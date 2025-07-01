package LinkDriver

import (
	"backend/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	databaseName   = "links"
	collectionName = "links"
)

type MongoLinkDriver struct {
	collection *mongo.Collection
	nowFunc    func() time.Time
}

func NewMongoLinkDriver(URI string) (*MongoLinkDriver, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	return &MongoLinkDriver{
		collection: collection,
		nowFunc: func() time.Time {
			return time.Now()
		},
	}, nil
}

func (m *MongoLinkDriver) GetLinkByID(ID string) (*models.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var link models.Link
	err := m.collection.FindOne(ctx, map[string]string{"_id": ID}).Decode(&link)
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (m *MongoLinkDriver) GetLinksForMember(Email string) ([]*models.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, map[string]string{"member_email": Email})
	if err != nil {
		return nil, err
	}

	var links []*models.Link
	for cursor.Next(ctx) {
		var link models.Link
		err := cursor.Decode(&link)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}

	return links, nil
}

func (m *MongoLinkDriver) GetQRsForMember(Email string) ([]*models.Link, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := m.collection.Find(ctx, bson.M{"member_email": Email, "has_qr": true})
	if err != nil {
		return nil, err
	}

	var links []*models.Link
	for cursor.Next(ctx) {
		var link models.Link
		err := cursor.Decode(&link)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}

	return links, nil
}

func (m *MongoLinkDriver) UpsertLink(Link *models.Link) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	_, err := m.collection.ReplaceOne(ctx, map[string]string{"_id": Link.ID}, Link, opts)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoLinkDriver) DeleteLink(ID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.collection.DeleteOne(ctx, map[string]string{"_id": ID})
	if err != nil {
		return err
	}

	return nil
}
