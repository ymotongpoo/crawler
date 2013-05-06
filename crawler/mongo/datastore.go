package crawler

import (
	"reflect"
	"strconv"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"crawler"
)

// MongoDB settings
var (
	// Configuration
	DefaultMongoHost = "localhost"
	DefaultMongoPort = 27017

	// Database names
	MainDatabase         = "main"
	DatDatabase          = "dat"
	BoardListCollection  = "board_list"
	ThreadListCollection = "thread_list"
	ThreadCollection     = "thread"
)

var ErrNotFound = mgo.ErrNotFound

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

// insertData inserts data in collection of database.
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

// getData returns collection in database.
func getCollection(session *mgo.Session, database, collection string) (*mgo.Collection, error) {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return nil, err
		}
	}
	return session.DB(database).C(collection), nil
}

// InsertDat inserts thread's dat data into collection for thread.
func InsertDat(session *mgo.Session, dats []*crawler.ThreadData) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
		defer session.Close()
	}
	return insertData(session, DatDatabase, ThreadCollection, dats)
}

// InsertBoards inserts board list into collection for board.
func InsertBoards(session *mgo.Session, boards []*crawler.Board) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
		defer session.Close()
	}
	return insertData(session, MainDatabase, BoardListCollection, boards)
}

// InsertThread inserts thread list into collection for threadlist.
func InsertThread(session *mgo.Session, thread *crawler.Thread) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
		defer session.Close()
	}
	return insertData(session, MainDatabase, ThreadListCollection, thread)
}

// UpdateThread updates target with thread.
func UpdateThread(session *mgo.Session, target, thread *crawler.Thread) error {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return err
		}
		defer session.Close()
	}
	c, err := GetThreadCollection(session)
	if err != nil {
		return err
	}
	query := bson.M{"URL": target.URL}

	err = c.Update(query, bson.M{"$set": thread})
	if err != nil {
		return err
	}
	return nil
}

// GetThreadListCollection returns collection for thread list.
func GetThreadListCollection(session *mgo.Session) (*mgo.Collection, error) {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return nil, err
		}
	}
	return getCollection(session, MainDatabase, ThreadListCollection)
}

// GetThreadListCollection returns collection for board list.
func GetBoardCollection(session *mgo.Session) (*mgo.Collection, error) {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		if err != nil {
			return nil, err
		}
	}
	return getCollection(session, MainDatabase, BoardListCollection)
}

// GetThreadListCollection returns collection for a thread.
func GetThreadCollection(session *mgo.Session) (*mgo.Collection, error) {
	var err error
	if session == nil {
		session, err = GetConnection(DefaultMongoHost, DefaultMongoPort)
		return nil, err
	}
	return getCollection(session, DatDatabase, ThreadCollection)
}

// GetBoards returns a list of boards in collection.
func GetBoards(session *mgo.Session) ([]*crawler.Board, error) {
	c, err := GetBoardCollection(session)
	if err != nil {
		return nil, err
	}
	var result []*crawler.Board
	err = c.Find(nil).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindThread returns thread's rescount.
func FindThread(session *mgo.Session, thread *crawler.Thread) (*crawler.Thread, error) {
	c, err := GetThreadListCollection(session)
	if err != nil {
		return nil, err
	}
	query := bson.M{"URL": thread.URL}
	var result *crawler.Thread
	err = c.Find(query).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}
	if result == nil {
		return nil, nil
	}
	return result, nil
}

// FindThreadsByName returns a slice of crawler.Thread. name can be PCRE.
func FindThreadsByName(session *mgo.Session, name string) ([]*crawler.Thread, error) {
	c, err := GetThreadListCollection(session)
	if err != nil {
		return nil, err
	}
	var result []*crawler.Thread
	query := bson.M{"Title": bson.M{"$regex": name}}
	err = c.Find(query).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
