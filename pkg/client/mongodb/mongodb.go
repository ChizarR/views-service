package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, user, password, database, authDB string) (*mongo.Database, error) {
	var mongoURL string
	var isAuth bool
	if user == "" && password == "" {
		isAuth = false
		mongoURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)
	}

	clientOptions := options.Client().ApplyURI(mongoURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDB,
			Username:   user,
			Password:   password,
		})

	}

	// Connect
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("CONNECTION to MONGO failed: %v", err)
	}

	// Ping
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("CONNECTION to MONGO failed: %v", err)
	}

	return client.Database(database), nil
}
