package bowyer_watson

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

const pc int = 6

func getRandomPointInCircle(radius float64) (float64, float64) {
	t := 2 * math.Pi * rand.Float64()
	u := rand.Float64() + rand.Float64()
	r := 0.0
	if u > 1 {
		r = 2.0 - u
	} else {
		r = u
	}
	return radius * r * math.Cos(t), radius * r * math.Sin(t)
}

func TestDelaunayTriangulation1(t *testing.T) {
	points := make([]Point, pc)

	for i := 0; i < len(points); i++ {
		x, y := getRandomPointInCircle(5)
		points[i] = Point{x, y}
	}

	super := Triangle{
		A: Point{0, 50},
		B: Point{50, -50},
		C: Point{-50, -50},
	}

	u := DelaunayTriangulation(points, super)

	t.Log(u)

	fmt.Println("number of triangles", len(u))
}

func TestDelaunayTriangulation2(t *testing.T) {
	N := 4
	points := make([]Point, N)

	for i := 0; i < N; i++ {
		x, y := getRandomPointInCircle(3.0)
		points[i] = Point{x, y}
	}

	super := Triangle{
		A: Point{0, 50},
		B: Point{50, -50},
		C: Point{-50, -50},
	}

	u := DelaunayTriangulation(points, super)

	if got, want := len(u), 2; got != want {
		t.Errorf("#triangles: got %v, want %v", got, want)
	}
}
