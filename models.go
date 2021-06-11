package go_utils

import (
	"fmt"
	"net/url"
	"time"
)

// QueryParam is an interface used for filter and sort parameters
type QueryParam interface {
	ToURLValues() (values url.Values)
}

// PaginationInput represents paging parameters
type PaginationInput struct {
	First  int    `json:"first"`
	Last   int    `json:"last"`
	After  string `json:"after"`
	Before string `json:"before"`
}

//IsEntity ...
func (p PaginationInput) IsEntity() {}

// SortInput is a generic container for strongly typed sorting parameters
type SortInput struct {
	SortBy []*SortParam `json:"sortBy"`
}

//IsEntity ...
func (s SortInput) IsEntity() {}

// SortParam represents a single field sort parameter
type SortParam struct {
	FieldName string    `json:"fieldName"`
	SortOrder SortOrder `json:"sortOrder"`
}

//IsEntity ...
func (s SortParam) IsEntity() {}

// FilterInput is s generic container for strongly type filter parameters
type FilterInput struct {
	Search   *string        `json:"search"`
	FilterBy []*FilterParam `json:"filterBy"`
}

//IsEntity ...
func (f FilterInput) IsEntity() {}

// FilterParam represents a single field filter parameter
type FilterParam struct {
	FieldName           string      `json:"fieldName"`
	FieldType           FieldType   `json:"fieldType"`
	ComparisonOperation Operation   `json:"comparisonOperation"`
	FieldValue          interface{} `json:"fieldValue"`
}

//IsEntity ...
func (f FilterParam) IsEntity() {}

// PhoneOptIn is used to persist and manage phone communication whitelists
type PhoneOptIn struct {
	MSISDN  string `json:"msisdn" firestore:"msisdn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

//IsEntity ...
func (p PhoneOptIn) IsEntity() {}

// USSDSessionLog is used to persist a log of USSD sessions
type USSDSessionLog struct {
	MSISDN    string `json:"msisdn" firestore:"msisdn"`
	SessionID string `json:"sessionID" firestore:"sessionID"`
}

//IsEntity ...
func (p USSDSessionLog) IsEntity() {}

// EmailOptIn is used to persist and manage email communication whitelists
type EmailOptIn struct {
	Email   string `json:"email" firestore:"optedIn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

//IsEntity ...
func (e EmailOptIn) IsEntity() {}

// FirebaseRefreshResponse is used to (de)serialize the results of a successful Firebase token refresh
type FirebaseRefreshResponse struct {
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

// PIN is used to store a PIN (Personal Identifiation Number) associated
// to a phone number sign up to Firebase
type PIN struct {
	UID     string `json:"uid" firestore:"uid"`
	MSISDN  string `json:"msisdn,omitempty" firestore:"msisdn"`
	PIN     string `json:"pin,omitempty" firestore:"pin"`
	IsValid bool   `json:"isValid,omitempty" firestore:"isValid"`
}

//IsEntity ...
func (p PIN) IsEntity() {}

// Upload represents a file uploaded to cloud storage
type Upload struct {
	ID          string    `json:"id" firestore:"id"`
	URL         string    `json:"url" firestore:"url"`
	Size        int       `json:"size" firestore:"size"`
	Hash        string    `json:"hash" firestore:"hash"`
	Creation    time.Time `json:"creation" firestore:"creation"`
	Title       string    `json:"title" firestore:"title"`
	ContentType string    `json:"contentType" firestore:"contentType"`
	Language    string    `json:"language" firestore:"language"`
	Base64data  string    `json:"base64data" firestore:"base64data"`
}

// IsEntity marks upload as an apollo federation entity
func (u Upload) IsEntity() {}

// IsNode marks upload as a relay node
func (u Upload) IsNode() {}

// SetID marks upload as a relay node
func (u Upload) SetID(id string) {
	u.ID = id
}

// GetID marks upload as a relay node
func (u Upload) GetID() ID {
	return IDValue(u.ID)
}

// UploadInput is used to send data for new uploads
type UploadInput struct {
	Title       string `json:"title"`
	ContentType string `json:"contentType"`
	Language    string `json:"language"`
	Base64data  string `json:"base64data"`
	Filename    string `json:"filename"`
}

// CustomError represents a custom error struct
// Reference https://blog.golang.org/error-handling-and-go
type CustomError struct {
	Err     error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
