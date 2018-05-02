// +build integration

package deptree

import (
	"io/ioutil"
	"testing"
)

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
		{name: "DateTime",
			input:     []string{"DateTime"},
			numberDep: 7,
		},
		{name: "DateTime-TimeZone",
			input:     []string{"DateTime-TimeZone"},
			numberDep: 6,
		},
		{name: "Params-ValidationCompiler",
			input:     []string{"Params-ValidationCompiler"},
			numberDep: 2,
		},
		{name: "Eval-Closure",
			input:     []string{"Eval-Closure"},
			numberDep: 0,
		},
		{name: "Exception-Class",
			input:     []string{"Exception-Class"},
			numberDep: 2,
		},
		{name: "Try-Tiny",
			input:     []string{"Try-Tiny"},
			numberDep: 0,
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

func BenchmarkResolveOneByOne(b *testing.B) {
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

func BenchmarkResolveAllAtOnce(b *testing.B) {
	path := "./cmd/deptree/data/"
	dt, err := New(path)
	if err != nil {
		b.Fatal(err)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		b.Fatal(err)
	}
	var distribs []string
	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		distribs = append(distribs, f.Name())

	}

	for n := 0; n < b.N; n++ {
		_, err := dt.Resolve(distribs...)
		if err != nil {
			b.Fatal(err)
		}

	}
}
