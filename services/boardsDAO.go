package services

import (
	"context"
	"fmt"
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/option"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
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
	now := time.Now().Unix()
	board := models.Board{
		BoardCreate: *boardCreate,
		ID:          uuid.NewV4().String(),
		CreatedAt:   now,
		ChangedAt:   now,
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

func (bd *BoardsDao) Update(update *models.BoardUpdate) (*models.Board, error) {
	update.ChangedAt = time.Now().Unix()
	res := bd.Collection().FindOneAndUpdate(
		nil,
		bson.NewDocument(bson.EC.String("_id", update.ID)),
		&SetWrapper{Set: update},
		findopt.OptReturnDocument(option.After),
	)
	var board models.Board
	if err := res.Decode(&board); err != nil {
		bd.Logger.Error("fail to update board", zap.Error(err), zap.Any("new value", update))
		return nil, err
	}
	return &board, nil
}

func (bd *BoardsDao) Delete(boardID string) (error) {
	res, err := bd.Collection().DeleteOne(nil, bson.NewDocument(bson.EC.String("_id", boardID)))
	if err != nil {
		bd.Logger.Error("fail to delete", zap.Error(err), zap.String("board_id", boardID))
		return err
	}
	if res.DeletedCount == 0 {
		bd.Logger.Error("board not found", zap.String("board_id", boardID))
		return fmt.Errorf("board not found")
	}
	return nil
}

func (bd *BoardsDao) AddCollaborators(boardID string, update models.BoardCollaborationUpdate) (models.Collaborators, error) {
	entries := bson.NewArray()
	for _, collaborator := range update.Collaboration {
		entries.Append(bson.VC.String(collaborator))
	}
	each := bson.EC.SubDocumentFromElements("collaborations", bson.EC.Array("$each", entries))
	addToSet := bson.EC.SubDocumentFromElements("$addToSet", each)
	updateChangedAt := bson.EC.SubDocumentFromElements("$set", bson.EC.Int64("changed_at", time.Now().Unix()))
	res := bd.Collection().FindOneAndUpdate(
		nil,
		bson.NewDocument(bson.EC.String("_id", boardID)),
		bson.NewDocument(addToSet, updateChangedAt),
		findopt.OptReturnDocument(option.After),
	)
	var board models.Board
	if err := res.Decode(&board); err != nil {
		bd.Logger.Error("fail to decode a board", zap.Error(err), zap.String("board_id", boardID))
		return nil, err
	}
	return board.Collaboration, nil
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
