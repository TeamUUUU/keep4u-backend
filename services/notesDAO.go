package services

import (
	"context"
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
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
	note := models.Note{
		NoteCreate: *noteCreate,
		ID:         uuid.NewV4().String(),
		CreatedAt:  time.Now().Unix(),
	}

	collection := nd.Collection()
	_, err := collection.InsertOne(context.Background(), &note)
	if err != nil {
		nd.Logger.Error("fail to insert note", zap.Error(err), zap.Any("note", *noteCreate))
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
