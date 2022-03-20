package main

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestRoomDAOFindByPlayer(t *testing.T) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	require.NoError(t, err)

	ctx := context.Background()
	dao := NewRoomDAO(session)
	defer dao.RemoveAll(ctx)

	r := &Room{
		Sides: []RoomSideInfo{{
			Name: "evgsol",
		}, {
			Name: "solarka",
		}},
	}

	_, err = dao.Insert(ctx, r)
	require.NoError(t, err)

	found, err := dao.FindOneByPlayer(ctx, "evgsol")
	require.NoError(t, err)
	require.NotNil(t, found)

	found, err = dao.FindOneByPlayer(ctx, "solarka")
	require.NoError(t, err)
	require.NotNil(t, found)

	found, err = dao.FindOneByPlayer(ctx, "miracle")
	require.Error(t, err)
	require.Nil(t, found)
}
