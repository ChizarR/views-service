package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChizarR/stats-service/internal/intaraction"
	"github.com/ChizarR/stats-service/pkg/error/mongoerr"
	"github.com/ChizarR/stats-service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ intaraction.Storage = &db{}

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) intaraction.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (db *db) GetOrCreate(ctx context.Context, date string) (intaraction.Intaraction, error) {
	result, err := db.FindOne(ctx, date)
	if err != nil {
		if errors.Is(err, mongoerr.ErrNotFound) {
			intr := intaraction.Intaraction{
				Id:    "",
				Date:  date,
				Views: map[string]int{},
			}
			_, err := db.Create(ctx, intr)
			if err != nil {
				return intaraction.Intaraction{}, err
			}
			result, err = db.FindOne(ctx, date)
			return result, nil
		}
		return intaraction.Intaraction{}, err
	}
	return result, nil
}

func (db *db) Create(ctx context.Context, intr intaraction.Intaraction) (string, error) {
	result, err := db.collection.InsertOne(ctx, intr)
	if err != nil {
		return "", fmt.Errorf("Failed to create new intaraction: %v", err)
	}
	objId, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return objId.Hex(), nil
	}
	return "", fmt.Errorf("Failed to convert ObjectID to Hex")
}

func (db *db) FindOne(ctx context.Context, date string) (intr intaraction.Intaraction, err error) {
	filter := bson.M{"date": date}
	result := db.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return intr, mongoerr.ErrNotFound
		}
		return intr, fmt.Errorf("Failed to find Intaractions for today: %s, %v", date, err)
	}

	if err := result.Decode(&intr); err != nil {
		return intr, fmt.Errorf("Failed to decode Intaraction due to: %v", err)
	}
	return intr, nil
}

func (db *db) FindAll(ctx context.Context) ([]intaraction.Intaraction, error) {
	var intaractions []intaraction.Intaraction
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return []intaraction.Intaraction{}, fmt.Errorf("Failed to find all docs: %v", err)
	}

	if err = cursor.All(ctx, &intaractions); err != nil {
		return intaractions, fmt.Errorf("Failed to read all documents from cursor: %v", err)
	}

	return intaractions, nil
}

func (db *db) Update(ctx context.Context, intr intaraction.Intaraction) error {
	objId, err := primitive.ObjectIDFromHex(intr.Id)
	if err != nil {
		return fmt.Errorf("Failed to convert Id to ObjectID, Id=%s", intr.Id)
	}

	filter := bson.M{"_id": objId}
	intrBytes, err := bson.Marshal(intr)
	if err != nil {
		return fmt.Errorf("Failed to Marshal Intaraction documet due to: %v", err)
	}

	var updateIntrObject bson.M
	err = bson.Unmarshal(intrBytes, &updateIntrObject)
	if err != nil {
		return fmt.Errorf("Failed to Unmarshal intrBytes due to: %v", err)
	}
	delete(updateIntrObject, "_id")

	update := bson.M{"$set": updateIntrObject}

	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Failed to update Intaraction due to: %v", err)
	}

	if result.MatchedCount == 0 {
		return mongoerr.ErrNotFound
	}

	db.logger.Tracef("Matched %d documents and modified %d documents\n", result.MatchedCount, result.ModifiedCount)
	return nil
}
