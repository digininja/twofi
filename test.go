package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
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
	s := make([]string, 4)

	s[0] = "RT @campuscodi: steel_con steel_con \"An employee of Deloitte's India office has been found to be the mastermind of a computer hacking gang that targeted Britis…"
	s[1] = "@fleetwoodmlyric @n00bznet I had a bad curry last night, you probably don't want to be with me EVERYWHERE."
	s[2] = "@nihsuyhp @Steel_Con @TimmehWimmy @SleepyEntropy @infosecmo @tamonten This looks good!  https://t.co/8ceByzw2r7"
	s[3] = "@HackingDave Is there anything useful in the video part of this? I've just skimmed through and don't see anything.… https://t.co/xBKML44PaR"

	min_length := 5
	show_count := false

	var results map[string]int
	results = make(map[string]int)

	fmt.Println()
	for word, count := range results {
		fmt.Printf("%s: %d\n", word, count)
	}
	fmt.Println()

	for _, str := range s {
		re, err := regexp.Compile(`[^\w \s \d]`)
		if err != nil {
			fmt.Printf(err.Error())
		}
		str = re.ReplaceAllString(str, " ")
		fmt.Println(str)

		words := strings.Fields(str)
		for _, word := range words {
			if len(word) >= min_length {
				fmt.Printf("word: %s\n", word)
				if _, ok := results[word]; ok {
					results[word] += 1
				} else {
					results[word] = 1
				}
			}
		}
	}
	keys := sortResults(results)
	for _, k := range keys {
		if show_count {
			fmt.Printf("%s: %d\n", k, results[k])
		} else {
			fmt.Printf("%s\n", k)
		}
	}
}
