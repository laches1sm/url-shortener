package main

import (
	"bytes"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/laches1sm/url_shortener/adapters"
	"github.com/laches1sm/url_shortener/http"
	"github.com/laches1sm/url_shortener/infrastructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mongodb"
	"os"
)

const (
	serviceRunning = `Url Shortener Service now running...`
)

func main() {
	var buf bytes.Buffer
	logger := log.New(io.MultiWriter(&buf), "parrot", 0)
	mongoSession, err := mongodb.NewMongoDBClientWithMSI("bb585a49-702f-4762-a11e-66866bo78729", "rebecca-audition", "rebecca-parrot-db", azure.PublicCloud)
	if err != nil{
		os.Exit(1)
	}
	mongoCollection := mongoSession.DB("rebecca-audition").C("parrots")
	opts := &options.ClientOptions{}
	mClient, err := mongo.NewClient(opts)
	if err != nil{
		os.Exit(1)
	}
	parrotInfra := infrastructure.UrlInfra{
		MongoClient: *mClient,
		MongoCollection: *mongoCollection,
	}
	urlAdapter := adapters.NewUrlShortenerAdapter(*logger, parrotInfra)

	server := http.NewParrotServer(*logger, &urlAdapter, &parrotInfra)
	server.SetupRoutes()

	logger.Println(serviceRunning)
	if err := server.Start(http.ServerPort); err != nil {
		logger.Println(err.Error())
	}
}