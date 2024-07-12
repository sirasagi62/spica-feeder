// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    rSSFeed, err := UnmarshalRSSFeed(bytes)
//    bytes, err = rSSFeed.Marshal()

package main

import "github.com/pelletier/go-toml/v2"

func UnmarshalRSSFeed(data []byte) (RSSFeed, error) {
	var r RSSFeed
	err := toml.Unmarshal(data, &r)
	return r, err
}

func (r *RSSFeed) Marshal() ([]byte, error) {
	return toml.Marshal(r)
}

type RSSFeed struct {
	Src []Src `json:"src"`
}

type Src struct {
	Main  *string `json:"main,omitempty"`
	Topic *Topic  `json:"topic,omitempty"`
	User  *Topic  `json:"user,omitempty"`
}

type Topic struct {
	URL       string   `json:"url"`
	Following []string `json:"following"`
}
