package deptree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	readFileFunc    func(string) ([]byte, error)
}

// New returns an instance of a Perl dependency tree resolver.
func New(path string) (Resolver, error) {
	r := &perlDepTreeResolver{
		path:         path,
		readFileFunc: ioutil.ReadFile,
		cache:        make(map[string]*Distribution),
	}

	distroMapPath := fmt.Sprintf("%s/%s", path, distroMapFile)
	data, err := r.readFileFunc(distroMapPath)
	if err != nil {
		return nil, fmt.Errorf("deptree: error reading the json file %s, %s", distroMapPath, err)
	}

	err = json.Unmarshal(data, &r.distributionMap)
	if err != nil {
		return nil, fmt.Errorf("deptree: error decoding the json file %s, %s", distroMapPath, err)
	}

	coreModulesPath := fmt.Sprintf("%s/%s", path, coreModulesFile)
	data, err = r.readFileFunc(coreModulesPath)
	if err != nil {
		return nil, fmt.Errorf("deptree: error reading the json file %s, %s", coreModulesPath, err)
	}

	err = json.Unmarshal(data, &r.coreModules)
	if err != nil {
		return nil, fmt.Errorf("deptree: error decoding the json file %s, %s", coreModulesPath, err)
	}

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
			return nil, err
		}

		deps, err := r.Resolve(dependencies...)
		if err != nil {
			return nil, err
		}

		dist.addDependencies(deps...)

		r.cache[d] = dist
	}

	return result, nil
}

// getDependencies returns the list of distributions required for a distribution.
// The function will ignore modules present in the core modules.
func (r *perlDepTreeResolver) getDependencies(distribution string) ([]string, error) {
	moduleMap, err := r.getRequiredModules(distribution)
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

// getDistribution returns the distribution name of a module.
func (r *perlDepTreeResolver) getDistribution(module string) (string, error) {
	if val, ok := r.distributionMap[module]; ok {
		return val, nil
	}
	return "", ModuleNotFoundError{module, fmt.Errorf("module %s not present on distribution map %s", module, distroMapFile)}
}

// getRequiredModules returns a map of required modules/version for a distribution.
func (r *perlDepTreeResolver) getRequiredModules(dist string) (map[string]string, error) {
	meta := &struct {
		Prereqs struct {
			Runtime struct {
				Requires map[string]string `json:"requires"`
			} `json:"runtime"`
		} `json:"prereqs"`
	}{}

	path := fmt.Sprintf("%s/%s/%s", r.path, dist, metaJSONFile)
	data, err := r.readFileFunc(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, DistributionNotFoundError{dist, err}
		}
		return nil, fmt.Errorf("deptree: could not open the file %s, %v", path, err)
	}

	err = json.Unmarshal(data, meta)
	if err != nil {
		return nil, fmt.Errorf("deptree: error decoding the json file %s, %s", path, err)
	}

	return meta.Prereqs.Runtime.Requires, nil
}

// filterCoreModules returns the list of modules passed to the method
// filtered by the modules present in the core modules.
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
