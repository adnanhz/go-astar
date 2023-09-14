package internal

import (
	"image"
	"math"
	"sync/atomic"

	"github.com/fzipp/astar"
)

type Node struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Id string  `json:"id"`
}

type Vertex struct {
	Start Node `json:"start"`
	End   Node `json:"end"`
}

type graph[Node comparable] map[Node][]Node

func (node *Node) ToImagePoint() image.Point {
	return image.Pt(int(node.X), int(node.Y))
}

func newGraph[Node comparable]() graph[Node] {
	return make(map[Node][]Node)
}

func (g graph[Node]) link(a, b Node) graph[Node] {
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
	return g
}

func (g graph[Node]) Neighbours(n Node) []Node {
	return g[n]
}

func nodeDist(p, q image.Point) float64 {
	d := q.Sub(p)
	return math.Sqrt(float64(d.X*d.X + d.Y*d.Y))
}

func CalculatePath(graph []Vertex, start Node, end Node, passes *int32) astar.Path[image.Point] {
	g := newGraph[image.Point]()
	for _, vertex := range graph {
		g.link(vertex.Start.ToImagePoint(), vertex.End.ToImagePoint())
	}
	path := astar.FindPath[image.Point](g, start.ToImagePoint(), end.ToImagePoint(), nodeDist, nodeDist)
	atomic.AddInt32(passes, 1)
	return path
}
