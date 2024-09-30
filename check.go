package main

import (
	"fmt"

	"github.com/rrgmc/whymodwhy/pkg"
)

// checkDirectDependency returns whether if a package is a direct dependency of the root module.
func checkDirectDependency(graph *pkg.Graph, p *pkg.Package) bool {
	for _, version := range p.Versions {
		for parent, _ := range version.Parents {
			if parent == graph.Root {
				if !graph.IsRootIndirectMod(p.Name) {
					return true
				}
			}
		}
	}
	return false
}

func checkRootPackages(graph *pkg.Graph, p *pkg.Package) ([]string, error) {
	visited := map[string]struct{}{}

	var ret []string

	for _, version := range p.Versions {
		for parentName, parentVersion := range version.Parents {
			if parentName == graph.Root {
				continue
			}
			pname := fmt.Sprintf("%s-%s", parentName, parentVersion)
			if _, ok := visited[pname]; ok {
				continue
			}
			visited[pname] = struct{}{}

			p, v := graph.GetPackageVersion(parentName, parentVersion)
			if p == nil {
				return nil, fmt.Errorf("internal error: package '%s' version '%s' not found", parentName, parentVersion)
			}

			f, err := getParentRootPackages(graph, p, v)
			if err != nil {
				return nil, err
			}
			ret = append(ret, f...)
		}
	}

	return dedupeSlice(ret), nil
}

func getParentRootPackages(graph *pkg.Graph, p *pkg.Package, v *pkg.PackageVersion) ([]string, error) {
	// fmt.Printf("!! %s -- %s\n", p.Name, v.Version)

	if checkDirectDependency(graph, p) {
		// fmt.Printf("@ DIRECT: %s\n", p.Name)
		return []string{p.Name}, nil
	}

	var ret []string

	for parentName, parentVersion := range v.Parents {
		if parentName == graph.Root {
			continue
		}

		parentp, parentv := graph.GetPackageVersion(parentName, parentVersion)
		if parentp == nil {
			return nil, fmt.Errorf("internal error: package '%s' version '%s' not found", parentName, parentVersion)
		}

		f, err := getParentRootPackages(graph, parentp, parentv)
		if err != nil {
			return nil, err
		}

		// fmt.Printf("@ %s(%s): %v\n", parentp.Name, parentv.Version, f)

		ret = append(ret, f...)
	}

	return ret, nil
}
