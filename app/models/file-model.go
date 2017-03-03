package models

import (
	//"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	FileModel struct {
		db *mgo.Database
	}

	File struct {
		Id        bson.ObjectId `bson:"_id"`
		Name      string        `bson:"name"`
		Original  bson.ObjectId `bson:"original"`
		Owner     bson.ObjectId `bson:"owner"`
		Extension string        `bson:"extension"`
		Content   []byte        `bson:"content"`
		IsPrivate bool          `bson:"isPrivate"`
	}
)

func NewFileModel(s *mgo.Session) *FileModel {
	return &FileModel{s.DB("repl")}
}
