package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rrgmc/whymodwhy/pkg"
)

var (
	printOnly       = flag.Bool("p", false, "print only")
	showLastVersion = flag.Bool("v", false, "show last version of returned packages")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	graph, err := pkg.RunGoModGraph()
	if err != nil {
		return err
	}

	if *printOnly {
		return runPrint(graph)
	}

	return runFind(graph)
}

func runFind(graph *pkg.Graph) error {
	if len(flag.Args()) == 0 {
		return errors.New("package name is required")
	}

	pkgs, ok := graph.Packages[flag.Arg(0)]
	if !ok {
		return fmt.Errorf("package '%s' not found in dependencies", flag.Arg(0))
	}

	return runFindPkg(graph, pkgs)
}

func runFindPkg(graph *pkg.Graph, p *pkg.Package) error {
	// printPackage(graph, p)

	// check if direct dependency
	if checkDirectDependency(graph, p) {
		fmt.Printf("'%s' is a direct dependency of the root module '%s', just run:\n", p.Name, graph.Root)
		fmt.Printf("go get %s\n", p.Name)
		return nil
	}

	pkgs, err := checkRootPackages(graph, p)
	if err != nil {
		return err
	}

	var lverr error

	fmt.Printf("to upgrade '%s' these packages must be upgraded:\n", p.Name)
	for _, fp := range pkgs {
		lastversion := ""

		if *showLastVersion {
			lv, err := GetLatestPackageVersion(fp)
			if err == nil {
				lastversion = fmt.Sprintf(" (last version: %s)", lv.Version)
			} else {
				lverr = errors.Join(lverr, fmt.Errorf("error getting last version of '%s': %w", fp, err))
			}
		}

		xp := graph.GetPackage(fp)
		if xp == nil {
			fmt.Printf("- %s%s\n", fp)
		} else {
			fmt.Printf("- %s (local version: %s)%s\n", fp, xp.LastVersion, lastversion)
		}
	}

	if lverr != nil {
		fmt.Printf("errors getting latest version of packages: %s", lverr)
	}

	return nil
}

func runPrint(graph *pkg.Graph) error {
	if len(flag.Args()) > 0 {
		pkgs, ok := graph.Packages[flag.Arg(0)]
		if !ok {
			return fmt.Errorf("unknown package %s", flag.Arg(0))
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
