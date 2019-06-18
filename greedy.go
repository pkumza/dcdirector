package dcdirector

import (
	"fmt"
	"log"
)

// NewGreedyDirector creates a new greedy director
func NewGreedyDirector(simplify, verbose bool) DCDirector {
	return GreedyDirector{
		ResSimplify: simplify,
		Verbose:     verbose,
	}
}

// GreedyDirector impls DCDirector
type GreedyDirector struct {
	ResSimplify bool
	Verbose     bool
}

type greedyNode struct {
	Vertex
	Diff    float64
	Weights []*Edge
}

// Route impls DCDirector
func (g GreedyDirector) Route(vertexes []Vertex) (*Graph, error) {
	// Step1: make a copy
	nodes := make([]*greedyNode, len(vertexes))
	for i, vertex := range vertexes {
		if vertex.Validate() != nil {
			return nil, fmt.Errorf("vertex.Validate() of %q reports %v", vertex.Name, vertex.Validate())
		}
		nodes[i] = &greedyNode{Vertex: vertex}
	}

	// Step2: transform in case egressSum is not equal to ingressSum
	var ingressSum, egressSum float64
	for _, node := range nodes {
		if node.Outdegree < 0 {
			return nil, fmt.Errorf("egress of %s is %f, less than 0", node.Name, node.Outdegree)
		}
		if node.Indegree < 0 {
			return nil, fmt.Errorf("ingress of %s is %f, less than 0", node.Name, node.Indegree)
		}
		egressSum += node.Outdegree
		ingressSum += node.Indegree
	}
	if egressSum == 0 {
		return nil, fmt.Errorf("egressSum is 0")
	}
	if ingressSum == 0 {
		return nil, fmt.Errorf("ingressSum is 0")
	}
	for _, node := range nodes {
		node.Outdegree *= ingressSum / egressSum
	}

	// Diff
	for _, node := range nodes {
		node.Diff = node.Outdegree - node.Indegree
	}

	// Print
	if g.Verbose {
		log.Println("===== 2 ==== ")
		for _, node := range nodes {
			log.Println(node)
		}
	}

	// Self Fill
	for _, node := range nodes {
		if g.ResSimplify {
			if node.Indegree == 0 || node.Outdegree == 0 {
				continue
			}
		}
		if node.Diff > 0 {
			node.Weights = append(node.Weights, &Edge{
				From:   node.Name,
				To:     node.Name,
				Weight: node.Indegree,
			})
		} else {
			node.Weights = append(node.Weights, &Edge{
				From:   node.Name,
				To:     node.Name,
				Weight: node.Outdegree,
			})
		}
	}

	// Overflow
	for _, convex := range nodes {
		if convex.Diff > 0 {
			for _, concavo := range nodes {
				if convex.Diff <= 0 {
					break
				}
				if concavo.Diff < 0 {
					if convex.Diff+concavo.Diff < 0 {
						convex.Weights = append(convex.Weights, &Edge{
							From:   convex.Name,
							To:     concavo.Name,
							Weight: convex.Diff,
						})
						concavo.Diff += convex.Diff
						convex.Diff = 0
					} else {
						convex.Weights = append(convex.Weights, &Edge{
							From:   convex.Name,
							To:     concavo.Name,
							Weight: -concavo.Diff,
						})
						convex.Diff += concavo.Diff
						concavo.Diff = 0
					}
				} else {
					continue
				}
			}
		} else {
			continue
		}
	}

	// Print
	if g.Verbose {
		log.Println("===== 3 ==== ")
		for _, node := range nodes {
			log.Println(node)
			for _, w := range node.Weights {
				log.Println(w.From, "->", w.To, w.Weight)
			}
		}
	}

	// Align max to 1.0
	for _, node := range nodes {
		var maxWeight float64
		for _, w := range node.Weights {
			if w.Weight > maxWeight {
				maxWeight = w.Weight
			}
		}
		if maxWeight == 0 {
			for _, w := range node.Weights {
				w.Weight = 0.0
			}
		} else {
			for _, w := range node.Weights {
				w.Weight /= maxWeight
			}
		}
	}

	// Print
	if g.Verbose {
		log.Println("===== 4 ==== ")
		for _, node := range nodes {
			log.Println(node)
			for _, w := range node.Weights {
				log.Println(w.From, "->", w.To, w.Weight)
			}
		}
	}

	graph := &Graph{
		edges: make([]Edge, 0),
	}
	for _, node := range nodes {
		for _, w := range node.Weights {
			if g.ResSimplify {
				if w.From == w.To && w.Weight >= 1.0 {
					continue
				}
				graph.edges = append(graph.edges, *w)
			} else {
				graph.edges = append(graph.edges, *w)
			}
		}
	}

	return graph, nil
}
