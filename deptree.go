package deptree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	distroMapFile   = "module-distro-map.json"
	coreModulesFile = "core-modules.json"
	metaJSONFile    = "META.json"
)

//Resolver defined the methods of a dependency tree resolver.
type Resolver interface {
	Resolve(distributions ...string) (DependencyTree, error)
}

type DependencyTree []*distribution

type distribution struct {
	Name         string         `json:",omitempty"`
	Dependencies DependencyTree `json:",omitempty"`
}

func (d DependencyTree) ToJSON() string {

	var buffer bytes.Buffer
	buffer.WriteString("{")
	for _, v := range d {
		buffer.WriteString(fmt.Sprintf("\"%s\":", v.Name))
		if v.Dependencies == nil || len(v.Dependencies) == 0 {
			buffer.WriteString("{}")
			break
		}
		buffer.WriteString(v.Dependencies.ToJSON())
	}
	buffer.WriteString("}")
	return buffer.String()
}

func (d distribution) contains(dist *distribution) bool {
	for _, dep := range d.Dependencies {
		if dep.Name == dist.Name {
			return true
		}
	}
	return false
}

type perlDepTreeResolver struct {
	path            string
	distributionMap map[string]string
	coreModules     []string
	cache           map[string]*distribution
}

//New returns an instance of a perl dependency tree resolver.
func New(path string) (Resolver, error) {
	r := &perlDepTreeResolver{
		path: path,
	}
	distroMapPath := fmt.Sprintf("%s/%s", path, distroMapFile)
	err := decodeFromFile(&r.distributionMap, distroMapPath)
	if err != nil {
		return nil, fmt.Errorf("error decoding the json file %s, %s", distroMapPath, err)
	}

	coreModulesPath := fmt.Sprintf("%s/%s", path, coreModulesFile)
	err = decodeFromFile(&r.coreModules, coreModulesPath)
	if err != nil {
		return nil, fmt.Errorf("error decoding the json file %s, %s", coreModulesPath, err)
	}
	r.cache = make(map[string]*distribution)
	return r, nil
}

func (r *perlDepTreeResolver) Resolve(distributions ...string) (DependencyTree, error) {
	var result []*distribution
	for _, d := range distributions {

		distro := &distribution{Name: d}

		if deps, ok := r.cache[d]; ok {
			distro = deps
			return result, nil
		}

		result = append(result, distro)

		dependencies, err := r.getDependencies(d)
		if err != nil {
			return nil, err
		}

		for _, dep := range dependencies {
			deps, err := r.Resolve(dep)
			for _, dep := range deps {
				if !distro.contains(dep) {
					distro.Dependencies = append(distro.Dependencies, dep)
				}
			}

			if err != nil {
				return nil, err
			}
		}
		r.cache[d] = distro
	}

	return result, nil
}

func (r *perlDepTreeResolver) getDependencies(dist string) ([]string, error) {
	var dependencies []string
	modules, err := r.getRequiresModules(dist)
	if err != nil {
		return nil, err
	}
	for m := range modules {
		i := sort.SearchStrings(r.coreModules, m)
		if i == len(r.coreModules) {
			continue
		}
		if val, ok := r.distributionMap[m]; ok {
			dependencies = append(dependencies, val)
		}
	}
	return dependencies, nil
}

func (r *perlDepTreeResolver) getRequiresModules(dist string) (map[string]string, error) {
	meta := &struct {
		Prereqs struct {
			Runtime struct {
				Requires map[string]string `json:"requires"`
			} `json:"runtime"`
		} `json:"prereqs"`
	}{}

	path := fmt.Sprintf("%s/%s/%s", r.path, dist, metaJSONFile)
	err := decodeFromFile(meta, path)
	if err != nil {
		return nil, err
	}

	return meta.Prereqs.Runtime.Requires, nil
}

func decodeFromFile(v interface{}, path string) error {
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
