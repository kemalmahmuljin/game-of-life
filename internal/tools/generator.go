package tools

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GenerateRandomBoard(rows, cols int, density float64, seed int64) (file string){
	// Guarantee reproducibility
	rand.Seed(seed)

	var den = int(100*density)

	file = fmt.Sprintf("%s%d%s%d%s%d%s%d%s", "setup/World_", rows, "_", cols, "_", den, "_", seed, ".txt")
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(file)
		check(err)
		defer f.Close()

		fmt.Fprintf(f, "%d,%d\n", rows, cols)

		for i:=0; i<rows; i++ {
			for j:=0; j<cols; j++ {
				if rand.Float64() < density {
					fmt.Fprintf(f, "%d,%d\n", i,j)
				}
			}
		}
	}
	return file
}
