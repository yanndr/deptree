package deptree_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"bitbucket.org/yanndr/deptree"
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

	dt, err := deptree.New("cmd/deptree/data")
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

func TestJson(t *testing.T) {

	tt := []struct {
		name   string
		input  []string
		indent string
		output string
	}{
		{"Specio", []string{"Specio"}, "\t", "{\n\t\"Specio\": {\n\t\t\"Devel-StackTrace\": {},\n\t\t\"Eval-Closure\": {},\n\t\t\"MRO-Compat\": {},\n\t\t\"Module-Runtime\": {},\n\t\t\"Role-Tiny\": {},\n\t\t\"Sub-Quote\": {},\n\t\t\"Test-Fatal\": {\n\t\t\t\"Try-Tiny\": {}\n\t\t},\n\t\t\"Try-Tiny\": {}\n\t}\n}"},
	}
	dt, err := deptree.New("cmd/deptree/data")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, err := dt.Resolve(tc.input...)
			if err != nil {
				t.Fatal(err)
			}
			js := r.ToJSON(tc.indent)
			if js != tc.output {
				t.Fatalf("expected: \n%s, got \n%s", tc.output, js)
			}
		})
	}
}

func BenchmarkResolveOneByOne(b *testing.B) {
	path := "./cmd/deptree/data/"
	dt, err := deptree.New(path)
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
	dt, err := deptree.New(path)
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

func BenchmarkResolveScale(b *testing.B) {
	path := "./cmd/deptree/data/"
	dt, err := deptree.New(path)
	if err != nil {
		b.Fatal(err)
	}

	distribution := "DateTime"

	for n := 1; n <= 1024; n *= 2 {
		var test []string
		for i := 0; i < n; i++ {
			test = append(test, distribution)
		}
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := dt.Resolve(test...)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
