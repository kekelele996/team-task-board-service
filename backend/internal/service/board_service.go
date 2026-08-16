package service

import (
	"fmt"
	"log/slog"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/dto"
	"github.com/gbkanban/gbkanban/internal/model"
	"github.com/gbkanban/gbkanban/internal/repository"
)

// BoardService handles board and column business logic.
type BoardService interface {
	Create(userID, workspaceID uint, req dto.CreateBoardRequest) (*model.Board, error)
	ListByWorkspace(workspaceID uint) ([]model.Board, error)
	Get(userID, workspaceID, boardID uint) (*model.Board, error)
	Update(userID, workspaceID, boardID uint, req dto.UpdateBoardRequest) (*model.Board, error)
	Delete(userID, workspaceID, boardID uint) error

	CreateColumn(userID, workspaceID, boardID uint, req dto.CreateColumnRequest) (*model.BoardColumn, error)
	UpdateColumn(userID, workspaceID, columnID uint, req dto.UpdateColumnRequest) (*model.BoardColumn, error)
	DeleteColumn(userID, workspaceID, columnID uint) error
	ListColumns(boardID uint) ([]model.BoardColumn, error)
}

type boardService struct {
	boards     repository.BoardRepo
	columns    repository.ColumnRepo
	workspaces repository.WorkspaceRepo
	workspace  WorkspaceService
	logger     *slog.Logger
}

// NewBoardService constructs a BoardService.
func NewBoardService(boards repository.BoardRepo, columns repository.ColumnRepo, workspaces repository.WorkspaceRepo, workspace WorkspaceService, logger *slog.Logger) BoardService {
	return &boardService{boards: boards, columns: columns, workspaces: workspaces, workspace: workspace, logger: logger}
}

func (s *boardService) Create(userID, workspaceID uint, req dto.CreateBoardRequest) (*model.Board, error) {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor); err != nil {
		return nil, err
	}
	board := &model.Board{WorkspaceID: workspaceID, Name: req.Name, Description: req.Description, CreatedBy: userID}
	if err := s.boards.Create(board); err != nil {
		return nil, fmt.Errorf("create board: %w", err)
	}
	for _, name := range constants.DefaultColumns {
		column := &model.BoardColumn{BoardID: board.ID, Name: name, Position: 0}
		if err := s.columns.Create(column); err != nil {
			return nil, fmt.Errorf("create default column %q: %w", name, err)
		}
	}
	board.Columns, _ = s.columns.ListByBoard(board.ID)
	return board, nil
}

func (s *boardService) ListByWorkspace(workspaceID uint) ([]model.Board, error) {
	return s.boards.ListByWorkspace(workspaceID)
}

func (s *boardService) Get(userID, workspaceID, boardID uint) (*model.Board, error) {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor, constants.RoleViewer); err != nil {
		return nil, err
	}
	board, err := s.boards.FindByID(boardID)
	if err != nil {
		return nil, err
	}
	if board.WorkspaceID != workspaceID {
		return nil, repository.ErrNotFound
	}
	return board, nil
}

func (s *boardService) Update(userID, workspaceID, boardID uint, req dto.UpdateBoardRequest) (*model.Board, error) {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor); err != nil {
		return nil, err
	}
	board, err := s.boards.FindByID(boardID)
	if err != nil {
		return nil, err
	}
	if board.WorkspaceID != workspaceID {
		return nil, repository.ErrNotFound
	}
	board.Name = req.Name
	board.Description = req.Description
	if err := s.boards.Update(board); err != nil {
		return nil, fmt.Errorf("update board: %w", err)
	}
	return board, nil
}

func (s *boardService) Delete(userID, workspaceID, boardID uint) error {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin); err != nil {
		return err
	}
	board, err := s.boards.FindByID(boardID)
	if err != nil {
		return err
	}
	if board.WorkspaceID != workspaceID {
		return repository.ErrNotFound
	}
	if err := s.boards.Delete(boardID); err != nil {
		return fmt.Errorf("delete board: %w", err)
	}
	return nil
}

func (s *boardService) ListColumns(boardID uint) ([]model.BoardColumn, error) {
	return s.columns.ListByBoard(boardID)
}

func (s *boardService) CreateColumn(userID, workspaceID, boardID uint, req dto.CreateColumnRequest) (*model.BoardColumn, error) {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor); err != nil {
		return nil, err
	}
	if err := s.ensureBoardInWorkspace(boardID, workspaceID); err != nil {
		return nil, err
	}
	max, err := s.columns.MaxPosition(boardID)
	if err != nil {
		return nil, fmt.Errorf("create column: %w", err)
	}
	column := &model.BoardColumn{BoardID: boardID, Name: req.Name, Position: max + 1, Color: req.Color}
	if err := s.columns.Create(column); err != nil {
		return nil, fmt.Errorf("create column: %w", err)
	}
	return column, nil
}

func (s *boardService) UpdateColumn(userID, workspaceID, columnID uint, req dto.UpdateColumnRequest) (*model.BoardColumn, error) {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor); err != nil {
		return nil, err
	}
	column, err := s.columns.FindByID(columnID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureBoardInWorkspace(column.BoardID, workspaceID); err != nil {
		return nil, err
	}
	column.Name = req.Name
	column.Color = req.Color
	if err := s.columns.Update(column); err != nil {
		return nil, fmt.Errorf("update column: %w", err)
	}
	return column, nil
}

func (s *boardService) DeleteColumn(userID, workspaceID, columnID uint) error {
	if err := s.workspace.RequireRole(workspaceID, userID, constants.RoleAdmin, constants.RoleEditor); err != nil {
		return err
	}
	column, err := s.columns.FindByID(columnID)
	if err != nil {
		return err
	}
	if err := s.ensureBoardInWorkspace(column.BoardID, workspaceID); err != nil {
		return err
	}
	if err := s.columns.Delete(columnID); err != nil {
		return fmt.Errorf("delete column: %w", err)
	}
	return nil
}

func (s *boardService) ensureBoardInWorkspace(boardID, workspaceID uint) error {
	board, err := s.boards.FindByID(boardID)
	if err != nil {
		return err
	}
	if board.WorkspaceID != workspaceID {
		return repository.ErrNotFound
	}
	return nil
}
