package repository_test

import (
	"fmt"
	"testing"

	"bou.ke/monkey"
	"github.com/STreeChin/contactapi/mocks"
	"github.com/STreeChin/contactapi/pkg/config"
	"github.com/STreeChin/contactapi/pkg/entities"
	"github.com/STreeChin/contactapi/pkg/log"
	"github.com/STreeChin/contactapi/pkg/repository"
	"github.com/STreeChin/contactapi/pkg/route/middleware/crpt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
)

var getContact = entities.Contact{
	ContactID: "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc43333",
	Email:     "StGr@gmail.com",
	FirstName: "St",
	LastName:  "Gr",
	Type:      "Contact",
	Phone:     "4159945916",
}
var insertThenDelContact = entities.Contact{
	ContactID: "person_AP2-9cbf7ac0-eec5-11e4-87bc-01insert1111",
	Email:     "PC@gmail.com",
	FirstName: "PI",
	LastName:  "CI",
	Type:      "Contact",
	Phone:     "4159941111",
}
var updateContact = entities.Contact{
	ContactID: "person_AP2-9cbf7ac0-eec5-11e4-87bc-01update0000",
	Email:     "PC@gmail.com",
	FirstName: "PU",
	LastName:  "CU",
	Type:      "Contact",
	Phone:     "4159940000",
}
var dbc = config.DatabaseConfig{URL: "mongodb://localhost:27017"}

/*var readDoc = bson.M{
	"contactid": "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc40000",
	"email":     "StGr@gmail.com",
	"firstName": "St",
	"lastName":  "Gr",
}*/

