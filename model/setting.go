package model

// Setting setting
type Setting struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	Copyright string `json:"copyright"`
	Keyword   string `json:"keyword"`
	EMail     string `json:"email"`
	ICP       string `json:"icp"`
}
