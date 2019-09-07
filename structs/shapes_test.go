package structs

import (
	"math"
	"testing"
)

func TestPerimeter(t *testing.T) {
	got := Perimeter(Rectangle{10.0, 10.0})
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f ", got, want)
	}
}

func TestArea(t *testing.T) {
	circle := Circle{10.0}
	rectangle := Rectangle{10.0, 10.0}
	triangle := Triangle{10.0, 2}

	areaTests := []struct {
		name    string
		shape   Shape
		hasArea float64
	}{
		{"rectangle", rectangle, 100.0},
		{"circle", circle, math.Pi * math.Pow(10, 2)},
		{"triangle", triangle, 10},
	}

	for _, tt := range areaTests {
		got := tt.shape.Area()
		want := tt.hasArea
		if got != want {
			t.Errorf("%#v got %.2f want %.2f", tt.shape, got, want)
		}
	}
}
