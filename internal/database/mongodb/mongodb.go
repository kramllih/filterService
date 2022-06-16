package mongodb

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kramllih/filterService/config"
	"github.com/kramllih/filterService/internal/database"
	"github.com/kramllih/filterService/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "filter"
)

type mongoDb struct {
	DB  *mongo.Client
	log *logger.Logger

	approvalCol *mongo.Collection
	rejectedCol *mongo.Collection
	messageCol  *mongo.Collection
}

func init() {
	database.RegisterType("mongodb", NewDB)
}

func NewDB(cfg *config.ConfigNamespace) (database.Client, error) {
	db := &mongoDb{
		log: logger.NewLogger("mongodb"),
	}

	config := mongoConfig{}

	if err := cfg.Config().UnpackRaw(&config); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Host))
	if err != nil {
		db.log.Error(err)
		return nil, err
	}

	db.DB = client
	db.approvalCol = client.Database(dbName).Collection("approvals")
	db.rejectedCol = client.Database(dbName).Collection("rejected")
	db.messageCol = client.Database(dbName).Collection("messages")

	return db, nil
}

func (c *mongoDb) StoreApproval(id string, approval []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data := bson.M{"_id": id, "message": approval}

	_, err := c.approvalCol.InsertOne(ctx, data)
	if err != nil {
		c.log.Errorf("Unable to insert approval to database: %s", err)
		return err
	}

	return nil

}
func (c *mongoDb) GetApproval(id string) (*database.Approval, error) {
	var approval *database.Approval

	type temp struct {
		Id      string `bson:"_id"`
		Message []byte `bson:"message"`
	}

	result := temp{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	res := c.approvalCol.FindOne(ctx, filter)

	res.Decode(&result)

	if result.Message != nil {
		json.Unmarshal(result.Message, &approval)
	}

	if approval != nil {
		return approval, nil
	}

	return nil, nil

}
func (c *mongoDb) GetAllApprovals() ([]*database.Approval, error) {

	approvals := []*database.Approval{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := c.approvalCol.Find(ctx, bson.M{})
	if err != nil {
		c.log.Errorf("Unable to find devices: %s", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		d := database.Approval{}

		type temp struct {
			Id      string `bson:"_id"`
			Message []byte `bson:"message"`
		}

		result := temp{}

		err := cur.Decode(&result)
		if err != nil {
			c.log.Errorf("bson decode error: %w", err)
			break
		}

		if result.Message != nil {
			json.Unmarshal(result.Message, &d)

			approvals = append(approvals, &d)
		}

	}
	if err := cur.Err(); err != nil {
		c.log.Errorf("Unable to devices: %+v", err)
		return nil, err
	}

	return approvals, nil
}
func (c *mongoDb) UpdateApprovals(id string, approval []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	data := bson.M{"_id": id, "message": approval}

	res := c.approvalCol.FindOneAndReplace(ctx, filter, data)
	if err := res.Err(); err != nil {
		c.log.Errorf("Unable to replace device: %+v", err)
		return err

	}

	return nil
}

func (c *mongoDb) DeleteApprovals(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	res := c.approvalCol.FindOneAndDelete(ctx, filter)
	if err := res.Err(); err != nil {
		c.log.Errorf("Unable to replace device: %+v", err)
		return err

	}

	return nil
}

func (c *mongoDb) StoreReject(id string, message []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data := bson.M{"_id": id, "message": message}

	_, err := c.rejectedCol.InsertOne(ctx, data)
	if err != nil {
		c.log.Errorf("Unable to insert rejected message to database: %s", err)
		return err
	}

	return nil

}

func (c *mongoDb) GetAllRejected() ([]*database.Message, error) {

	rejected := []*database.Message{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := c.rejectedCol.Find(ctx, bson.M{})
	if err != nil {
		c.log.Errorf("Unable to find devices: %s", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		d := database.Message{}

		type temp struct {
			Id      string `bson:"_id"`
			Message []byte `bson:"message"`
		}

		result := temp{}

		err := cur.Decode(&result)
		if err != nil {
			c.log.Errorf("bson decode error: %w", err)
			break
		}

		if result.Message != nil {
			json.Unmarshal(result.Message, &d)

			rejected = append(rejected, &d)
		}

	}
	if err := cur.Err(); err != nil {
		c.log.Errorf("Unable to devices: %+v", err)
		return nil, err
	}

	return rejected, nil
}

func (c *mongoDb) StoreMessage(id string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data := bson.M{"_id": id, "message": message}

	_, err := c.messageCol.InsertOne(ctx, data)
	if err != nil {
		c.log.Errorf("Unable to insert rejected message to database: %s", err)
		return err
	}

	return nil

}

func (c *mongoDb) GetMessage(id string) (*database.Message, error) {

	var message *database.Message

	type temp struct {
		Id      string `bson:"_id"`
		Message []byte `bson:"message"`
	}

	result := temp{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	res := c.messageCol.FindOne(ctx, filter)

	res.Decode(&result)

	if result.Message != nil {
		json.Unmarshal(result.Message, &message)
	}

	if message != nil {
		return message, nil
	}

	return nil, nil

}

func (c *mongoDb) UpdateMessage(id string, message []byte) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	data := bson.M{"_id": id, "message": message}

	res := c.messageCol.FindOneAndReplace(ctx, filter, data)
	if err := res.Err(); err != nil {
		c.log.Errorf("Unable to replace device: %+v", err)
		return err

	}

	return nil
}

func (c *mongoDb) GetAllMessages() ([]*database.Message, error) {

	messages := []*database.Message{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := c.messageCol.Find(ctx, bson.M{})
	if err != nil {
		c.log.Errorf("Unable to find devices: %s", err)
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {

		d := database.Message{}

		type temp struct {
			Id      string `bson:"_id"`
			Message []byte `bson:"message"`
		}

		result := temp{}

		err := cur.Decode(&result)
		if err != nil {
			c.log.Errorf("bson decode error: %w", err)
			break
		}

		if result.Message != nil {
			json.Unmarshal(result.Message, &d)

			messages = append(messages, &d)
		}

	}
	if err := cur.Err(); err != nil {
		c.log.Errorf("Unable to devices: %+v", err)
		return nil, err
	}

	return messages, nil
}
