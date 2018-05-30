package main

import (
	"fmt"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats"
	"os"
	"github.com/peterhellberg/hn"
)

var nc *nats.Conn
var news []string

func main()  {
	uri := os.Getenv("NATS_URI")

	var err error

	nc, err = nats.Connect(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to NATS server " + uri)

	nc.QueueSubscribe("HackerNewsByText", "HackerNewsTeller", publishHackerNewsFromStream)
	select {} // Block forever
}

func publishHackerNewsFromStream(m *nats.Msg) {
	hn := hn.DefaultClient
	ids, err := hn.TopStories()
	if err != nil {
		panic(err)
	}

	myHackerNews := tr.HackerNews{}

	for _, id := range ids[:10] {
		item, err := hn.Item(id)
		if err != nil {
			panic(err)
		}
		news = append(news, item.Title)

		myHackerNews.News = news

		data, err := proto.Marshal(&myHackerNews)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		nc.Publish(m.Reply, data)

		nc.Flush()
	}
}
