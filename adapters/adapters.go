package adapters

import (
	"context"
	"encoding/json"
	"github.com/laches1sm/url_shortener/infrastructure"
	"github.com/laches1sm/url_shortener/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type UrlShortenerAdapter struct{
	Logger log.Logger
	Infra infrastructure.UrlInfra
}

func NewUrlShortenerAdapter(log log.Logger, infra infrastructure.UrlInfra) UrlShortenerAdapter{
	return UrlShortenerAdapter{
		Logger: log,
		Infra: infra,
	}
}

func (adapter *UrlShortenerAdapter) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		adapter.Logger.Printf(`Not a valid POST Url Shortner Request, aborting with error...`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		adapter.Logger.Printf(`Error while reading request body: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// unmarshal the body into a request struct
	requestBody := &models.UrlRequest{}
	if err = json.Unmarshal(body, requestBody); err != nil{
		adapter.Logger.Printf(`Error while trying to unmarshal request: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	// read resp body...
	resp, err := adapter.Infra.CreateURLDocument(context.Background(), requestBody)
	if err != nil {
		return
	}
	urlGet, err := json.Marshal(resp)
	if err != nil {
		return
	}
	writeResponse(w, urlGet, http.StatusOK)
}

func (adapter *UrlShortenerAdapter) RedirectUrl(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		adapter.Logger.Printf(`Not a valid POST Url Shortner Request, aborting with error...`)
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		adapter.Logger.Printf(`Error while reading request body: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// unmarshal the body into a request struct
	requestBody := &models.ShortURLRequest{}
	if err = json.Unmarshal(body, requestBody); err != nil{
		adapter.Logger.Printf(`Error while trying to unmarshal request: %s`, err.Error())
		_ = marshalAndWriteErrorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	// read resp body...
	resp, err := adapter.Infra.GetUrlDocument(context.Background(), requestBody)
	if err != nil {
		return
	}

	// Basic check to see if the long-form url is not malformed in some way.
    if resp.LongUrl != ""{
		http.Redirect(w, r, resp.LongUrl, 301)
	}else{
		marshalAndWriteErrorResponse(w, `error while trying to redirect: long url not found`, 400)
	}


}
func marshalAndWriteErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) error {
	msg := map[string]string{
		"message": errorMessage,
	}
	responseBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := writeResponse(w, responseBody, statusCode); err != nil {
		return err
	}

	return nil
}

func writeResponse(w http.ResponseWriter, body []byte, statusCode int) error {
	w.Header().Set(`Content-Type`, `application/json`)
	w.Header().Set(`Content-Length`, strconv.Itoa(len(body)))

	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
		return err
	}

	return nil
}
