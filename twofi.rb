#!/usr/bin/env ruby
# encoding: utf-8

#
# This script takes a list of twitter usernames or search terms and generates a
# word list based on them. For usernames it requests the last 500 tweets from
# that user, for a search term it requests 500 tweets including that term.
#
# The script is based on an original idea from the
# "7 Habits of Highly Effective Hackers" blog
# http://7habitsofhighlyeffectivehackers.blogspot.com.au/2012/05/using-twitter-to-build-password.html
#
# Author:: Robin Wood (robin@digininja.org)
# Copyright:: Copyright (c) Robin Wood 2014
# Licence:: Creative Commons Attribution-Share Alike 2.0
# 
# CHANGELOG:
# * 2014-07 (@femoltor): 
#     Add optional ignore usernames feature. 
#     For non english users: Fix problem with latin an cirilyc characters.
#     Developing dictionaries to ignore non meaningfull words (subset of http://www.link-assistant.com/seo-stop-words.html)  
#
# TODO: Ignore shortened URLs
# TODO: Optional ignore hashtags (sometimes they are useful, sometimes not)

require 'yaml'
require 'twitter'
require 'getoptlong'

opts = GetoptLong.new(
  [ '--help', '-h', GetoptLong::NO_ARGUMENT ],
  [ '--config', GetoptLong::REQUIRED_ARGUMENT ],
  [ '--count', '-c', GetoptLong::NO_ARGUMENT ],
  [ '--min_word_length', "-m" , GetoptLong::REQUIRED_ARGUMENT ],
  [ '--term_file', "-T" , GetoptLong::REQUIRED_ARGUMENT ],
  [ '--terms', "-t" , GetoptLong::REQUIRED_ARGUMENT ],
  [ '--user_file', "-U" , GetoptLong::REQUIRED_ARGUMENT ],
  [ '--users', "-u" , GetoptLong::REQUIRED_ARGUMENT ],
  [ '--ignore-usernames', "-i" , GetoptLong::NO_ARGUMENT ],
  [ '--verbose', "-v" , GetoptLong::NO_ARGUMENT ]
)

def sample_config
  puts "The config file \"#{@config_file}\" is missing or invalid, please create a config file in the format:"
  puts "options:
  api_key: <YOUR KEY>
  api_secret: <YOUR SECRET>

To get your keys you must register with Twitter at: https://apps.twitter.com/
"
  exit
end

def usage
  puts 'twoif 2.0-beta Robin Wood (robin@digininja.org) (www.digininja.org)
twoif - Twitter Words of Interest

Usage: twoif [OPTIONS]
  --help, -h: show help
  --config <file>: config file, default is twofi.yml
  --count, -c: include the count with the words
  --min_word_length, -m: minimum word length
  --term_file, -T <file>: a file containing a list of terms
  --terms, -t: comma separated search terms
    quote words containing spaces, no space after commas
  --user_file, -U <file>: a file containing a list of users
  --users, -u: comma separated usernames
    quote words containing spaces, no space after commas
  --ignore-usernames, -i: Ignore the usernames mentioned in the tweets
  --verbose, -v: verbose

'
  exit
end

# Default this to nil and it is then created
# when first needed in the search

@twitter_client = nil

def twitter_search(query)
  if @twitter_client.nil?
    @twitter_client = Twitter::REST::Client.new do |config|
      config.consumer_key = @api_key
      config.consumer_secret = @api_secret
      unless @bearer_token.nil?
        config.bearer_token = @bearer_token
      end
    end
  end

  begin
    data = @twitter_client.search(query, :result_type => "recent")
  rescue Twitter::Error::RequestTimeout
    puts "There was a timeout trying to connect to Twitter."
    puts "Please check your network connection and try again.\n\n"
    exit
  rescue Twitter::Error::Forbidden, Twitter::Error::Unauthorized
    puts "The authentication with Twitter failed, please check your API keys."
    puts "If there is a bearer_token entry in your config file try removing that.\n\n"
    exit
  end

  return data
end

def is_username(word)
  return !/^@[^\s]{3,15}$/.match(word).nil?
end

users=[]
terms=[]
min_word_length=3
show_count=false
ignoreusernames=false
@config_file = "twofi.yml"

