package soundcloud_test

import (
	"testing"

	"github.com/AYehia0/soundcloud-dl/pkg/soundcloud"
)

func TestIsValidUrl(t *testing.T) {
	validUrls := []string{
		"https://soundcloud.com/something/2092-jfl-2jlafj",
		"https://soundcloud.com/ahmad9kamal/ftrod7evucus?si=dcb16d100afe402f89ca4b196c1dc756&utm_source=clipboard&utm_medium=text&utm_campaign=social_sharing",
		"https://soundcloud.com/mahmoud-said-451898307/sets/rtutgnky5you?si=57e5fb4b81e148338e6d90f57e6eebad",
		"https://soundcloud.com/search/sets?q=%D8%A7%D9%84%D8%B4%D9%8A%D8%AE+%D8%AD%D8%B3%D9%86+%D8%B5%D8%A7%D9%84%D8%AD",
	}

	invalidUrls := []string{
		"https",
		"",
		"./what/the/heck/",
		"https/soundcloud/hello/ljas020-33",
		"google.com/me",
	}

	for _, url := range validUrls {
		if !soundcloud.IsValidUrl(url) {
			t.Errorf("Expected %v to be a valid Url!", url)
		}
	}

	for _, url := range invalidUrls {
		if soundcloud.IsValidUrl(url) {
			t.Errorf("Expected %v to be a invalid Url!", url)
		}
	}

}
