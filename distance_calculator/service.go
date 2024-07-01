package main

import (
	"fmt"
	"math"

	"github.com/64bitAryan/go-microservice/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUDATA) (float64, error)
}

type CalculatorService struct {
	points [][]float64
}

func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{
		points: make([][]float64, 0),
	}
}

func (s *CalculatorService) CalculateDistance(data types.OBUDATA) (float64, error) {
	distance := 0.0
	if len(s.points) > 0 {
		prevPoint := s.points[len(s.points)-1]
		distance = calculateDistance(prevPoint[0], data.Lat, prevPoint[1], data.Long)
	}
	s.points = append(s.points, []float64{data.Lat, data.Long})
	fmt.Println("Calculating the distance")
	return distance, nil
}

func calculateDistance(x1, x2, y1, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
