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
	flag.Var(&names, "name", "Distribition names to resolve; you can define this flag multile time.")
	flag.Parse()

	fmt.Printf("%s \n", names)

	_, err := deptree.New(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//fmt.Println(dt)
}
