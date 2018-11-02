package services

type SetWrapper struct {
	Set interface{} `bson:"$set,omitempty"`
}
