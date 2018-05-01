package deptree

import (
	"testing"
)

func TestToJSON(t *testing.T) {
	dt := Distributions{
		&Distribution{
			Name:         "test",
			Dependencies: nil,
		},
	}

	dtwithDep := Distributions{
		&Distribution{
			Name:         "test With dep",
			Dependencies: dt,
		},
	}

	multipleDep := Distributions{
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

	multipleRoot := Distributions{
		&Distribution{
			Name:         "dist1",
			Dependencies: nil,
		},
		&Distribution{
			Name:         "dist2",
			Dependencies: nil,
		},
	}

	complex := Distributions{
		&Distribution{
			Name:         "dist1",
			Dependencies: multipleDep,
		},
		&Distribution{
			Name:         "dist2",
			Dependencies: multipleDep,
		},
	}

	tt := []struct {
		name   string
		input  Distributions
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
