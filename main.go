package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/rrgmc/whymodwhy/pkg"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	graph, err := pkg.RunGoModGraph()
	if err != nil {
		return err
	}

	for _, pkgs := range graph.Packages {
		fmt.Printf("%s %s (%s) %s\n", strings.Repeat("=", 5), pkgs.Name, pkgs.LastVersion, strings.Repeat("=", 5))
		for _, version := range pkgs.Versions {
			versionextra := ""
			if version.Version == pkgs.LastVersion {
				versionextra = " (last)"
			}
			fmt.Printf("\tVersion: %s%s\n", version.Version, versionextra)
			if len(version.Parents) > 0 {
				fmt.Printf("\t\t%s Parents %s\n", strings.Repeat("-", 5), strings.Repeat("-", 5))
				for parentPkg, parentVersion := range version.Parents {
					if parentPkg == graph.Root {
						pextra := ""
						if slices.Contains(graph.RootIndirectMods, pkgs.Name) {
							pextra = " (indirect)"
						}
						fmt.Printf("\t\t%s%s\n", parentPkg, pextra)
					} else {
						fmt.Printf("\t\t%s (%s)\n", parentPkg, parentVersion)
					}
				}
			}
			if len(version.Deps) > 0 {
				fmt.Printf("\t\t%s Deps %s\n", strings.Repeat("-", 5), strings.Repeat("-", 5))
				for parentPkg, parentVersion := range version.Deps {
					fmt.Printf("\t\t%s (%s)\n", parentPkg, parentVersion)
				}
			}
		}
	}

	// spew.Dump(graph)

	return nil
}
