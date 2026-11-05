package news

import (
	"bytes"
	"encoding/xml"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io"
	"net/url"
	"strings"
	"time"
	"unicode"
)

type Source struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Region   string `json:"region"`
	Category string `json:"category"`
	Custom   bool   `json:"custom"`
	Updated  string `json:"updated"`
	Error    string `json:"error"`
}
type Story struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Source      string `json:"source"`
	SourceID    string `json:"sourceId"`
	Region      string `json:"region"`
	Category    string `json:"category"`
	Published   string `json:"published"`
	Summary     string `json:"summary"`
	Text        string `json:"text"`
	ContentKind string `json:"contentKind"`
	Image       string `json:"image"`
	SavedAt     string `json:"savedAt,omitempty"`
}
type item struct {
	Title string `xml:"title"`
	Links []struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
	} `xml:"link"`
	Description string `xml:"description"`
	Summary     string `xml:"summary"`
	Encoded     string `xml:"encoded"`
	Content     string `xml:"content"`
	PubDate     string `xml:"pubDate"`
	Published   string `xml:"published"`
	Updated     string `xml:"updated"`
	Thumb       struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
}

// Feed/article markup is rendered exclusively as text, never as trusted HTML.
func plain(raw string) string {
	z := html.NewTokenizer(strings.NewReader(raw))
	var b strings.Builder
	skip := 0
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		switch tt {
		case html.StartTagToken, html.EndTagToken:
			t := z.Token()
			switch t.Data {
			case "script", "style", "noscript", "iframe", "svg":
				if tt == html.StartTagToken {
					skip++
				} else if skip > 0 {
					skip--
				}
			case "p", "div", "br", "li", "h1", "h2", "h3":
				if skip == 0 {
					b.WriteString("\n\n")
				}
			}
		case html.TextToken:
			if skip == 0 {
				b.Write(z.Text())
			}
		}
	}
	lines := strings.Split(b.String(), "\n")
	out := []string{}
	for _, l := range lines {
		l = strings.Join(strings.FieldsFunc(l, unicode.IsSpace), " ")
		if l != "" {
			out = append(out, l)
		}
	}
	s := strings.Join(out, "\n\n")
	r := []rune(s)
	if len(r) > 100000 {
		s = string(r[:100000])
	}
	return s
}
func safeLink(raw, base string) string {
	u, e := url.Parse(strings.TrimSpace(raw))
	if e != nil {
		return ""
	}
	b, e := url.Parse(base)
	if e != nil {
		return ""
	}
	u = b.ResolveReference(u)
	if _, e = validURL(u.String()); e != nil {
		return ""
	}
	u.Fragment = ""
	q := u.Query()
	for k := range q {
		if strings.HasPrefix(k, "utm_") || k == "fbclid" {
			q.Del(k)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}
func parseFeed(data []byte, s Source) ([]Story, error) {
	d := xml.NewDecoder(bytes.NewReader(data))
	d.CharsetReader = charset.NewReaderLabel
	out := []Story{}
	seen := map[string]bool{}
	root := ""
	for {
		tok, e := d.Token()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, errors.New("This URL did not return a valid RSS or Atom feed")
		}
		st, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if root == "" {
			root = st.Name.Local
			if root != "rss" && root != "feed" && root != "RDF" {
				return nil, errors.New("This URL is a web page, not an RSS or Atom feed")
			}
		}
		if st.Name.Local != "item" && st.Name.Local != "entry" {
			continue
		}
		var it item
		if e = d.DecodeElement(&it, &st); e != nil {
			return nil, e
		}
		link := ""
		for _, l := range it.Links {
			if l.Text != "" {
				link = l.Text
			}
			if l.Href != "" && (l.Rel == "alternate" || l.Rel == "") {
				link = l.Href
				break
			}
		}
		link = safeLink(link, s.URL)
		title := plain(it.Title)
		if link == "" || title == "" || seen[link] {
			continue
		}
		seen[link] = true
		summary := plain(it.Description)
		if summary == "" {
			summary = plain(it.Summary)
		}
		body := plain(it.Encoded)
		if body == "" {
			body = plain(it.Content)
		}
		kind := "Feed excerpt"
		if body != "" {
			kind = "Feed content"
		} else {
			body = summary
		}
		if summary == "" {
			summary = body
		}
		if r := []rune(summary); len(r) > 350 {
			summary = string(r[:350]) + "…"
		}
		published := ""
		for _, v := range []string{it.PubDate, it.Published, it.Updated} {
			for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC822Z, time.RFC822, "2006-01-02T15:04:05-0700"} {
				if t, e := time.Parse(layout, strings.TrimSpace(v)); e == nil {
					published = t.UTC().Format(time.RFC3339)
					break
				}
			}
			if published != "" {
				break
			}
		}
		out = append(out, Story{ID: apphost.Hash(link), Title: title, URL: link, Source: s.Name, SourceID: s.ID, Region: s.Region, Category: s.Category, Published: published, Summary: summary, Text: body, ContentKind: kind, Image: safeLink(it.Thumb.URL, s.URL)})
		if len(out) >= 100 {
			break
		}
	}
	if root == "" {
		return nil, errors.New("Empty feed")
	}
	return out, nil
}
func articleText(data []byte) string {
	doc, e := html.Parse(bytes.NewReader(data))
	if e != nil {
		return ""
	}
	var best string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "article" || n.Data == "main") {
			var b bytes.Buffer
			_ = html.Render(&b, n)
			p := plain(b.String())
			if len(p) > len(best) {
				best = p
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return best
}
