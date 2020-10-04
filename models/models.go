package models

import "time"

// UrlRequest is a representation of the long-form url sent to the API.
type UrlRequest struct {
	UrlToBeShortened string `json:"url"`
}

// ShortURLRequest is a representation of when the user receives the shortened url and wants to retrieve the long form URL from the API
type ShortURLRequest struct{
	ShortUrl string `json:"url"`

}

// UrlResponse is a representation of the response the user should expect to get if they send a valid request to the API.
// It contains both the shortered version of the URL and the orginal form of the URL.
type UrlResponse struct{
	ShortUrl string `json:"url"`
	LongUrl string `json:"long_url"`
	Clicks int `json:"clicks"`
	CreatedAt time.Time `json:"createdAt"`
	URLId string `json:"id"`
}

// ErrorResponse is a representation of the response the user should expect to get if they send an invalid response to the API,
// or of an error that has occurred when handling their request.
type ErrorResponse struct {
	Message string `json:"err_message"`
}
