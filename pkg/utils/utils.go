package utils

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
)

func ParseTestCase(path string) (result TestCase, err error) {
	file, err := os.Open(path)
	if err != nil {
		return TestCase{}, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var line string
	line, err = reader.ReadString('\n')
	if err != nil {
		return TestCase{}, err
	}

	line = strings.TrimSpace(line)

	var initialPositions []InitialPosition
	for line != "" {
		lineParts := strings.Split(line, " ")
		initialPositions = append(initialPositions, InitialPosition{lineParts[0], lineParts[1]})
		line, err = reader.ReadString('\n')
		if err != nil {
			return TestCase{}, err
		}

		line = strings.TrimSpace(line)
	}

	line, _ = reader.ReadString('\n')
	line = strings.TrimSpace(line)
	whiteCatures := strings.Split(line[1:len(line)-1], " ")

	line, _ = reader.ReadString('\n')
	line = strings.TrimSpace(line)
	blackCatures := strings.Split(line[1:len(line)-1], " ")

	line, _ = reader.ReadString('\n')
	line = strings.TrimSpace(line)

	var moves []string
	for line != "" {
		line = strings.TrimSpace(line)
		moves = append(moves, line)
		line, _ = reader.ReadString('\n')
	}

	return TestCase{
		WhiteCapture:     whiteCatures,
		BlackCapture:     blackCatures,
		Moves:            moves,
		InitialPositions: initialPositions,
	}, nil
}

func StringifyBoard(board [][]string) string {

	row := len(board)
	col := len(board[0])

	var buffer bytes.Buffer

	for i := 0; i < col; i++ {
		if i == 0 {
			buffer.WriteString("   ")
		}
		colLetter := (string)(i + 'a')
		buffer.WriteString(" " + colLetter + "  ")
	}
	buffer.WriteString("\n")

	for i := row - 1; i >= 0; i-- {
		buffer.WriteString(strconv.Itoa(i + 1))
		buffer.WriteString(" |")

		for j := 0; j < col; j++ {
			buffer.WriteString(stringifySquare(board[len(board)-i-1][j]))
		}

		buffer.WriteString(" " + strconv.Itoa(i+1))
		buffer.WriteString("\n")

		if i != 0 {
			buffer.WriteString("\n")
		}
	}

	for i := 0; i < col; i++ {
		if i == 0 {
			buffer.WriteString("   ")
		}
		colLetter := (string)(i + 'a')
		buffer.WriteString(" " + colLetter + "  ")
	}
	buffer.WriteString("\n")

	return buffer.String()
}

func stringifySquare(sq string) string {
	if len(sq) == 0 {
		return " _ |"
	}

	return " " + sq + " |"
}
