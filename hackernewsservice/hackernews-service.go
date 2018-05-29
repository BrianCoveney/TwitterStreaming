package main

import (
	"fmt"
	"github.com/nats-io/nats"
	"os"
	"github.com/peterhellberg/hn"
	tr "github.com/BrianCoveney/TwitterStreaming/transport"
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

	nc.QueueSubscribe("TwitterByText", "TwitterTeller", publishHackerNewsFromStream)
	select {} // Block forever
}

func publishHackerNewsFromStream(m *nats.Msg) {
	hn := hn.DefaultClient
	ids, err := hn.TopStories()
	if err != nil {
		panic(err)
	}

	myHackerNews := &tr.HackerNews{}
	news = myHackerNews.News

	for i, id := range ids[:10] {
		item, err := hn.Item(id)
		if err != nil {
			panic(err)
		}

		news = append(news, item.Text)
		fmt.Println("NEWS: ", news)

		fmt.Println(i, "–", item.Title, "\n   ", item.URL, "\n")
	}
}
