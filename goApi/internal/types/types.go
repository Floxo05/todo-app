package types

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
	ID        *int    `json:"id"`
}

type ShareToUserRequest struct {
	Username string `json:"username"`
	TodoID   int    `json:"id"`
}
