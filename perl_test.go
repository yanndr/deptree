package deptree

import (
	"fmt"
	"testing"
)

func fakeReadFile(path string) ([]byte, error) {

	if path == "/distrib1/META.json" {
		return []byte(`{
			"prereqs": {
			  "runtime": {
				"requires": {
				  "module2": ""
				}
			  }
			}
		  }
		  `), nil
	} else if path == "/distrib2/META.json" {
		return []byte(`{
			"prereqs": {
			  "runtime": {
				"requires": {
				}
			  }
			}
		  }
		  `), nil
	}
	return nil, nil
}
func TestResolve(t *testing.T) {
	r := perlDepTreeResolver{
		coreModules:     []string{"module1", "module2"},
		distributionMap: map[string]string{"module1": "distrib1", "module2": "distrib2"},
		cache:           make(map[string]*Distribution),
		readFileFunc:    fakeReadFile,
	}

	res, err := r.Resolve("distrib1")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	fmt.Println(len(res))
}

func TestFilterCoreModules(t *testing.T) {
	tt := []struct {
		name   string
		input  map[string]string
		output []string
	}{
		{"Present", map[string]string{"module1": ""}, nil},
		{"Not Present", map[string]string{"module3": ""}, []string{"module3"}},
		{"One present", map[string]string{"test": "", "module1": ""}, []string{"test"}},
		{"Perl module", map[string]string{"perl": "", "module3": ""}, []string{"module3"}},
	}

	r := perlDepTreeResolver{
		coreModules: []string{"module1", "module2"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := r.filterCoreModules(tc.input)

			if len(result) != len(tc.output) {
				t.Fatalf("error expected %v result(s) got %v result(s)", len(tc.output), len(result))
			}
			for k, v := range result {
				if v != tc.output[k] {
					t.Fatalf("error expected %v  got %v", tc.output[k], v)
				}
			}
		})
	}
}

func TestGetDistribution(t *testing.T) {
	tt := []struct {
		name   string
		input  string
		output string
		err    error
	}{
		{"Present", "module1", "distrib1", nil},
		{"not Present", "module3", "", ModuleNotFoundError{"module3", fmt.Errorf("module %s not present on distribution map %s", "module3", distroMapFile)}},
	}

	r := perlDepTreeResolver{
		distributionMap: map[string]string{"module1": "distrib1", "module2": "distrib2"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := r.getDistribution(tc.input)

			if err != nil && tc.err == nil {
				t.Errorf("expected no error, got error \"%v\"", err)
			} else if tc.err != nil && err == nil {
				t.Errorf("expected error \"%v\", got no error", tc.err)
			} else if tc.err != nil && err.Error() != tc.err.Error() {
				t.Errorf("expected error \"%v\", got error \"%v\"", tc.err, err)
			}

			if result != tc.output {
				t.Fatalf("error expected %v  got %v", tc.output, result)
			}
		})
	}
}
