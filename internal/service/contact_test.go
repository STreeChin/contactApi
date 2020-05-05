package service_test

import (
	"testing"

	"bou.ke/monkey"
	"github.com/STreeChin/contactapi/internal/service"
	"github.com/STreeChin/contactapi/pkg/entities"
	"github.com/STreeChin/contactapi/pkg/log"
	"github.com/STreeChin/contactapi/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

/*var postJsonStr = []byte(
	`{
				"getContact": {
                 	"FirstName": "Slarty",
                 	"LastName": "Bartfast",
                	"Email": "test@slarty.com",
                 	"custom": {
                    	 "integer--Test--Field": "1024"
                    }
				}
			}`)*/

//var postContact = model.Contact{FirstName: "Slarty", LastName: "Bartfast", Email: "test@slarty.com", Custom: map[string]interface{}{"Test Field": 1024}}
var getContact = entities.Contact{
	ContactID:  "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23",
	Email:      "StGr@gmail.com",
	FirstName:  "St",
	LastName:   "Gr",
	Type:       "Contact",
	Phone:      "4159945916",
	CreatTime:  "2015-04-29T23:15:25.347Z",
	UpdateTime: "2015-04-29T23:15:25.347Z",
	LeadSource: "Autopilot",
	Status:     "Testing",
	Company:    "Magpie API",
	Lists:      []string{"contactlist_9EAF39E4-9AEC-4134-964A"},
}
var postContact = entities.Contact{
	ContactID: "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23",
	Email:     "StGr@gmail.com",
	FirstName: "St",
	LastName:  "Gr",
	Type:      "Contact",
	Phone:     "4159945916",
}

//TestGetOneContact
func TestGetOneContact(t *testing.T) {
	Convey("TestGetOneContact", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		mockCache := mocks.NewMockCache(ctl)
		mockRep := mocks.NewMockRepository(ctl)
		cSrv := service.NewContactService(logger, mockCfg, mockCache, mockRep)

		expected := &getContact
		keyContactID, valueContactID := "contactid", getContact.ContactID
		//errCache := errors.Wrap(redis.ErrNil, "email")
		Convey("UT Normal Case", func() {
			Convey("UT Normal Case1: Get from cache by email", func() {
				key, value := "email", "StGr@gmail.com"
				gomock.InOrder(
					mockCache.EXPECT().GetOneContact(getContact.Email).Return(&getContact, nil),
				)

				cont, err := cSrv.GetOneContact(key, value)
				So(err, ShouldEqual, nil)
				So(cont.ContactID, ShouldEqual, expected.ContactID)
				So(cont.Email, ShouldEqual, expected.Email)
			})

			Convey("UT Normal Case2: Get from cache by contactID", func() {
				gomock.InOrder(
					mockCache.EXPECT().GetEmailByContactID(valueContactID).Return(getContact.Email, nil),
					mockCache.EXPECT().GetOneContact(getContact.Email).Return(&getContact, nil),
				)

				cont, err := cSrv.GetOneContact(keyContactID, valueContactID)
				So(err, ShouldEqual, nil)
				So(cont.ContactID, ShouldEqual, expected.ContactID)
				So(cont.Email, ShouldEqual, expected.Email)
			})

			Convey("UT Normal Case3: Get from db by contactId", func() {
				gomock.InOrder(
					mockCache.EXPECT().GetEmailByContactID(valueContactID).Return("", errors.New("redigo: nil returned")),
					mockRep.EXPECT().GetOneContact(keyContactID, valueContactID).Return(&getContact, nil),
					mockCache.EXPECT().SetOneContact(getContact.Email, &getContact).Return(nil),
					mockCache.EXPECT().SetEmailByContactID(getContact.ContactID, getContact.Email).Return(nil),
				)
				defer monkey.UnpatchAll()
				monkey.Patch(errors.Cause, func(error) error {
					return redis.ErrNil
				})
				cont, err := cSrv.GetOneContact(keyContactID, valueContactID)
				So(err, ShouldEqual, nil)
				So(cont.ContactID, ShouldEqual, expected.ContactID)
				So(cont.Email, ShouldEqual, expected.Email)
			})
		})

		Convey("AbNormal Case:", func() {
			gomock.InOrder(
				mockCache.EXPECT().GetEmailByContactID(valueContactID).Return("", redis.ErrNil),
			)

			Convey("AbNormal Case1: Get from db fail", func() {
				err := errors.New("read db fail")
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact(keyContactID, valueContactID).Return(nil, err),
				)

				cont, resErr := cSrv.GetOneContact(keyContactID, valueContactID)
				So(errors.Cause(resErr), ShouldEqual, err)
				So(cont, ShouldEqual, nil)
			})

			Convey("Normal Case2: Set cache fail", func() {
				err := errors.New("set cache fail")
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact(keyContactID, valueContactID).Return(&getContact, nil),
					mockCache.EXPECT().SetOneContact(getContact.Email, &getContact).Return(err),
					mockCache.EXPECT().SetEmailByContactID(getContact.ContactID, getContact.Email).Return(err),
				)

				cont, resErr := cSrv.GetOneContact(keyContactID, valueContactID)
				So(resErr, ShouldEqual, nil)
				So(cont.ContactID, ShouldEqual, cont.ContactID)
			})
		})
	})
}

