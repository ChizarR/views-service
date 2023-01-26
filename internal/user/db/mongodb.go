package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ChizarR/stats-service/internal/user"
	"github.com/ChizarR/stats-service/pkg/error/mongoerr"
	"github.com/ChizarR/stats-service/pkg/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ user.Storage = &db{}

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (db *db) GetOrCreate(ctx context.Context, tgId int) (user.User, error) {
	result, err := db.FindOne(ctx, tgId)
	if err != nil {
		if errors.Is(err, mongoerr.ErrNotFound) {
			u := user.User{
				Id:           "",
				TgId:         tgId,
				Intaractions: []user.Intaractions{},
			}
			_, err := db.Create(ctx, u)
			if err != nil {
				return user.User{}, err
			}
			result, err = db.FindOne(ctx, tgId)
			return result, nil
		}
		return user.User{}, err
	}
	return result, nil
}

func (db *db) Create(ctx context.Context, user user.User) (string, error) {
	result, err := db.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("Failed to create new user: %v", err)
	}
	objId, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return objId.Hex(), nil
	}
	return "", fmt.Errorf("Failed to convert ObjectID to Hex")
}

func (db *db) FindOne(ctx context.Context, tgId int) (user.User, error) {
	filter := bson.M{"tg_id": tgId}
	result := db.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return user.User{}, mongoerr.ErrNotFound
		}
		return user.User{}, fmt.Errorf("Failed to find User with tg_id: %s", tgId)
	}

	var u user.User
	if err := result.Decode(&u); err != nil {
		return u, fmt.Errorf("Failed to decode User due to: %v", err)
	}
	return u, nil
}

func (db *db) FindAll(ctx context.Context) ([]user.User, error) {
	var users []user.User
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return users, fmt.Errorf("Failed to find all docs: %v", err)
	}

	if err = cursor.All(ctx, &users); err != nil {
		return users, fmt.Errorf("Failed to read all documents from cursor: %v", err)
	}

	return users, nil
}

func (db *db) Update(ctx context.Context, user user.User) error {
	objId, err := primitive.ObjectIDFromHex(user.Id)
	if err != nil {
		return fmt.Errorf("Failed to convert Id to ObjectID, Id=%s", user.Id)
	}

	filter := bson.M{"_id": objId}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("Failed to Marshal User documet due to: %v", err)
	}

	var updateUserObject bson.M
	err = bson.Unmarshal(userBytes, &updateUserObject)
	if err != nil {
		return fmt.Errorf("Failed to Unmarshal userBytes due to: %v", err)
	}
	delete(updateUserObject, "_id")

	update := bson.M{"$set": updateUserObject}

	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Failed to update User due to: %v", err)
	}

	if result.MatchedCount == 0 {
		return mongoerr.ErrNotFound
	}

	db.logger.Tracef("Matched %d documents and modified %d documents\n", result.MatchedCount, result.ModifiedCount)
	return nil
}
