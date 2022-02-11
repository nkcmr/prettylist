package main // import "code.nkcmr.net/prettylist"

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/spf13/cobra"
)

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s: error: "+format+"\n", append([]interface{}{filepath.Base(os.Args[0])}, a...))
	os.Exit(1)
}

type countWriter int64

func (c *countWriter) Write(b []byte) (int, error) {
	n, err := os.Stdout.Write(b)
	atomic.AddInt64((*int64)(c), int64(n))
	return n, err
}

func (c *countWriter) Sync() error {
	return os.Stdout.Sync()
}

func (c *countWriter) Reset() {
	atomic.StoreInt64((*int64)(c), 0)
}

func (c *countWriter) Len() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func rootCommand() *cobra.Command {
	var f struct {
		linLen, nleadspace       int
		qchar                    string
		delim                    string
		oneline                  bool
		doubleQuote, singleQuote bool
	}
	cmd := &cobra.Command{
		Use:   "prettylist",
		Short: "take each line over STDIN and format into a comma-delimited list",
		Args:  cobra.NoArgs,
		PreRun: func(*cobra.Command, []string) {
			if f.oneline {
				f.linLen = math.MaxInt
			}
			if f.qchar == "" {
				if f.doubleQuote {
					f.qchar = `"`
				} else if f.singleQuote {
					f.qchar = `'`
				}
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			if f.nleadspace >= f.linLen {
				fatal("too many leading spaces")
				return
			}
			s := bufio.NewScanner(os.Stdin)
			cw := new(countWriter)
			io.WriteString(cw, strings.Repeat(" ", f.nleadspace))
			first := true
			for s.Scan() {
				item := s.Text()
				prefix := ""
				if !first {
					prefix = f.delim
				} else {
					first = false
				}
				// if the next addition is going to overflow the linLen, get a new line
				if cw.Len()+int64(len(prefix))+int64(len(f.qchar)*2)+int64(len(item)) >= int64(f.linLen) {
					fmt.Fprintln(cw, f.delim)
					cw.Sync()
					cw.Reset()
					io.WriteString(cw, strings.Repeat(" ", f.nleadspace))
					prefix = ""
					first = false
				}
				io.WriteString(cw, prefix)
				io.WriteString(cw, f.qchar)
				io.WriteString(cw, item)
				io.WriteString(cw, f.qchar)
			}
			if cw.Len() > int64(f.nleadspace) {
				cw.Sync()
			}
			if err := s.Err(); err != nil {
				fatal("failed to read stdin: %s", err.Error())
			}
		},
	}
	cmd.Flags().StringVarP(&f.delim, "delimiter", "d", ", ", "the string that should separate each item")
	cmd.Flags().BoolVar(&f.oneline, "oneline", false, "puts all input on one line")
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
