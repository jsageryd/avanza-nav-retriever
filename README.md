# Avanza NAV retriever

Automatically fetches the latest price for an Avanza instrument and appends it
to a ledger-compatible price-db file, then commits and pushes the update.

This is a rewrite of https://github.com/jsageryd/avanza_zero_nav_retriever.

## Installation
If you have $GOPATH set:

    go get github.com/jsageryd/avanza-nav-retriever

If you don't care about Go and just want to build the thing:

    git clone https://github.com/jsageryd/avanza-nav-retriever.git
    cd avanza-nav-retriever
    go build

If you prefer to just grab a pre-built binary:

http://gobuild.io/github.com/jsageryd/avanza-nav-retriever

## Usage example
Fetch the latest price for Avanza Zero and store it in `$HOME/foo/`

    avanza-nav-retriever \
      --repo "$HOME/foo/" \
      --pricedb 'price-db' \
      --url 'https://www.avanza.se/fonder/om-fonden.html/41567/avanza-zero' \
      --ident 'ZERO'

## Licence
Copyright (c) 2014 Johan Sageryd <j@1616.se>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
