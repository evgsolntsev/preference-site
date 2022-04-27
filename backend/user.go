package main

import (
	"context"
	"errors"

	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserDatabaseName   = "preference"
	UserCollectionName = "users"
)

type UserDAO struct {
	collection *mgo.Collection
}

func NewUserDAO(session *mgo.Session) *UserDAO {
	return &UserDAO{
		collection: session.DB(UserDatabaseName).C(UserCollectionName),
	}
}

func (d *UserDAO) Insert(ctx context.Context, new *User) error {
	if err := d.collection.Insert(new); err != nil {
		return err
	}

	return nil
}

func (d *UserDAO) FindOneByLogin(ctx context.Context, login string) (*User, error) {
	var result User
	if err := d.collection.FindId(login).One(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *UserDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

type UserManager struct {
	dao *UserDAO
}

func NewUserManager(dao *UserDAO) *UserManager {
	return &UserManager{
		dao: dao,
	}
}

func (m *UserManager) Create(ctx context.Context, login, password string) error {
	_, err := m.dao.FindOneByLogin(ctx, login)
	if err == nil {
		return errors.New("login already exist")
	}

	if err.Error() != "not found" {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &User{
		Login:        login,
		PasswordHash: hash,
	}

	return m.dao.Insert(ctx, newUser)
}

func (m *UserManager) Check(ctx context.Context, login, password string) error {
	u, err := m.dao.FindOneByLogin(ctx, login)
	if err != nil {
		if err.Error() == "not found" {
			return errors.New("login not found")
		}
		return err
	}

	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
}
