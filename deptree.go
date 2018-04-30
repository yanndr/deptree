package deptree

import (
	"fmt"
	"sort"

	"bitbucket.org/yanndr/deptree/json"
)

const (
	distroMapFile   = "module-distro-map.json"
	coreModulesFile = "core-modules.json"
	metaJSONFile    = "META.json"
)

//Resolver defined the methods of a dependency tree resolver.
type Resolver interface {
	Resolve(distributions ...string) ([]*Distribution, error)
}

type Distribution struct {
	Name         string          `json:",omitempty"`
	Dependencies []*Distribution `json:",omitempty"`
}

type perlDepTreeResolver struct {
	path            string
	distributionMap map[string]string
	coreModules     []string
}

//New returns an instance of a perl dependency tree resolver.
func New(path string) (Resolver, error) {
	dt := &perlDepTreeResolver{
		path: path,
	}
	distroMapPath := fmt.Sprintf("%s/%s", path, distroMapFile)
	err := json.DecodeFromFile(&dt.distributionMap, distroMapPath)
	if err != nil {
		return nil, fmt.Errorf("error decoding the json file %s, %s", distroMapPath, err)
	}

	coreModulesPath := fmt.Sprintf("%s/%s", path, coreModulesFile)
	err = json.DecodeFromFile(&dt.coreModules, coreModulesPath)
	if err != nil {
		return nil, fmt.Errorf("error decoding the json file %s, %s", coreModulesPath, err)
	}
	return dt, nil
}

func (r *perlDepTreeResolver) Resolve(distributions ...string) ([]*Distribution, error) {
	var result []*Distribution
	for _, d := range distributions {

		distro := &Distribution{Name: d}
		result = append(result, distro)

		dependencies, err := r.getDependencies(d)
		if err != nil {
			return nil, err
		}

		for _, dep := range dependencies {

			deps, err := r.Resolve(dep)
			for _, dep := range deps {
				distro.Dependencies = append(distro.Dependencies, dep)
			}

			if err != nil {
				return nil, err
			}

		}
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

func (dt *perlDepTreeResolver) getRequiresModules(dist string) (map[string]interface{}, error) {
	meta := &struct {
		Prereqs struct {
			Runtime struct {
				Requires map[string]interface{} `json:"requires"`
			} `json:"runtime"`
		} `json:"prereqs"`
	}{}

	path := fmt.Sprintf("%s/%s/%s", dt.path, dist, metaJSONFile)
	err := json.DecodeFromFile(meta, path)
	if err != nil {
		return nil, err
	}

	return meta.Prereqs.Runtime.Requires, nil
}
