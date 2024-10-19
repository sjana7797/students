package types

type Student struct {
	Id    int    `json:"id"`
	Name  string `validate:"required" json:"name"`
	Email string `validate:"required,email" json:"email"`
	Age   int    `validate:"required,number" json:"age"`
}

type UpdateStudent struct {
	Name  *string `json:"name,omitempty"`
	Email *string `validate:"email" json:"email,omitempty"`
	Age   *int    `validate:"number" json:"age,omitempty"`
}
