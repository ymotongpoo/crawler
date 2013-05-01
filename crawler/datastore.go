package crawler

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	DefaultServer = "localhost"
)

func InsertDat(dats []*ThreadData) error {
	session, err := mgo.Dial(DefaultServer)
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("crawler").C("dat")
	for _, d := range dats {
		err = c.Insert(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertBoards(boards []*Board) error {
	session, err := mgo.Dial(DefaultServer)
	if err != nil {
		return err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("crawler").C("board")
	for _, b := range boards {
		err = c.Insert(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertThread(thread *Thread) (int, error) {
	session, err := mgo.Dial(DefaultServer)
	if err != nil {
		return 0, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("crawler").C("thread")
	err = c.Insert(thread)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func GetBoards() ([]*Board, error) {
	session, err := mgo.Dial(DefaultServer)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	c := session.DB("crawler").C("board")
	result := []*Board{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
