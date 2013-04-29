package crawler

import (
	_ "errors"
	"fmt"
	"regexp"

	_ "github.com/PuerkitoBio/goquery"
)

const (
	SourceCodec = "Shift_JIS"
	MenuURL     = "http://menu.2ch.net/bbsmenu.html"
)

var (
	UrlRe = regexp.MustCompile(".*/([0-9]+)/.*")
)

func GetBoardList() ([]*BoardLink, error) {
	doc, err := newDocument(MenuURL, SourceCodec)
	if err != nil {
		return nil, err
	}
	links := ParseMenu(*doc)
	return links, nil
}

func Run() {
	fmt.Println("gopher")
}
