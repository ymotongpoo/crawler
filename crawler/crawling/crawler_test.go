package crawling

import (
	"testing"

	"crawler"
)

func TestGetBoardList(t *testing.T) {
	list, err := GetBoardList()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(list) == 0 {
		t.Errorf("no boards.")
	}
	for _, b := range list {
		t.Logf("%v -> %v", b.Title, b.URL)
	}
}

func TestGetThreadList(t *testing.T) {
	in := []*crawler.Board{
		&crawler.Board{
			Title: "一人暮らし",
			URL:   "http://uni.2ch.net/homealone/",
		},
		&crawler.Board{
			Title: "PC初心者",
			URL:   "http://kohada.2ch.net/pcqa/",
		},
	}
	for _, b := range in {
		list, err := GetThreadList(b)
		if err != nil {
			t.Errorf("%v\n", err)
		}
		if len(list) == 0 {
			t.Errorf("%v -> no threads.", b.Title)
		}
		for _, l := range list {
			t.Logf("%v -> %v, %v, %v", b.Title, l.Title, l.URL, l.ResCount)
		}
	}
}

func TestGetThreadData(t *testing.T) {
	board := &crawler.Board{
		Title: "一人暮らし",
		URL:   "http://uni.2ch.net/homealone/",
	}
	thread := &crawler.Thread{
		Title:    "とりあえず部屋にあるもの挙げてみろ",
		Board:    board,
		URL:      "http://uni.2ch.net/test/read.cgi/homealone/1329046861/",
		ResCount: 67,
	}
	list, err := GetThreadData(thread)
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(list) == 0 {
		t.Errorf("%v -> no thread data.")
	}
	for _, l := range list {
		t.Logf("%v : %v\n", l.No, l.Comment)
	}
}
