package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

const (
	SourceCodec = "Shift_JIS"
	MenuURL     = "http://menu.2ch.net/bbsmenu.html"
	ChannelSize = 1000
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
func GetThreadData(t *Thread) ([]*ThreadData, error) {
	match := UrlRe.FindStringSubmatch(t.URL)
	bg20 := `http://bg20.2ch.net/test/r.so/`
	if len(match) > 0 {
		index := match[1]
		bgurl := fmt.Sprintf("%s%s%s/", bg20, t.Board.URL[7:], index)
		r, err := decodeData(bgurl, SourceCodec)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		if !bytes.Contains(data, []byte("ERROR = 5656")) {
			return ParseThread(t, data), nil
		}
	}
	return nil, errors.New("No thread data found.")
}

func CrawlThread(threads <-chan *Thread) {
	for t := range threads {
		dats, err := GetThreadData(t)
		if err != nil {
			log.Printf(t.Board.URL)
		}
		if len(dats) > 0 {
			var old_count int
			// TODO(ymotongpoo): Change interface to return (int, error)
			err := InsertThread(t)
			if err != nil {
				log.Printf("Failed to store thread %v\n", t.Title)
				continue
			}
			dats = dats[old_count:]
			InsertDat(dats)
		} else {
			log.Printf("bg20 is dead. %v", t.Title)
		}
	}
}

func CrawlBoard(boards <-chan *Board) {
	for b := range boards {
		log.Println(b.Title)
		threads, err := GetThreadList(b)
		if err != nil {
			log.Printf("Error on fetching thread list in %v", b.URL)
		}
		if len(threads) > 0 {
			tasks := make(chan *Thread, ChannelSize)
			go func() {
				for _, t := range threads {
					tasks <- t
				}
				close(tasks)
			}()
			ExecCrawl(tasks, 8)
		}
	}
}

func ExecCrawl(threads <-chan *Thread, maxWorkers int) {
	for w := 0; w < maxWorkers; w++ {
		go CrawlThread(threads)
	}

}

func Run(maxWorkers int) {
	boards, err := GetBoardList()
	if err != nil {
		log.Fatalf("Error on fetching board list")
	}
	InsertBoards(boards)
	// boards = GetBoards()
	tasks := make(chan *Board, ChannelSize)
	go func() {
		for _, b := range boards {
			tasks <- b
		}
		close(tasks)
	}()
	for w := 0; w < maxWorkers; w++ {
		go CrawlBoard(tasks)
	}
}
