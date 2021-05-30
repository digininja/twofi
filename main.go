package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/dghubble/go-twitter/twitter"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	//fmt.Println("Starting")

	flags := struct {
		consumerKey    string
		consumerSecret string
		term           string
		records        int
	}{}

	flag.StringVar(&flags.consumerKey, "consumer-key", "", "Twitter Consumer Key")
	flag.StringVar(&flags.consumerSecret, "consumer-secret", "", "Twitter Consumer Secret")
	flag.StringVar(&flags.term, "term", "", "Search Term")
	flag.IntVar(&flags.records, "records", 5, "Number of search results")
	flag.Parse()
	//flagutil.SetFlagsFromEnv(flag.CommandLine, "TWITTER")

	if flags.consumerKey == "" || flags.consumerSecret == "" {
		log.Fatal("Application Access Token required")
	}
	if flags.term == "" {
		log.Fatal("No search term provided")
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     flags.consumerKey,
		ClientSecret: flags.consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// search tweets
	// example
	// https://github.com/dghubble/go-twitter/blob/master/twitter/search.go
	searchTweetParams := &twitter.SearchTweetParams{
		Query: flags.term,
		// TweetMode: "extended",
		Count: flags.records,
	}

	search, _, _ := client.Search.Tweets(searchTweetParams)
	// Print all
	// fmt.Printf("%+v\n", search.Statuses)

	var re = regexp.MustCompile(`[^\w \s \d]`)

	for _, res := range search.Statuses {
		fmt.Printf("XXXXXXXXXXXXXXXX\n")
		fmt.Printf("%s\n", res.Text)
		clean := re.ReplaceAllString(res.Text, ` `)
		fmt.Printf("%s\n", clean)
		words := strings.Fields(clean)
		for i, word := range words {
			fmt.Printf("%d, %s\n", i, word)
		}

		// fmt.Printf("%+v\n", res)
	}
	//fmt.Printf("SEARCH TWEETS:\n%+v\n", search)
	//fmt.Printf("SEARCH METADATA:\n%+v\n", search.Metadata)

	// Search Tweets
	/*
		search, resp, err := client.Search.Tweets(&twitter.SearchTweetParams{
			Query: "gopher",
		})
		fmt.Printf("search: %s\n", search)
		fmt.Println("Done")
	*/
}
