/*
Check if the URL passed is a valid URL by matching a regex.
*/
package soundcloud

import (
	"log"
	"regexp"
)

// check if the url is a soundcloud url
func IsValidUrl(url string) bool {

	/*
		   ^ - start of string
		   (?:https?://)? - an optional http:// or https://
		   (?:[^/\s]+\.)* - zero or more repetitions of :
			   [^/.\s]+ - one or more chars other than /, . and whitespace
			   \. - a dot
		   google\.com - an escaped keyword
		   (?:/[^/\s]+)* - zero or more repetitions of a / and then one or more chars other than / and whitespace chars
		   /? - an optional /
		   $ - end of string
	*/
	pattern := `^(?:https?://)?(?:[^/.\s]+\.)*soundcloud\.com(?:/[^/\s]+)*/?$`
	matched, err := regexp.MatchString(pattern, url)
	if err != nil {
		log.Fatalf("Something went wrong while parsing the URL : %s", err)
	}

	if matched {
		return true
	}
	return false
}
