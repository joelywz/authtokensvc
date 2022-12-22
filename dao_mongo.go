package authtokensvc

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ Dao = (*MongoDao)(nil)

type MongoDao struct {
	collection *mongo.Collection
}

func NewMongoDao(collection *mongo.Collection) *MongoDao {
	return &MongoDao{
		collection: collection,
	}
}

func (dao *MongoDao) DeleteAll(userId string) error {

	filter := bson.D{{Key: "user_id", Value: userId}}

	_, err := dao.collection.DeleteMany(context.Background(), filter)

	return err
}

func (dao *MongoDao) Delete(tokenId string) error {

	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "access_token_id", Value: tokenId}},
		bson.D{{Key: "refresh_token_id", Value: tokenId}},
	}}}

	_, err := dao.collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil
}

func (dao *MongoDao) SaveToken(token *Token) error {

	_, err := dao.collection.InsertOne(context.Background(), token)

	if err != nil {
		return err
	}

	return nil
}

func (dao *MongoDao) GetRefreshToken(refreshTokenId string) (*Token, error) {
	filter := bson.M{"refresh_token_id": refreshTokenId}

	result := dao.collection.FindOne(context.Background(), filter)

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}

	var token Token

	result.Decode(&token)

	return &token, nil
}

func (dao *MongoDao) GetAccessToken(accessTokenId string) (*Token, error) {
	filter := bson.M{"access_token_id": accessTokenId}

	result := dao.collection.FindOne(context.Background(), filter)

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, nil
	}

	var token Token

	result.Decode(&token)

	return &token, nil
}

func (dao *MongoDao) HasToken(id string) (bool, error) {
	filter := bson.M{"$or": bson.A{
		bson.M{"access_token_id": id},
		bson.M{"refresh_token_id": id},
	}}

	count, err := dao.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