//TestGetOneContact
func TestGetOneContact(t *testing.T) {
	Convey("TestGetOneContact", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		gomock.InOrder(
			mockCfg.EXPECT().GetDBConfig().Return(&dbc).AnyTimes(),
		)
		rep := repository.NewRepository(logger, mockCfg)

		key, value := "email", getContact.Email
		_ = rep.InsertOneContact(&getContact)

		Convey("UT Normal Case", func() {
			Convey("UT Normal Case1: Get from cache by email", func() {
				cont, err := rep.GetOneContact(key, value)
				expected := &getContact
				So(err, ShouldEqual, nil)
				So(cont.ContactID, ShouldEqual, expected.ContactID)
				So(cont.Email, ShouldEqual, expected.Email)
			})
		})

		Convey("UT AbNormal Case", func() {
			Convey("UT AbNormal Case1: mongo fail", func() {
				cont, err := rep.GetOneContact(key, "abnormal.com")
				So(errors.Cause(err), ShouldEqual, mongo.ErrNoDocuments)
				So(cont, ShouldEqual, nil)
			})

			Convey("UT AbNormal Case2: AesEncrypt fail", func() {
				errResult := errors.New("AesEncrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesEncrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				cont, err := rep.GetOneContact(key, value)
				So(errors.Cause(err), ShouldEqual, errResult)
				So(cont, ShouldEqual, nil)
			})

			Convey("UT AbNormal Case3: AesDecrypt fail", func() {
				errResult := errors.New("AesDecrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesDecrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				cont, err := rep.GetOneContact(key, value)
				So(errors.Cause(err), ShouldEqual, errResult)
				So(cont, ShouldEqual, nil)
			})

			/*			Convey("UT AbNormal Case4: bson.Marshal fail", func() {
						errMarshal := errors.New("bson.Marshal fail")
						defer monkey.UnpatchAll()
						monkey.Patch(bson.Marshal, func(interface{}) ([]byte, error) {
							return nil, errMarshal
						})

						cont, errM := rep.GetOneContact(key, value)
						So(errors.Cause(errM), ShouldEqual, errMarshal)
						So(cont, ShouldEqual, nil)
					})*/
		})
	})
}

//TestInsertOneContact
func TestInsertOneContact(t *testing.T) {
	Convey("TestInsertOneContact", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		gomock.InOrder(
			mockCfg.EXPECT().GetDBConfig().Return(&dbc).AnyTimes(),
		)
		rep := repository.NewRepository(logger, mockCfg)

		key, value := "email", insertThenDelContact.Email

		Convey("UT Normal Case", func() {
			Convey("UT Normal Case1: insert one", func() {
				err := rep.DeleteOneContact(key, value)
				if err != nil {
					panic(err)
				}

				err = rep.InsertOneContact(&insertThenDelContact)
				So(err, ShouldEqual, nil)

				expected := &insertThenDelContact
				result, err := rep.GetOneContact(key, value)
				if err != nil {
					panic(err)
				}

				So(result.ContactID, ShouldEqual, expected.ContactID)
				So(result.Email, ShouldEqual, expected.Email)
			})
		})

		Convey("UT AbNormal Case", func() {
			Convey("UT AbNormal Case1: mongo fail: dup email", func() {
				err := rep.InsertOneContact(&insertThenDelContact)
				fmt.Println(err)
				So(errors.Cause(err), ShouldNotEqual, nil)
			})

			Convey("UT AbNormal Case2: AesEncrypt fail", func() {
				errResult := errors.New("AesEncrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesEncrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				err := rep.InsertOneContact(&insertThenDelContact)
				So(errors.Cause(err), ShouldEqual, errResult)
			})

			Convey("UT AbNormal Case3: AesDecrypt fail", func() {
				err := rep.DeleteOneContact(key, value)
				if err != nil {
					panic(err)
				}

				errResult := errors.New("AesDecrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesDecrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				err = rep.InsertOneContact(&insertThenDelContact)
				So(errors.Cause(err), ShouldEqual, errResult)
			})

		})
	})
}

//TestUpdateOneContact
func TestUpdateOneContact(t *testing.T) {
	Convey("TestUpdateOneContact", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		gomock.InOrder(
			mockCfg.EXPECT().GetDBConfig().Return(&dbc).AnyTimes(),
		)
		rep := repository.NewRepository(logger, mockCfg)

		key, value := "email", updateContact.Email
		_ = rep.InsertOneContact(&insertThenDelContact)

		Convey("UT Normal Case", func() {
			Convey("UT Normal Case1: update one", func() {
				err := rep.UpdateOneContact(&updateContact)
				So(err, ShouldEqual, nil)

				expected := &updateContact
				result, err := rep.GetOneContact(key, value)
				if err != nil {
					panic(err)
				}
				So(result.ContactID, ShouldEqual, expected.ContactID)
				So(result.Email, ShouldEqual, expected.Email)
				So(result.FirstName, ShouldEqual, expected.FirstName)
			})
		})

		Convey("UT AbNormal Case", func() {
			Convey("UT AbNormal Case1: AesEncrypt fail", func() {
				errResult := errors.New("AesEncrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesEncrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				err := rep.UpdateOneContact(&updateContact)
				So(errors.Cause(err), ShouldEqual, errResult)
			})

			Convey("UT AbNormal Case2: AesDecrypt fail", func() {
				err := rep.DeleteOneContact(key, value)
				if err != nil {
					panic(err)
				}

				errResult := errors.New("AesDecrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesDecrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				err = rep.UpdateOneContact(&updateContact)
				So(errors.Cause(err), ShouldEqual, errResult)
			})

		})
	})
}

//TestInsertOneContact
func TestDeleteOneContact(t *testing.T) {
	Convey("TestDeleteOneContact", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		gomock.InOrder(
			mockCfg.EXPECT().GetDBConfig().Return(&dbc).AnyTimes(),
		)
		rep := repository.NewRepository(logger, mockCfg)

		key, value := "email", insertThenDelContact.Email
		_ = rep.InsertOneContact(&insertThenDelContact)

		Convey("UT Normal Case", func() {
			Convey("UT Normal Case1: delete one", func() {
				err := rep.DeleteOneContact(key, value)
				So(err, ShouldEqual, nil)
				result, err := rep.GetOneContact(key, value)
				So(errors.Cause(err), ShouldEqual, mongo.ErrNoDocuments)
				So(result, ShouldEqual, nil)
			})
		})

		Convey("UT AbNormal Case", func() {
			Convey("UT AbNormal Case1: mongo fail: dup email", func() {
				err := rep.DeleteOneContact(key, "norecord@g.com")
				fmt.Println(err)
				So(errors.Cause(err), ShouldEqual, nil)
			})

			Convey("UT AbNormal Case2: AesEncrypt fail", func() {
				_ = rep.InsertOneContact(&insertThenDelContact)
				errResult := errors.New("AesEncrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesEncrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				err := rep.DeleteOneContact(key, value)
				So(errors.Cause(err), ShouldEqual, errResult)
			})
		})
	})
}

//TestGetOneContact
func TestGetContactIdByApiKey(t *testing.T) {
	Convey("TestGetContactIdByApiKey", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		gomock.InOrder(
			mockCfg.EXPECT().GetDBConfig().Return(&dbc).AnyTimes(),
		)
		rep := repository.NewRepository(logger, mockCfg)
		apiKey, contactID := "65263027fab7d440ba4c5f3b834fb800", "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23"

		Convey("UT Normal Case", func() {
			Convey("UT Normal Case1: Get from cache by email", func() {
				result, err := rep.GetContactIDByAPIKey(apiKey)
				So(err, ShouldEqual, nil)
				So(result, ShouldEqual, contactID)
			})
		})

		Convey("UT AbNormal Case", func() {
			Convey("UT AbNormal Case1: mongo fail", func() {
				result, err := rep.GetContactIDByAPIKey("invalidKey")
				So(errors.Cause(err), ShouldEqual, mongo.ErrNoDocuments)
				So(result, ShouldEqual, "")
			})

			Convey("UT AbNormal Case2: AesEncrypt fail", func() {
				errResult := errors.New("AesEncrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesEncrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				result, err := rep.GetContactIDByAPIKey(apiKey)
				So(errors.Cause(err), ShouldEqual, errResult)
				So(result, ShouldEqual, "")
			})

			Convey("UT AbNormal Case3: AesDecrypt fail", func() {
				errResult := errors.New("AesDecrypt fail")
				defer monkey.UnpatchAll()
				monkey.Patch(crpt.AesDecrypt, func([]byte) ([]byte, error) {
					return nil, errResult
				})

				result, err := rep.GetContactIDByAPIKey(apiKey)
				So(errors.Cause(err), ShouldEqual, errResult)
				So(result, ShouldEqual, "")
			})
		})
	})
}
