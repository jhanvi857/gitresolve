package history

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	funcDefRegex  = regexp.MustCompile(`(?:func|def|function)\s+([A-Za-z0-9_]+)`)
	methodGoRegex = regexp.MustCompile(`func\s*\([^)]+\)\s*([A-Za-z0-9_]+)`)
)

// BuildSymbolIndex parses source files in the repository to populate SymbolCallers.
// It analyzes Go files using standard AST analysis and source files using syntax patterns.
func (idx *Index) BuildSymbolIndex(repoPath string, files []string) {
	for _, file := range files {
		absPath := filepath.Join(repoPath, file)
		if strings.HasSuffix(file, ".go") {
			idx.indexGoSymbols(absPath, file)
		} else if isSourceFile(file) {
			idx.indexGenericSymbols(absPath, file)
		}
	}
}

// indexGoSymbols parses a Go file and extracts all function calls and their caller contexts.
func (idx *Index) indexGoSymbols(absPath, relPath string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, absPath, nil, 0)
	if err != nil {
		return
	}

	var currentFunc string
	ast.Inspect(node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			currentFunc = fn.Name.Name
		case *ast.CallExpr:
			calleeName := ""
			switch fun := fn.Fun.(type) {
			case *ast.Ident:
				calleeName = fun.Name
			case *ast.SelectorExpr:
				calleeName = fun.Sel.Name
			}

			if calleeName != "" {
				line := fset.Position(fn.Pos()).Line
				sym := SymbolID{Name: calleeName}
				idx.SymbolCallers[sym] = append(idx.SymbolCallers[sym], CallerRef{
					File: filepath.ToSlash(relPath),
					Line: line,
					Name: currentFunc,
				})
			}
		}
		return true
	})
}

// indexGenericSymbols provides basic call extraction for non-Go source files.
func (idx *Index) indexGenericSymbols(absPath, relPath string) {
	data, err := os.ReadFile(filepath.Clean(absPath))
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	callRegex := regexp.MustCompile(`\b([A-Za-z0-9_]+)\s*\(`)

	for i, line := range lines {
		matches := callRegex.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			if len(m) > 1 {
				name := m[1]
				if isCommonKeyword(name) {
					continue
				}
				sym := SymbolID{Name: name}
				idx.SymbolCallers[sym] = append(idx.SymbolCallers[sym], CallerRef{
					File: filepath.ToSlash(relPath),
					Line: i + 1,
					Name: "",
				})
			}
		}
	}
}

// CallersOfSymbol returns all caller references matching the symbol name.
func (idx *Index) CallersOfSymbol(name string) []CallerRef {
	var result []CallerRef
	for sym, callers := range idx.SymbolCallers {
		if sym.Name == name {
			result = append(result, callers...)
		}
	}
	return result
}

// ExtractSymbolsFromLines detects symbol/function names declared or modified in lines.
func ExtractSymbolsFromLines(lines []string) []string {
	var symbols []string
	seen := make(map[string]bool)

	for _, line := range lines {
		if m := methodGoRegex.FindStringSubmatch(line); len(m) > 1 {
			name := m[1]
			if !seen[name] && !isCommonKeyword(name) {
				seen[name] = true
				symbols = append(symbols, name)
			}
		} else if m := funcDefRegex.FindStringSubmatch(line); len(m) > 1 {
			name := m[1]
			if !seen[name] && !isCommonKeyword(name) {
				seen[name] = true
				symbols = append(symbols, name)
			}
		}
	}
	return symbols
}

// ExtractEnclosingSymbol tries to find what function encloses a line in a file.
func ExtractEnclosingSymbol(repoPath, relPath string, targetLine int) string {
	if !strings.HasSuffix(relPath, ".go") {
		return ""
	}
	absPath := filepath.Join(repoPath, relPath)
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, absPath, nil, 0)
	if err != nil {
		return ""
	}

	var found string
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			start := fset.Position(fn.Pos()).Line
			end := fset.Position(fn.End()).Line
			if targetLine >= start && targetLine <= end {
				found = fn.Name.Name
				return false
			}
		}
		return true
	})
	return found
}

func isSourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".c", ".cpp", ".rs":
		return true
	default:
		return false
	}
}

func isCommonKeyword(name string) bool {
	switch name {
	case "if", "for", "while", "switch", "return", "catch", "import", "make", "new", "append", "len", "cap", "panic", "recover", "func", "function", "var", "const", "type":
		return true
	default:
		return false
	}
}
