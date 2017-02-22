package models

import (
	"golang.org/x/crypto/bcrypt"
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
		Role     string        `bson:"role"`
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

func (um *UserModel) FindOne(email, pwd string) (result *User, err error) {
	result = new(User)

	um.db.C("users").Find(bson.M{"email": email}).One(result)
	if result.Id == "" {
		return
	}

	err = bcrypt.CompareHashAndPassword(result.Password, []byte(pwd))
	return
}
