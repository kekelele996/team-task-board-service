package dto

// CreateBoardRequest creates a board and its default columns.
type CreateBoardRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=128"`
	Description string `json:"description" validate:"max=2000"`
}

// UpdateBoardRequest updates board metadata.
type UpdateBoardRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=128"`
	Description string `json:"description" validate:"max=2000"`
}

// CreateColumnRequest adds a custom column.
type CreateColumnRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=64"`
	Color string `json:"color" validate:"max=32"`
}

// UpdateColumnRequest renames or recolors a column.
type UpdateColumnRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=64"`
	Color string `json:"color" validate:"max=32"`
}
