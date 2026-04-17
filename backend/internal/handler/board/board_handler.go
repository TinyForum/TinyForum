package board

import (
	boardService "tiny-forum/internal/service/board"
)

type BoardHandler struct {
	boardSvc *boardService.BoardService
}

func NewBoardHandler(boardSvc *boardService.BoardService) *BoardHandler {
	return &BoardHandler{boardSvc: boardSvc}
}
