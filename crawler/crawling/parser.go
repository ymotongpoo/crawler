package crawling

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/mahonia"
	gq "github.com/PuerkitoBio/goquery"

	"crawler"
)

var (
	CgiRe   = regexp.MustCompile(`(.*read.cgi/\w+/\d+/).*`)
	ResRe   = regexp.MustCompile(`.*\((\d+)\).*`)
	TitleRe = regexp.MustCompile(`\d+:\s+(.*)\((\d+)\)`)
)

func decodeData(url, srcCodec string) (*mahonia.Reader, error) {
	// Load URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// Convert character settings
	decode := mahonia.NewDecoder(srcCodec)
	if decode == nil {
		return nil, errors.New("Cannot create decoder.")
	}
	return decode.NewReader(resp.Body), nil
}

// newDocument loads a page on url in srcCodec and returns
// goquery Document.
func newDocument(r *mahonia.Reader) (*gq.Document, error) {
	root, err := html.Parse(r)
	if err != nil {
		return &gq.Document{}, err
	}
	doc := gq.NewDocumentFromNode(root)
	return doc, nil
}

// ParseMenu parses 2ch BBS menu page and returns a slice of board list.
func ParseMenu(d *gq.Document) (result []*crawler.Board) {
	d.Find("a").Each(func(_ int, s *gq.Selection) {
		href, exist := s.Attr("href")
		if !exist {
			return
		}
		if strings.HasPrefix(href, "http") && !strings.HasSuffix(href, "bbsmenu") {
			if !strings.HasSuffix(href, "php") &&
				!strings.HasSuffix(href, ".net/") &&
				!strings.HasSuffix(href, ".jp/") {
				result = append(result, &crawler.Board{
					Title: s.Text(),
					URL:   href,
				})
			}
		}
	})
	return result
}

// Parsecrawler.ThreadList returns a thread list in a board b.
func ParseThreadList(b *crawler.Board, d *gq.Document) (result []*crawler.Thread) {
	var base_url string
	d.Find("base").Each(func(_ int, s *gq.Selection) {
		var exist bool
		base_url, exist = s.Attr("href")
		if !exist {
			return
		}
	})
	if len(base_url) == 0 {
		return nil
	}

	d.Find("a").Each(func(_ int, s *gq.Selection) {
		v := s.Text()
		href, exist := s.Attr("href")
		if !exist {
			return
		}
		url := base_url + href
		if strings.HasSuffix(url, "50") {
			match := TitleRe.FindStringSubmatch(v)
			if len(match) > 0 {
				title := strings.TrimSpace(match[1])
				res_cnt, err := strconv.Atoi(match[2])
				if err != nil {
					return
				}
				result = append(result, &crawler.Thread{
					Title:    title,
					Board:    b,
					URL:      url[:len(url)-3],
					ResCount: res_cnt,
				})
			}
		}
	})
	return result
}

// Parsecrawler.Thread loads all dat file data into a slice of crawler.ThreadData.
func ParseThread(t *crawler.Thread, r []byte) (result []*crawler.ThreadData) {
	buffer := bytes.NewBuffer(r)
	i := 1
	for {
		line, err := buffer.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil
		}
		cols := bytes.Split(line, []byte("<>"))
		if len(cols) > 4 {
			td := &crawler.ThreadData{
				Board:   t.Board,
				Thread:  t,
				Handle:  string(cols[0]),
				MailTo:  string(cols[1]),
				Date:    string(cols[2]),
				Comment: string(cols[3]),
				Other:   string(cols[4]),
				No:      i,
			}
			result = append(result, td)
			i++
		}
	}
	return result
}
