package feed

import (
	"encoding/xml"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Title         string       `xml:"title"`
	Link          string       `xml:"link"`
	Description   string       `xml:"description"`
	Language      string       `xml:"language"`
	PubDate       string       `xml:"pubDate"`
	LastBuildDate string       `xml:"lastBuildDate"`
	ITunesImage   *ITunesImage `xml:"itunes:image,omitempty"`
	Items         []*Item      `xml:"item"`
}

// Episode describes a single episode in the feed. The exported fields will be referenced in the actual feed,
type Episode struct {
	ContentFile      string
	AudioFile        string
	content          string
	Summary          string
	summaryFile      string
	IllustrationFile string
	Created          time.Time
	Updated          time.Time
}

type Link struct {
	Href, Rel, Type, Length string
}
type Author struct {
	Name, Email string
}
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type Item struct {
	Title       string       `xml:"title"`
	Link        string       `xml:"link"`
	Description string       `xml:"description"`
	PubDate     string       `xml:"pubDate"`
	Guid        string       `xml:"guid"`
	Enclosure   *Enclosure   `xml:"enclosure,omitempty"`
	ITunesImage *ITunesImage `xml:"itunes:image,omitempty"`
}

type Feed struct {
	Title       string
	Link        *Link
	Description string
	Author      *Author
	Updated     time.Time
	Created     time.Time
	Id          string
	Subtitle    string
	Items       []*Item
	Copyright   string
	Image       *Image
}

type Image struct {
	Url, Title, Link string
	Width, Height    int
}

type ITunesImage struct {
	Href string `xml:"href,attr"`
}
