package internal

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/dchest/uniuri"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"labix.org/v2/mgo/bson"
)

type Database struct {
	db *mongo.Database
}

type User struct {
	ID      primitive.ObjectID `bson:"_id, omitempty"`
	Details struct {
		Username string
		Name     string
		Email    string
	}
	Profile struct {
		Colour string
		Status string
		Bio    string
		URL    string
	}
	Agora struct {
		Score struct {
			Posts int32
			Stars int32
		}
	}

	ConnectedApps struct {
		Discord string
	}
}

func NewDatabase(mongoDatabase *mongo.Database) *Database {
	database := new(Database)

	database.db = mongoDatabase

	return database
}

func (database *Database) GetMemberCount() (int64, error) {
	return database.db.Collection("users").CountDocuments(context.TODO(), bson.M{})
}

func (database *Database) GetUserByEmail(email string) (User, error) {
	email, err := govalidator.NormalizeEmail(email)
	if err != nil {
		return User{}, err
	}

	filter := bson.M{"details.email": email}

	var user User

	err = database.db.Collection("users").FindOne(
		context.Background(),
		filter).Decode(&user)

	return user, err
}

func (database *Database) GetUserByUsername(username string) (User, error) {
	filter := bson.M{"details.username": username}

	var user User

	collation := new(options.Collation)
	collation.Locale = "en"
	collation.Strength = 2

	queryOptions := new(options.FindOneOptions)
	queryOptions.Collation = collation

	err := database.db.Collection("users").FindOne(
		context.Background(),
		filter, queryOptions).Decode(&user)

	return user, err
}

func (database *Database) GetUserByID(id string) (User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err
	}

	filter := bson.M{"_id": objectID}

	var user User

	err = database.db.Collection("users").FindOne(
		context.Background(),
		filter).Decode(&user)

	return user, err
}

func (database *Database) GetUsersByID(ids primitive.M) ([]User, error) {
	findOptions := options.Find()
	findOptions.Projection = bson.M{
		"details.name":   true,
		"profile.status": true,
	}

	cur, err := database.db.Collection("users").Find(context.TODO(), ids, findOptions)
	if err != nil {
		return []User{}, err
	}

	var users []User

	cur.All(context.Background(), &users)

	err = cur.Close(context.TODO())

	return users, err
}

func (database *Database) GetUserByDiscordID(discordID string) (User, error) {
	filter := bson.M{"connectedApps.discord": discordID}

	var user User

	err := database.db.Collection("users").FindOne(
		context.Background(),
		filter).Decode(&user)

	return user, err
}

func (database *Database) InsertNotification(notification primitive.M) error {
	_, err := database.db.Collection("notifications").InsertOne(context.TODO(), notification)

	return err
}

func (database *Database) ConnectDUAccount(discordToken string, discordID string) (string, error) {
	filter := bson.M{"discordToken": discordToken}
	update := bson.M{
		"$unset": bson.M{
			"discordToken": true,
		}, "$set": bson.M{
			"connectedApps.discord": discordID,
		},
	}

	queryOptions := new(options.FindOneAndUpdateOptions)
	var returnDocument options.ReturnDocument = options.After
	queryOptions.ReturnDocument = &returnDocument

	var user User

	err := database.db.Collection("users").FindOneAndUpdate(
		context.Background(),
		filter,
		update).Decode(&user)

	return user.Details.Email, err
}

func (database *Database) AddDiscordVerificationToken(id string) (string, error) {
	token := uniuri.NewLen(16)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"discordToken": token,
		},
	}

	_, err = database.db.Collection("users").UpdateOne(
		context.TODO(),
		filter,
		update)

	return token, err
}
