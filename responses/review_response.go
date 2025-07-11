package responses

type UserMiniResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CommentResponse struct {
	ID      uint             `json:"id"`
	User    UserMiniResponse `json:"user"`
	Content string           `json:"content"`
}

type LikeResponse struct {
	ID   uint             `json:"id"`
	User UserMiniResponse `json:"user"`
}

type ReviewResponse struct {
	ID       uint              `json:"id"`
	Rating   int               `json:"rating"`
	Content  string            `json:"content"`
	User     UserMiniResponse  `json:"user"`
	Comments []CommentResponse `json:"comments"`
	Likes    []LikeResponse    `json:"likes"`
}
