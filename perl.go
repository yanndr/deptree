package deptree

import (
	"fmt"
	"os"
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
	cache           map[string]*Distribution
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
	r.cache = make(map[string]*Distribution)
	return r, nil
}

func (r *perlDepTreeResolver) Resolve(distributions ...string) (Distributions, error) {
	var result Distributions
	for _, d := range distributions {

		if v, ok := r.cache[d]; ok {
			result = append(result, v)
			continue
		}

		dist := &Distribution{Name: d}
		result = append(result, dist)
		dependencies, err := r.getDependencies(d)
		if err != nil {
			if _, ok := err.(DistributionNotFoundError); ok {
				return nil, err
			}
			return nil, fmt.Errorf("resolve: can't get dependencies of %s: %v", d, err)
		}

		deps, err := r.Resolve(dependencies...)
		if err != nil {
			return nil, err
		}

		for _, dep := range deps {
			if !dist.contains(dep) {
				dist.Dependencies = append(dist.Dependencies, dep)
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

	var distributions []string
	for _, m := range modules {
		d, err := r.getDistribution(m)
		if err != nil {
			return nil, err
		}
		distributions = append(distributions, d)
	}
	return distributions, nil
}

//getDistribution returns the distribution name of a module.
func (r *perlDepTreeResolver) getDistribution(module string) (string, error) {

	if val, ok := r.distributionMap[module]; ok {
		return val, nil
	}
	return "", DistributionNotFoundError{name: "modules3"}
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
		if os.IsNotExist(err) {
			return nil, DistributionNotFoundError{dist, err}
		}
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
