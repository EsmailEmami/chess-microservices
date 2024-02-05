package chessboard

func (c *Chessboard) IsValidMove(piece *Piece, to Position) bool {
	// Check if the 'to' position is within the board boundaries
	if to.Row < 0 || to.Row >= 8 || to.Col < 0 || to.Col >= 8 {
		return false
	}

	from := piece.Position

	// Check if there is a piece at the 'from' position
	if c.GetPiece(from.Row, from.Col) != piece {
		return false
	}

	// Check if there is no piece at the 'to' position or if the piece at 'to' is of the opposite color
	if !c.isEmptyPiece(to.Row, to.Col) && c.GetPiece(to.Row, to.Col).Color == piece.Color {
		return false
	}

	switch piece.Type {
	case Pawn:
		return c.isValidMovePawn(piece, from, to)

	case Rook:
		return c.isValidMoveRook(from, to)

	case Knight:
		return c.isValidMoveKnight(from, to)

	case Bishop:
		return c.isValidMoveBishop(from, to)

	case Queen:
		return c.isValidMoveQueen(from, to)

	case King:
		return c.isValidMoveKing(from, to)
	default:
		return false
	}
}

func (c *Chessboard) isValidMovePawn(piece *Piece, from, to Position) bool {
	var (
		isToEmptyPiece = c.isEmptyPiece(to.Row, to.Col)
		toPiece        = c.GetPiece(to.Row, to.Col)
	)

	if piece.Color == White {
		// Moving one square forward
		if to.Row == from.Row-1 && to.Col == from.Col && isToEmptyPiece {
			return true
		}
		// Moving two squares forward on the first move
		if from.Row == 6 && to.Row == 4 && to.Col == from.Col && c.isEmptyPiece(5, to.Col) && isToEmptyPiece {
			return true
		}
		// Capturing diagonally
		if to.Row == from.Row-1 && abs(to.Col-from.Col) == 1 &&
			!isToEmptyPiece && toPiece.Color != piece.Color {
			return true
		}
	} else {
		// Moving one square forward
		if to.Row == from.Row+1 && to.Col == from.Col && isToEmptyPiece {
			return true
		}
		// Moving two squares forward on the first move
		if from.Row == 1 && to.Row == 3 && to.Col == from.Col && c.isEmptyPiece(2, to.Col) && isToEmptyPiece {
			return true
		}
		// Capturing diagonally
		if to.Row == from.Row+1 && abs(to.Col-from.Col) == 1 &&
			!isToEmptyPiece && toPiece.Color != piece.Color {
			return true
		}
	}

	return false
}

func (c *Chessboard) isValidMoveRook(from, to Position) bool {
	return (to.Row == from.Row || to.Col == from.Col) && !c.isPathBlocked(from, to)
}

func (c *Chessboard) isValidMoveKnight(from, to Position) bool {
	return (abs(to.Row-from.Row) == 2 && abs(to.Col-from.Col) == 1) ||
		(abs(to.Row-from.Row) == 1 && abs(to.Col-from.Col) == 2)
}

func (c *Chessboard) isValidMoveBishop(from, to Position) bool {
	return abs(to.Row-from.Row) == abs(to.Col-from.Col) && !c.isPathBlocked(from, to)
}

func (c *Chessboard) isValidMoveQueen(from, to Position) bool {
	return (to.Row == from.Row || to.Col == from.Col || abs(to.Row-from.Row) == abs(to.Col-from.Col)) && !c.isPathBlocked(from, to)
}
func (c *Chessboard) isValidMoveKing(from, to Position) bool {
	return abs(to.Row-from.Row) <= 1 && abs(to.Col-from.Col) <= 1
}

func (c *Chessboard) isPathBlocked(from, to Position) bool {
	if from.Row == to.Row {
		start, end := min(from.Col, to.Col)+1, max(from.Col, to.Col)
		for i := start; i < end; i++ {
			if !c.isEmptyPiece(from.Row, i) {
				return true
			}
		}
	} else if from.Col == to.Col {
		start, end := min(from.Row, to.Row)+1, max(from.Row, to.Row)
		for i := start; i < end; i++ {
			if !c.isEmptyPiece(i, from.Col) {
				return true
			}
		}
	} else {
		// Check for diagonal path
		rowInc := 1
		if to.Row < from.Row {
			rowInc = -1
		}
		colInc := 1
		if to.Col < from.Col {
			colInc = -1
		}
		for i, j := from.Row+rowInc, from.Col+colInc; i != to.Row && j != to.Col; i, j = i+rowInc, j+colInc {
			if !c.isEmptyPiece(i, j) {
				return true
			}
		}
	}

	return false
}
