package domain

type CreateUsersResponse struct {
	Data *UserResponse `json:"data,omitempty"`
}

type GetUserResponse struct {
	Data *UserResponse `json:"data,omitempty"`
}

type GetUsersResponse struct {
	Data []UserResponse `json:"data"`
}

type UpdateUserResponse struct {
	Data *UserResponse `json:"data,omitempty"`
}

type DeleteUserResponse struct {
	Data *UserResponse `json:"data,omitempty"`
}

// SymbolResponse represents the response version of an Symbol.
type UserResponse struct {
	ID        string `json:"id"`
	SurName   string `json:"surname"`
	LastName  string `json:"lastname"`
	Age       string `json:"age"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt *int64 `json:"updated_at,omitempty"`

	// New Fields
	FatherName string `json:"father_name,omitempty" validate:"omitempty"`
	MotherName string `json:"mother_name,omitempty" validate:"omitempty"`
}
