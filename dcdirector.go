package dcdirector

import (
	"fmt"
	"strings"
)

// == Definitions ==

// Vertex is a node that represents a data center
type Vertex struct {
	Name      string   // name of data center, should be unique, should not contains "->" or ";" or "*"
	Outdegree float64  // how many clients deployed in this data center
	Indegree  float64  // how many servers deployed in this data center
	Neighbors []string // data centers' name that close to current data center.
}

// Graph is a result calling relationship
type Graph struct {
	edges []Edge
}

// Edge is a directed line between two vertex
type Edge struct {
	From   string
	To     string
	Weight float64
}

// DCDirector is an interface for implementation
type DCDirector interface {
	// Route accepts a bunch of vertexes and returns a graph
	Route(vertexes []Vertex) (*Graph, error)
}

// == Implementions

// Validate makes sure that vertex is valid
func (v Vertex) Validate() error {
	if strings.Contains(v.Name, "->") || strings.Contains(v.Name, ";") || strings.Contains(v.Name, "*") {
		return fmt.Errorf(`vertex name %q should not contains "->" or ";" or "*"`, v.Name)
	}
	if v.Outdegree < 0 {
		return fmt.Errorf(`vertex outdegress %f is less than 0`, v.Outdegree)
	}
	if v.Indegree < 0 {
		return fmt.Errorf(`vertex indegress %f is less than 0`, v.Indegree)
	}
	return nil
}

func (g Graph) String() string {
	ret := ""
	for _, edge := range g.edges {
		if edge.Weight <= 0 {
			ret += edge.From + "!>" + edge.To + ";"
		} else {
			if edge.Weight >= 1 {
				ret += edge.From + "->" + edge.To + ";"
			} else {
				ret += fmt.Sprintf("%s->%s*%.3f;", edge.From, edge.To, edge.Weight)
			}
		}
	}
	return ret
}
