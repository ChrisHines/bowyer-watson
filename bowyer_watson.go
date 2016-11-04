// Package bowyer_watson provides an implementation of the Bowyer-Watson
// algorithm to create a Delaunay Triangulation. Given a set of points, it
// will output a set of Triangles.
package bowyer_watson

import (
	"math"
	"sort"
)

// Point represents a basic x,y coordinate.
type Point struct {
	X, Y float64
}

type pointsByX []Point

func (s pointsByX) Len() int           { return len(s) }
func (s pointsByX) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s pointsByX) Less(i, j int) bool { return s[i].X < s[j].X }

// Triangle contains three points that form a triangle.
type Triangle struct {
	A, B, C         Point
	center          Point
	radius, radius2 float64
}

// CalcCircumCircle calculates t's circumcircle and caches the results in t.
// It must be called before using CircumcircleContains.
func (t *Triangle) CalcCircumCircle() {
	ab := sqr(t.A.X) + sqr(t.A.Y)
	cd := sqr(t.B.X) + sqr(t.B.Y)
	ef := sqr(t.C.X) + sqr(t.C.Y)

	t.center.X = (ab*(t.C.Y-t.B.Y) + cd*(t.A.Y-t.C.Y) + ef*(t.B.Y-t.A.Y)) / (t.A.X*(t.C.Y-t.B.Y) + t.B.X*(t.A.Y-t.C.Y) + t.C.X*(t.B.Y-t.A.Y)) / 2
	t.center.Y = (ab*(t.C.X-t.B.X) + cd*(t.A.X-t.C.X) + ef*(t.B.X-t.A.X)) / (t.A.Y*(t.C.X-t.B.X) + t.B.Y*(t.A.X-t.C.X) + t.C.Y*(t.B.X-t.A.X)) / 2
	t.radius2 = sqr(t.A.X-t.center.X) + sqr(t.A.Y-t.center.Y)
	t.radius = math.Sqrt(t.radius2)
}

// HasVertex determine if p is one of t's vertices.
func (t *Triangle) HasVertex(p Point) bool {
	return t.A == p || t.B == p || t.C == p
}

// CircumcircleContains determines if p is contained within the circumcircle
// of t. A circumcircle is the circle whose circumference contains all 3
// vertices of a triangle.
func (t *Triangle) CircumcircleContains(p Point) bool {
	dist2 := sqr(p.X-t.center.X) + sqr(p.Y-t.center.Y)
	return dist2 <= t.radius2
}

// Edge is a line segment.
type Edge struct {
	A, B Point
}

// isEqual returns true if e2 is an equivalent edge to e1.
func (e1 Edge) isEqual(e2 Edge) bool {
	return (e1.A == e2.A && e1.B == e2.B || e1.A == e2.B && e1.B == e2.A)
}

// DelaunayTriangulation returns the triangles in the Delaunay triangulation
// of points. All elements of points must lie within super. Source for
// algorithm: paulbourke.net/papers/triangulate
func DelaunayTriangulation(points []Point, super Triangle) []Triangle {
	super.CalcCircumCircle()
	ts := []Triangle{super}

	pts := make([]Point, len(points))
	copy(pts, points)
	sort.Sort(pointsByX(pts))

	var result []Triangle
	var edges []Edge
	for _, p := range pts {
		edges = edges[:0]

		for i := 0; i < len(ts); {
			t := &ts[i]
			if p.X > t.center.X+t.radius {
				result = append(result, *t)
				n := len(ts) - 1
				ts[i] = ts[n]
				ts = ts[:n]
			} else if t.CircumcircleContains(p) {
				edges = append(edges,
					Edge{t.A, t.B},
					Edge{t.A, t.C},
					Edge{t.B, t.C},
				)
				n := len(ts) - 1
				ts[i] = ts[n]
				ts = ts[:n]
			} else {
				i++
			}
		}

		for j := 0; j < len(edges); j++ {
			for i := j + 1; i < len(edges); {
				if !edges[j].isEqual(edges[i]) {
					i++
					continue
				}
				n := len(edges) - 1
				edges[i] = edges[n]
				edges = edges[:n]
			}
		}

		for _, e := range edges {
			t := Triangle{A: e.A, B: e.B, C: p}
			t.CalcCircumCircle()
			ts = append(ts, t)
		}
	}

	result = append(result, ts...)

	//Remove any triangles using the Points of the super
	for i := 0; i < len(result); {
		t := result[i]
		if !t.HasVertex(super.A) && !t.HasVertex(super.B) && !t.HasVertex(super.C) {
			i++
			continue
		}
		n := len(result) - 1
		result[i] = result[n]
		result = result[:n]
	}

	return result
}

func sqr(x float64) float64 {
	return x * x
}
