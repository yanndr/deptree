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
	flag.Var(&names, "name", "Distribition names to resolve; you can define this flag multiple time.")
	flag.Usage = usage

	flag.Parse()

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

	fmt.Println(deps.ToJSON("  "))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -name distribution \n\n", os.Args[0])
	fmt.Fprint(os.Stderr, "This command displays the tree of dependency of one or multiple Perl distributions.\neg. deptree -name DateTime -name Specio\n\n")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")
}
