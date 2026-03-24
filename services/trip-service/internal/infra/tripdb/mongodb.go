package tripdb

import (
	"context"
	"fmt"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	tripsCollection     = "trips"
	rideFaresCollection = "ride_fares"
)

type mongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(db *mongo.Database) *mongoRepository {
	return &mongoRepository{db: db}
}

func (r *mongoRepository) CreateTrip(ctx context.Context, t *tripdomain.Trip) (*tripdomain.Trip, error) {
	result, err := r.db.Collection(tripsCollection).InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return nil, fmt.Errorf("unexpected inserted trip id type %T", result.InsertedID)
	}

	t.ID = insertedID
	return t, nil
}

func (r *mongoRepository) GetTripByID(ctx context.Context, id string) (*tripdomain.Trip, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var t tripdomain.Trip
	if err := r.db.Collection(tripsCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&t); err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *mongoRepository) SaveRideFare(ctx context.Context, fare *tripdomain.RideFare) error {
	result, err := r.db.Collection(rideFaresCollection).InsertOne(ctx, fare)
	if err != nil {
		return err
	}

	insertedID, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("unexpected inserted ride fare id type %T", result.InsertedID)
	}

	fare.ID = insertedID
	return nil
}

func (r *mongoRepository) GetRideFareByID(ctx context.Context, id string) (*tripdomain.RideFare, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var fare tripdomain.RideFare
	if err := r.db.Collection(rideFaresCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&fare); err != nil {
		return nil, err
	}

	return &fare, nil
}
