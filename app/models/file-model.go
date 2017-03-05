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

func (fm *FileModel) IsThereUserFile(fileName string, gId, uId bson.ObjectId) bool {
	result := struct {
		Id bson.ObjectId `bson:"_id"`
	}{}

	fm.db.C("groups").Find(bson.M{
		"_id":         gId,
		"files.owner": uId,
		"files.name":  fileName,
	}).Select(bson.M{"_id": 1}).One(&result)

	return result.Id != ""
}

func (fm *FileModel) AddFile(file *File, gId bson.ObjectId) error {
	return fm.db.C("groups").Update(
		bson.M{"_id": gId},
		bson.M{
			"$push": bson.M{"files": file},
		},
	)
}
