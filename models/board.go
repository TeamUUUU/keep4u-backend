package models

type Board struct {
	ID            string        `json:"id,omitempty" bson:"_id,omitempty"`
	Collaborators Collaborators `json:"collaborators,omitempty" bson:"collaborators,omitempty"`
	CreatedAt     int64         `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ChangedAt     int64         `json:"changed_at,omitempty" bson:"changed_at,omitempty"`
	BoardCreate   `bson:",inline"`
}

type BoardCreate struct {
	Title       string `json:"title" bson:"title"`
	OwnerID     string `json:"owner_id,omitempty" bson:"owner_id,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type BoardUpdate struct {
	ID                       string `json:"id,omitempty" bson:"_id, omitempty"`
	Title                    string `json:"title,omitempty" bson:"title,omitempty"`
	Description              string `json:"description,omitempty" bson:"description,omitempty"`
	ChangedAt                int64  `json:"changed_at,omitempty" bson:"changed_at,omitempty"`
	BoardCollaborationUpdate `bson:",inline"`
}

type BoardCollaborationUpdate struct {
	Collaboration Collaborators `json:"collaborators,omitempty" bson:"collaborators,omitempty"`
}

type Boards []*Board

type Collaborators []string
