package twitterapi_client

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
	pb "github.com/BrianCoveney/TwitterStreaming/twitter_route"
	"fmt"
	"reflect"
	"github.com/BrianCoveney/stringutil/credentials"
)

var client *twitter.Client

func Auth() {
	// The twitter package provides a Client for accessing the Twitter API.
	// Here we create a new twitter client using the credentials in the config file
	config := oauth1.NewConfig(viper.GetString(credentials.ConsumerKey),
		viper.GetString(credentials.ConsumerSecret))
	token := oauth1.NewToken(viper.GetString(credentials.AccessToken),
		viper.GetString(credentials.AccessSecret))

	// Authentication is handled by the http.Client which is
	// passed to twitter NewClient as Oauth1
	httpClient := config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)

}

// GetStream : returns new stream with the given params
func GetStream(params *pb.Params) (*twitter.Stream, error) {

	// Convert request params to a twitter specific struct
	tp := &twitter.StreamFilterParams{
		Track:         params.Track,
		StallWarnings: twitter.Bool(params.StallWarnings),
	}

	stream, err := client.Streams.Filter(tp)
	if err != nil { // Handle errors reading the config file
		return nil, err
	}

	fmt.Print(reflect.TypeOf(stream))

	return stream, nil

}
