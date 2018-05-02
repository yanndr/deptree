// +build unit

package deptree

import (
	"testing"
)

func TestToJSON(t *testing.T) {

	tt := []struct {
		name   string
		input  Distributions
		indent string
		output string
	}{
		{name: "simple", input: dt, output: `{"test": {}}`},
		{name: "with dependencies", input: dtwithDep, output: `{"test With dep": {"test": {}}}`},
		{name: "with multiple dependencies", input: multipleDep, output: `{"test With multiple dep": {"test": {},"test2": {}}}`},
		{name: "with multiple root", input: multipleRoot, output: `{"dist1": {},"dist2": {}}`},
		{name: "complex", input: complex, output: `{"dist1": {"test With multiple dep": {"test": {},"test2": {}}},"dist2": {"test With multiple dep": {"test": {},"test2": {}}}}`},
		{name: "with dependencies with indent", input: dtwithDep, output: "{\n\t\"test With dep\": {\n\t\t\"test\": {}\n\t}\n}", indent: "\t"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToJSON(tc.indent)
			if result != tc.output {
				t.Fatalf("error, expected \n%s, got \n%s", tc.output, result)
			}
		})
	}
}

var (
	dt = Distributions{
		&Distribution{
			Name:         "test",
			Dependencies: nil,
		},
	}

	dtwithDep = Distributions{
		&Distribution{
			Name:         "test With dep",
			Dependencies: dt,
		},
	}

	multipleDep = Distributions{
		&Distribution{
			Name: "test With multiple dep",
			Dependencies: Distributions{
				&Distribution{
					Name:         "test",
					Dependencies: nil,
				},
				&Distribution{
					Name:         "test2",
					Dependencies: nil,
				},
			},
		},
	}

	multipleRoot = Distributions{
		&Distribution{
			Name:         "dist1",
			Dependencies: nil,
		},
		&Distribution{
			Name:         "dist2",
			Dependencies: nil,
		},
	}

	complex = Distributions{
		&Distribution{
			Name:         "dist1",
			Dependencies: multipleDep,
		},
		&Distribution{
			Name:         "dist2",
			Dependencies: multipleDep,
		},
	}
)
