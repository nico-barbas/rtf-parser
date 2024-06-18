package main

// TODO(nico): Need a modifiable base font size somewhere
const baseFontSize int = 12

type MeasuringUnit int

const (
	MeasuringUnitPoint MeasuringUnit = iota
	MeasuringUnitTwip
	MeasuringUnitEm
	MeasuringUnitPixel
)

func ConvertUnits(value int, from MeasuringUnit, to MeasuringUnit) (result int) {
	switch from {
	case MeasuringUnitPoint:
		result = convertPoints(value, to)
	case MeasuringUnitTwip:
		pt := value / 20
		result = convertPoints(pt, to)
	case MeasuringUnitEm:
		// result = convertPoints(value, to)
	case MeasuringUnitPixel:
		// result = convertPoints(value, to)
	}

	return
}

func convertPoints(value int, to MeasuringUnit) (result int) {
	switch to {
	case MeasuringUnitPoint:
		result = value
	case MeasuringUnitEm:
		result = value / baseFontSize
	case MeasuringUnitPixel:
		result = value
	}
	return
}
