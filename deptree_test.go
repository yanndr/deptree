package deptree

// func TestToJSON(t *testing.T) {

// 	tt := []struct {
// 		name   string
// 		input  Distributions
// 		indent string
// 		output string
// 	}{
// 		{name: "simple", input: dt, output: `{"test": {}}`},
// 		{name: "with dependencies", input: dtwithDep, output: `{"test With dep": {"test": {}}}`},
// 		{name: "with multiple dependencies", input: multipleDep, output: `{"test With multiple dep": {"test": {},"test2": {}}}`},
// 		{name: "with multiple root", input: multipleRoot, output: `{"dist1": {},"dist2": {}}`},
// 		{name: "complex", input: complex, output: `{"dist1": {"test With multiple dep": {"test": {},"test2": {}}},"dist2": {"test With multiple dep": {"test": {},"test2": {}}}}`},
// 		{name: "with dependencies with indent", input: dtwithDep, output: "{\n\t\"test With dep\": {\n\t\t\"test\": {}\n\t}\n}", indent: "\t"},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			//result := tc.input.ToJSON(tc.indent)
// 			if result != tc.output {
// 				t.Fatalf("error, expected \n%s, got \n%s", tc.output, result)
// 			}
// 		})
// 	}
// }

var (
	dt = Distributions{"test": nil}

	dtwithDep = Distributions{"test With dep": dt}

	multipleDep = Distributions{"test With multiple dep": Distributions{"test": nil}, "test2": nil}

	multipleRoot = Distributions{"dist1": nil, "dist2": nil}

	complex = Distributions{"dist1": multipleDep, "dist2": multipleDep}
)
