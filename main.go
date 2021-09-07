package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: error: "+format+"\n", append([]interface{}{filepath.Base(os.Args[0])}, a...))
	os.Exit(1)
}

func rootCommand() *cobra.Command {
	var f struct {
		linLen, nleadspace int
		qchar              string

		doubleQuote, singleQuote bool
	}
	cmd := &cobra.Command{
		Use:   "prettylist",
		Short: "take each line over STDIN and format into a comma-delimited list",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			if f.qchar == "" {
				if f.doubleQuote {
					f.qchar = `"`
				} else if f.singleQuote {
					f.qchar = `'`
				}
			}

			if f.nleadspace >= f.linLen {
				fatal("too many leading spaces")
				return
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
			if err := s.Err(); err != nil {
				fatal("failed to read stdin: %s", err.Error())
			}
		},
	}
	cmd.Flags().IntVar(&f.linLen, "maxlen", 80, "sets the max line length")
	cmd.Flags().StringVar(&f.qchar, "quote", "", "defines a surrounding quote character for each item")
	cmd.Flags().BoolVar(&f.doubleQuote, "doublequote", false, `shorthand for --quote '"'`)
	cmd.Flags().BoolVar(&f.singleQuote, "singlequote", false, `shorthand for --quote "'"`)
	cmd.Flags().IntVar(&f.nleadspace, "nspaces", 0, "set the number of space characters each line should have prefixed")
	return cmd
}

func main() {
	if err := rootCommand().Execute(); err != nil {
		fatal(err.Error())
	}
}
