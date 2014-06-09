package models

import (
	"regexp"
)

var MentionRegex = regexp.MustCompile(`([\#\@])([A-Za-z0-9]+)`)

// Datastore objects

type Like struct {
	LikeId int
	PostId int
	UserId int
}

type Follower struct {
	FollowerId   int
	UserId       int
	FollowUserId int
}

// JSON response objects

type SimpleJSONResponse struct {
	Status  string
	Message string
}

func FormatContentMentions(content []byte) []byte {

	formattedContent := MentionRegex.ReplaceAllFunc(content, func(m []byte) []byte {
		match := string(m)
		parts := MentionRegex.FindSubmatch(m)

		mentionType := string(parts[1])
		mentionValue := string(parts[2])

		if mentionType == "@" {
			// '@' mention
			return []byte("<a href=\"/" + mentionValue + "\">" + match + "</a>")
		}
		// '#' hashtag
		return []byte("<a href=\"/search?hashtag=" + mentionValue + "\">" + match + "</a>")
	})

	return formattedContent

}
