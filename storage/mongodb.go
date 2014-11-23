package storage

import (
	"log"

	"github.com/ian-kent/Go-MailHog/data"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// MongoDB represents MongoDB backed storage backend
type MongoDB struct {
	Session    *mgo.Session
	Collection *mgo.Collection
}

// CreateMongoDB creates a MongoDB backed storage backend
func CreateMongoDB(uri, db, coll string) *MongoDB {
	log.Printf("Connecting to MongoDB: %s\n", uri)
	session, err := mgo.Dial(uri)
	if err != nil {
		log.Printf("Error connecting to MongoDB: %s", err)
		return nil
	}
	return &MongoDB{
		Session:    session,
		Collection: session.DB(db).C(coll),
	}
}

// Store stores a message in MongoDB and returns its storage ID
func (mongo *MongoDB) Store(m *data.Message) (string, error) {
	err := mongo.Collection.Insert(m)
	if err != nil {
		log.Printf("Error inserting message: %s", err)
		return "", err
	}
	return string(m.ID), nil
}

// List returns a list of messages by index
func (mongo *MongoDB) List(start int, limit int) (*data.Messages, error) {
	messages := &data.Messages{}
	err := mongo.Collection.Find(bson.M{}).Skip(start).Limit(limit).Select(bson.M{
		"id":              1,
		"_id":             1,
		"from":            1,
		"to":              1,
		"content.headers": 1,
		"content.size":    1,
		"created":         1,
	}).All(messages)
	if err != nil {
		log.Printf("Error loading messages: %s", err)
		return nil, err
	}
	return messages, nil
}

// DeleteOne deletes an individual message by storage ID
func (mongo *MongoDB) DeleteOne(id string) error {
	_, err := mongo.Collection.RemoveAll(bson.M{"id": id})
	return err
}

// DeleteAll deletes all messages stored in MongoDB
func (mongo *MongoDB) DeleteAll() error {
	_, err := mongo.Collection.RemoveAll(bson.M{})
	return err
}

// Load loads an individual message by storage ID
func (mongo *MongoDB) Load(id string) (*data.Message, error) {
	result := &data.Message{}
	err := mongo.Collection.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Printf("Error loading message: %s", err)
		return nil, err
	}
	return result, nil
}