package domain

type CreateUserRequest struct {
	Surname  string `json:"surname" validate:"required"`
	Lastname string `json:"lastname" validate:"required"`
	Age      string `json:"age" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`

	// New Fields
	FatherName string `json:"father_name,omitempty" validate:"omitempty"`
	MotherName string `json:"mother_name,omitempty" validate:"omitempty"`
}

type UpdateUserRequest struct {
	Surname  *string `json:"surname" validate:"omitempty"`
	Lastname *string `json:"lastname" validate:"omitempty"`
	Age      *string `json:"age" validate:"omitempty"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Phone    *string `json:"phone" validate:"omitempty"`

	// New Fields
	FatherName *string `json:"father_name,omitempty" validate:"omitempty"`
	MotherName *string `json:"mother_name,omitempty" validate:"omitempty"`
}
