package crawler

import (
	"testing"
)

func TestGetBoardList(t *testing.T) {
	list, err := GetBoardList()
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(list) == 0 {
		t.Errorf("no boards.")
	}
	for _, l := range list {
		t.Logf("%v -> %v", l.Title, l.URL)
	}
}
