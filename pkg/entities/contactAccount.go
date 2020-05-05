package entities

//Contact  a struct to parse the request, the difference between with the Contact is the custom field
type Contact struct {
	// used internally
	ContactID         string   `json:"contact_id"`
	Email             string   `json:"Email"`
	Twitter           string   `json:"Twitter"`
	FirstName         string   `json:"FirstName"`
	LastName          string   `json:"LastName"`
	Salutation        string   `json:"Salutation"`
	Company           string   `json:"Company"`
	NumberOfEmployees string   `json:"NumberOfEmployees"`
	Title             string   `json:"Title"`
	Industry          string   `json:"Industry"`
	Phone             string   `json:"Phone"`
	MobilePhone       string   `json:"MobilePhone"`
	Fax               string   `json:"Fax"`
	Website           string   `json:"Website"`
	MailingStreet     string   `json:"MailingStreet"`
	MailingCity       string   `json:"MailingCity"`
	MailingState      string   `json:"MailingState"`
	MailingPostalCode string   `json:"MailingPostalCode"`
	MailingCountry    string   `json:"MailingCountry"`
	LeadSource        string   `json:"LeadSource"`
	Status            string   `json:"Status"`
	LinkedIn          string   `json:"LinkedIn"`
	Lists             []string `json:"lists"`
	// used internally
	Type string `json:"type"`
	// used internally
	CreatTime string `json:"created_at"`
	// used internally
	UpdateTime string `json:"updated_at"`
	// used internally
	OwnerName string `json:"owner_name"`
	// used internally
	Unsubscribed bool `json:"unsubscribed"`
	// used internally
	Custom map[string]interface{} `json:"custom"`
	// used internally
	AutopilotSessionID string `json:"_autopilot_session_id"`
	// used internally
	AutopilotList string `json:"_autopilot_list"`
	// used internally
	Notify string `json:"notify"`
}

//ReqContact request
type ReqContact struct {
	Contact Contact `json:"contact"`
}
