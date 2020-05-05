package repository

import (
	"github.com/STreeChin/contactapi/pkg/config"
	"github.com/STreeChin/contactapi/pkg/database"
	"github.com/STreeChin/contactapi/pkg/entities"
	"github.com/STreeChin/contactapi/pkg/route/middleware/crpt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//DBHandler db
type DBHandler interface {
	FindOne(db, coll, key string, value interface{}) (bson.M, error)
	InsertOne(db, coll string, value interface{}) error
	UpdateOne(db, coll string, key string, value interface{}, update interface{}) error
	DeleteOne(db, coll, key string, value interface{}) error
}

type repository struct {
	log       *logrus.Logger
	DbHandler DBHandler
}

//NewRepository instance
func NewRepository(log *logrus.Logger, cfg config.Config) *repository {
	dbHandler := database.NewDataStore(log, cfg)
	// debug
	database.InitMongoDB(dbHandler)
	return &repository{log, dbHandler}
}

func (r *repository) GetOneContact(key, value string) (*entities.Contact, error) {
	var err error
	var readDoc bson.M
	var contact *entities.Contact

	value, err = r.repEncrypt(value)
	if err != nil {
		return nil, errors.WithMessage(err, "rep getOneContact")
	}

	readDoc, err = r.DbHandler.FindOne("contact", "contactInfo", key, value)
	if err != nil {
		return nil, errors.WithMessage(err, "rep getOneContact")
	}

	contact, err = r.bsonToContact(readDoc)
	if err != nil {
		return nil, errors.WithMessage(err, "rep getOneContact")
	}

	if contact.ContactID != "" {
		contact.ContactID, err = r.repDecrypt(contact.ContactID)
		if err != nil {
			return nil, errors.WithMessage(err, "rep getOneContact")
		}
	}
	if contact.Email != "" {
		contact.Email, err = r.repDecrypt(contact.Email)
	}

	return contact, errors.WithMessage(err, "rep getOneContact")
}

func (r *repository) InsertOneContact(contact *entities.Contact) error {
	var err error
	contact.Email, err = r.repEncrypt(contact.Email)
	if err != nil {
		return errors.WithMessage(err, "rep insertOneContact")
	}
	contact.ContactID, err = r.repEncrypt(contact.ContactID)
	if err != nil {
		return errors.WithMessage(err, "rep insertOneContact")
	}

	err = r.DbHandler.InsertOne("contact", "contactInfo", contact)
	if err != nil {
		return errors.WithMessage(err, "rep insertOneContact")
	}

	contact.Email, err = r.repDecrypt(contact.Email)
	if err != nil {
		return errors.WithMessage(err, "rep insertOneContact")
	}

	contact.ContactID, err = r.repDecrypt(contact.ContactID)
	return errors.WithMessage(err, "rep insertOneContact")
}

func (r *repository) UpdateOneContact(contact *entities.Contact) error {
	var err error
	contact.Email, err = r.repEncrypt(contact.Email)
	if err != nil {
		return errors.WithMessage(err, "rep updateOneContact")
	}
	contact.ContactID, err = r.repEncrypt(contact.ContactID)
	if err != nil {
		return errors.WithMessage(err, "rep updateOneContact")
	}

	err = r.DbHandler.UpdateOne("contact", "contactInfo", "email", contact.Email, contact)
	if err != nil {
		return errors.WithMessage(err, "rep UpdateOneContact")
	}

	contact.Email, err = r.repDecrypt(contact.Email)
	if err != nil {
		return errors.WithMessage(err, "rep updateOneContact")
	}

	contact.ContactID, err = r.repDecrypt(contact.ContactID)
	return errors.WithMessage(err, "rep updateOneContact")
}

func (r *repository) DeleteOneContact(key, value string) error {
	value, err := r.repEncrypt(value)
	if err != nil {
		return errors.WithMessage(err, "rep deleteOneContact")
	}

	err = r.DbHandler.DeleteOne("contact", "contactInfo", key, value)
	if err != nil {
		return errors.WithMessage(err, "rep deleteOneContact")
	}

	return errors.WithMessage(err, "rep deleteOneContact")
}

func (r *repository) GetContactIDByAPIKey(apiKey string) (string, error) {
	query, err := crpt.AesEncrypt([]byte(apiKey))
	if err != nil {
		return "", errors.WithMessage(err, "rep getContactId encrypt")
	}

	readDoc, err := r.DbHandler.FindOne("contact", "apiKey", "apikey", query)
	if err != nil {
		return "", errors.WithMessage(err, "rep getContactId find")
	}

	contactID, err := crpt.AesDecrypt(readDoc["contactid"].(primitive.Binary).Data)
	return string(contactID), errors.WithMessage(err, "rep getContactId decrypt")
}

func (r *repository) bsonToContact(doc bson.M) (*entities.Contact, error) {
	bsonBytes, err := bson.Marshal(doc)
	if err != nil {
		return nil, errors.Wrap(err, "rep bsonToContact")
	}

	contact := new(entities.Contact)
	err = bson.Unmarshal(bsonBytes, contact)

	return contact, errors.Wrap(err, "rep bsonToContact")
}

func (r *repository) repEncrypt(key string) (string, error) {
	encryptKey, err := crpt.AesEncrypt([]byte(key))
	return string(encryptKey), errors.WithMessage(err, "rep repEncrypt")
}

func (r *repository) repDecrypt(key string) (string, error) {
	decryptKey, err := crpt.AesDecrypt([]byte(key))
	return string(decryptKey), errors.WithMessage(err, "rep repDecrypt")
}
