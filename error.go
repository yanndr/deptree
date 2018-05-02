package deptree

import "fmt"

//DistributionNotFoundError is an error when the distribution could not be found.
type DistributionNotFoundError struct {
	name string
	err  error
}

func (e DistributionNotFoundError) Error() string {
	return fmt.Sprintf("distribution %s not found: %v", e.name, e.err)
}

//ModuleNotFoundError is an error when a module could not be found.
type ModuleNotFoundError struct {
	name string
	err  error
}

func (e ModuleNotFoundError) Error() string {
	return fmt.Sprintf("module %s not found: %v", e.name, e.err)
}
