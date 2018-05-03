package deptree

import (
	"errors"
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
	} else if path == "/distrib3/META.json" {
		return []byte(`{
			"prereqs": {
			  "runtime": {
				"requires": {
					"module4":""
				}
			  }
			}
		  }
		  `), nil
	}
	return []byte("{}"), nil
}

var distrib2 = Distributions{"distrib2": nil}

var distrib1 = Distributions{"distrib1": distrib2}

func TestResolve(t *testing.T) {
	tt := []struct {
		name   string
		input  []string
		output Distributions
		err    error
	}{
		{"normal", []string{"distrib1"}, distrib1, nil},
		{"error module", []string{"distrib3"}, nil, ModuleNotFoundError{"module4", errors.New("module module4 not present on distribution map module-distro-map.json")}},
	}

	r := perlDepTreeResolver{
		coreModules:     []string{"module3"},
		distributionMap: map[string]string{"module1": "distrib1", "module2": "distrib2"},
		cache:           make(Distributions),
		readFileFunc:    fakeReadFile,
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := r.Resolve(tc.input...)
			handleError(t, err, tc.err)

			// if res.ToJSON("") != tc.output.ToJSON("") {
			// 	t.Fatalf("error got %s, expected %s", res.ToJSON(""), tc.output.ToJSON(""))
			// }
		})
	}

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

			handleError(t, err, tc.err)

			if result != tc.output {
				t.Fatalf("error expected %v  got %v", tc.output, result)
			}
		})
	}
}

func handleError(t *testing.T, err, expectedError error) {
	if err != nil && expectedError == nil {
		t.Errorf("expected no error, got error \"%v\"", err)
	} else if expectedError != nil && err == nil {
		t.Errorf("expected error \"%v\", got no error", expectedError)
	} else if expectedError != nil && err.Error() != expectedError.Error() {
		t.Errorf("expected error \"%v\", got error \"%v\"", expectedError, err)
	}
}
