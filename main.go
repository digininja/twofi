package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strings"

	"github.com/dghubble/go-twitter/twitter"

	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"os"
)

type Configuration struct {
	Term      string
	Key       string
	Secret    string
	MinLength int
	NumTweets int
	IncCounts bool
	IncTweets bool
	Debug     bool
}

func getConfig(configFile string) Configuration {
	if _, err := os.Stat(configFile); err != nil {
		fmt.Println("Config file not found")
		os.Exit(1)
	}
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("Error opening config file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("Error parsing config file: %s\n", err)
	}
	log.Debug("Key: %s\n", configuration.Key)
	log.Debug("Secret: %s\n", configuration.Secret)
	log.Debug("Min Length: %d\n", configuration.MinLength)

	return configuration
}

func main() {
	flags := struct {
		apiKey    string
		apiSecret string
		term      string
		numTweets int
		minLength int
		incCounts bool
		incTweets bool
		debug     bool
	}{}

	// Use this to help set default values for integer and bool params
	// https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go

	flag.StringVar(&flags.apiKey, "api-key", "", "Twitter API Key")
	flag.StringVar(&flags.apiSecret, "api-secret", "", "Twitter API Secret")
	flag.StringVar(&flags.term, "term", "", "Search Term - For tweets from a specific user prefix the user with 'from:'")
	flag.IntVar(&flags.numTweets, "num-tweets", 5, "Number of tweets to look at")
	flag.IntVar(&flags.minLength, "min-length", 5, "Minimum word length")
	flag.BoolVar(&flags.incCounts, "inc-counts", false, "Include counts")
	flag.BoolVar(&flags.incTweets, "inc-tweets", false, "Include tweets")
	flag.BoolVar(&flags.debug, "debug", false, "Debug")
	flag.Parse()
	//flagutil.SetFlagsFromEnv(flag.CommandLine, "TWITTER")

	configFile := "config.json"
	configuration := getConfig(configFile)

	if flags.apiKey != "" {
		configuration.Key = flags.apiKey
	}
	if flags.apiSecret != "" {
		configuration.Secret = flags.apiSecret
	}

	if configuration.Key == "" || configuration.Secret == "" {
		log.Fatal("Application Access Tokens required.\nSee README for more information.")
	}

	/*
				if flags.numTweets != "" {
					configuration.MinLength = flags.numTweets
				}

			if configuration.NumTweets == "" {
				configuration.NumTweets = 10
			}
			/*

				if flags.minLength != "" {
					configuration.MinLength = flags.minLength
				}

		if configuration.MinLength == "" {
			configuration.MinLength = 5
		}
	*/

	if flags.term != "" {
		configuration.Term = flags.term
	}

	if configuration.Term == "" {
		log.Fatal("No search term provided")
	}

	/*
		if flags.debug != "" {
			configuration.Debug = flags.debug
		}
	*/

	if configuration.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     configuration.Key,
		ClientSecret: configuration.Secret,
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
		Query: configuration.Term,
		// TweetMode: "extended", - When used, the result is in FullText not in Text
		Count: configuration.NumTweets,
	}

	search, _, _ := client.Search.Tweets(searchTweetParams)
	// Print all
	log.Debugf("%+v\n", search.Statuses)

	// Regexp to remove any non-letter, symbol or number
	var re = regexp.MustCompile(`[^\w \s \d]`)

	// Structure of a tweet - called a status
	// https://github.com/dghubble/go-twitter/blob/master/twitter/statuses.go
	// Twitter API docs
	// https://developer.twitter.com/en/docs/twitter-api/tweets/search/quick-start/recent-search
	data := make(map[string]int)

	log.Debugf("Got %d tweets returned", len(search.Statuses))

	for _, status := range search.Statuses {
		// Dump the full response
		log.Debugf("%+v\n", status)
		if flags.incTweets {
			fmt.Printf("%s\n", status.Text)
			fmt.Println(strings.Repeat("-", 30))
			// fmt.Printf("%s\n", status.FullText)
		}
		clean := re.ReplaceAllString(status.Text, ` `)
		log.Debugf("%s\n", clean)
		words := strings.Fields(clean)
		for _, word := range words {
			if len(word) >= flags.minLength {
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
		if flags.incCounts {
			fmt.Printf("%s, %d\n", kv.Key, kv.Value)
		} else {
			fmt.Printf("%s\n", kv.Key)
		}
	}
}
