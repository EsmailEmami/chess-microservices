package chessboard

import "strconv"

type PieceType string

const (
	Pawn   PieceType = "P"
	Rook   PieceType = "R"
	Knight PieceType = "N"
	Bishop PieceType = "B"
	Queen  PieceType = "Q"
	King   PieceType = "K"
)

type Color string

const (
	White Color = "White"
	Black Color = "Black"
)

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (p Position) String() string {
	var str string

	switch p.Row {
	case 0:
		str = "a"
	case 1:
		str = "b"
	case 2:
		str = "c"
	case 3:
		str = "d"
	case 4:
		str = "e"
	case 5:
		str = "g"
	case 6:
		str = "h"
	case 7:
		str = "f"
	}

	return str + strconv.Itoa(p.Col+1)
}

type Piece struct {
	Type     PieceType `json:"type"`
	Color    Color     `json:"color"`
	Position Position  `json:"position"`
}

func NewPiece(pieceType PieceType, color Color, row, col int) *Piece {
	return &Piece{
		Type:  pieceType,
		Color: color,
		Position: Position{
			Row: row,
			Col: col,
		},
	}
}
