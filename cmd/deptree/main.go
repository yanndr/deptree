// Command line interface to resolve Perl distributions dependency tree.
package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/yanndr/deptree"
)

type nameFlag []string

func (n *nameFlag) String() string {
	return fmt.Sprintf("%v", *n)
}

func (n *nameFlag) Set(value string) error {
	*n = append(*n, value)
	return nil
}

func main() {
	var (
		path  string
		names nameFlag
	)

	flag.StringVar(&path, "path", "./data", "The path to the CPAN folder.")
	flag.Var(&names, "name", "Distribition name to resolve; this flag is mandatory. You need to define it once but you can also define it multiple times.")
	flag.Usage = usage

	flag.Parse()

	if len(names) == 0 {
		fmt.Print("No Distribution name provided.\n\n")
		usage()
		os.Exit(2)
	}

	r, err := deptree.New(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	deps, err := r.Resolve(names...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(deps.ToJSON("\t"))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s  [OPTIONS] \n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "This command displays the dependency tree of one or multiple Perl distributions.\neg: %s -name DateTime -name Specio\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
}
