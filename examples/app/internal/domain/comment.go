package domain

type Comment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	UserID  string `json:"userId"`
	PostID  string `json:"postId"`
	User    *User  `json:"user,omitempty"`
	Post    *Post  `json:"post,omitempty"`
}
