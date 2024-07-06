// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    zennSearchResult, err := UnmarshalZennSearchResult(bytes)
//    bytes, err = zennSearchResult.Marshal()

package main

import (
	"encoding/json"
	"time"
)

func UnmarshalZennSearchResult(data []byte) (ZennSearchResult, error) {
	var r ZennSearchResult
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ZennSearchResult) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ViewerResult struct {
	Title string
	URL   string
	Date  time.Time
}

type ZennSearchResult struct {
	Articles []ArticleElement `json:"articles"`
	NextPage int64            `json:"next_page"`
}

type ArticleElement struct {
	ID                  int64        `json:"id"`
	PostType            PostType     `json:"post_type"`
	Title               string       `json:"title"`
	Slug                string       `json:"slug"`
	CommentsCount       int64        `json:"comments_count"`
	LikedCount          int64        `json:"liked_count"`
	BodyLettersCount    int64        `json:"body_letters_count"`
	ArticleType         ArticleType  `json:"article_type"`
	Emoji               string       `json:"emoji"`
	IsSuspendingPrivate bool         `json:"is_suspending_private"`
	PublishedAt         time.Time    `json:"published_at"`
	BodyUpdatedAt       time.Time    `json:"body_updated_at"`
	SourceRepoUpdatedAt *time.Time   `json:"source_repo_updated_at"`
	Pinned              bool         `json:"pinned"`
	Path                string       `json:"path"`
	User                User         `json:"user"`
	Publication         *Publication `json:"publication"`
}

type Publication struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	DisplayName      string `json:"display_name"`
	AvatarSmallURL   string `json:"avatar_small_url"`
	Pro              bool   `json:"pro"`
	AvatarRegistered bool   `json:"avatar_registered"`
}

type User struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarSmallURL string `json:"avatar_small_url"`
}

type ArticleType string

const (
	Idea ArticleType = "idea"
	Tech ArticleType = "tech"
)

type PostType string

const (
	Article PostType = "Article"
)