begin
  opts.each do |opt, arg|
    case opt
    when "--config"
      @config_file = arg
    when '--count'
      show_count = true
    when '--help'
      usage
    when "--user_file"
      begin
        File.new(arg, 'r').each_line do |line|
          username = 'from:' + line.chomp.sub(/^@/, '')
          terms << username
        end
      rescue
        puts "Unable to read the users file\n"
        exit
      end
    when "--term_file"
      begin
        File.new(arg, 'r').each_line do |line|
          terms << line.chomp
        end
      rescue
        puts "Unable to read the terms file\n"
        exit
      end
    when '--terms'
      arg.split(',').each do |term|
        terms << term
      end
    when '--users'
      arg.split(',').each do |user|
        username = 'from:' + user.chomp.sub(/^@/, '')
        terms << username
      end
    when '--min_word_length'
      min_word_length=arg.to_i
      if min_word_length<1
        usage
      end
    when '--ignore-usernames'
      ignoreusernames=true
    when '--verbose'
      verbose=true
    when '--write'
      outfile=arg
    end
  end
rescue => e
  usage
end

if terms.count == 0
  puts 'You must specify at least one search term or username'
  puts
  usage
end

# Check the config file exits then parse out of it
# the stuff that we need

if File.exists?(@config_file)
  config = YAML.load_file(@config_file)
  if config == false
    sample_config
  end
else
  sample_config
end

@api_key = nil
@api_secret = nil
@bearer_token = nil

if config.include?"options"
  if config["options"].include?"api_key" and config["options"].include?"api_secret"
    @api_key = config["options"]["api_key"]
    @api_secret = config["options"]["api_secret"]
  else
    sample_config
  end

  if @api_key == "<YOUR KEY>"
    sample_config
  end

  if config["options"].include?"bearer_token"
    @bearer_token = config["options"]["bearer_token"]
  else
    @bearer_token = nil
  end
else
  sample_config
end

results = []

terms.each do |term|
  data = twitter_search(term)
  results += data.to_a
end

if results.count == 0
  puts "No search results"
else
  wordlist = {}
  results.each do |result|
    # have to .dup the text as it comes in frozen
    text = result.full_text.dup
    # Strip any non word type characters and substitute accents and other Latin and cirilyc chars
    text.tr!(
    "ÀÁÂÃÄÅàáâãäåĀāĂăĄąÇçĆćĈĉĊċČčÐðĎďĐđÈÉÊËèéêëĒēĔĕĖėĘęĚěĜĝĞğĠġĢģĤĥĦħÌÍÎÏìíîïĨĩĪīĬĭĮįİıĴĵĶķĸĹĺĻļĽľĿŀŁłÑñŃńŅņŇňŉŊŋÒÓÔÕÖØòóôõöøŌōŎŏŐőŔŕŖŗŘřŚśŜŝŞşŠšſŢţŤťŦŧÙÚÛÜùúûüŨũŪūŬŭŮůŰűŲųŴŵÝýÿŶŷŸŹźŻżŽž",
    "AAAAAAaaaaaaAaAaAaCcCcCcCcCcDdDdDdEEEEeeeeEeEeEeEeEeGgGgGgGgHhHhIIIIiiiiIiIiIiIiIiJjKkkLlLlLlLlLlNnNnNnNnnNnOOOOOOooooooOoOoOoRrRrRrSsSsSsSssTtTtTtUUUUuuuuUuUuUuUuUuUuWwYyyYyYZzZzZz"
    )
    text.gsub!(/[^\w \s \d \@]/, ' ')
    text.gsub!("@"," ") if !ignoreusernames
    words = text.split(/\s/)
    words.each do |word|
      #Empty or shorter than required
      if word == '' or word.length < min_word_length or (is_username(word) and ignoreusernames)
        next
      end
      if wordlist.key?(word)
        wordlist[word] += 1
      else
        wordlist[word] = 1
      end
    end
  end

  sorted_wordlist = wordlist.sort_by do |word, count| -count end
  sorted_wordlist.each do |word, count|
    if show_count
      puts word + ', ' + count.to_s
    else
      puts word
    end
  end
end

# Write out the bearer token, this saves making unnecessary
# requests next time
unless @twitter_client.bearer_token.nil?
  config['options']["bearer_token"] = @twitter_client.bearer_token.to_s
  File.open(@config_file,'w') do |h| 
    h.write config.to_yaml
  end
end
