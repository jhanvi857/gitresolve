//go:build !treesitter || !cgo

package analysis

import "fmt"

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
	return nil, fmt.Errorf("ParseFile: AST analysis requires cgo-enabled build")
}

func FindFunction(result *ParseResult, funcName string) *Node {
	if result == nil || result.Root == nil {
		return nil
	}
	return findNodeByName(result.Root, funcName, "function_declaration", "method_declaration")
}

func FindChangedNodes(result *ParseResult, startLine, endLine int) []*Node {
	if result == nil || result.Root == nil {
		return nil
	}
	var affected []*Node
	collectNodesInRange(result.Root, startLine, endLine, &affected)
	return affected
}

func CompareNodes(a, b *Node) bool {
	if a == nil || b == nil {
		return false
	}
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
