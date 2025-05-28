package domain

type User struct {
	ID        string `json:"id"`
	Surname   string `json:"surname"`
	Lastname  string `json:"lastname"`
	Age       string `json:"age"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt *int64 `json:"updated_at,omitempty"`

	// New Field
	FatherName string `json:"father_name,omitempty"`
	MotherName string `json:"mother_name,omitempty"`
}
