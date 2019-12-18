package gitlab

// User models a gitlab user entity
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	State    string `json:"state"`
	Avatar   string `json:"avatar_url"`
	URL      string `json:"web_url"`
}

