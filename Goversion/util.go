package main

import (
	"errors"
	"math"
	"sync"
)

type indexPair struct {
	First  int
	Second int
}

type histogram2D struct {
	BinsX int
	BinsY int
	MinX  float64
	MaxX  float64
	MinY  float64
	MaxY  float64
	Data  [][]int
	Mutex sync.Mutex
}

func NewHistogram2D(binsX, binsY int, minX, maxX, minY, maxY float64) *histogram2D {
	data := make([][]int, binsX)
	for i := range data {
		data[i] = make([]int, binsY)
	}
	return &histogram2D{
		BinsX: binsX,
		BinsY: binsY,
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		Data:  data,
	}
}

func (h *histogram2D) Increment(x, y float64) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	indexX := int((x - h.MinX) / (h.MaxX - h.MinX) * float64(h.BinsX))
	if indexX == h.BinsX {
		indexX--
	}

	indexY := int((y - h.MinY) / (h.MaxY - h.MinY) * float64(h.BinsY))
	if indexY == h.BinsY {
		indexY--
	}

	h.Data[indexX][indexY]++
}

func (h *histogram2D) CalculateMutualInformation() float64 {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	total := 0
	for i := 0; i < h.BinsX; i++ {
		for j := 0; j < h.BinsY; j++ {
			total += h.Data[i][j]
		}
	}

	var hx, hy float64
	for i := 0; i < h.BinsX; i++ {
		px := float64(0)
		for j := 0; j < h.BinsY; j++ {
			px += float64(h.Data[i][j]) / float64(total)
		}
		if px != 0 {
			hx -= px * math.Log2(px)
		}
	}

	for j := 0; j < h.BinsY; j++ {
		py := float64(0)
		for i := 0; i < h.BinsX; i++ {
			py += float64(h.Data[i][j]) / float64(total)
		}
		if py != 0 {
			hy -= py * math.Log2(py)
		}
	}

	var hxy float64
	for i := 0; i < h.BinsX; i++ {
		for j := 0; j < h.BinsY; j++ {
			p := float64(h.Data[i][j]) / float64(total)
			if p != 0 {
				hxy -= p * math.Log2(p)
			}
		}
	}

	return hx + hy - hxy
}

func CalculateIndices1D(bins int, min, max float64, data []float64) ([]int, error) {
	if min >= max {
		return nil, errors.New("min has to be smaller than max")
	}
	if bins < 1 {
		return nil, errors.New("there must be at least one bin")
	}

	indices := make([]int, len(data))
	for i, value := range data {
		if value < min || value > max {
			indices[i] = -1 // Indicates out of range
			continue
		}
		index := int((value - min) / (max - min) * float64(bins))
		if index == bins {
			index--
		}
		indices[i] = index
	}

	return indices, nil
}

func CalculateIndices2D(binsX, binsY int, minX, maxX, minY, maxY float64, dataX, dataY []float64) ([]indexPair, error) {
	if minX >= maxX {
		return nil, errors.New("minX has to be smaller than maxX")
	}
	if minY >= maxY {
		return nil, errors.New("minY has to be smaller than maxY")
	}
	if binsX < 1 {
		return nil, errors.New("there must be at least one binX")
	}
	if binsY < 1 {
		return nil, errors.New("there must be at least one binY")
	}
	if len(dataX) != len(dataY) {
		return nil, errors.New("dataX and dataY must have the same size")
	}

	indices := make([]indexPair, len(dataX))
	for i := range dataX {
		if dataX[i] < minX || dataX[i] > maxX || dataY[i] < minY || dataY[i] > maxY {
			indices[i] = indexPair{First: -1, Second: -1} // Indicates out of range
			continue
		}
		indexX := int((dataX[i] - minX) / (maxX - minX) * float64(binsX))
		if indexX == binsX {
			indexX--
		}
		indexY := int((dataY[i] - minY) / (maxY - minY) * float64(binsY))
		if indexY == binsY {
			indexY--
		}
		indices[i] = indexPair{First: indexX, Second: indexY}
	}

	return indices, nil
}

func ShiftedMutualInformation(shiftFrom, shiftTo, binsX, binsY int, minX, maxX, minY, maxY float64, dataX, dataY []float64, shiftStep int) ([]float64, error) {
	if shiftFrom >= shiftTo {
		return nil, errors.New("shiftFrom has to be smaller than shiftTo")
	}
	if minX >= maxX {
		return nil, errors.New("minX has to be smaller than maxX")
	}
	if minY >= maxY {
		return nil, errors.New("minY has to be smaller than maxY")
	}
	if binsX < 1 || binsY < 1 {
		return nil, errors.New("there must be at least one binX and one binY")
	}
	if len(dataX) != len(dataY) {
		return nil, errors.New("dataX and dataY must have the same size")
	}
	if shiftStep < 1 {
		return nil, errors.New("shiftStep must be greater or equal 1")
	}

	var wg sync.WaitGroup
	numShifts := (shiftTo-shiftFrom)/shiftStep + 1
	mi := make([]float64, numShifts)

	for i := shiftFrom; i <= shiftTo; i += shiftStep {
		wg.Add(1)
		go func(shift int) {
			defer wg.Done()

			hist := NewHistogram2D(binsX, binsY, minX, maxX, minY, maxY)

			for j := 0; j < len(dataX); j++ {
				x := dataX[j]
				y := dataY[j]

				if shift < 0 {
					if j < -shift {
						continue
					}
					x = dataX[j+shift]
					y = dataY[j]
				} else if shift > 0 {
					if j >= len(dataX)-shift {
						continue
					}
					x = dataX[j]
					y = dataY[j+shift]
				}

				hist.Increment(x, y)
			}

			mi[(shift-shiftFrom)/shiftStep] = hist.CalculateMutualInformation()
		}(i)
	}

	wg.Wait()
	return mi, nil
}
