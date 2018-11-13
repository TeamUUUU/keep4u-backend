package models

type Access struct {
	UserID    string   `bson:"_id"`
	Documents []string `bson:"documents"`
}
