package services

import (
	"context"
	"github.com/TeamUUUU/keep4u-backend/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/updateopt"
	"go.uber.org/zap"
)

type DocumentAccessService struct {
	Db             *mongo.Client
	Logger         *zap.Logger
	Database       string
	CollectionName string
}

func (das *DocumentAccessService) Collection() *mongo.Collection {
	return das.Db.Database(das.Database).Collection(das.CollectionName)
}

func (das *DocumentAccessService) UpdateAccess(access *models.Access) error {
	documents := bson.NewArray()
	for _, doc := range access.Documents {
		documents.Append(bson.VC.String(doc))
	}
	_, err := das.Collection().UpdateOne(nil,
		bson.NewDocument(bson.EC.String("_id", access.UserID)),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$addToSet",
				bson.EC.SubDocumentFromElements("documents",
					bson.EC.Array("$each", documents),
				),
			),
			bson.EC.SubDocumentFromElements("$setOnInsert",
				bson.EC.String("_id", access.UserID),
			),
		),
		updateopt.Upsert(true),
	)
	if err != nil {
		das.Logger.Error("fail to update access", zap.Error(err))
	}
	return err
}

func (das *DocumentAccessService) DropAccess(access *models.Access) error {
	documents := bson.NewArray()
	for _, doc := range access.Documents {
		documents.Append(bson.VC.String(doc))
	}
	_, err := das.Collection().UpdateOne(nil,
		nil,
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$pull",
				bson.EC.SubDocumentFromElements("documents",
					bson.EC.Array("$in", documents),
				),
			),
			bson.EC.SubDocumentFromElements("$setOnInsert",
				bson.EC.String("_id", access.UserID),
			),
		),
	)
	if err != nil {
		das.Logger.Error("fail to drop access", zap.Error(err))
	}
	return err
}

func (das *DocumentAccessService) CheckAccess(access *models.Access) (bool, error) {
	ctx := context.Background()
	docs := bson.NewArray()
	for _, doc := range access.Documents {
		docs.Append(bson.VC.String(doc))
	}
	res, err := das.Collection().Find(ctx,
		bson.NewDocument(bson.EC.String("_id", access.UserID),
			bson.EC.SubDocumentFromElements("documents",
				bson.EC.Array("$all", docs),
			),
		),
	)
	if err == mongo.ErrNoDocuments {
		res.Close(ctx)
		return false, nil
	}
	if err != nil {
		das.Logger.Error("fail to check access", zap.Error(err))
		return false, err
	}
	res.Close(ctx)
	return true, nil
}
