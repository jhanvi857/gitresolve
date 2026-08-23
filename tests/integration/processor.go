package main

import (
	"fmt"
	"math"
)

type Result struct {
	Value float64
	Tag   string
}

// ComputeFactor executes a complex logical transformation on inputs.
// This is used for Rule 5 (Logic Conflict) and TypeSignature testing.
func ComputeFactor(inputs []float64, weight float64) ([]Result, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("no inputs provided")
	}

	results := make([]Result, len(inputs))
	for i, v := range inputs {
		res := math.Pow(v, 2) * weight
		results[i] = Result{
			Value: res,
			Tag:   "computed",
		}
	}
	return results, nil
}

func main() {}
