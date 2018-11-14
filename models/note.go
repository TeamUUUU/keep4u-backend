package models

type Note struct {
	ID         string `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ChangedAt  int64  `json:"changed_at,omitempty" bson:"changed_at,omitempty"`
	NoteCreate `bson:",inline"`
}

type NoteCreate struct {
	Title       string      `json:"title" bson:"title"`
	Content     string      `json:"content" bson:"content"`
	Attachments Attachments `json:"attachments,omitempty" bson:"attachments,omitempty"`
	BoardID     string      `json:"board_id, omitempty" bson:"board_id,omitempty"`
	OwnerID     string      `json:"-" bson:"owner_id,omitempty"`
}

type NoteUpdate struct {
	ID        string `json:"id,omitempty" bson:"-"`
	Title     string `json:"title,omitempty" bson:"title,omitempty"`
	Content   string `json:"content,omitempty" bson:"content,omitempty"`
	ChangedAt int64  `json:"changed_at,omitempty" bson:"changed_at,omitempty"`
}

type Notes []*Note
