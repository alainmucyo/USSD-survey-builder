package dto

type USSDRequest struct {
	SessionId   string `json:"sessionId"`
	TrackId     string `json:"trackId,omitempty"`
	Text        string `json:"text"`
	PhoneNumber string `json:"phoneNumber"`
	ServiceCode string `json:"serviceCode" binding:"-"`
	Action      string `json:"action" binding:"-"`
	Response    string `json:"response,omitempty" binding:"-"`
}
