// +build unit

package deptree

import (
	"testing"
)

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
		{"not Present", "module3", "", DistributionNotFoundError{name: "modules3"}},
	}

	r := perlDepTreeResolver{
		distributionMap: map[string]string{"module1": "distrib1", "module2": "distrib2"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := r.getDistribution(tc.input)

			if err != tc.err {
				t.Errorf("expected error \"%v\", got error \"%v\"", tc.err, err)
			}

			if result != tc.output {
				t.Fatalf("error expected %v  got %v", tc.output, result)
			}
		})
	}
}
