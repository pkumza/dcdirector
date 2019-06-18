package dcdirector

import "testing"

func TestEastToWest(t *testing.T) {
	nodes := []Vertex{
		{Name: "us-east1", Outdegree: 50.0, Indegree: 0.0},
		{Name: "asia-south1", Outdegree: 60.0, Indegree: 60.0},
		{Name: "us-west1", Outdegree: 0, Indegree: 50.0},
		{Name: "us-central1", Outdegree: 20, Indegree: 20.0},
	}
	dcTrafficGraph, err := NewGreedyDirector(false, false).Route(nodes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dcTrafficGraph %q", dcTrafficGraph)
}

func TestCross(t *testing.T) {
	nodes := []Vertex{
		{Name: "us-east1", Outdegree: 50.0, Indegree: 0.0},
		{Name: "asia-south1", Outdegree: 70.0, Indegree: 60.0},
		{Name: "us-west1", Outdegree: 0, Indegree: 50.0},
		{Name: "us-central1", Outdegree: 10, Indegree: 20.0},
	}
	dcTrafficGraph, err := NewGreedyDirector(true, false).Route(nodes)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dcTrafficGraph %q", dcTrafficGraph)
}
