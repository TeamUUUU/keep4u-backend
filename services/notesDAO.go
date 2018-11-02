package services

import (
	"context"
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

type NotesDAO struct {
	Db             *mongo.Client
	Logger         *zap.Logger
	Database       string
	CollectionName string
}

func (nd *NotesDAO) Collection() *mongo.Collection {
	return nd.Db.Database(nd.Database).Collection(nd.CollectionName)
}

func (nd *NotesDAO) Create(noteCreate *models.NoteCreate) (*models.Note, error) {
	now := time.Now().Unix()
	note := models.Note{
		NoteCreate: *noteCreate,
		ID:         uuid.NewV4().String(),
		CreatedAt:  now,
		ChangedAt:  now,
	}
	collection := nd.Collection()
	_, err := collection.InsertOne(context.Background(), &note)
	if err != nil {
		nd.Logger.Error("fail to insert note", zap.Error(err), zap.Any("note", *noteCreate))
		return nil, err
	}
	return &note, nil
}

func (nd *NotesDAO) GetNote(noteID string) (*models.Note, error) {
	res := nd.Collection().FindOne(nil, bson.NewDocument(bson.EC.String("_id", noteID)))
	var note models.Note
	if err := res.Decode(&note); err != nil {
		nd.Logger.Error("fail to find note by note id", zap.Error(err), zap.String("note_id", noteID))
		return nil, err
	}
	return &note, nil
}

func (nd *NotesDAO) GetNotesForBoard(boardid string) (models.Notes, error) {
	cur, err := nd.Collection().Find(nil, bson.NewDocument(bson.EC.String("board_id", boardid)))
	if err != nil {
		nd.Logger.Error("fail to find notes by board id", zap.Error(err), zap.String("note_id", boardid))
		return nil, err
	}
	notes := make(models.Notes, 0)
	ctx := context.Background()
	defer cur.Close(ctx)
	for cur.Next(nil) {
		var note models.Note
		if err := cur.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}
	if err := cur.Err(); err != nil {
		nd.Logger.Error("cursor errored", zap.Error(err))
		return nil, err
	}
	return notes, nil
}

func (nd *NotesDAO) Update(noteUpdate *models.NoteUpdate) (*models.Note, error) {
	noteUpdate.ChangedAt = time.Now().Unix()
	updateParam := SetWrapper{Set: noteUpdate}
	res := nd.Collection().FindOneAndUpdate(nil,
		bson.NewDocument(bson.EC.String("_id", noteUpdate.ID)),
		updateParam,
		findopt.OptReturnDocument(option.After),
	)

	var note models.Note
	if err := res.Decode(&note); err != nil {
		nd.Logger.Error("fail to perform update", zap.Error(err), zap.Any("new_value", noteUpdate))
		return nil, err
	}
	return &note, nil
}

func (nd *NotesDAO) Delete(noteID string) (error) {
	res, err := nd.Collection().DeleteOne(nil, bson.NewDocument(bson.EC.String("_id", noteID)))
	if err != nil {
		nd.Logger.Error("fail to delete", zap.Error(err), zap.String("note_id", noteID))
		return err
	}
	if res.DeletedCount == 0 {
		nd.Logger.Error("note not found", zap.String("note_id", noteID))
		return mongo.ErrNoDocuments
	}
	return nil
}
