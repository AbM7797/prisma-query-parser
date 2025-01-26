package domain

type User struct {
	ID      string    `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Post    []Post    `json:"post"`
	Comment []Comment `json:"comment"`
}
