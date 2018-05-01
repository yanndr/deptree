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
		{"not Present", "module3", "", distributionNotFoundError{name: "module3"}},
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
