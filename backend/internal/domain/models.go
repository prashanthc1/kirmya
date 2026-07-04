package domain

type User struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  string   `json:"role"`
	Tags  []string `json:"tags"`
}

type Job struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type Workspace struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
	Type  string `json:"type"`
}

type ChatRoom struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Meeting struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	StartsAt string `json:"startsAt"`
}

type Notification struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Idea struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
