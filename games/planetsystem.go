package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const SYSTEMSIZE int = 80

func main() {
	//rand.Seed(1)
	rand.Seed(time.Now().UnixNano())

	// Make a clean map
	system := [SYSTEMSIZE + 1][SYSTEMSIZE + 1]string{}
	for y := 0; y < SYSTEMSIZE+1; y++ {
		for x := 0; x < SYSTEMSIZE+1; x++ {
			system[y][x] = " "
		}
	}
	// Add corner posts
	//system[0][0] = "#"
	//system[0][SYSTEMSIZE] = "#"
	//system[SYSTEMSIZE][0] = "#"
	//system[SYSTEMSIZE][SYSTEMSIZE] = "#"

	// Add the system's star
	system[SYSTEMSIZE/2][SYSTEMSIZE/2] = "0"

	// Add asteroid belts
	belts := randIntRange(1, 5)
	for b := 0; b < belts; b++ {
		radius := float64(randIntRange(4, SYSTEMSIZE/2))
		//fmt.Printf("belt %d, range %.1f\n", b, radius)
		for i := 0; i < 100; i++ {
			randpi := randFLoatRange(2 * math.Pi)
			jitter := radius + float64(randIntRange(-1, 1))
			x := int(jitter*math.Sin(randpi)) + SYSTEMSIZE/2
			y := int(jitter*math.Cos(randpi)) + SYSTEMSIZE/2
			system[y][x] = "."
		}
	}

	// Add planets
	planets := randIntRange(2, 15)
	for p := 0; p < planets; p++ {
		x := randIntRange(3, SYSTEMSIZE)
		y := randIntRange(3, SYSTEMSIZE)
		// Add moons?
		//if rand.Float64() > 0.75 {
		radius := float64(randIntRange(1, 4))
		moons := randIntRange(0, 10)
		for i := 0; i < moons; i++ {
			randpi := randFLoatRange(2 * math.Pi)
			//jitter := radius + float64(randIntRange(-1, 1))
			jitter := radius
			rx := x + int(jitter*math.Sin(randpi)) // + SYSTEMSIZE/2
			ry := y + int(jitter*math.Cos(randpi)) // + SYSTEMSIZE/2
			if rx < 0 || rx > SYSTEMSIZE || ry < 0 || ry > SYSTEMSIZE {
				continue
			}
			system[ry][rx] = ","
		}
		//}
		system[y][x] = "o"
	}

	// Show map
	for y := 0; y < SYSTEMSIZE+1; y++ {
		for x := 0; x < SYSTEMSIZE+1; x++ {
			fmt.Printf(system[y][x])
		}
		fmt.Printf("\n")
	}
}

func randFLoatRange(max float64) float64 {
	return rand.Float64() * max
}

func randIntRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}
