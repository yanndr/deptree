package deptree

import (
	"testing"
)

func TestToJSON(t *testing.T) {
	dt := DependencyTree{
		&distribution{
			Name:         "test",
			Dependencies: nil,
		},
	}

	dtwithDep := DependencyTree{
		&distribution{
			Name:         "test With dep",
			Dependencies: dt,
		},
	}

	multipleDep := DependencyTree{
		&distribution{
			Name: "test With multiple dep",
			Dependencies: DependencyTree{
				&distribution{
					Name:         "test",
					Dependencies: nil,
				},
				&distribution{
					Name:         "test2",
					Dependencies: nil,
				},
			},
		},
	}

	multipleRoot := DependencyTree{
		&distribution{
			Name:         "dist1",
			Dependencies: nil,
		},
		&distribution{
			Name:         "dist2",
			Dependencies: nil,
		},
	}

	complex := DependencyTree{
		&distribution{
			Name:         "dist1",
			Dependencies: multipleDep,
		},
		&distribution{
			Name:         "dist2",
			Dependencies: multipleDep,
		},
	}

	tt := []struct {
		name   string
		input  DependencyTree
		output string
	}{
		{name: "simple", input: dt, output: `{"test":{}}`},
		{name: "with dependencies", input: dtwithDep, output: `{"test With dep":{"test":{}}}`},
		{name: "with multiple dependencies", input: multipleDep, output: `{"test With multiple dep":{"test":{},"test2":{}}}`},
		{name: "with multiple root", input: multipleRoot, output: `{"dist1":{},"dist2":{}}`},
		{name: "complex", input: complex, output: `{"dist1":{"test With multiple dep":{"test":{},"test2":{}}},"dist2":{"test With multiple dep":{"test":{},"test2":{}}}}`},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToJSON("")
			if result != tc.output {
				t.Fatalf("error, expected %s, got %s", tc.output, result)
			}
		})
	}
}
