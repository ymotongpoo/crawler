package crawler

import (
	"errors"
	"net/http"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/mahonia"
	gq "github.com/PuerkitoBio/goquery"
)

type BoardLink struct {
	Title string
	URL   string
}

func newDocument(url, srcCodec string) (*gq.Document, error) {
	// Load URL
	resp, err := http.Get(url)
	if err != nil {
		return &gq.Document{}, err
	}
	defer resp.Body.Close()

	// Convert character settings
	decode := mahonia.NewDecoder(srcCodec)
	if decode == nil {
		return &gq.Document{}, errors.New("Cannot create decoder.")
	}
	r := decode.NewReader(resp.Body)

	root, err := html.Parse(r)
	if err != nil {
		return &gq.Document{}, err
	}
	doc := gq.NewDocumentFromNode(root)
	return doc, nil
}

func ParseMenu(d gq.Document) (result []*BoardLink) {
	d.Find("a").Each(func(_ int, s *gq.Selection) {
		href, exist := s.Attr("href")
		if !exist {
			return
		}
		if strings.HasPrefix(href, "http") && !strings.HasSuffix(href, "bbsmenu") {
			if !strings.HasSuffix(href, "php") &&
				!strings.HasSuffix(href, ".net/") &&
				!strings.HasSuffix(href, ".jp/") {
				result = append(result, &BoardLink{
					Title: s.Text(),
					URL:   href,
				})
			}
		}
	})
	return result
}
