package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// Generate test data
	const (
		numPoints  = 1000
		binsX      = 10
		binsY      = 10
		minX, maxX = 0.0, 10.0
		minY, maxY = 0.0, 10.0
	)

	dataX := make([]float64, numPoints)
	dataY := make([]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		dataX[i] = rand.Float64()*(maxX-minX) + minX
		dataY[i] = rand.Float64()*(maxY-minY) + minY
	}

	// Calculate indices
	_, err := CalculateIndices1D(binsX, minX, maxX, dataX)
	if err != nil {
		fmt.Println("Error calculating indices for X:", err)
		return
	}
	//fmt.Println(indicesX)
	_, err = CalculateIndices1D(binsY, minY, maxY, dataY)
	if err != nil {
		fmt.Println("Error calculating indices for Y:", err)
		return
	}
	//fmt.Println(indicesY)

	// Calculate mutual information
	shiftFrom, shiftTo := -2, 2
	shiftStep := 1
	mi, err := ShiftedMutualInformation(shiftFrom, shiftTo, binsX, binsY, minX, maxX, minY, maxY, dataX, dataY, shiftStep)
	if err != nil {
		fmt.Println("Error calculating mutual information:", err)
		return
	}

	fmt.Println("Mutual Information for each shift:")
	for i, val := range mi {
		fmt.Printf("Shift %d: %.6f\n", shiftFrom+i*shiftStep, val)
	}
}
