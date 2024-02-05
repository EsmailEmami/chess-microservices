package chessboard

import (
	"errors"
	"strconv"
)

func getOpponentColor(color Color) Color {
	if color == White {
		return Black
	}
	return White
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func GetPosition(position string) (*Position, error) {
	if len(position) != 2 {
		return nil, errors.New("invalid position format")
	}

	row, err := strconv.Atoi(string(position[1]))
	if err != nil || row < 1 || row > 8 {
		return nil, errors.New("invalid row value")
	}

	col := int(position[0] - 'a')
	if col < 0 || col > 7 {
		return nil, errors.New("invalid col value")
	}

	return &Position{
		Row: row - 1,
		Col: col,
	}, nil
}

func isValidArea(row, col int) bool {
	return !(row > 7 || row < 0 || col > 7 || col < 0)
}
