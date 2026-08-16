package handler

import (
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/middleware"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gbkanban/gbkanban/internal/service"
	"github.com/gin-gonic/gin"
)

// BoardHandler exposes board and column endpoints.
type BoardHandler struct {
	boards service.BoardService
}

// NewBoardHandler constructs a BoardHandler.
func NewBoardHandler(boards service.BoardService) *BoardHandler {
	return &BoardHandler{boards: boards}
}

// List lists boards in a workspace.
func (h *BoardHandler) List(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boards, err := h.boards.ListByWorkspace(workspaceID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, boards)
}

// Create creates a board with default columns.
func (h *BoardHandler) Create(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	var req dto.CreateBoardRequest
	if !bindJSON(c, &req) {
		return
	}
	board, err := h.boards.Create(middleware.CurrentUserID(c), workspaceID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, board)
}

// Get returns board details with columns.
func (h *BoardHandler) Get(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	board, err := h.boards.Get(middleware.CurrentUserID(c), workspaceID, boardID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, board)
}

// Update updates board metadata.
func (h *BoardHandler) Update(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	var req dto.UpdateBoardRequest
	if !bindJSON(c, &req) {
		return
	}
	board, err := h.boards.Update(middleware.CurrentUserID(c), workspaceID, boardID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, board)
}

// Delete deletes a board.
func (h *BoardHandler) Delete(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	if err := h.boards.Delete(middleware.CurrentUserID(c), workspaceID, boardID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}

// ListColumns lists columns in a board.
func (h *BoardHandler) ListColumns(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	board, err := h.boards.Get(middleware.CurrentUserID(c), workspaceID, boardID)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, board.Columns)
}

// CreateColumn adds a custom column.
func (h *BoardHandler) CreateColumn(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	boardID, ok := parseUintParam(c, "board_id")
	if !ok {
		return
	}
	var req dto.CreateColumnRequest
	if !bindJSON(c, &req) {
		return
	}
	column, err := h.boards.CreateColumn(middleware.CurrentUserID(c), workspaceID, boardID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.Created(c, column)
}

// UpdateColumn renames or recolors a column.
func (h *BoardHandler) UpdateColumn(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	columnID, ok := parseUintParam(c, "column_id")
	if !ok {
		return
	}
	var req dto.UpdateColumnRequest
	if !bindJSON(c, &req) {
		return
	}
	column, err := h.boards.UpdateColumn(middleware.CurrentUserID(c), workspaceID, columnID, req)
	if err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, column)
}

// DeleteColumn deletes a column.
func (h *BoardHandler) DeleteColumn(c *gin.Context) {
	workspaceID, ok := parseUintParam(c, "workspace_id")
	if !ok {
		return
	}
	columnID, ok := parseUintParam(c, "column_id")
	if !ok {
		return
	}
	if err := h.boards.DeleteColumn(middleware.CurrentUserID(c), workspaceID, columnID); err != nil {
		respondError(c, err)
		return
	}
	response.OK(c, gin.H{"deleted": true})
}
