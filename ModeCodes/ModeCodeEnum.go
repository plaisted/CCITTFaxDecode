package modecodes

type Mode uint8

const (
	Pass Mode = iota + 1
	Horizontal
	VerticalZero
	VerticalR1
	VerticalR2
	VerticalR3
	VerticalL1
	VerticalL2
	VerticalL3
	Extension
)
