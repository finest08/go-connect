package store

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/finest08/go-connect/internal/gen/messaging/v1"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = errors.New("document not found")
)

type Storer interface {
	CreateMsg(ctx context.Context, msg *pb.MessageThread) error
	QueryMsg(ctx context.Context, qr *pb.QueryRequest) ([]*pb.MessageThread, int64, error)
	GetMsg(ctx context.Context, id string) (*pb.MessageThread, error)
	UpdateMsg(id string, ctx context.Context, msg *pb.MessageThread) error
	DeleteMsg(id string) error
	CreateUser(ctx context.Context, usr *pb.User) error
	GetUser(ctx context.Context, id string) (*pb.User, error)
	DeleteUser(id string) error
}

func (s Store) CreateUser(ctx context.Context, usr *pb.User) error {
	_, err := s.locaCollUser.InsertOne(ctx, usr)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s Store) GetUser(ctx context.Context, id string) (*pb.User, error) {
	var usr pb.User

	if err := s.locaCollUser.FindOne(context.Background(), bson.M{"id": id}).Decode(&usr); err != nil {
		if err == mongo.ErrNoDocuments {
			return &usr, err
		}
		return &usr, err
	}

	return &usr, nil
}

func (s Store) DeleteUser(id string) error {
	if _, err := s.locaCollUser.DeleteOne(context.Background(), bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}


func (s Store) CreateMsg(ctx context.Context, msg *pb.MessageThread) error {
	_, err := s.locaCollMsg.InsertOne(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (s Store) QueryMsg(ctx context.Context, qr *pb.QueryRequest) ([]*pb.MessageThread, int64, error) {

	filter := bson.M{}
	if qr.SearchText != "" {
		filter = bson.M{"$text": bson.M{"$search": `"` + qr.SearchText + `"`}}
	}

	opt := options.FindOptions{
		Skip:  &qr.Offset,
		Limit: &qr.Limit,
		Sort:  bson.M{"date": -1},
	}

	cursor, err := s.locaCollMsg.Find(ctx, filter, &opt)
	if err != nil {
		return nil, 0, err
	}

	var msgs []*pb.MessageThread
	if err := cursor.All(ctx, &msgs); err != nil {
		return nil, 0, err
	}

	matches, err := s.locaCollMsg.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return msgs, matches, err
}

func (s Store) GetMsg(ctx context.Context, id string) (*pb.MessageThread, error) {
	var msg pb.MessageThread

	if err := s.locaCollMsg.FindOne(context.Background(), bson.M{"id": id}).Decode(&msg); err != nil {
		if err == mongo.ErrNoDocuments {
			return &msg, err
		}
		return &msg, err
	}

	return &msg, nil
}

func (s Store) UpdateMsg(id string, ctx context.Context, msg *pb.MessageThread) error {
	insertResult, err := s.locaCollMsg.ReplaceOne(ctx, bson.M{"id": id}, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nInserted a Single Document: %v\n", insertResult)

	return err
}

func (s Store) DeleteMsg(id string) error {
	if _, err := s.locaCollMsg.DeleteOne(context.Background(), bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}
