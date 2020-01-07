# prettylist

a tiny unix pipe utitlity that reads in newline delimited items and organizes those items in lines that do not exceed a configurable amount of characters

## installation

requires Golang

```
$ go install -u -v github.com/nkcmr/prettylist
```

## usage

i deal a lot with CSV files so if I am trying to prepare a large list of things to search for in SQL it usually looks like this:

```
$ cat data.csv | \
	q 'select id_column from - where thing = "foo"' | \
	xsv sample 10 | \
	tail -n+2 | \
	prettylist -quote "'" -maxlen 120 -nspaces 4
```

(`q` is https://github.com/harelba/q, `xsv` is https://github.com/BurntSushi/xsv)

and the output from `prettylist` should look like this:

```
    '493f8050dd7849d581d243e744297bff', '8ac65ae3a7774633afcb47d7a8a14636',
    '9d2e61a1cca74ffcb118a2d6ad36df0f', '68054d0e4c0245a3bb171829a3a63a37',
    '20094123d6094d66bf29a24cc134285d', '10724a73a1d8461caddb84729bb5bc2b',
    '0e30c31ace4847038daa8a33333b32aa', '2e7e94795f2c4e4a8fdfd9e1baf4bd3e',
    'b196843ea1f04f3993841e7e1da37aee', '8220daacac4e4347b9463ec452a4e2ab'
```

ready to be pasted into an `IN (...)` for a query!

## license

Copyright 2020 Nicholas Comer

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.