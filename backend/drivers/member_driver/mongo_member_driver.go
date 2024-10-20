package member_driver

import (
	"backend/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	databaseName   = "members"
	collectionName = "members"
)

type MongoMemberDriver struct {
	collection *mongo.Collection
	nowFunc    func() time.Time
}

func NewMongoMemberDriver(URI string) (*MongoMemberDriver, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	return &MongoMemberDriver{
		collection: collection,
		nowFunc: func() time.Time {
			return time.Now()
		},
	}, nil
}

func (m *MongoMemberDriver) CountMembersWithEmail(Email string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := m.collection.CountDocuments(ctx, map[string]string{"email": Email})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (m *MongoMemberDriver) GetMemberByEmail(Email string) (*models.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var member models.Member
	err := m.collection.FindOne(ctx, map[string]string{"email": Email}).Decode(&member)
	if err != nil {
		return nil, err
	}

	return &member, nil
}

func (m *MongoMemberDriver) UpsertMember(Member *models.Member) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Replace().SetUpsert(true)
	_, err := m.collection.ReplaceOne(ctx, map[string]string{"email": Member.Email}, Member, opts)
	return err
}

func (m *MongoMemberDriver) DeleteMember(Email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := m.collection.DeleteOne(ctx, map[string]string{"email": Email})
	return err
}
