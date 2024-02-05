package chessboard

import "fmt"

type ChessboardPiece struct {
	Color     Color
	PieceType PieceType
	Row, Col  int
}

type ChessBoardMove struct {
	From Position
	To   Position
}
type Chessboard struct {
	Pieces     [8][8]*Piece
	MovesCount [8][8]int
}

func NewDefault() *Chessboard {
	board := new(Chessboard)
	board.setupDefult()
	return board
}

func New(pieces []*ChessboardPiece, moves []*ChessBoardMove) *Chessboard {
	board := &Chessboard{}

	if len(pieces) == 0 {
		board.setupDefult()
		return board
	}

	for _, piece := range pieces {
		board.Pieces[piece.Row][piece.Col] = NewPiece(piece.PieceType, piece.Color, piece.Row, piece.Col)
	}

	for _, move := range moves {
		board.increaseMovesCount(move.From, move.To)
	}

	return board
}

func (c *Chessboard) increaseMovesCount(from, to Position) {
	c.MovesCount[from.Row][from.Col]++
	c.MovesCount[to.Row][to.Col]++
}

func (c *Chessboard) PlacePiece(piece *Piece, position Position) error {
	isValidMove := false
	validMoves := c.GetValidMoves(piece)

	for _, move := range validMoves {
		if move.Row == position.Row && move.Col == position.Col {
			isValidMove = true
			break
		}
	}

	if !isValidMove {
		return fmt.Errorf("this is not a valid move from [%d %d] to [%d %d]", piece.Position.Row, position.Col, position.Row, position.Col)
	}

	c.increaseMovesCount(piece.Position, position)

	c.Pieces[piece.Position.Row][piece.Position.Col] = nil
	c.Pieces[position.Row][position.Col] = piece
	piece.Position = position

	return nil
}

func (c *Chessboard) PlacePieceFromPosition(from, to Position) error {
	piece := c.GetPiece(from.Row, from.Col)

	if piece == nil {
		return fmt.Errorf("there is no piece in [%d %d]", from.Row, from.Col)
	}

	return c.PlacePiece(piece, to)
}

func (c *Chessboard) GetPiece(row, col int) *Piece {
	if !isValidArea(row, col) {
		return nil
	}

	return c.Pieces[row][col]
}

func (c *Chessboard) isEmptyPiece(row, col int) bool {
	return c.GetPiece(row, col) == nil
}

func (c *Chessboard) setupDefult() {
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

func (c *Chessboard) GetPieces() []*ChessboardPiece {
	positions := make([]*ChessboardPiece, 0)

	for i, row := range c.Pieces {
		for j, piece := range row {

			if piece != nil {
				positions = append(positions, &ChessboardPiece{
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

func (c *Chessboard) PrintChessboard() {
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

func (c *Chessboard) PrintPieceMoveableParts(position Position) {
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
