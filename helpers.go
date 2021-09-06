package main

import (
	"math/rand"
	"time"
)

// seed random number
var R = rand.New(rand.NewSource(time.Now().UnixNano()))

// ApplyJitter return a random number
func ApplyJitter(input int) (output int) {
	deviation := int(0.25 * float64(input))
	return input - deviation + rand.Intn(2*deviation)
}

func FindMinAndMax(a []int) (min int, max int) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func Sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}
