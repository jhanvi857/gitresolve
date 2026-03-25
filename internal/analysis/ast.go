//go:build cgo && treesitter

package analysis

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
)

// Node represents one element in the syntax tree

type Node struct {
	Type      string
	Name      string
	StartLine int
	EndLine   int
	Children  []*Node
}

type ParseResult struct {
	Root     *Node
	Language string
	Source   []byte
}

func ParseFile(filePath string, content []byte) (*ParseResult, error) {
	lang, err := detectLanguage(filePath)
	if err != nil {
		return nil, fmt.Errorf("ParseFile: %w", err)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, fmt.Errorf("ParseFile: parsing: %w", err)
	}
	root := convertNode(tree.RootNode(), content)

	return &ParseResult{
		Root:     root,
		Language: filepath.Ext(filePath),
		Source:   content,
	}, nil
}

func FindFunction(result *ParseResult, funcName string) *Node {
	return findNodeByName(result.Root, funcName, "function_declaration", "method_declaration")
}

func FindChangedNodes(result *ParseResult, startLine, endLine int) []*Node {
	var affected []*Node
	collectNodesInRange(result.Root, startLine, endLine, &affected)
	return affected
}

func CompareNodes(a, b *Node) bool {
	if a.Type != b.Type {
		return false
	}

	if len(a.Children) != len(b.Children) {
		return false
	}

	for i := range a.Children {
		if a.Children[i].Type != b.Children[i].Type {
			return false
		}
	}

	return true
}

func detectLanguage(filePath string) (*sitter.Language, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return golang.GetLanguage(), nil
	case ".js", ".jsx", ".ts", ".tsx":
		return javascript.GetLanguage(), nil
	default:
		return nil, fmt.Errorf("detectLanguage: unsupported language %s", ext)
	}
}

func convertNode(n *sitter.Node, source []byte) *Node {
	if n == nil {
		return nil
	}

	node := &Node{
		Type:      n.Type(),
		StartLine: int(n.StartPoint().Row),
		EndLine:   int(n.EndPoint().Row),
	}
	for i := 0; i < int(n.ChildCount()); i++ {
		child := n.Child(i)
		if child.Type() == "identifier" {
			node.Name = child.Content(source)
			break
		}
	}

	for i := 0; i < int(n.ChildCount()); i++ {
		child := convertNode(n.Child(i), source)
		if child != nil {
			node.Children = append(node.Children, child)
		}
	}

	return node
}

func findNodeByName(node *Node, name string, types ...string) *Node {
	if node == nil {
		return nil
	}

	if node.Name == name {
		for _, t := range types {
			if node.Type == t {
				return node
			}
		}
	}

	for _, child := range node.Children {
		if found := findNodeByName(child, name, types...); found != nil {
			return found
		}
	}

	return nil
}

func collectNodesInRange(node *Node, startLine, endLine int, result *[]*Node) {
	if node == nil {
		return
	}
	nodeOverlaps := node.StartLine <= endLine && node.EndLine >= startLine
	if nodeOverlaps {
		*result = append(*result, node)
	}
	for _, child := range node.Children {
		collectNodesInRange(child, startLine, endLine, result)
	}
}
