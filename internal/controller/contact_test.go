package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	"github.com/STreeChin/contactapi/internal/controller"
	"github.com/STreeChin/contactapi/pkg/entities"
	"github.com/STreeChin/contactapi/pkg/log"
	"github.com/STreeChin/contactapi/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
)

var mockContact = entities.Contact{
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

var postJSONStr = []byte(
	`{
				"contact": {
                 	"FirstName": "Slarty",
                 	"LastName": "Bartfast",
                	"Email": "Slarty@test.com",
                 	"custom": {
                    	 "integer--Test--Field": "1024"
                    }
				}
			}`)

var postContact = entities.Contact{FirstName: "Slarty", LastName: "Bartfast", Email: "Slarty@test.com", Custom: map[string]interface{}{"Test Field": 1024}}

func formHTTTest(act, url string, body []byte) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(act, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("autopilotapikey", "65263027fab7d440ba4c5f3b834fb800")

	w := httptest.NewRecorder()
	return req, w
}

//TestGetOneContactCtrl Use GoConvey test framework
func TestGetOneContactCtrl(t *testing.T) {
	Convey("GetOneContactCtrl", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		mockSrv := mocks.NewMockContactService(ctl)
		cCtrl := controller.NewContactController(logger, mockSrv)

		act := "Get"
		contactIDURL := "/v1/contact/person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23"
		emailURL := "/v1/contact/StGr@gmail.com"
		Convey("UT Normal Case", func() {
			Convey(" Normal Case1: 200, Get by contactId", func() {
				/*req := httptest.NewRequest("POST", "/v1/contact/person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("autopilotapikey", "65263027fab7d440ba4c5f3b834fb800")*/
				req, w := formHTTTest(act, contactIDURL, nil)
				defer monkey.UnpatchAll()
				monkey.Patch(mux.Vars, func(r *http.Request) map[string]string {
					return map[string]string{"contact_id_or_email": "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23"}
				})
				gomock.InOrder(
					mockSrv.EXPECT().GetOneContact("contactid", "person_AP2-9cbf7ac0-eec5-11e4-87bc-6df09cc44d23").Return(&mockContact, nil),
				)

				cCtrl.GetOneContactCtrl(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				result := new(entities.Contact)
				_ = json.NewDecoder(w.Body).Decode(result)
				expected := &mockContact
				So(result.ContactID, ShouldEqual, expected.ContactID)
				So(result.Email, ShouldEqual, expected.Email)
			})

			Convey("UT Normal Case2 : 200, Get by email", func() {
				req, w := formHTTTest(act, emailURL, nil)
				defer monkey.UnpatchAll()
				monkey.Patch(mux.Vars, func(r *http.Request) map[string]string {
					return map[string]string{"contact_id_or_email": "StGr@gmail.com"}
				})
				gomock.InOrder(
					mockSrv.EXPECT().GetOneContact("email", "StGr@gmail.com").Return(&mockContact, nil),
				)

				cCtrl.GetOneContactCtrl(w, req)

				So(w.Code, ShouldEqual, http.StatusOK)
				result := new(entities.Contact)
				_ = json.NewDecoder(w.Body).Decode(result)
				expected := &mockContact
				So(result.ContactID, ShouldEqual, expected.ContactID)
				So(result.Email, ShouldEqual, expected.Email)
			})

		})

		Convey("UT AbNormal Case1 : 404 contact could not be found", func() {
			req, w := formHTTTest(act, emailURL, nil)

			defer monkey.UnpatchAll()
			monkey.Patch(mux.Vars, func(r *http.Request) map[string]string {
				return map[string]string{"contact_id_or_email": "StGr@gmail.com"}
			})
			gomock.InOrder(
				mockSrv.EXPECT().GetOneContact("email", "StGr@gmail.com").Return(nil, mongo.ErrNoDocuments),
			)

			cCtrl.GetOneContactCtrl(w, req)

			So(w.Code, ShouldEqual, http.StatusNotFound)
			result := map[string]string{}
			_ = json.NewDecoder(w.Body).Decode(&result)
			expected := map[string]string{"error": "Not Found", "message": "Contact could not be found."}
			So(result["error"], ShouldEqual, expected["error"])
			So(result["message"], ShouldEqual, expected["message"])
		})

		Convey("UT AbNormal Case2 : 400 Invalid email", func() {
			req, w := formHTTTest(act, "/v1/contact/StGrgmail.com", nil)

			defer monkey.UnpatchAll()
			monkey.Patch(mux.Vars, func(r *http.Request) map[string]string {
				return map[string]string{"contact_id_or_email": "StGrgmail.com"}
			})

			cCtrl.GetOneContactCtrl(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
			result := map[string]string{}
			_ = json.NewDecoder(w.Body).Decode(&result)
			expected := map[string]string{"error": "Bad Request", "message": "Invalid contact_id_or_email value provided."}
			So(result["error"], ShouldEqual, expected["error"])
			So(result["message"], ShouldEqual, expected["message"])
		})

		Convey("UT AbNormal Case3 : 400 Invalid contactId", func() {
			req, w := formHTTTest(act, "/v1/contact/---9cbf7ac0-eec5-11e4-87bc-6df09cc44d23", nil)
			defer monkey.UnpatchAll()
			monkey.Patch(mux.Vars, func(r *http.Request) map[string]string {
				return map[string]string{"contact_id_or_email": "---9cbf7ac0-eec5-11e4-87bc-6df09cc44d23"}
			})

			cCtrl.GetOneContactCtrl(w, req)

			So(w.Code, ShouldEqual, http.StatusBadRequest)
			result := map[string]string{}
			_ = json.NewDecoder(w.Body).Decode(&result)
			expected := map[string]string{"error": "Bad Request", "message": "Invalid contact_id_or_email value provided."}
			So(result["error"], ShouldEqual, expected["error"])
			So(result["message"], ShouldEqual, expected["message"])
		})

		Convey("UT AbNormal Case4 : 500 Internal Server Error", func() {
			req, w := formHTTTest(act, emailURL, nil)
			defer monkey.UnpatchAll()
			monkey.Patch(mux.Vars, func(r *http.Request) map[string]string {
				return map[string]string{"contact_id_or_email": "StGr@gmail.com"}
			})
			gomock.InOrder(
				mockSrv.EXPECT().GetOneContact("email", "StGr@gmail.com").Return(nil, errors.New("other error")),
			)

			cCtrl.GetOneContactCtrl(w, req)

			So(w.Code, ShouldEqual, http.StatusInternalServerError)
			result := map[string]string{}
			_ = json.NewDecoder(w.Body).Decode(&result)
			expected := map[string]string{"error": "Internal Server Error", "message": "other error"}
			So(result["error"], ShouldEqual, expected["error"])
			So(result["message"], ShouldEqual, expected["message"])
		})

		Convey("UT Abnormal Case5: req and w is nil", func() {
			var w http.ResponseWriter
			var req *http.Request
			cCtrl.GetOneContactCtrl(w, req)
			So(w, ShouldEqual, nil)
		})
	})
}

//TestAddUpdateContactCtrl Use GoConvey test framework
func TestAddOrUpdateContactCtrl(t *testing.T) {
	Convey("AddOrUpdateContactCtrl", t, func() {
		ctl := gomock.NewController(t)
		defer ctl.Finish()
		mockCfg := mocks.NewMockConfig(ctl)
		logger := log.NewLogger(mockCfg)
		mockSrv := mocks.NewMockContactService(ctl)
		cCtrl := controller.NewContactController(logger, mockSrv)
		act := "POST"
		contactURL := "/v1/contact"

		Convey("UT Normal Case", func() {
			Convey("UT Case1 Normal: 200", func() {
				req, w := formHTTTest(act, contactURL, postJSONStr)
				gomock.InOrder(
					mockSrv.EXPECT().AddOrUpdateContact(&postContact).Return("person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7", nil),
				)

				cCtrl.AddOrUpdateContactCtrl(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
				result := map[string]string{}
				_ = json.Unmarshal(w.Body.Bytes(), &result)
				So(result["contact_id"], ShouldNotEqual, nil)
				expected := map[string]string{"contact_id": "person_9EAF39E4-9AEC-4134-964A-D9D8D54162E7"}
				So(result["contact_id"], ShouldEqual, expected["contact_id"])
			})

			Convey("UT Case2 Normal: 200", func() {
				var postJSONStrBool = []byte(
					`{
				"contact": {
                 	"FirstName": "Slarty",
                 	"LastName": "Bartfast",
                	"Email": "Slarty@test.com",
                 	"custom": {
                    	 "boolean--Test--Field": "true"
                    }
				}
			}`)
				var postContactBool = entities.Contact{FirstName: "Slarty", LastName: "Bartfast", Email: "Slarty@test.com", Custom: map[string]interface{}{"Test Field": true}}
				req, w := formHTTTest(act, contactURL, postJSONStrBool)
				gomock.InOrder(
					mockSrv.EXPECT().AddOrUpdateContact(&postContactBool).Return("person_9EAF39E4-9AEC-4134-964A-D9D8D5410002", nil),
				)

				cCtrl.AddOrUpdateContactCtrl(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
				result := map[string]string{}
				_ = json.Unmarshal(w.Body.Bytes(), &result)
				So(result["contact_id"], ShouldNotEqual, nil)
				expected := map[string]string{"contact_id": "person_9EAF39E4-9AEC-4134-964A-D9D8D5410002"}
				So(result["contact_id"], ShouldEqual, expected["contact_id"])
			})

			Convey("UT Case3 Normal: 200", func() {
				var postJSONStrBool = []byte(
					`{
				"contact": {
                 	"FirstName": "Slarty",
                 	"LastName": "Bartfast",
                	"Email": "Slarty@test.com",
                 	"custom": {
                    	 "float--Test--Field": "3.0"
                    }
				}
			}`)
				var postContactBool = entities.Contact{FirstName: "Slarty", LastName: "Bartfast", Email: "Slarty@test.com", Custom: map[string]interface{}{"Test Field": 3.0}}
				req, w := formHTTTest(act, contactURL, postJSONStrBool)
				gomock.InOrder(
					mockSrv.EXPECT().AddOrUpdateContact(&postContactBool).Return("person_9EAF39E4-9AEC-4134-964A-D9D8D5410002", nil),
				)

				cCtrl.AddOrUpdateContactCtrl(w, req)
				So(w.Code, ShouldEqual, http.StatusOK)
				result := map[string]string{}
				_ = json.Unmarshal(w.Body.Bytes(), &result)
				So(result["contact_id"], ShouldNotEqual, nil)
				expected := map[string]string{"contact_id": "person_9EAF39E4-9AEC-4134-964A-D9D8D5410002"}
				So(result["contact_id"], ShouldEqual, expected["contact_id"])
			})
		})

		Convey("AbNormal Case:", func() {
			Convey("UT Abnormal Case1: 400, No email provided.", func() {
				var jsonStr = []byte(
					`{"contact": {
				"FirstName": "Slarty",
				"LastName": "Bartfast",
				"custom": {
					"integer--Test--Field": "1024"
				}
			  }
			}`)
				req, w := formHTTTest(act, contactURL, jsonStr)

				cCtrl.AddOrUpdateContactCtrl(w, req)

				So(w.Code, ShouldEqual, http.StatusBadRequest)
				result := map[string]string{}
				_ = json.Unmarshal(w.Body.Bytes(), &result)
				expected := map[string]string{"error": "Bad Request", "message": "No contact details provided."}
				So(result["error"], ShouldEqual, expected["error"])
				So(result["message"], ShouldEqual, expected["message"])
			})

			Convey("UT Abnormal Case2: 500, AddOrUpdateContact return error", func() {
				req, w := formHTTTest(act, contactURL, postJSONStr)
				gomock.InOrder(
					mockSrv.EXPECT().AddOrUpdateContact(&postContact).Return("", errors.New("AddOrUpdateContact return error")),
				)
				cCtrl.AddOrUpdateContactCtrl(w, req)

				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				result := map[string]string{}
				_ = json.Unmarshal(w.Body.Bytes(), &result)
				expected := map[string]string{"error": "Internal Server Error", "message": "AddOrUpdateContact return error"}
				So(result["error"], ShouldEqual, expected["error"])
				So(result["message"], ShouldEqual, expected["message"])
			})

			Convey("UT Abnormal Case3: 500, parseContactFromReq return error", func() {
				var jsonStr = []byte(`{
					"contact": {
						"FirstName": "Slarty",
						"LastName": "Bartfast",
						"custom": {
							"invalidType--Test--Field": "1024"
						}
					}
				}`)
				req, w := formHTTTest(act, contactURL, jsonStr)

				cCtrl.AddOrUpdateContactCtrl(w, req)
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				result := map[string]string{}
				_ = json.Unmarshal(w.Body.Bytes(), &result)
				expected := map[string]string{"error": "Internal Server Error", "message": "Internal Error"}
				So(result["error"], ShouldEqual, expected["error"])
				So(result["message"], ShouldEqual, expected["message"])
			})

			Convey("UT Abnormal Case4: req and w is nil", func() {
				var w http.ResponseWriter
				var req *http.Request
				cCtrl.AddOrUpdateContactCtrl(w, req)
				So(w, ShouldEqual, nil)
			})
		})
	})
}
