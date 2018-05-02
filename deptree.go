package deptree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

//Resolver defined the methods of a dependency tree resolver.
type Resolver interface {
	Resolve(distributions ...string) (Distributions, error)
}

//Distribution or package is the representation of a distrubition.
// the name of the distribution its dependencies.
type Distribution struct {
	Name         string
	Dependencies Distributions
}

func (d Distribution) contains(dist *Distribution) bool {
	for _, dep := range d.Dependencies {
		if dep.Name == dist.Name {
			return true
		}
	}
	return false
}

//Distributions is an array of distribution that represent the tree of dependecies for distributions.
type Distributions []*Distribution

//ToJSON export the dependency tree to a JSON format
//with the indentantion ident. if the ident is empty, the
//JSON will be render in one line.
func (d Distributions) ToJSON(indent string) string {
	var buffer bytes.Buffer
	d.toJSON(&buffer, indent, 0)
	return buffer.String()
}

func (d Distributions) toJSON(dst *bytes.Buffer, indent string, depth int) {
	dst.WriteString("{")
	depth++
	newline(dst, indent, depth)
	for k, v := range d {
		dst.WriteString(fmt.Sprintf("\"%s\": ", v.Name))
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

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return fmt.Errorf("deptree: error finding the file %s , %s", path, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("deptree: error opening the file %s , %s", path, err)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("deptree: error decoding the the file %s, %s", file.Name(), err)
	}

	return nil
}
