package main

import (
	"fmt"
	"math"
)

type Point struct {
	x int
	y int
}

func (p Point) distance() float64 {
	return math.Sqrt(float64(p.x*p.x + p.y*p.y))
}

func (p Point) noop(d int) {
	p.x += d
	p.y += d
}

func (p *Point) shift(d int) {
	p.x += d
	p.y += d
}

func main() {
	p1 := Point{x: 0, y: 0}
	fmt.Println(p1)

	p1.shift(2)
	fmt.Println(p1)

	p1.noop(3)
	fmt.Println(p1)
}
