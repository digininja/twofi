# twofi, Twitter Words of Interest

Copyright(c) 2021, Robin Wood <robin@digi.ninja>

When attempting to crack passwords custom word lists are very useful additions
to standard dictionaries. An interesting idea originally released on the "7
Habits of Highly Effective Hackers" blog was to use Twitter to help generate
those lists based on searches for keywords related to the list that is being
cracked. I've expanded this idea into twofi which will take multiple search
terms and return a word list sorted by most common first.

The original blog post is at:

http://7habitsofhighlyeffectivehackers.blogspot.com.au/2012/05/using-twitter-to-build-password.html

A second option, suggested by @pentest4dummies, was to look at what specific
users have been saying and use their own tweets to build up the list so I've
added that as well. Given a list of twitter usernames the script will bring back
as many tweets for each user as the API will allow and use those to create the
list.

Installation
============

The only ruby gem that probably isn't installed by default is the twitter one, to
install this run:

bundle install

Then you can run twofi by either using ruby

ruby twofi.rb

or making it executable then running it directly

chmod a+x twofi.rb
./twofi.rb

Version 1 of Twofi used the now removed Twitter search feature which did not
require any authentication. Version 2 now uses the new API which requires you to
have a Twitter account and apply for API keys. The process is simple and
instant, no cash, no waiting for human approval, so no big deal. You need to go
to:

https://apps.twitter.com/

And fill in your details. This will give you a pair of keys which you then need
to put into the twofi.yml config file.


At the moment the script expects the config file to be in the same directory as
twofi is being ran from, if this is not the case you can tell it where the
config file is by using the --config parameter.

Usage
=====

Usage: twofi [OPTIONS]
	--help, -h: show help
	--count, -c: include the count with the words
	--config <file>: config file, default is twofi.yml
	--min_word_length, -m: minimum word length
	--term_file, -T <file>: a file containing a list of terms
	--terms, -t: comma separated search terms
		quote words containing spaces, no space after commas
	--user_file, -U <file>: a file containing a list of users
	--users, -u: comma separated search terms
		quote words containing spaces, no space after commas
	--verbose, -v: verbose

Usage is fairly simple, you can specify search terms or usernames either on the
command line as comma separated lists or through files which you pass in. If you
are specifying the terms or users on the command line you cannot have a space
between the comma and the words, i.e. this is good:

term1,term2,term3

and this is bad:

term1, term2, term3

This is because of the way the command line arguments are parsed, the space
is taken to mean a new parameter.

If you are using files each term/username should be on its own line.

When specifying usernames you do not need the @ symbol, if you pass it it will
be stripped off when used anyway so save yourself some typing.

The --count option allows you to request the number of times each word is used.
This might help if you only have a limited number of attempts to use the words
and so need to decide which are really worth trying.

At the moment there is nothing for the script to be verbose about so the verbose
flag does nothing. I've included it for future versions.

Change Log
==========

2.0-beta - Updated to use the new authenticated API
1.0 - Initial release

Licence
=======
This project released under the Creative Commons Attribution-Share Alike 2.0
UK: England & Wales

( http://creativecommons.org/licenses/by-sa/2.0/uk/ )
