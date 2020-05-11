package game

import (
	"math"
)

func DegreesToX(Direction float64) (XVal float64) {
	XVal = math.Cos(float64(Direction * (math.Pi / 180)))
	//	YVal = math.sin(Direction)
	return (XVal)
}

func DegreesToY(Direction float64) (YVal float64) {
	YVal = math.Sin(float64(Direction * (math.Pi / 180)))
	//	YVal = math.sin(Direction)
	return (YVal)
}