//TestAddOrUpdateContact
func TestAddOrUpdateContact(t *testing.T) {
	Convey("TestAddOrUpdateContact", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		mockCache := mocks.NewMockCache(ctl)
		mockRep := mocks.NewMockRepository(ctl)
		cCtrl := service.NewContactService(logger, mockCfg, mockCache, mockRep)

		Convey("Normal Case", func() {
			expected := postContact.ContactID
			Convey("Normal Case1: Update to db", func() {
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact("email", postContact.Email).Return(&postContact, nil),
					mockRep.EXPECT().UpdateOneContact(&postContact).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.Email).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.ContactID).Return(nil),
				)

				contactID, err := cCtrl.AddOrUpdateContact(&postContact)
				So(err, ShouldEqual, nil)
				So(contactID, ShouldEqual, expected)
			})

			Convey("UT Normal Case2: Insert to db", func() {
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact("email", postContact.Email).Return(&postContact, errors.New("find fail")),
					mockRep.EXPECT().InsertOneContact(&postContact).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.Email).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.ContactID).Return(nil),
				)

				contactID, err := cCtrl.AddOrUpdateContact(&postContact)
				So(err, ShouldEqual, nil)
				So(contactID, ShouldEqual, expected)
			})
		})

		Convey("AbNormal Case", func() {
			Convey("AbNormal Case1: insert to db fail", func() {
				err := errors.New("insert to db fail")
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact("email", postContact.Email).Return(&postContact, errors.New("find fail")),
					mockRep.EXPECT().InsertOneContact(&postContact).Return(err),
				)

				_, resErr := cCtrl.AddOrUpdateContact(&postContact)
				So(errors.Cause(resErr), ShouldEqual, err)
			})
			Convey("AbNormal Case2: update to db fail", func() {
				err := errors.New("update to db fail")
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact("email", postContact.Email).Return(&postContact, nil),
					mockRep.EXPECT().UpdateOneContact(&postContact).Return(err),
				)

				_, resErr := cCtrl.AddOrUpdateContact(&postContact)
				So(errors.Cause(resErr), ShouldEqual, err)
			})
			Convey("AbNormal Case3: del cache fail", func() {
				err := errors.New("del cache fail")
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact("email", postContact.Email).Return(&postContact, errors.New("find fail")),
					mockRep.EXPECT().InsertOneContact(&postContact).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.Email).Return(err),
				)

				_, resErr := cCtrl.AddOrUpdateContact(&postContact)
				So(errors.Cause(resErr), ShouldEqual, err)
			})
			Convey("AbNormal Case4: del cache fail", func() {
				err := errors.New("del cache fail")
				gomock.InOrder(
					mockRep.EXPECT().GetOneContact("email", postContact.Email).Return(&postContact, errors.New("find fail")),
					mockRep.EXPECT().InsertOneContact(&postContact).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.Email).Return(nil),
					mockCache.EXPECT().DelOneContact(postContact.ContactID).Return(err),
				)

				_, resErr := cCtrl.AddOrUpdateContact(&postContact)
				So(errors.Cause(resErr), ShouldEqual, err)
			})
		})
	})
}
