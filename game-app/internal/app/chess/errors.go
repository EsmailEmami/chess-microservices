package chess

import (
	"errors"
)

var (
	ErrInvalidGame              = errors.New("this is not your game")
	ErrInvalidTurn              = errors.New("it is not tour turn to move")
	ErrGameNoPlayers            = errors.New("the game has no players")
	ErrGameNotFound             = errors.New("game not found")
	ErrGameWaitingStatus        = errors.New("game is in waiting status")
	ErrGameIsNotInWaitingStatus = errors.New("game is not in waiting status, you can not play")
	ErrGameIsOver               = errors.New("game is over")
)
