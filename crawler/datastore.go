package crawler

import (
	"reflect"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	DefaultMongoServer = "localhost"
	DatabaseName       = "crawler"
)

func insertData(collection string, data interface{}) error {
	session, err := mgo.Dial(DefaultMongoServer)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DatabaseName).C(collection)
	r := reflect.ValueOf(data)
	t := r.Type()
	k := t.Kind()
	switch k {
	case reflect.Slice:
		for i := 0; i < r.Len(); i++ {
			int := r.Index(i).Interface()
			if err = c.Insert(int); err != nil {
				return err
			}
		}
	default:
		if err = c.Insert(r); err != nil {
			return err
		}
	}
	return nil
}

func InsertDat(dats []*ThreadData) error {
	return insertData("dat", dats)
}

func InsertBoards(boards []*Board) error {
	return insertData("board", boards)
}

func InsertThread(thread *Thread) error {
	return insertData("thread", thread)
}

func GetBoards() ([]*Board, error) {
	session, err := mgo.Dial(DefaultMongoServer)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	c := session.DB(DatabaseName).C("board")
	result := []*Board{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
