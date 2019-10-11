package graph

/*
Graph models a graph structure.
*/
type Graph struct {
	Vertices []*Vertex
}

/*
Vertex models a vertix in a graph structure.
*/
type Vertex struct {
	ID             string
	ParentVertices []*Vertex
}

/*
VertexMap sorts the vertices into a map keyed by their id's.
*/
func VertexMap(vertices []Vertex) map[string]*Vertex {

	result := make(map[string]*Vertex, len(vertices))

	for _, v := range vertices {
		result[v.ID] = &v
	}

	return result

}
