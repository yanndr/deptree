//Package deptree is a package that resolves the dependency tree for package/distributions
//
// So far only a Perl implementation is available but extensions could be added.
package deptree

// Resolver defines the method of a dependency tree resolver.
type Resolver interface {
	//Resolve returns the distribution list with dependencies.
	Resolve(distributions ...string) (Distributions, error)
}

// Distributions is an array of distribution that represents the tree of dependecies for distributions.
type Distributions map[string]Distributions
