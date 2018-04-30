package deptree

import (
	"io/ioutil"
	"testing"
)

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
