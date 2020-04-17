package model

// CommentView 注释视图
type CommentView struct {
	ID         int           `json:"id"`
	Content    string        `json:"content"`
	Creater    string        `json:"creater"`
	CreateDate string        `json:"createDate"`
	Reply      []interface{} `json:"reply"`
}
