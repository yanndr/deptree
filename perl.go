package deptree

import (
	"fmt"
	"sort"
)

const (
	distroMapFile   = "module-distro-map.json"
	coreModulesFile = "core-modules.json"
	metaJSONFile    = "META.json"
)

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
