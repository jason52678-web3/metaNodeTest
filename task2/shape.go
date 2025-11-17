package main

import (
	"fmt"
	"math"
)

type Share interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	height float64
	width  float64
}

func (r *Rectangle) Area() float64 {
	return r.height * r.width
}
func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.height + r.width)
}

type Circle struct {
	radius float64
}

func (c *Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c *Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func main() {
	rectangle := Rectangle{height: 10, width: 20}
	circle := Circle{radius: 5}

	fmt.Println("rectangle Area: ", rectangle.Area())
	fmt.Println("Rectangle Perimeter: ", rectangle.Perimeter())
	fmt.Println("circle Area: ", circle.Area())
	fmt.Println("circle Perimeter: ", circle.Perimeter())
}
