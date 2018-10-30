package models

type Attachment struct {
	ID   string `json:"id" bson:"_id"`
	Kind string `json:"kind" bson:"kind"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type Attachments []Attachment
