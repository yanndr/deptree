package deptree

import (
	"fmt"

	"bitbucket.org/yanndr/deptree/json"
)

const (
	distroMapFile   = "module-distro-map.json"
	coreModulesFile = "core-modules.json"
	metaJSONFile    = "META.json"
)

//Resolver defined the methods of a dependency tree resolver.
type Resolver interface {
	Resolve()
}

type perlDepTreeResolver struct {
	path            string
	distributionMap map[string]string
	coreModules     []string
}

func (*perlDepTreeResolver) Resolve() {}

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
