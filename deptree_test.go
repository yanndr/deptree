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

	tt := []struct {
		name   string
		input  DependencyTree
		output string
	}{
		{name: "simple", input: dt, output: `{"test":{}}`},
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
