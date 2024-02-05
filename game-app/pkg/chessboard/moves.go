package chessboard

// GetValidMoves returns valid moves for a piece at the specified position
func (c *Chessboard) GetValidMoves(piece *Piece) []Position {
	moves := c.calculateValidMoves(piece)

	if c.IsInCheck(piece.Color) {
		var newValidMoves []Position

		for i := 0; i < len(moves); i++ {
			if !c.wouldMoveResultInCheck(piece.Color, piece.Position, moves[i]) {
				newValidMoves = append(newValidMoves, moves[i])
			}
		}

		return newValidMoves
	} else {
		return moves
	}
}

func (c *Chessboard) GetValidMovesFromPosition(position Position) []Position {
	piece := c.GetPiece(position.Row, position.Col)

	if piece == nil {
		return []Position{}
	}

	return c.GetValidMoves(piece)
}

func (c *Chessboard) calculateValidMoves(piece *Piece) []Position {
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

func (c *Chessboard) getValidMovesForPawn(piece *Piece) []Position {
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
	if c.IsValidMove(piece, newPosition) && c.isEmptyPiece(newPosition.Row, newPosition.Col) {
		validMoves = append(validMoves, newPosition)
	}

	// Move two squares forward on the first move
	if (piece.Color == White && position.Row == 6) || (piece.Color == Black && position.Row == 1) {
		newPosition := Position{position.Row + 2*direction, position.Col}
		if c.IsValidMove(piece, newPosition) && c.isEmptyPiece(newPosition.Row, newPosition.Col) {
			validMoves = append(validMoves, newPosition)
		}
	}

	// Capture diagonally
	capturePositions := []Position{{position.Row + direction, position.Col - 1}, {position.Row + direction, position.Col + 1}}
	for _, capturePos := range capturePositions {
		if c.IsValidMove(piece, capturePos) && !c.isEmptyPiece(capturePos.Row, capturePos.Col) &&
			c.GetPiece(capturePos.Row, capturePos.Col).Color != piece.Color {
			validMoves = append(validMoves, capturePos)
		}
	}

	return validMoves
}

func (c *Chessboard) getValidMovesForRook(piece *Piece) []Position {
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

func (c *Chessboard) getValidMovesForKnight(piece *Piece) []Position {
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

func (c *Chessboard) getValidMovesForBishop(piece *Piece) []Position {
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

func (c *Chessboard) getValidMovesForQueen(piece *Piece) []Position {
	rookMoves := c.getValidMovesForRook(piece)
	bishopMoves := c.getValidMovesForBishop(piece)
	return append(rookMoves, bishopMoves...)
}

func (c *Chessboard) getValidMovesForKing(piece *Piece) []Position {
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

	// Check for castling moves
	castlingMoves := c.getValidCastlingMoves(piece.Color, position)
	validMoves = append(validMoves, castlingMoves...)

	return validMoves
}

// getValidCastlingMoves returns valid castling moves for the king
func (c *Chessboard) getValidCastlingMoves(color Color, kingPos Position) []Position {
	validMoves := make([]Position, 0)

	// Check for kingside castling
	if c.canCastleKingside(color) {
		kingsidePos := Position{kingPos.Row, kingPos.Col + 2}
		if !c.isSquareAttacked(kingPos, getOpponentColor(color)) &&
			!c.isSquareAttacked(kingsidePos, getOpponentColor(color)) &&
			c.Pieces[kingPos.Row][kingPos.Col+1] == nil &&
			c.Pieces[kingPos.Row][kingPos.Col+2] == nil {
			validMoves = append(validMoves, kingsidePos)
		}
	}

	// Check for queenside castling
	if c.canCastleQueenside(color) {
		queensidePos := Position{kingPos.Row, kingPos.Col - 2}
		if !c.isSquareAttacked(kingPos, getOpponentColor(color)) &&
			!c.isSquareAttacked(queensidePos, getOpponentColor(color)) &&
			c.Pieces[kingPos.Row][kingPos.Col-1] == nil &&
			c.Pieces[kingPos.Row][kingPos.Col-2] == nil &&
			c.Pieces[kingPos.Row][kingPos.Col-3] == nil {
			validMoves = append(validMoves, queensidePos)
		}
	}

	return validMoves
}

// canCastleKingside checks if kingside castling is allowed for the specified color
func (c *Chessboard) canCastleKingside(color Color) bool {
	kingPos := c.findKingPosition(color)
	if (color == White && kingPos.Row != 0) || (color == Black && kingPos.Row != 7) {
		return false // King is not in the starting position
	}
	return c.canCastleKingsideWhite() || c.canCastleKingsideBlack()
}

// canCastleQueenside checks if queenside castling is allowed for the specified color
func (c *Chessboard) canCastleQueenside(color Color) bool {
	kingPos := c.findKingPosition(color)
	if (color == White && kingPos.Row != 0) || (color == Black && kingPos.Row != 7) {
		return false // King is not in the starting position
	}
	return c.canCastleQueensideWhite() || c.canCastleQueensideBlack()
}

// canCastleKingsideWhite checks if kingside castling is allowed for white
func (c *Chessboard) canCastleKingsideWhite() bool {
	// Conditions for kingside castling for white
	if c.hasPieceMoved(Position{0, 4}) || c.hasPieceMoved(Position{0, 7}) {
		return false // King or kingside rook has moved
	}
	if c.Pieces[0][5] != nil || c.Pieces[0][6] != nil {
		return false // Squares between king and rook are not empty
	}
	if c.isSquareAttacked(Position{0, 4}, Black) || c.isSquareAttacked(Position{0, 5}, Black) || c.isSquareAttacked(Position{0, 6}, Black) {
		return false // King or castling squares are under attack
	}
	return true
}

// canCastleQueensideWhite checks if queenside castling is allowed for white
func (c *Chessboard) canCastleQueensideWhite() bool {
	// Conditions for queenside castling for white
	if c.hasPieceMoved(Position{0, 4}) || c.hasPieceMoved(Position{0, 0}) {
		return false // King or queenside rook has moved
	}
	if c.Pieces[0][1] != nil || c.Pieces[0][2] != nil || c.Pieces[0][3] != nil {
		return false // Squares between king and rook are not empty
	}
	if c.isSquareAttacked(Position{0, 4}, Black) || c.isSquareAttacked(Position{0, 3}, Black) || c.isSquareAttacked(Position{0, 2}, Black) {
		return false // King or castling squares are under attack
	}
	return true
}

// canCastleKingsideBlack checks if kingside castling is allowed for black
func (c *Chessboard) canCastleKingsideBlack() bool {
	// Conditions for kingside castling for black
	if c.hasPieceMoved(Position{7, 4}) || c.hasPieceMoved(Position{7, 7}) {
		return false // King or kingside rook has moved
	}
	if c.Pieces[7][5] != nil || c.Pieces[7][6] != nil {
		return false // Squares between king and rook are not empty
	}
	if c.isSquareAttacked(Position{7, 4}, White) || c.isSquareAttacked(Position{7, 5}, White) || c.isSquareAttacked(Position{7, 6}, White) {
		return false // King or castling squares are under attack
	}
	return true
}

// canCastleQueensideBlack checks if queenside castling is allowed for black
func (c *Chessboard) canCastleQueensideBlack() bool {
	// Conditions for queenside castling for black
	if c.hasPieceMoved(Position{7, 4}) || c.hasPieceMoved(Position{7, 0}) {
		return false // King or queenside rook has moved
	}
	if c.Pieces[7][1] != nil || c.Pieces[7][2] != nil || c.Pieces[7][3] != nil {
		return false // Squares between king and rook are not empty
	}
	if c.isSquareAttacked(Position{7, 4}, White) || c.isSquareAttacked(Position{7, 3}, White) || c.isSquareAttacked(Position{7, 2}, White) {
		return false // King or castling squares are under attack
	}
	return true
}

// hasPieceMoved checks if a piece at the specified position has moved
func (c *Chessboard) hasPieceMoved(pos Position) bool {
	return c.MovesCount[pos.Row][pos.Col] > 0
}

func (c *Chessboard) isSquareAttacked(square Position, byColor Color) bool {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := c.Pieces[i][j]
			if piece != nil && piece.Color == byColor {
				validMoves := c.calculateValidMoves(piece)
				for _, move := range validMoves {
					if move == square {
						return true
					}
				}
			}
		}
	}
	return false
}

func (c *Chessboard) findKingPosition(color Color) Position {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := c.Pieces[i][j]
			if piece != nil && piece.Type == King && piece.Color == color {
				return Position{i, j}
			}
		}
	}
	return Position{-1, -1}
}

