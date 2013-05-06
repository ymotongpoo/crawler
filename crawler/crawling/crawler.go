package crawling

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"runtime"
	"time"

	"crawler"
	datastore "crawler/mongo"
)

const (
	SourceCodec = "Shift_JIS"
	MenuURL     = "http://menu.2ch.net/bbsmenu.html"
	ChannelSize = 1000
)

var (
	UrlRe         = regexp.MustCompile(`.*/(\d+)/.*`)
	NumCPU        = runtime.NumCPU()
	CrawlWaitTime = 15 * time.Second
)

// GetBoardList returns 2ch board list
func GetBoardList() ([]*crawler.Board, error) {
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
func GetThreadList(b *crawler.Board) ([]*crawler.Thread, error) {
	url := b.URL + "subback.html"
	r, err := decodeData(url, SourceCodec)
	if err != nil {
		return nil, err
	}
	doc, err := newDocument(r)
	if err != nil {
		return nil, err
	}
	links := ParseThreadList(b, doc)
	return links, nil
}

// GetThreadData returns a list of thread comments in a thread t of a board b.
func GetThreadData(t *crawler.Thread) ([]*crawler.ThreadData, error) {
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

// CrawlThread go through all threads in channel 'threads' and
// store thread data into datastore.
func CrawlThread(threads <-chan *crawler.Thread) {
	for t := range threads {
		dats, err := GetThreadData(t)
		if err != nil {
			log.Printf("%v: %v", t.URL, err)
		}
		if len(dats) > 0 {
			r, err := datastore.FindThread(nil, t)
			if err != nil && err != datastore.ErrNotFound {
				log.Printf("Failed to find thread: %v", t.Title)
				continue
			}
			var old_count int
			if r != nil {
				old_count = r.ResCount
				if t.ResCount > old_count {
					err := datastore.UpdateThread(nil, r, t)
					if err != nil {
						log.Printf("Failed to update thread: %v", t.Title)
						continue
					}
					dats = dats[old_count:]
					datastore.InsertDat(nil, dats)
				}
			} else {
				err := datastore.InsertThread(nil, t)
				if err != nil {
					log.Printf("Failed to store thread %v\n", t.Title)
					continue
				}
			}
		} else {
			log.Printf("bg20 is dead. %v", t.Title)
		}
		time.Sleep(CrawlWaitTime)
	}
}

// CrawlBoard go through all boardss in channel 'boards' and
// store board data into datastore.
func CrawlBoard(boards <-chan *crawler.Board) {
	for b := range boards {
		log.Println(b.Title)
		threads, err := GetThreadList(b)
		if err != nil {
			log.Printf("Error on fetching thread list in %v", b.URL)
		}
		if len(threads) > 0 {
			tasks := make(chan *crawler.Thread, ChannelSize)
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

// ExecCrawl runs workers for CrawlThread.
func ExecCrawl(threads <-chan *crawler.Thread, maxWorkers int) {
	for w := 0; w < maxWorkers; w++ {
		go CrawlThread(threads)
	}
}

// Run launches workers for CrawlBoard.
func Run(maxWorkers int) {
	runtime.GOMAXPROCS(NumCPU)
	boards, err := GetBoardList()
	if err != nil {
		log.Fatalf("Error on fetching board list")
	}
	datastore.InsertBoards(nil, boards)
	// boards = datastore.GetBoards()
	tasks := make(chan *crawler.Board, ChannelSize)
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
