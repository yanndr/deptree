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
	err := decodeJSONFromFile(&r.distributionMap, distroMapPath)
	if err != nil {
		return nil, fmt.Errorf("error decoding the json file %s, %s", distroMapPath, err)
	}

	coreModulesPath := fmt.Sprintf("%s/%s", path, coreModulesFile)
	err = decodeJSONFromFile(&r.coreModules, coreModulesPath)
	if err != nil {
		return nil, fmt.Errorf("error decoding the json file %s, %s", coreModulesPath, err)
	}
	r.cache = make(map[string]*distribution)
	return r, nil
}

func (r *perlDepTreeResolver) Resolve(distributions ...string) (DependencyTree, error) {
	var result DependencyTree
	for _, d := range distributions {
		if v, ok := r.cache[d]; ok {
			return append(result, v), nil
		}

		dist := &distribution{Name: d}
		result = append(result, dist)

		dependencies, err := r.getDependencies(d)
		if err != nil {
			return nil, err
		}

		for _, dep := range dependencies {
			deps, err := r.Resolve(dep)
			for _, dep := range deps {
				if !dist.contains(dep) {
					dist.Dependencies = append(dist.Dependencies, dep)
				}
			}

			if err != nil {
				return nil, err
			}
		}
		r.cache[d] = dist
	}

	return result, nil
}

//getDependencies returns the list of distributions requires for a distribution.
// the function will ignore modules present in the core modules.
func (r *perlDepTreeResolver) getDependencies(distribution string) ([]string, error) {
	moduleMap, err := r.getRequiresModules(distribution)
	if err != nil {
		return nil, err
	}
	modules := r.filterCoreModules(moduleMap)

	return r.getDistributions(modules)
}

func (r *perlDepTreeResolver) getDistributions(modules []string) ([]string, error) {
	var distributions []string
	for _, m := range modules {
		if val, ok := r.distributionMap[m]; ok {
			distributions = append(distributions, val)
		} else {
			return nil, distributionNotFoundError{name: m}
		}
	}
	return distributions, nil
}

//getRequiresModules returns a map of requires modules/version for a distribution.
func (r *perlDepTreeResolver) getRequiresModules(dist string) (map[string]string, error) {
	meta := &struct {
		Prereqs struct {
			Runtime struct {
				Requires map[string]string `json:"requires"`
			} `json:"runtime"`
		} `json:"prereqs"`
	}{}

	path := fmt.Sprintf("%s/%s/%s", r.path, dist, metaJSONFile)
	err := decodeJSONFromFile(meta, path)
	if err != nil {
		return nil, fmt.Errorf("get modules error: could not decode the file %s, %v", path, err)
	}
	return meta.Prereqs.Runtime.Requires, nil
}

func (r *perlDepTreeResolver) filterCoreModules(modules map[string]string) []string {
	var result []string
	for m := range modules {
		if m == "perl" {
			continue
		}
		i := sort.SearchStrings(r.coreModules, m)
		if i < len(r.coreModules) && r.coreModules[i] == m {
			continue
		}

		result = append(result, m)
	}

	return result
}
