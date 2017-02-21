package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	UserModel struct {
		db *mgo.Database
	}

	User struct {
		Id       bson.ObjectId `bson:"_id"`
		Name     string        `bson:"name"`
		Email    string        `bson:"email"`
		Password []byte        `bson:"password"`
	}
)

func NewUserModel(s *mgo.Session) *UserModel {
	return &UserModel{s.DB("repl")}
}

func (um *UserModel) IsThere(email string) bool {
	result := new(User)

	um.db.C("users").Find(bson.M{"email": email}).One(result)

	return result.Id != ""
}

func (um *UserModel) AddUser(user *User) error {
	return um.db.C("users").Insert(*user)
}
