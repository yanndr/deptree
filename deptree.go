//Package deptree is a package that resolve dependecy tree for package/distributions
//
// So far only a perl implementation is available.
package deptree

import (
	"bytes"
	"fmt"
	"sort"
)

//Resolver defined the methods of a dependency tree resolver.
type Resolver interface {
	//Resolve returns the distribution list with their dependencies.
	Resolve(distributions ...string) (Distributions, error)
}

//Distribution is the representation of a distribution.
// It contains the name of the distribution and a list of its dependencies.
type Distribution struct {
	Name         string
	Dependencies Distributions
}

//addDependencies add a list of dependencies to the distribution
//it inserts the dependencies ordered by name for a more efficient search/insert.
func (d *Distribution) addDependencies(distributions ...*Distribution) {
	for _, dis := range distributions {

		i := sort.Search(len(d.Dependencies), func(i int) bool { return d.Dependencies[i].Name >= dis.Name })

		if i == len(d.Dependencies) { //not found and it is the biggest value
			d.Dependencies = append(d.Dependencies, dis)
			continue
		}
		if d.Dependencies[i].Name == dis.Name { // already in the list
			continue
		}

		d.Dependencies = append(d.Dependencies, &Distribution{})
		copy(d.Dependencies[i+1:], d.Dependencies[i:])
		d.Dependencies[i] = dis
	}
}

//Distributions is an array of distribution that represent the tree of dependecies for distributions.
type Distributions []*Distribution

//ToJSON export the dependency tree to a JSON format
//with the indentantion passed as parameter ident.
//if the ident is empty, the JSON will be render in one line.
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
