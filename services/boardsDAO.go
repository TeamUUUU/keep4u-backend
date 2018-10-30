package services

import (
	"context"
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type BoardsDao struct {
	Db             *mongo.Client
	Logger         *zap.Logger
	Database       string
	CollectionName string
}

func (bd *BoardsDao) Collection() *mongo.Collection {
	return bd.Db.Database(bd.Database).Collection(bd.CollectionName)
}

func (bd *BoardsDao) Create(boardCreate *models.BoardCreate) (*models.Board, error) {
	board := models.Board{
		BoardCreate: *boardCreate,
		ID:          uuid.NewV4().String(),
	}
	collection := bd.Collection()
	_, err := collection.InsertOne(context.Background(), &board)
	if err != nil {
		bd.Logger.Error("fail to insert note", zap.Error(err), zap.Any("note", *boardCreate))
		return nil, err
	}
	return &board, nil
}

func (bd *BoardsDao) GetBoardById(id string) (*models.Board, error) {
	res := bd.Collection().FindOne(nil, bson.NewDocument(bson.EC.String("_id", id)))
	var board models.Board
	if err := res.Decode(&board); err != nil {
		bd.Logger.Error("fail to get board by id", zap.Error(err), zap.String("id", id))
		return nil, err
	}
	return &board, nil
}

func (bd *BoardsDao) GetBoardsForUser(ownerid string) (models.Boards, error) {

	cur, err := bd.Collection().Find(nil, bson.NewDocument(bson.EC.String("owner_id", ownerid)))
	if err != nil {
		bd.Logger.Error("fail to find boards by owner id", zap.Error(err), zap.String("owner_id", ownerid))
		return nil, err
	}
	ctx := context.Background()
	defer cur.Close(ctx)

	boards := make(models.Boards, 0)

	for cur.Next(nil) {
		var board models.Board
		if err := cur.Decode(&board); err != nil {
			return nil, err
		}
		boards = append(boards, &board)
	}
	if err := cur.Err(); err != nil {
		bd.Logger.Error("cursor error", zap.Error(err))
		return nil, err
	}
	return boards, nil
}
