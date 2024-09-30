package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"maps"
	"os/exec"
	"slices"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

type Edge struct {
	From string
	To   string
}

type GraphV1 struct {
	Root        string
	Edges       []Edge
	MvsPicked   []string
	MvsUnpicked []string
}

type Graph struct {
	Root             string
	RootIndirectMods []string
	Packages         map[string]*Package
}

func (g *Graph) IsRootIndirectMod(name string) bool {
	return slices.Contains(g.RootIndirectMods, name)
}

func (g *Graph) GetPackage(name string) *Package {
	p, exists := g.Packages[name]
	if !exists {
		return nil
	}
	return p
}

func (g *Graph) GetPackageVersion(name string, version string) (*Package, *PackageVersion) {
	p, exists := g.Packages[name]
	if !exists {
		return nil, nil
	}
	v, exists := p.Versions[version]
	if !exists {
		return nil, nil
	}
	return p, v
}

func (g *Graph) SortedPackages() []*Package {
	pkgs := slices.Collect(maps.Keys(g.Packages))
	slices.Sort(pkgs)
	var ret []*Package
	for _, p := range pkgs {
		ret = append(ret, g.Packages[p])
	}
	return ret
}

type Package struct {
	Name        string
	LastVersion string
	Versions    map[string]*PackageVersion
}

func (p *Package) SortedVersions() []*PackageVersion {
	versions := slices.Collect(maps.Keys(p.Versions))
	slices.SortFunc(versions, func(a, b string) int {
		return semver.Compare(a, b)
	})
	slices.Reverse(versions)
	var ret []*PackageVersion
	for _, v := range versions {
		ret = append(ret, p.Versions[v])
	}
	return ret
}

type PackageVersion struct {
	Version string
	Parents map[string]string
	Deps    map[string]string
}

type PackageItem struct {
	Name    string
	Version string
}

func RunGoModGraph() (*Graph, error) {
	indirectMods, err := indirects()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("go", "mod", "graph")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	g, err := ParseGoModGraph(strings.NewReader(out.String()))
	if err != nil {
		return nil, err
	}
	g.RootIndirectMods = indirectMods
	return g, nil
}

func ParseGoModGraph(r io.Reader) (*Graph, error) {
	scanner := bufio.NewScanner(r)
	g := &Graph{
		Packages: map[string]*Package{},
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 words in line, but got %d: %s", len(parts), line)
		}

		from := parseVersion(parts[0])
		to := parseVersion(parts[1])

		if skipPackage(from.Name) || skipPackage(to.Name) {
			continue
		}

		for _, node := range []PackageItem{from, to} {
			if _, ok := g.Packages[node.Name]; !ok {
				g.Packages[node.Name] = &Package{
					Name:     node.Name,
					Versions: map[string]*PackageVersion{},
				}
			}

			if _, ok := g.Packages[node.Name].Versions[node.Version]; !ok {
				g.Packages[node.Name].Versions[node.Version] = &PackageVersion{
					Version: node.Version,
					Parents: map[string]string{},
					Deps:    map[string]string{},
				}
			}

			if g.Packages[node.Name].LastVersion == "" || semver.Compare(g.Packages[node.Name].LastVersion, node.Version) < 0 {
				g.Packages[node.Name].LastVersion = node.Version
			}
			if node.Version == "" {
				g.Root = node.Name
			}
		}

		// if from.Version != "" && to.Version != "" {
		g.Packages[from.Name].Versions[from.Version].Deps[to.Name] = to.Version
		g.Packages[to.Name].Versions[to.Version].Parents[from.Name] = from.Version
		// }
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return g, nil
}

func ParseGoModGraphV1(r io.Reader) (*GraphV1, error) {
	scanner := bufio.NewScanner(r)
	var g GraphV1
	seen := make(map[string]bool)
	mvsPicked := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 words in line, but got %d: %s", len(parts), line)
		}

		from, to := parts[0], parts[1]
		g.Edges = append(g.Edges, Edge{From: from, To: to})

		for _, node := range []string{from, to} {
			if seen[node] {
				continue
			}
			seen[node] = true

			var module, version string
			if i := strings.IndexByte(node, '@'); i >= 0 {
				module, version = node[:i], node[i+1:]
			} else {
				g.Root = node
				continue
			}

			if maxVersion, exists := mvsPicked[module]; exists {
				if semver.Compare(maxVersion, version) < 0 {
					g.MvsUnpicked = append(g.MvsUnpicked, module+"@"+maxVersion)
					mvsPicked[module] = version
				} else {
					g.MvsUnpicked = append(g.MvsUnpicked, node)
				}
			} else {
				mvsPicked[module] = version
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for module, version := range mvsPicked {
		g.MvsPicked = append(g.MvsPicked, module+"@"+version)
	}

	sort.Strings(g.MvsPicked)
	return &g, nil
}

func skipPackage(name string) bool {
	return false
	//
	// if name == "go" || name == "toolchain" {
	// 	return false
	// }
	// return strings.HasPrefix(name, "golang.org/x/")
}

// "// indirect" packages are listed as direct dependencies; this is how go mod
// graph outputs it.
func indirects() ([]string, error) {
	out, err := exec.Command("go", "list", "-f", "{{.Indirect}} {{.Path}}", "-m", "all").CombinedOutput()
	if err != nil {
		return nil, err
	}

	in := make([]string, 0, 8)
	for _, line := range strings.Split(string(out), "\n") {
		indir, pkg, _ := strings.Cut(line, " ")
		if indir == "true" {
			in = append(in, pkg)
		}
	}
	return in, nil
}

func parseVersion(node string) PackageItem {
	if i := strings.IndexByte(node, '@'); i >= 0 {
		return PackageItem{Name: node[:i], Version: node[i+1:]}
	} else {
		return PackageItem{Name: node}
	}
}
