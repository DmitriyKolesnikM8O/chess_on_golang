package utils

import (
	"bufio"
	"os"
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
