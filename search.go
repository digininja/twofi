package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func sortResults(results map[string]int) []string {
	keys := make([]string, 0, len(results))

	for key := range results {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return results[keys[i]] > results[keys[j]]
	})

	return keys
}

func main() {
	min_length := 5
	search_count := 10
	term := "from:digininja"

	flags := struct {
		consumerKey    string
		consumerSecret string
	}{}

	flag.StringVar(&flags.consumerKey, "consumer-key", "", "Twitter Consumer Key")
	flag.StringVar(&flags.consumerSecret, "consumer-secret", "", "Twitter Consumer Secret")
	flag.Parse()
	flagutil.SetFlagsFromEnv(flag.CommandLine, "TWITTER")

	if flags.consumerKey == "" || flags.consumerSecret == "" {
		log.Fatal("Application Access Token required")
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
	searchTweetParams := &twitter.SearchTweetParams{
		Query: term,
		Count: search_count,
	}
	search, _, _ := client.Search.Tweets(searchTweetParams)
	//fmt.Printf("SEARCH TWEETS:\n%+v\n", search)

	//fmt.Printf("\n\n")
	//fmt.Printf("Result count = %d\n\n", len(search.Statuses))
	//fmt.Printf("\n\n")

	re, err := regexp.Compile(`[^\w \s \d]`)
	if err != nil {
		log.Fatal(err.Error())
	}

	//fmt.Printf("\n\n")

	var results map[string]int
	results = make(map[string]int)

	show_count := true

	for _, status := range search.Statuses {
		//fmt.Printf("Text: %s\n", status.Text)
		//fmt.Printf("Full Text: %s\n", status.FullText)

		str := re.ReplaceAllString(status.Text, " ")
		//fmt.Println(str)

		words := strings.Fields(str)
		for _, word := range words {
			if len(word) >= min_length {
				//fmt.Printf("word: %s\n", word)
				if _, ok := results[word]; ok {
					results[word] += 1
				} else {
					results[word] = 1
				}
			}
		}
	}
	keys := sortResults (results)
	for _, k := range keys{
		if show_count {
			fmt.Printf("%s: %d\n", k, results[k])
		} else {
			fmt.Printf("%s\n", k)
		}
	}
}
