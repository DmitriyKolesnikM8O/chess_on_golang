package utils

type TestCase struct {
	WhiteCapture     []string
	BlackCapture     []string
	Moves            []string
	InitialPositions []InitialPosition
}

type InitialPosition struct {
	Sign     string
	Position string
}
