package main

import (
	"flag"
	"fmt"
	"log"
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
	}{}

	flag.StringVar(&flags.apiKey, "api-key", "", "Twitter API Key")
	flag.StringVar(&flags.apiSecret, "api-secret", "", "Twitter API Secret")
	flag.StringVar(&flags.term, "term", "", "Search Term")
	flag.IntVar(&flags.records, "records", 5, "Number of search results")
	flag.IntVar(&flags.min_length, "min_length", 5, "Minimum word length")
	flag.Parse()
	//flagutil.SetFlagsFromEnv(flag.CommandLine, "TWITTER")

	if flags.apiKey == "" || flags.apiSecret == "" {
		log.Fatal("Application Access Token required")
	}
	if flags.term == "" {
		log.Fatal("No search term provided")
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
		// TweetMode: "extended",
		Count: flags.records,
	}

	search, _, _ := client.Search.Tweets(searchTweetParams)
	// Print all
	// fmt.Printf("%+v\n", search.Statuses)

	var re = regexp.MustCompile(`[^\w \s \d]`)

	data := make(map[string]int)

	for _, res := range search.Statuses {
		//fmt.Printf("XXXXXXXXXXXXXXXX\n")
		//	fmt.Printf("%s\n", res.Text)
		clean := re.ReplaceAllString(res.Text, ` `)
		//	fmt.Printf("%s\n", clean)
		words := strings.Fields(clean)
		for _, word := range words {
			if len(word) >= flags.min_length {
				if _, ok := data[word]; !ok {
					data[word] = 0
				}
				data[word]++

				//		fmt.Printf("%d, %s\n", i, word)
			}
		}
	}

	/*
		To sort the results alphabetically
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	*/

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
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	}

	/*
		for _, k := range keys {
			fmt.Println(k, data[k])
		}
	*/

}
