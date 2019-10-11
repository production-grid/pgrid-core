package graph

/*
TopographicSort returns the id's of all vertices topographically sorted.
*/
func TopographicSort(vertices []Vertex) []string {

	resultIds := make([]string, 0, 16)

	seen := make(map[string]bool)

	var reliefValve = 0

	for (reliefValve < 100) && (len(resultIds) < len(vertices)) {
		reliefValve++
		for _, v := range vertices {
			_, ok := seen[v.ID]
			if !ok {
				if len(v.ParentVertices) == 0 {
					seen[v.ID] = true
					resultIds = append(resultIds, v.ID)
				} else {
					if seenAllParents(seen, v) {
						seen[v.ID] = true
						resultIds = append(resultIds, v.ID)
					}
				}
			}
		}
	}

	return resultIds
}

func seenAllParents(seenMap map[string]bool, vertex Vertex) bool {

	if len(vertex.ParentVertices) == 0 {
		return true
	}

	for _, v := range vertex.ParentVertices {
		_, parentSeen := seenMap[v.ID]
		if !parentSeen {
			return false
		}
	}

	return true

}
