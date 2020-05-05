package database

import (
	"fmt"

	"github.com/STreeChin/contactapi/pkg/entities"
	"github.com/STreeChin/contactapi/pkg/route/middleware/crpt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//InitMongoDB init db
func InitMongoDB(m *mongoDB) {
	//InitAPIKey(m)
	InitContactInfo(m)
}

//InitAPIKey init apiKey
func InitAPIKey(mgoDB *mongoDB) {
	apiKey := "65263027fab7d440ba4c5f3b834fb800"
	contactID := "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23"
	encryptAPIKey, _ := crpt.AesEncrypt([]byte(apiKey))
	encryptUserID, _ := crpt.AesEncrypt([]byte(contactID))
	/*	oneDoc := APIKeyAccount{
		ContactID: string(encryptUserID),
		APIKey:    string(encryptApiKey),
	}*/
	type APIKeyAccountByte struct {
		ContactID []byte `json:"contactid"`
		APIKey    []byte `json:"apikey"`
	}

	oneDoc := APIKeyAccountByte{
		ContactID: encryptUserID,
		APIKey:    encryptAPIKey,
	}

	contact := new(entities.Contact)
	//dbDoc, err = mgoDB.FindOne("contact", "apiKey", "apikey", string(encryptApiKey))
	dbDoc, err := mgoDB.FindOne("contact", "apiKey", "apikey", encryptAPIKey)
	if err != nil {
		fmt.Println("Init DB apiKey: no document in DB apiKey")
		if err := mgoDB.InsertOne("contact", "apiKey", oneDoc); err != nil {
			fmt.Println("Init DB apiKey: insert DB  apiKey fail.")
		}
		fmt.Println("Init DB apiKey: insert  document to DB apiKey success. ")

		readDoc, err := mgoDB.FindOne("contact", "apiKey", "apikey", encryptAPIKey)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(readDoc["contactid"].(primitive.Binary).Data)
		contactID, _ := crpt.AesDecrypt(readDoc["contactid"].(primitive.Binary).Data)
		fmt.Println("Init DB apiKey: read DB apiKey success. ", string(contactID))

		return
	}
	bsonBytes, _ := bson.Marshal(dbDoc)
	_ = bson.Unmarshal(bsonBytes, contact)
	//ts := dbDoc["contactid"].(string)
	dbContactID, err := crpt.AesDecrypt([]byte(contact.ContactID))
	if err != nil {
		logrus.Error("Init DB apiKey: : AesDecrypt error.")
	}
	fmt.Println("Init DB apiKey: read apiKey success. ", string(dbContactID), contactID)

}

//InitContactInfo init contactInfo
func InitContactInfo(mgoDB *mongoDB) {
	contact := new(entities.Contact)
	contactID := "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23"
	email := "repositoryInit@gmail.com"
	contact.ContactID = contactID
	contact.Email = email
	contact.FirstName = "rep"
	contact.LastName = "InitDB"
	contact.Phone = "4159945916"

	enContactID, _ := crpt.AesEncrypt([]byte(contact.ContactID))
	enEmail, _ := crpt.AesEncrypt([]byte(contact.Email))

	contact.ContactID = string(enContactID)
	contact.Email = string(enEmail)
	//readDoc, err := mgoDB.FindOne("contact", "contactInfo", "email", "init@test.com")
	readDoc, err := mgoDB.FindOne("contact", "contactInfo", "email", string(enEmail))
	if err != nil {
		fmt.Println("Init DB contactInfo: read fail.", err)
		if err := mgoDB.InsertOne("contact", "contactInfo", contact); err != nil {
			fmt.Println("Init DB contactInfo: insert DB fail.", err)
			return
		}
		fmt.Println("Init DB contactInfo: insert DB success. ", contactID, email)
		return
	}

	dbContactID, _ := crpt.AesDecrypt([]byte(readDoc["contactid"].(string)))
	dbEmail, _ := crpt.AesDecrypt([]byte(readDoc["email"].(string)))
	fmt.Println("Init DB contactInfo: read success. ", string(dbEmail), string(dbContactID), contact.FirstName)
}
