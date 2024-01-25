package chessboard

import "fmt"

type Board struct {
	Pieces [][]*Piece `json:"pieces"`
}

func NewBoard(pieces ...*BoardPiece) *Board {
	board := make([][]*Piece, 8)
	for i := range board {
		board[i] = make([]*Piece, 8)
	}
	b := &Board{
		Pieces: board,
	}

	if len(pieces) == 0 {
		b.defaultPieces()
		return b
	}

	// place the pieces
	for _, piece := range pieces {
		b.Pieces[piece.Row][piece.Col] = NewPiece(piece.PieceType, piece.Color, piece.Row, piece.Col)
	}

	return b
}

func (c *Board) PlacePiece(piece *Piece, position Position) error {

	if !c.IsValidMove(piece, position) {
		return fmt.Errorf("this is not a valid move from [%d %d] to [%d %d]", piece.Position.Row, position.Col, position.Row, position.Col)
	}

	c.Pieces[piece.Position.Row][piece.Position.Col] = nil
	c.Pieces[position.Row][position.Col] = piece
	piece.Position = position

	return nil
}

func (c *Board) PlacePieceFromPosition(from, to Position) error {
	piece := c.GetPiece(from.Row, from.Col)

	if piece == nil {
		return fmt.Errorf("there is no piece in [%d %d]", from.Row, from.Col)
	}

	return c.PlacePiece(piece, to)
}

func (c *Board) GetPiece(row, col int) *Piece {
	if row > 7 || row < 0 || col > 7 || col < 0 {
		return nil
	}

	return c.Pieces[row][col]
}

func (c *Board) isPieceEmpty(row, col int) bool {
	return c.GetPiece(row, col) == nil
}

func (c *Board) PrintBoard() {
	for i, row := range c.Pieces {
		for j, piece := range row {
			if piece == nil {
				fmt.Printf("[  %d  %d  ]", i, j)
			} else {
				fmt.Printf("[ %s ]", string(piece.Color)+string(piece.Type))
			}
		}
		fmt.Println()
	}
}

func (c *Board) PrintPieceMoveableParts(position Position) {
	piece := c.Pieces[position.Row][position.Col]

	var moveablePositions []Position

	if piece != nil {
		moveablePositions = c.GetValidMoves(piece)
	}

	fmt.Println("Moves:", piece.Type)

	for i, row := range c.Pieces {
		for j, piece := range row {
			isTaken := false

			for _, takenPosition := range moveablePositions {
				if takenPosition.Row == i && takenPosition.Col == j {
					isTaken = true
					break
				}
			}

			if piece == nil {
				if isTaken {
					fmt.Printf("[--%d--%d--]", i, j)
				} else {
					fmt.Printf("[  %d  %d  ]", i, j)
				}
			} else {
				if isTaken {
					fmt.Printf("[-%s-]", string(piece.Color)+string(piece.Type))
				} else {
					fmt.Printf("[ %s ]", string(piece.Color)+string(piece.Type))
				}
			}
		}
		fmt.Println()
	}
}

func (c *Board) defaultPieces() {
	// Black Pieces

	c.Pieces[0][0] = NewPiece(Rook, Black, 0, 0)
	c.Pieces[0][1] = NewPiece(Knight, Black, 0, 1)
	c.Pieces[0][2] = NewPiece(Bishop, Black, 0, 2)
	c.Pieces[0][3] = NewPiece(Queen, Black, 0, 3)
	c.Pieces[0][4] = NewPiece(King, Black, 0, 4)
	c.Pieces[0][5] = NewPiece(Bishop, Black, 0, 5)
	c.Pieces[0][6] = NewPiece(Knight, Black, 0, 6)
	c.Pieces[0][7] = NewPiece(Rook, Black, 0, 7)

	for i := 0; i < 8; i++ {
		c.Pieces[1][i] = NewPiece(Pawn, Black, 1, i)
	}

	// White Pieces
	c.Pieces[7][0] = NewPiece(Rook, White, 7, 0)
	c.Pieces[7][1] = NewPiece(Knight, White, 7, 1)
	c.Pieces[7][2] = NewPiece(Bishop, White, 7, 2)
	c.Pieces[7][3] = NewPiece(Queen, White, 7, 3)
	c.Pieces[7][4] = NewPiece(King, White, 7, 4)
	c.Pieces[7][5] = NewPiece(Bishop, White, 7, 5)
	c.Pieces[7][6] = NewPiece(Knight, White, 7, 6)
	c.Pieces[7][7] = NewPiece(Rook, White, 7, 7)

	for i := 0; i < 8; i++ {
		c.Pieces[6][i] = NewPiece(Pawn, White, 6, i)
	}

}

type BoardPiece struct {
	Color     Color
	PieceType PieceType
	Row, Col  int
}

func (c *Board) GetPiecesPositions() []*BoardPiece {
	positions := make([]*BoardPiece, 0)

	for i, row := range c.Pieces {
		for j, piece := range row {

			if piece != nil {
				positions = append(positions, &BoardPiece{
					Color:     piece.Color,
					PieceType: piece.Type,
					Row:       i,
					Col:       j,
				})
			}
		}
	}

	return positions
}
