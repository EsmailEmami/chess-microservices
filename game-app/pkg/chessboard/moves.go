package chessboard

func (c *Board) GetValidMovesFromPosition(position Position) []Position {
	piece := c.GetPiece(position.Row, position.Col)

	if piece == nil {
		return []Position{}
	}

	return c.GetValidMoves(piece)
}

func (c *Board) GetValidMoves(piece *Piece) []Position {
	switch piece.Type {
	case Pawn:
		return c.getValidMovesForPawn(piece)

	case Rook:
		return c.getValidMovesForRook(piece)

	case Knight:
		return c.getValidMovesForKnight(piece)

	case Bishop:
		return c.getValidMovesForBishop(piece)

	case Queen:
		return c.getValidMovesForQueen(piece)

	case King:
		return c.getValidMovesForKing(piece)

	default:
		return nil
	}
}

func (c *Board) getValidMovesForPawn(piece *Piece) []Position {
	var (
		position   = piece.Position
		validMoves = make([]Position, 0)
	)

	// Determine the direction of movement based on the pawn's color
	direction := 1
	if piece.Color == White {
		direction = -1
	}

	// Move one square forward
	newPosition := Position{position.Row + direction, position.Col}
	if c.IsValidMove(piece, newPosition) && c.isPieceEmpty(newPosition.Row, newPosition.Col) {
		validMoves = append(validMoves, newPosition)
	}

	// Move two squares forward on the first move
	if (piece.Color == White && position.Row == 6) || (piece.Color == Black && position.Row == 1) {
		newPosition := Position{position.Row + 2*direction, position.Col}
		if c.IsValidMove(piece, newPosition) && c.isPieceEmpty(newPosition.Row, newPosition.Col) {
			validMoves = append(validMoves, newPosition)
		}
	}

	// Capture diagonally
	capturePositions := []Position{{position.Row + direction, position.Col - 1}, {position.Row + direction, position.Col + 1}}
	for _, capturePos := range capturePositions {
		if c.IsValidMove(piece, capturePos) && !c.isPieceEmpty(capturePos.Row, capturePos.Col) &&
			c.GetPiece(capturePos.Row, capturePos.Col).Color != piece.Color {
			validMoves = append(validMoves, capturePos)
		}
	}

	return validMoves
}

func (c *Board) getValidMovesForRook(piece *Piece) []Position {
	var (
		position   = piece.Position
		validMoves = make([]Position, 0)
	)

	// Check horizontally
	for i := 0; i < 8; i++ {
		newPosition := Position{position.Row, i}
		if c.IsValidMove(piece, newPosition) {
			validMoves = append(validMoves, newPosition)
		}
	}

	// Check vertically
	for i := 0; i < 8; i++ {
		newPosition := Position{i, position.Col}
		if c.IsValidMove(piece, newPosition) {
			validMoves = append(validMoves, newPosition)
		}
	}

	return validMoves
}

func (c *Board) getValidMovesForKnight(piece *Piece) []Position {
	var (
		position   = piece.Position
		validMoves = make([]Position, 0)
	)

	possibleMoves := []Position{
		{position.Row - 2, position.Col - 1}, {position.Row - 2, position.Col + 1},
		{position.Row - 1, position.Col - 2}, {position.Row - 1, position.Col + 2},
		{position.Row + 1, position.Col - 2}, {position.Row + 1, position.Col + 2},
		{position.Row + 2, position.Col - 1}, {position.Row + 2, position.Col + 1},
	}

	for _, move := range possibleMoves {
		if c.IsValidMove(piece, move) {
			validMoves = append(validMoves, move)
		}
	}

	return validMoves
}

func (c *Board) getValidMovesForBishop(piece *Piece) []Position {
	var (
		position   = piece.Position
		validMoves = make([]Position, 0)
	)

	// Check diagonally
	for i := 0; i < 8; i++ {
		newPosition1 := Position{position.Row + i, position.Col + i}
		newPosition2 := Position{position.Row + i, position.Col - i}
		newPosition3 := Position{position.Row - i, position.Col + i}
		newPosition4 := Position{position.Row - i, position.Col - i}

		if c.IsValidMove(piece, newPosition1) {
			validMoves = append(validMoves, newPosition1)
		}
		if c.IsValidMove(piece, newPosition2) {
			validMoves = append(validMoves, newPosition2)
		}
		if c.IsValidMove(piece, newPosition3) {
			validMoves = append(validMoves, newPosition3)
		}
		if c.IsValidMove(piece, newPosition4) {
			validMoves = append(validMoves, newPosition4)
		}
	}

	return validMoves
}

func (c *Board) getValidMovesForQueen(piece *Piece) []Position {
	rookMoves := c.getValidMovesForRook(piece)
	bishopMoves := c.getValidMovesForBishop(piece)
	return append(rookMoves, bishopMoves...)
}

func (c *Board) getValidMovesForKing(piece *Piece) []Position {
	var (
		position   = piece.Position
		validMoves = make([]Position, 0)
	)

	// Check in all eight directions
	directions := []Position{
		{1, 0}, {1, 1}, {0, 1}, {-1, 1},
		{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
	}

	for _, dir := range directions {
		newPosition := Position{position.Row + dir.Row, position.Col + dir.Col}
		if c.IsValidMove(piece, newPosition) {
			validMoves = append(validMoves, newPosition)
		}
	}

	return validMoves
}
