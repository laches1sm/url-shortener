package infrastructure

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/globalsign/mgo/dbtest"
	"github.com/laches1sm/url_shortener/models"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func NewTestDBServer() dbtest.DBServer{
	return dbtest.DBServer{}
}

func NewMockMongoClient() mongo.Client{
	return mongo.Client{}
}
func TestUrlInfra_CreateURLDocument(t *testing.T) {
	Convey("Given I have a valid URL request", t, func(){
		req := models.UrlRequest{UrlToBeShortened: "parrots.com"}
		testSvr := NewTestDBServer()
		testCollection := testSvr.Session().DB("TEST").C("test")
		testClient := NewMockMongoClient()
		testInfra := UrlInfra{
			MongoClient:    testClient,
			MongoCollection: *testCollection,
		}
		Convey("When I send a request to shorten my URL", func(){
			resp, err := testInfra.CreateURLDocument(context.Background(), &req)

			Convey("I should get a valid response back from the server", func(){
				So(resp, ShouldNotBeNil)
			})
			Convey("And the error shall be nil", func(){
				So(err, ShouldBeNil)
			})
		})
    testSvr.Stop()
	})

}
func TestUrlInfra_GetUrlDocument(t *testing.T) {
	Convey("Given I have a valid short URL request", t, func(){
		short := models.ShortURLRequest{ShortUrl: "p.com"}
		testSvr := NewTestDBServer()
		testCollection := testSvr.Session().DB("TEST").C("test")
		testClient := NewMockMongoClient()
		testInfra := UrlInfra{
			MongoClient:    testClient,
			MongoCollection: *testCollection,
		}
		_ = testCollection.Insert(bson.M{"short_url": "p.com"})

		Convey("When I send a valid request ", func() {
			resp, err := testInfra.GetUrlDocument(context.Background(), &short)

			Convey("The result should be valid", func(){
				So(resp, ShouldNotBeNil)
			})
			Convey("And the error should be nil", func(){
				So(err, ShouldBeNil)
			})
		})
	})
	Convey("Given I have a short URL request that does not exist in the database", t, func(){
		short := models.ShortURLRequest{ShortUrl: "notaparrot"}
		testSvr := NewTestDBServer()
		testCollection := testSvr.Session().DB("TEST").C("test")
		testClient := NewMockMongoClient()
		testInfra := UrlInfra{
			MongoClient:    testClient,
			MongoCollection: *testCollection,
		}
		Convey("When I send this request", func(){
			resp, err := testInfra.GetUrlDocument(context.Background(), &short)

			Convey("The response should be nil", func(){
				So(resp, ShouldBeNil)
			})
			Convey("And the error will not be nil", func(){
				So(err, ShouldNotBeNil)
			})
		})
	})
}

