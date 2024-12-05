package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Images   []string           `json:"images" bson:"images"`
	Fullname string             `json:"fullname" bson:"fullname"`
}

func (u *User) GetByID(ctx context.Context, client *mongo.Client, id string) error {
	// Get the user by ID.
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("err in GetByID %s %s", id, err)
	}

	// Get the user collection.
	collection := client.Database(os.Getenv("MONGODBUSERSDBNAME")).Collection(os.Getenv("MONGODBUSERSCOLLECTION"))

	// Find the user by ID.
	filter := bson.D{{"_id", objID}}
	err = collection.FindOne(ctx, filter).Decode(u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("user with ID %s not found", objID.Hex())
		}
		return err
	}

	return nil
}
