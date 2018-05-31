package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"time"
	"io/ioutil"
	"github.com/nats-io/nats"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"sync"
	"os"
	"html/template"
)

const (
	hosts      = "sentiment-mongodb:27017"
	database   = "sentiment_db"
	username   = ""
	password   = ""
	collection = "sentiment"
)

type Sent struct {
	Date         		time.Time
	TwitterScore 		int32 `json:"score"`
	HackerNewsScore 	int32 `json:"score"`
}

type MongoStore struct {
	session *mgo.Session
}

var mongoStore = MongoStore{}

var nc *nats.Conn

func main() {
	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to NATS server " + uri)

	//Create MongoDB session
	session := initialiseMongo()
	mongoStore.session = session

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sentiment", sentGetHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":9090", router))
}

func initialiseMongo() (session *mgo.Session) {

	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	return
}

func sentGetHandler(w http.ResponseWriter, r *http.Request) {

	col := mongoStore.session.DB(database).C(collection)

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Get twitter sentiment
	tweetSentiment := tr.Sentiment{}
	go func() {
		msg, err := nc.Request("TwitterSentimentByText", data, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with Twitter SentimentByText. Waiting 2 seconds before retrying:", err)
			return
		}
		receivedSentiment := tr.Sentiment{}
		err = proto.Unmarshal(msg.Data, &receivedSentiment)
		if err != nil {
			log.Print("ERROR ", err)
		}
		tweetSentiment = receivedSentiment
		wg.Done()
	}()

	// Get hacker news sentiment
	hackerNewsSentiment := tr.Sentiment{}
	go func() {
		msg, err := nc.Request("HackerNewsSentimentByText", nil, 3000*time.Millisecond)
		if err != nil {
			fmt.Println("Something went wrong with HNews SentimentByText. Waiting 2 seconds before retrying:", err)
			return
		}
		receivedSentiment := tr.Sentiment{}
		err = proto.Unmarshal(msg.Data, &receivedSentiment)
		if err != nil {
			log.Print("ERROR ", err)
		}
		hackerNewsSentiment = receivedSentiment
		wg.Done()
	}()
	wg.Wait()

	s := Sent{
		Date:         		time.Now().UTC(),
		TwitterScore: 		tweetSentiment.Score,
		HackerNewsScore:	hackerNewsSentiment.Score,
	}

	// Insert tweetSentiment with timestamp into MongoDB
	err = col.Insert(s)
	if err != nil {
		panic(err)
	}

	// Create map to hold variables to pass into html template
	m := map[string]interface{}{
		"Date":       s.Date,
		"TweetScore": s.TwitterScore,
		"HackerNewsScore": s.HackerNewsScore,
	}

	t, _ := template.ParseFiles("view.html")
	t.Execute(w, m)
}
