package main

import "math"

func hypot(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func sub(x, y int) (z int) {
	z = x - y
	return
}

func first(x, _ int) int {
	return x
}

func zero(int, int) int {
	return 0
}