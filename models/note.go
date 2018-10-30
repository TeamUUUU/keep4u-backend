package models

type Note struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt  int64  `json:"created_at" bson:"string"`
	NoteCreate `bson:",inline"`
}

type NoteCreate struct {
	Title       string      `json:"title" title:"title"`
	Content     string      `json:"content" bson:"content"`
	Attachments Attachments `json:"attachments,omitempty" bson:"attachments,omitempty"`
	BoardID     string      `json:"board_id, omitempty" bson:"board_id,omitempty"`
}

type Notes []*Note
