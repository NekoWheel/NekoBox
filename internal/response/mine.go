package response

import "time"

type MineQuestionsItem struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Content    string    `json:"content"`
	IsAnswered bool      `json:"isAnswered"`
	IsPrivate  bool      `json:"isPrivate"`
}

type MineQuestions struct {
	Total     int64                `json:"total"`
	Cursor    string               `json:"cursor"`
	Questions []*MineQuestionsItem `json:"questions"`
}

type MineProfile struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