func (c *Chessboard) IsInCheck(color Color) bool {
	kingPos := c.findKingPosition(color)
	opponentColor := getOpponentColor(color)

	// Check if any opponent piece can attack the king
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := c.Pieces[i][j]
			if piece != nil && piece.Color == opponentColor {
				validMoves := c.calculateValidMoves(piece)
				for _, move := range validMoves {
					if move == kingPos {
						return true
					}
				}
			}
		}
	}

	return false
}

func (c *Chessboard) cloneBoard() *Chessboard {
	clone := new(Chessboard)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := c.Pieces[i][j]
			if piece != nil {
				clone.Pieces[i][j] = &Piece{
					Type:  piece.Type,
					Color: piece.Color,
					Position: Position{
						Row: piece.Position.Row,
						Col: piece.Position.Col,
					},
				}
			}

			clone.MovesCount[i][j] = c.MovesCount[i][j]
		}
	}
	return clone
}

func (c *Chessboard) wouldMoveResultInCheck(color Color, from, to Position) bool {
	simulatedBoard := c.cloneBoard()

	fromPiece := simulatedBoard.GetPiece(from.Row, from.Col)
	fromPiece.Position = to

	// move the piece
	simulatedBoard.Pieces[to.Row][to.Col] = simulatedBoard.Pieces[from.Row][from.Col]
	simulatedBoard.Pieces[from.Row][from.Col] = nil
	simulatedBoard.increaseMovesCount(from, to)

	kingPos := simulatedBoard.findKingPosition(color)
	return simulatedBoard.isSquareAttacked(kingPos, getOpponentColor(color))
}

func (c *Chessboard) IsCheckmate(color Color) bool {
	if !c.IsInCheck(color) {
		return false
	}

	// Iterate through all pieces of the checked color and check if any valid moves can escape check
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			piece := c.Pieces[i][j]
			if piece != nil && piece.Color == color {
				for x := 0; x < 8; x++ {
					for y := 0; y < 8; y++ {
						if c.IsValidMove(piece, Position{x, y}) {
							// Try making the move and check if the king is still in check
							simulatedBoard := c.cloneBoard()
							simulatedBoard.PlacePiece(piece, Position{x, y})
							if !simulatedBoard.IsInCheck(color) {
								return false // King can escape check, not checkmate
							}
						}
					}
				}
			}
		}
	}

	return true // No valid moves to escape check, it's checkmate
}
