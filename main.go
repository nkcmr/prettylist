package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func _main(args []string) error {
	var f struct {
		linLen, nleadspace int
		qchar              string
	}
	fs := flag.NewFlagSet("prettylist", flag.ContinueOnError)
	fs.IntVar(&f.linLen, "maxlen", 80, "sets the maximum line length")
	fs.StringVar(&f.qchar, "quote", "", "defines a surrounding quote character for each item")
	fs.IntVar(&f.nleadspace, "nspaces", 0, "set the number of space characters each line should have prepended")
	if err := fs.Parse(args[1:]); err != nil {
		return errors.Wrap(err, "failed to parse arguments")
	}
	if f.nleadspace >= f.linLen {
		return errors.New("too many leading spaces")
	}
	s := bufio.NewScanner(os.Stdin)
	buf := strings.Repeat(" ", f.nleadspace)
	first := true
	for s.Scan() {
		item := s.Text()
		prefix := ""
		if !first {
			prefix = ", "
		} else {
			first = false
		}
		// if the next addition is going to overflow the linLen, get a new line
		if len(buf)+len(prefix)+(len(f.qchar)*2)+len(item) >= f.linLen {
			fmt.Fprintf(os.Stdout, "%s,\n", buf)
			os.Stdout.Sync()
			buf = strings.Repeat(" ", f.nleadspace)
			prefix = ""
			first = false
		}
		buf += prefix + f.qchar + item + f.qchar
	}
	if len(buf) > f.nleadspace {
		fmt.Fprintln(os.Stdout, buf)
	}
	return errors.Wrap(s.Err(), "failed to read stdin")
}

func main() {
	if err := _main(os.Args); err != nil {
		log.Printf("%s: error: %s", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}
