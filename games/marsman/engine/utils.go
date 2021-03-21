package engine

import (
	"log"
	"math/rand"
	"os"

	osimplex "github.com/ojrac/opensimplex-go"
)

var logfile *os.File
var noise *osimplex.OpenSimplexNoise

func check_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OpenLog(file string) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening logfile: %v\n", err)
	}
	logfile = f
	log.SetOutput(logfile)
	log.Println("LOG STARTED")
}

func CloseLog() {
	if logfile == nil {
		return
	}
	log.Println("LOG STOPPED")
	logfile.Close()
}

func RandRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func Prob(chance int) bool {
	return RandRange(0, 100) <= chance
}

func InitSeeds(seed int64) {
	rand.Seed(seed)
	noise = osimplex.NewOpenSimplexWithSeed(seed)
}

func Noise(x, y float64) float64 {
	scale := 25.0
	sx, sy := x/scale, y/scale
	return noise.Eval2(sx, sy)
}
