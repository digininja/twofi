package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strings"

	"github.com/dghubble/go-twitter/twitter"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	//fmt.Println("Starting")

	flags := struct {
		apiKey     string
		apiSecret  string
		term       string
		records    int
		min_length int
		inc_counts bool
		inc_tweets bool
		debug      bool
	}{}

	flag.StringVar(&flags.apiKey, "api-key", "", "Twitter API Key")
	flag.StringVar(&flags.apiSecret, "api-secret", "", "Twitter API Secret")
	flag.StringVar(&flags.term, "term", "", "Search Term - For tweets from a specific user prefix the user with 'from:'")
	flag.IntVar(&flags.records, "records", 5, "Number of search results")
	flag.IntVar(&flags.min_length, "min-length", 5, "Minimum word length")
	flag.BoolVar(&flags.inc_counts, "inc-counts", false, "Include counts")
	flag.BoolVar(&flags.inc_tweets, "inc-tweets", false, "Include tweets")
	flag.BoolVar(&flags.debug, "debug", false, "Debug")
	flag.Parse()
	//flagutil.SetFlagsFromEnv(flag.CommandLine, "TWITTER")

	if flags.apiKey == "" || flags.apiSecret == "" {
		log.Fatal("Application Access Token required")
	}
	if flags.term == "" {
		log.Fatal("No search term provided")
	}

	if flags.debug {
		log.SetLevel(log.DebugLevel)
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     flags.apiKey,
		ClientSecret: flags.apiSecret,
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
		// TweetMode: "extended", - When used, the result is in FullText not in Text
		Count: flags.records,
	}

	search, _, _ := client.Search.Tweets(searchTweetParams)
	// Print all
	log.Debugf("%+v\n", search.Statuses)

	var re = regexp.MustCompile(`[^\w \s \d]`)

	// Structure of a tweet - called a status
	// https://github.com/dghubble/go-twitter/blob/master/twitter/statuses.go
	// Twitter API docs
	// https://developer.twitter.com/en/docs/twitter-api/tweets/search/quick-start/recent-search
	data := make(map[string]int)

	for _, status := range search.Statuses {
		// Dump the full response
		log.Debugf("%+v\n", status)
		if flags.inc_tweets {
			fmt.Printf("%s\n", status.Text)
			// fmt.Printf("%s\n", status.FullText)
		}
		clean := re.ReplaceAllString(status.Text, ` `)
		log.Debugf("%s\n", clean)
		words := strings.Fields(clean)
		for _, word := range words {
			if len(word) >= flags.min_length {
				if _, ok := data[word]; !ok {
					data[word] = 0
				}
				data[word]++
			}
		}
	}

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range data {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		if flags.inc_counts {
			fmt.Printf("%s, %d\n", kv.Key, kv.Value)
		} else {
			fmt.Printf("%s\n", kv.Key)
		}
	}
}
