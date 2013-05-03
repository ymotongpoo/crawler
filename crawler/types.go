package crawler

type Board struct {
	Title string
	URL   string
}

type Thread struct {
	Title    string
	Board    *Board
	URL      string
	ResCount int
}

type ThreadData struct {
	Board   *Board
	Thread  *Thread
	URL     string
	Handle  string
	MailTo  string
	Date    string
	Comment string
	Other   string
	No      int
}
