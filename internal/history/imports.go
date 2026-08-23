package history

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BuildImportEdges parses Go source files to extract import edges for
// cycle detection. Only Go files are supported in this phase — JS/TS
// import cycle detection is deferred to a future release to avoid
// flaky regex-based parsing producing false cycles.
func (idx *Index) BuildImportEdges(repoPath string, files []string) {
	for _, file := range files {
		if !strings.HasSuffix(file, ".go") {
			continue
		}

		absPath := filepath.Join(repoPath, file)
		imports, err := extractGoImports(absPath)
		if err != nil {
			continue // skip unparseable files silently
		}

		// Use the package directory as the node (not individual files)
		pkg := filepath.Dir(file)
		for _, imp := range imports {
			idx.ImportEdges[pkg] = appendUnique(idx.ImportEdges[pkg], imp)
		}
	}
}

// extractGoImports returns the import paths from a Go source file.
func extractGoImports(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []string
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		imports = append(imports, path)
	}
	return imports, nil
}

// BuildImportEdgesFromSource parses Go source content directly.
// Used in tests where files don't exist on disk.
func (idx *Index) BuildImportEdgesFromSource(sources map[string]string) {
	fset := token.NewFileSet()
	for file, src := range sources {
		if !strings.HasSuffix(file, ".go") {
			continue
		}

		f, err := parser.ParseFile(fset, file, src, parser.ImportsOnly)
		if err != nil {
			continue
		}

		pkg := filepath.Dir(file)
		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			idx.ImportEdges[pkg] = appendUnique(idx.ImportEdges[pkg], path)
		}
	}
}

// ImportCycles runs Tarjan's strongly connected components algorithm on
// the import graph and returns all cycles (SCCs with size > 1).
func (idx *Index) ImportCycles() []Cycle {
	return tarjanSCC(idx.ImportEdges)
}

// FileInImportCycle checks if the given file's package participates in
// any detected import cycle.
// FileInImportCycle checks if the given file's package participates in
// any detected import cycle.
func (idx *Index) FileInImportCycle(file string) (bool, Cycle) {
	pkg := filepath.ToSlash(filepath.Dir(file))
	cycles := idx.ImportCycles()
	for _, cycle := range cycles {
		for _, node := range cycle {
			if node == pkg {
				return true, cycle
			}
		}
	}
	return false, nil
}

// tarjanSCC implements Tarjan's algorithm for finding strongly connected
// components in a directed graph. Returns components with size > 1 (cycles).
func tarjanSCC(graph map[string][]string) []Cycle {
	var (
		index    int
		stack    []string
		onStack  = make(map[string]bool)
		indices  = make(map[string]int)
		lowlinks = make(map[string]int)
		defined  = make(map[string]bool)
		result   []Cycle
	)

	// Collect all nodes
	nodes := make(map[string]bool)
	for k, vs := range graph {
		nodes[k] = true
		for _, v := range vs {
			nodes[v] = true
		}
	}

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indices[v] = index
		lowlinks[v] = index
		defined[v] = true
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range graph[v] {
			if !defined[w] {
				strongConnect(w)
				if lowlinks[w] < lowlinks[v] {
					lowlinks[v] = lowlinks[w]
				}
			} else if onStack[w] {
				if indices[w] < lowlinks[v] {
					lowlinks[v] = indices[w]
				}
			}
		}

		// Root of an SCC
		if lowlinks[v] == indices[v] {
			var scc []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}

			// Include SCCs with >1 members (multi-node cycles) or
			// single-node SCCs with a self-edge (self-loops).
			isCycle := len(scc) > 1
			if len(scc) == 1 {
				for _, w := range graph[scc[0]] {
					if w == scc[0] {
						isCycle = true
						break
					}
				}
			}
			if isCycle {
				// Sort for deterministic output
				sort.Strings(scc)
				result = append(result, Cycle(scc))
			}
		}
	}

	// Sort node iteration for deterministic results
	sortedNodes := make([]string, 0, len(nodes))
	for n := range nodes {
		sortedNodes = append(sortedNodes, n)
	}
	sort.Strings(sortedNodes)

	for _, v := range sortedNodes {
		if !defined[v] {
			strongConnect(v)
		}
	}

	return result
}

// appendUnique appends val to slice if not already present.
func appendUnique(slice []string, val string) []string {
	for _, s := range slice {
		if s == val {
			return slice
		}
	}
	return append(slice, val)
}

// ListGoFiles returns all .go files under repoPath (non-test, non-vendor).
func ListGoFiles(repoPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == ".git" || base == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			rel, err := filepath.Rel(repoPath, path)
			if err == nil {
				files = append(files, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	return files, err
}
