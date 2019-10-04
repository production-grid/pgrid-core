package graph

/*
CycleCheckVertex checks the given vertix for cycles.
*/
func CycleCheckVertex(vertex *Vertex) bool {

	seenMap := make(map[string]bool)

	if len(vertex.ParentVertices) == 0 {
		return false
	}

	return cycleCheckParents(seenMap, vertex)

}

func cycleCheckParents(seen map[string]bool, vertex *Vertex) bool {

	_, found := seen[vertex.ID]

	if found {
		return true
	}

	seen[vertex.ID] = true

	for _, v := range vertex.ParentVertices {
		if cycleCheckParents(seen, v) {
			return true
		}
	}

	return false

}

/*
CycleCheckGraph checks the set of vertices for cycles.
*/
func CycleCheckGraph(vertices []Vertex) bool {

	for _, v := range vertices {
		if CycleCheckVertex(&v) {
			return true
		}
	}
	return false

}
