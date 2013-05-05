package crawler

import (
	"reflect"
	"strconv"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"crawler"
)

var (
	DefaultMongoHost = "localhost"
	DefaultMongoPort = 27017
	DatabaseName     = "crawler"
)

const (
	MainDatabase         = "main"
	DatDatabase          = "dat"
	BoardListCollection  = "board_list"
	ThreadListCollection = "thread_list"
	ThreadCollection     = "thread"
)

// GetConnection returns a pointer of mgo.Session based on host and port.
func GetConnection(host string, port int) (*mgo.Session, error) {
	if host == "" {
		host = DefaultMongoHost
	}
	if port == 0 {
		port = DefaultMongoPort
	}
	portStr := strconv.Itoa(port)
	return mgo.Dial("mongodb://" + host + ":" + portStr)
}

// insertData inserts data in collection
func insertData(session *mgo.Session, database, collection string, data interface{}) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(database).C(collection)
	r := reflect.ValueOf(data)
	t := r.Type()
	k := t.Kind()
	switch k {
	case reflect.Slice:
		for i := 0; i < r.Len(); i++ {
			itf := r.Index(i).Interface()
			if err = c.Insert(itf); err != nil {
				return err
			}
		}
	default:
		if err = c.Insert(r.Interface()); err != nil {
			return err
		}
	}
	return nil
}

func getData(session *mgo.Session, database, collection string) (interface{}, error) {
	return nil, nil
}

func InsertDat(session *mgo.Session, dats []*crawler.ThreadData) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
	}
	return insertData(session, MainDatabase, ThreadCollection, dats)
}

func InsertBoards(session *mgo.Session, boards []*crawler.Board) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
	}
	return insertData(session, MainDatabase, BoardListCollection, boards)
}

func InsertThread(session *mgo.Session, thread *crawler.Thread) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
	}
	return insertData(session, MainDatabase, ThreadListCollection, thread)
}

func GetBoards(session *mgo.Session) ([]*crawler.Board, error) {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return nil, err
		}
	}
	defer session.Close()
	c := session.DB(DatabaseName).C("board")
	result := []*crawler.Board{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
