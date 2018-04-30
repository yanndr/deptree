package deptree

import (
	"io/ioutil"
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

	tt := []struct {
		name   string
		input  DependencyTree
		output string
	}{
		{name: "simple", input: dt, output: `{"test":{}}`},
		{name: "with dependencies", input: dtwithDep, output: `{"test With dep":{"test":{}}}`},
		{name: "with multiple dependencies", input: multipleDep, output: `{"test With multiple dep":{"test":{},"test2":{}}}`},
		{name: "with multiple root", input: multipleRoot, output: `{"dist1":{},"dist2":{}}`},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ToJSON()
			if result != tc.output {
				t.Fatalf("error, expected %s, got %s", tc.output, result)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	tt := []struct {
		name      string
		input     []string
		numberDep int
	}{
		{name: "Specio",
			input:     []string{"Specio"},
			numberDep: 8,
		},
	}

	dt, err := New("./cmd/deptree/data/")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := dt.Resolve(tc.input...)
			if err != nil {
				t.Fatal(err)
			}

			for _, v := range result {
				if len(v.Dependencies) != tc.numberDep {
					t.Fatalf("result length wrong, expected %v got %v", tc.numberDep, len(v.Dependencies))
				}
			}
		})
	}

}

func BenchmarkResolve(b *testing.B) {
	path := "./cmd/deptree/data/"
	dt, err := New(path)
	if err != nil {
		b.Fatal(err)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		for _, f := range files {
			if !f.IsDir() {
				continue
			}
			_, err := dt.Resolve(f.Name())
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
