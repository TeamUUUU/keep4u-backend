package models

type Note struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt  int64  `json:"created_at" bson:"created_at"`
	ChangedAt  int64  `json:"changed_at" bson:"changed_at"`
	NoteCreate `bson:",inline"`
}

type NoteCreate struct {
	Title       string      `json:"title" bson:"title" binding:"required"`
	Content     string      `json:"content" bson:"content" binding:"required"`
	Attachments Attachments `json:"attachments,omitempty" bson:"attachments,omitempty"`
	BoardID     string      `json:"board_id, omitempty" bson:"board_id,omitempty"`
}

type NoteUpdate struct {
	ID        string `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title,omitempty"`
	Content   string `json:"content" bson:"content,omitempty"`
	ChangedAt int64  `json:"changed_at" bson:"changed_at"`
}

type Notes []*Note
