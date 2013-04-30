package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
)

const (
	SourceCodec = "Shift_JIS"
	MenuURL     = "http://menu.2ch.net/bbsmenu.html"
)

var (
	UrlRe = regexp.MustCompile(`.*/(\d+)/.*`)
)

// GetBoardList returns 2ch board list
func GetBoardList() ([]*Board, error) {
	r, err := decodeData(MenuURL, SourceCodec)
	if err != nil {
		return nil, err
	}
	doc, err := newDocument(r)
	if err != nil {
		return nil, err
	}
	links := ParseMenu(doc)
	return links, nil
}

// GetThreadList returns 2ch thread list in a board b.
func GetThreadList(b *Board) ([]*Thread, error) {
	url := b.URL + "subback.html"
	r, err := decodeData(url, SourceCodec)
	if err != nil {
		return nil, err
	}
	doc, err := newDocument(r)
	if err != nil {
		return nil, err
	}
	links := ParseThreadList(b, *doc)
	return links, nil
}

// GetThreadData returns a list of thread comments in a thread t of a board b.
func GetThreadData(b *Board, t *Thread) ([]*ThreadData, error) {
	match := UrlRe.FindStringSubmatch(t.URL)
	bg20 := `http://bg20.2ch.net/test/r.so/`
	if len(match) > 0 {
		index := match[1]
		bgurl := fmt.Sprintf("%s%s%s/", bg20, b.URL[7:], index)
		r, err := decodeData(bgurl, SourceCodec)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		if !bytes.Contains(data, []byte("ERROR = 5656")) {
			return ParseThread(b, t, data), nil
		}
	}
	return nil, errors.New("No thread data found.")
}

func Run() {
	fmt.Println("gopher")
}
