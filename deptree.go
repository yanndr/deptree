package deptree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

//Resolver defined the methods of a dependency tree resolver.
type Resolver interface {
	Resolve(distributions ...string) (DependencyTree, error)
}

type distribution struct {
	Name         string
	Dependencies DependencyTree
}

func (d distribution) contains(dist *distribution) bool {
	for _, dep := range d.Dependencies {
		if dep.Name == dist.Name {
			return true
		}
	}
	return false
}

//DependencyTree is the tree of dependecies for distributions.
type DependencyTree []*distribution

type distributionNotFoundError struct {
	name string
}

func (e distributionNotFoundError) Error() string {
	return fmt.Sprintf("distribution %s not found:", e.name)
}

//ToJSON export the dependency tree to a JSON format
//with the indentantion ident. if the ident is empty, the
//JSON will be render in one line.
func (d DependencyTree) ToJSON(indent string) string {
	var buffer bytes.Buffer
	d.toJSON(&buffer, indent, 0)
	return buffer.String()
}

func (d DependencyTree) toJSON(dst *bytes.Buffer, indent string, depth int) {
	dst.WriteString("{")
	depth++
	newline(dst, indent, depth)
	for k, v := range d {
		dst.WriteString(fmt.Sprintf("\"%s\":", v.Name))
		if v.Dependencies != nil && len(v.Dependencies) > 0 {
			v.Dependencies.toJSON(dst, indent, depth)
		} else {
			dst.WriteString("{}")
		}
		if k < len(d)-1 {
			dst.WriteRune(',')
			newline(dst, indent, depth)
		}
	}
	newline(dst, indent, depth-1)
	dst.WriteString("}")
}

func newline(dst *bytes.Buffer, indent string, depth int) {
	if indent == "" {
		return
	}
	dst.WriteByte('\n')
	for i := 0; i < depth; i++ {
		dst.WriteString(indent)
	}
}

func decodeJSONFromFile(v interface{}, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening the file %s , %s", path, err)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding the the file %s, %s", file.Name(), err)
	}

	return nil
}
