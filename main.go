package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

var (
	repo    string
	url     string
	ident   string
	pricedb string
)

func main() {
	initConf()
	updatePrice()
}

func initConf() {
	flag.StringVar(&repo, "repo", "", "path to Git workdir containing price-db file")
	flag.StringVar(&url, "url", "", "URL to Avanza page from which to fetch data")
	flag.StringVar(&ident, "ident", "", "price identifier to use in price-db")
	flag.StringVar(&pricedb, "pricedb", "price-db", "the name of the price-db file")
	flag.Parse()

	repo = strings.TrimSpace(repo)
	url = strings.TrimSpace(url)
	ident = strings.TrimSpace(ident)
	pricedb = strings.TrimSpace(pricedb)

	if repo == "" {
		log.Fatal("Need repo, something like '/path/to/valid/git/repo/'")
	}
	repo = strings.TrimRight(repo, "/") + "/"
	if url == "" {
		log.Fatal("Need URL, something like 'https://www.avanza.se/fonder/om-fonden.html/41567/avanza-zero'")
	}
	if ident == "" {
		log.Fatal("Need ident, something like 'ZERO'")
	}
	if pricedb == "" {
		log.Fatal("Need ident, something like 'price-db'")
	}
}

func updatePrice() {
	os.MkdirAll(repo, 0755)
	file, err := os.OpenFile(
		repo+pricedb,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	theirs := theirPrice()
	ours := ourPrice(file)

	if ours != theirs {
		addPrice(theirs, file)
	}
}

type price struct {
	date     string
	ident    string
	amount   string
	currency string
}

func (p price) String() string {
	return fmt.Sprintf("P %s %s %s %s", p.date, p.ident, p.amount, p.currency)
}

func addPrice(price price, file *os.File) {
	file.WriteString(fmt.Sprintln(price))

	var cmd *exec.Cmd

	cmd = exec.Command("git", "init")
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		log.Fatalf("git init: %v", err)
	}

	cmd = exec.Command("git", "add", pricedb)
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		log.Fatalf("git add: %v", err)
	}

	commitMessage := fmt.Sprintf("Update %s", price.date)
	cmd = exec.Command("git", "commit", "-m", commitMessage)
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		log.Fatalf("git commit: %v", err)
	}

	cmd = exec.Command("git", "push")
	cmd.Dir = repo
	if err := cmd.Run(); err != nil {
		log.Fatalf("git push: %v", err)
	}
}

func ourPrice(file *os.File) price {
	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}
	return stringToPrice(line)
}

func stringToPrice(line string) price {
	s := strings.Split(line, " ")
	if line == "" || len(s) != 5 {
		return price{}
	}
	return price{
		date:     s[1],
		ident:    s[2],
		amount:   s[3],
		currency: s[4],
	}
}

func theirPrice() price {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	productNode := nodeByAttr("itemtype", "http://schema.org/Product", doc)

	// Date
	reviewNode := nodeByAttr("itemtype", "http://schema.org/Review", productNode)
	dateNode := nodeByAttr("itemprop", "datePublished", reviewNode)
	date := strings.TrimSpace(dateNode.FirstChild.Data)
	if date == "" {
		log.Fatal("Cannot find their date")
	}

	// Amount and currency
	offerNode := nodeByAttr("itemtype", "http://schema.org/Offer", productNode)
	priceNode := nodeByAttr("itemprop", "price", offerNode)
	currencyNode := nodeByAttr("itemprop", "priceCurrency", offerNode)
	amount := strings.TrimSpace(attrValue("content", priceNode))
	if amount == "" {
		log.Fatal("Cannot find their amount")
	}
	currency := strings.TrimSpace(attrValue("content", currencyNode))
	if amount == "" {
		log.Fatal("Cannot find their currency")
	}
	return price{
		date:     date,
		ident:    ident,
		amount:   amount,
		currency: currency,
	}
}

func nodeByAttr(attrKey string, attrVal string, node *html.Node) *html.Node {
	if node.Type == html.ElementNode {
		if hasAttr(attrKey, attrVal, node) {
			return node
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res := nodeByAttr(attrKey, attrVal, c)
		if res != nil {
			return res
		}
	}
	return nil
}

func hasAttr(attrKey string, attrVal string, node *html.Node) bool {
	for _, attr := range node.Attr {
		if (attr.Key == attrKey) && (attr.Val == attrVal) {
			return true
		}
	}
	return false
}

func attrValue(attrKey string, node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == attrKey {
			return attr.Val
		}
	}
	return ""
}
