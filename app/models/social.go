package models

import (
	"regexp"
)

var MentionRegex = regexp.MustCompile(`([\#\@])([A-Za-z0-9]+)`)

// Datastore objects

type Like struct {
	LikeId          int
	PostId	        int
	UserId          int
}

type Follower struct {
	FollowerId      int
	UserId	        int
	FollowUserId    int
}

// JSON response objects

type SimpleJSONResponse struct {
  Status      string
  Message     string
}


func FormatContentMentions(content string) string {

	formattedContent := MentionRegex.ReplaceAllStringFunc(content, func(m string) string {
		parts := MentionRegex.FindStringSubmatch(m)

		if parts[1] == "@" {
			// '@' mention
			return "<a href=\"/" + parts[2] + "\">" + m + "</a>"
		}
		// '#' hashtag
		return "<a href=\"/search?hashtag=" + parts[2] + "\">" + m + "</a>"
	})

	return formattedContent

}