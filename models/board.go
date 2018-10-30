package models

type Board struct {
	ID            string         `json:"id,omitempty" bson:"_id,omitempty"`
	Collaboration Collaborations `json:"collaborations,omitempty" bson:"collaborations"`
	BoardCreate `bson:",inline"`
}

type BoardCreate struct {
	Title       string `json:"title" bson:"title"`
	OwnerID     string `json:"owner_id" binding:"required" bson:"owner_id"`
	Description string `json:"description" bson:"description"`
}

type Boards []*Board

type Collaborations []string
