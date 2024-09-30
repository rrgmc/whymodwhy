package main

import (
	"fmt"
	"log"
	"os"
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

	if len(os.Args) > 1 {
		pkgs, ok := graph.Packages[os.Args[1]]
		if !ok {
			return fmt.Errorf("unknown package %s", os.Args[1])
		}
		printPackage(graph, pkgs)
		return nil
	}

	for _, pkgs := range graph.SortedPackages() {
		printPackage(graph, pkgs)
	}

	return nil
}

func printPackage(graph *pkg.Graph, pkgs *pkg.Package) {
	fmt.Printf("%s %s (%s) %s\n", strings.Repeat("=", 5), pkgs.Name, pkgs.LastVersion, strings.Repeat("=", 5))
	for _, version := range pkgs.SortedVersions() {
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
					if graph.IsRootIndirectMod(pkgs.Name) {
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
