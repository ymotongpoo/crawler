package crawling

import (
	"bytes"
	"reflect"
	"testing"

	"code.google.com/p/go.net/html"
	gq "github.com/PuerkitoBio/goquery"

	"crawler"
)

func TestParseMenu(t *testing.T) {
	in := `<HTML>
<HEAD>
<TITLE>BBS MENU for 2ch</TITLE>
<BASE TARGET="cont">
</HEAD>

<BODY TEXT="#CC3300" BGCOLOR="#FFFFFF" link="#0000FF" alink="#ff0000" vlink="#660099">
<A HREF="http://www.download.co.jp/" TARGET="_blank">
<IMG BORDER=0 WIDTH=75 HEIGHT=75 SRC="http://www2.2ch.net/new2ch.gif"></A>
<BR>
<font size=2>

<A HREF=http://www.2ch.net/>2chの入り口</A><br>
<A HREF=http://www.2ch.net/guide/>2ch総合案内</A>

<BR><BR><B>地震</B><BR>
<A HREF=http://aa2.2ch.net/eq/>臨時地震</A><br>
<A HREF=http://aa2.2ch.net/eqplus/>臨時地震+</A>

<BR><BR><B>おすすめ</B><br>
<A HREF=http://book.2ch.net/bizplus/>ビジネスnews+</A><br>
<A HREF=http://news6.2ch.net/mnewsplus/>芸スポ速報+</A><br>
<A HREF=http://life.2ch.net/sale/>バーゲンセール</A><br>
<A HREF=http://dempa.gozans.com/>電波2ch</A>

<BR><BR><B>特別企画</B><BR>
<A HREF=http://2ch.tora3.net/>2chビューア</A>
`
	reader := bytes.NewReader([]byte(in))

	want := []*crawler.Board{
		&crawler.Board{"臨時地震", "http://aa2.2ch.net/eq/"},
		&crawler.Board{"臨時地震+", "http://aa2.2ch.net/eqplus/"},
		&crawler.Board{"ビジネスnews+", "http://book.2ch.net/bizplus/"},
		&crawler.Board{"芸スポ速報+", "http://news6.2ch.net/mnewsplus/"},
		&crawler.Board{"バーゲンセール", "http://life.2ch.net/sale/"},
		&crawler.Board{"電波2ch", "http://dempa.gozans.com/"},
	}

	root, err := html.Parse(reader)
	if err != nil {
		t.Errorf("wrong HTML %v", err)
	}
	out := ParseMenu(gq.NewDocumentFromNode(root))

	for i, w := range want {
		var found bool
		for _, o := range out {
			if reflect.DeepEqual(w, o) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%v: %v not found", i, w)
			for _, o := range out {
				t.Errorf("\t%v\t%v\t%v\t%v", o.Title, o.URL, []byte(o.Title), []byte(o.URL))
			}
		}
	}
}

func TestParseThreadList(t *testing.T) {
	board := &crawler.Board{"臨時地震", "http://aa2.2ch.net/eq/"}
	in := `<html lang="ja"><head><title>地震があったので臨時＠2ch掲示板＠スレッド一覧</title><meta http-equiv="Content-Type" content="text/html; charset=Shift_JIS"><base href="http://hayabusa.2ch.net/test/read.cgi/eq/" target="body"><script type="text/javascript" src="http://www2.2ch.net/snow/index.js" defer></script><style type="text/css"><!--
a { margin-right: 1em; }div.floated { border: 1px outset honeydew; float: left; height: 20em; line-height: 1em; margin: 0 0 .5em 0; padding: .5em; }div.floated, div.block { background-color: honeydew; }div.floated a, div.block a { display: block; margin-right: 0; text-decoration: none; white-space: nowrap; }div.floated a:hover, div.block a:hover { background-color: cyan; }div.floated a:active, div.block a:active { background-color: gold; }div.right { clear: left; text-align: right; }div.right a { margin-right: 0; }div.right a.js { background-color: dimgray; border: 1px outset dimgray; color: palegreen; text-decoration: none; }
--></style></head><body><div><small id="trad">
<a href="1366513738/l50">2: ひまつぶし2 (276)</a>
<a href="1366440298/l50">3: 世界・海外の地震スレ39～地球は揺れる～ (483)</a>
<a href="1364476302/l50">5: 山梨県民専用県内総合情報15 (519)</a>
<a href="1366460859/l50">7: 緊急警告！北日本・関東・近畿大震災迫る！！ (314)</a>
<a href="1367580117/l50">8: 【KiK-net】強震モニタを見守るスレ649 (182)</a>
</small></div><div class="right"><small><a href="javascript:changeSubbackStyle();" target="_self" class="js">表示スタイル切替</a>&nbsp;
<a href="javascript:switchReadJsMode();" target="_self" class="js">read.cgi モード切替</a>&nbsp;
<a href="../../../eq/kako/"><b>過去ログ倉庫はこちら</b></a></small></div>
</body></html>`

	want := []*crawler.Thread{
		&crawler.Thread{"ひまつぶし2", board,
			`http://hayabusa.2ch.net/test/read.cgi/eq/1366513738/`, 276},
		&crawler.Thread{"世界・海外の地震スレ39～地球は揺れる～", board,
			`http://hayabusa.2ch.net/test/read.cgi/eq/1366440298/`, 483},
		&crawler.Thread{"山梨県民専用県内総合情報15", board,
			`http://hayabusa.2ch.net/test/read.cgi/eq/1364476302/`, 519},
		&crawler.Thread{"緊急警告！北日本・関東・近畿大震災迫る！！", board,
			`http://hayabusa.2ch.net/test/read.cgi/eq/1366460859/`, 314},
		&crawler.Thread{"【KiK-net】強震モニタを見守るスレ649", board,
			`http://hayabusa.2ch.net/test/read.cgi/eq/1367580117/`, 182},
	}

	reader := bytes.NewReader([]byte(in))
	root, err := html.Parse(reader)
	if err != nil {
		t.Errorf("wrong HTML%v", err)
	}
	out := ParseThreadList(board, gq.NewDocumentFromNode(root))
	if len(out) == 0 {
		t.Errorf("no threads")
	}

	for i, w := range want {
		var found bool
		for _, o := range out {
			if reflect.DeepEqual(w, o) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%v: %v, not found", i, w)
			for _, o := range out {
				t.Errorf("\t%v", o)
			}
		}
	}
}
