package infrastructure

import (
	"context"
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/laches1sm/url_shortener/models"
	mongo "github.com/globalsign/mgo"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"strings"
	"time"
"github.com/satori/go.uuid"
	"crypto/md5"
	"encoding/hex"

)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)


type UrlInfra struct{
	MongoClient mongo2.Client
	MongoCollection mongo.Collection
}

// CreateURLDocument is a function that will take a url that the user has provided and place it within the MongoDB
 func (url *UrlInfra) CreateURLDocument (ctx context.Context, request *models.UrlRequest) (*models.UrlResponse, error){
 	
	 // Create a random string that will serve as the short form of the url.
	 shortUrl := createShortUrl(3)

	 // Check if randomly generated string already exists in the db. Unlikely, but best to safeguard against this.

	 query := url.MongoCollection.Find(shortUrl).Iter()
	 var results []bson.M
	 if err := query.All(&results); err != nil{
	 	return nil, err
	 }
	 for _, result := range results{
	 	for _, resultVal := range result{
	 		if shortUrl == resultVal{
	 			// Create new short url if it already exists in the db
	 			shortUrl = createShortUrl(3)
			}
		}
	 }

	 // Create new uuid for the URL ID.
	 urlID := uuid.NewV4()
	 newUrl := &models.UrlResponse{
		 ShortUrl:  shortUrl,
		 LongUrl:   request.UrlToBeShortened,
		 Clicks:    0,
		 CreatedAt: time.Now(),
		 URLId: urlID.String(),
	 }

	 err := url.MongoCollection.Insert(newUrl)
 	if err != nil {
 		return nil, errors.New("error while trying to insert url into mongo")
	}

 	return newUrl, nil
 }


 // I found this good solution to making a random string here: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
 // Highly recommend reading, the answerer explains things really well and I've learnt a lot from this post!
func createShortUrl(n int) string{
	var src = rand.NewSource(time.Now().UnixNano())
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	hash := md5.Sum([]byte(sb.String()))
	// Reduce the hash down
	return hex.EncodeToString(hash[:n])

}

func (url *UrlInfra) GetUrlDocument(ctx context.Context, shortURL *models.ShortURLRequest) (*models.UrlResponse, error){
	result := models.UrlResponse{}
	filter := bson.D{{"url", shortURL.ShortUrl}}
	err := url.MongoCollection.Find(filter).One(&result)
	if err != nil{
		return nil, err
	}
	return &result, nil
}