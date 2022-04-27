package main

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite

	Ctx     context.Context
	DAO     *UserDAO
	Manager *UserManager
}

func TestUser(t *testing.T) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	require.NoError(t, err)

	ctx := context.Background()
	dao := NewUserDAO(session)
	defer dao.RemoveAll(ctx)

	manager := NewUserManager(dao)

	suite.Run(t, &UserSuite{
		Ctx:     ctx,
		DAO:     dao,
		Manager: manager,
	})
}

func (s *UserSuite) TestCheckUnexisting() {
	err := s.Manager.Check(s.Ctx, "unexisting", "pass")
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "login not found")
}

func (s *UserSuite) TestCheckWrongPassword() {
	require.NoError(s.T(), s.Manager.Create(s.Ctx, "user1", "pass", "some@mail.com"))
	err := s.Manager.Check(s.Ctx, "user1", "wrong pass")
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "hashedPassword is not the hash of the given password")
}

func (s *UserSuite) TestCheckOK() {
	require.NoError(s.T(), s.Manager.Create(s.Ctx, "user2", "pass", "some@mail.com"))
	require.NoError(s.T(), s.Manager.Check(s.Ctx, "user2", "pass"))
}

func (s *UserSuite) TestCreateExisting() {
	require.NoError(s.T(), s.Manager.Create(s.Ctx, "user3", "pass", "some@mail.com"))
	err := s.Manager.Create(s.Ctx, "user3", "pass", "some@mail.com")
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "login already exist")
}

func (s *UserSuite) TestCreateWithInvalidMail() {
	err := s.Manager.Create(s.Ctx, "user4", "pass", "broken email")
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "invalid email")
}
