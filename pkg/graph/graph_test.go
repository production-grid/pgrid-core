package graph

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {

	vertices := make([]Vertex, 0)

	top1 := Vertex{ID: "top_level_1"}
	top2 := Vertex{ID: "top_level_2"}

	mid1 := Vertex{ID: "mid_level_1", ParentVertices: []*Vertex{&top1, &top2}}
	mid2 := Vertex{ID: "mid_level_2", ParentVertices: []*Vertex{&top2}}
	mid3 := Vertex{ID: "mid_level_3", ParentVertices: []*Vertex{&top2}}
	mid4 := Vertex{ID: "mid_level_4", ParentVertices: []*Vertex{&top1}}
	bot1 := Vertex{ID: "bottom_level_1", ParentVertices: []*Vertex{&mid1}}
	bot2 := Vertex{ID: "bottom_level_2", ParentVertices: []*Vertex{&mid2}}
	bot3 := Vertex{ID: "bottom_level_3", ParentVertices: []*Vertex{&mid2}}
	bot4 := Vertex{ID: "bottom_level_4", ParentVertices: []*Vertex{&mid2}}
	bot5 := Vertex{ID: "bottom_level_5", ParentVertices: []*Vertex{&mid2}}
	bot6 := Vertex{ID: "bottom_level_6", ParentVertices: []*Vertex{&mid3}}
	bot7 := Vertex{ID: "bottom_level_7", ParentVertices: []*Vertex{&mid4}}
	bot8 := Vertex{ID: "bottom_level_8", ParentVertices: []*Vertex{&mid4}}

	vertices = append(vertices, bot1, bot2, bot3, bot4, bot5, bot6, bot7, bot8)
	vertices = append(vertices, top1, top2)
	vertices = append(vertices, mid1, mid2, mid3, mid4)

	cycleCheckResult := CycleCheckGraph(vertices)

	if cycleCheckResult {
		t.Error("Invalid Cycle Check Result")
	}

	top1.ParentVertices = []*Vertex{&bot1}

	cycleCheckResult = CycleCheckGraph(vertices)

	if !cycleCheckResult {
		t.Error("Invalid Cycle Check Result")
	}

}

func TestTopographicSort(t *testing.T) {

	vertices := make([]Vertex, 0)

	top1 := Vertex{ID: "top_level_1"}
	top2 := Vertex{ID: "top_level_2"}

	mid1 := Vertex{ID: "mid_level_1", ParentVertices: []*Vertex{&top1, &top2}}
	mid2 := Vertex{ID: "mid_level_2", ParentVertices: []*Vertex{&top2}}
	mid3 := Vertex{ID: "mid_level_3", ParentVertices: []*Vertex{&top2}}
	mid4 := Vertex{ID: "mid_level_4", ParentVertices: []*Vertex{&top1}}
	bot1 := Vertex{ID: "bottom_level_1", ParentVertices: []*Vertex{&mid1}}
	bot2 := Vertex{ID: "bottom_level_2", ParentVertices: []*Vertex{&mid2}}
	bot3 := Vertex{ID: "bottom_level_3", ParentVertices: []*Vertex{&mid2}}
	bot4 := Vertex{ID: "bottom_level_4", ParentVertices: []*Vertex{&mid2}}
	bot5 := Vertex{ID: "bottom_level_5", ParentVertices: []*Vertex{&mid2}}
	bot6 := Vertex{ID: "bottom_level_6", ParentVertices: []*Vertex{&mid3}}
	bot7 := Vertex{ID: "bottom_level_7", ParentVertices: []*Vertex{&mid4, &top2}}
	bot8 := Vertex{ID: "bottom_level_8", ParentVertices: []*Vertex{&mid4}}

	vertices = append(vertices, bot1, bot2, bot3, bot4, bot5, bot6, bot7, bot8)
	vertices = append(vertices, top1, top2)
	vertices = append(vertices, mid1, mid2, mid3, mid4)

	tables := TopographicSort(vertices)

	if len(tables) < 1 {
		t.Error("no results")
	}

	if len(tables) < len(vertices) {
		t.Error("not enough results returned")
	}

	for _, v := range tables {
		fmt.Println(v)
	}

}
