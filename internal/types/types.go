package types

type Student struct {
	Id    int    `json:"id"`
	Name  string `validate:"required" json:"name"`
	Email string `validate:"required,email" json:"email"`
	Age   int    `validate:"required,number" json:"age"`
}
