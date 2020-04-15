package pool

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pisign/pisign-backend/api/text"
	"github.com/pisign/pisign-backend/types"
)

func NewFakePool() Pool {
	return Pool{
		registerChan:   nil,
		unregisterChan: nil,
		Map:            nil,
		name:           "test name",
		ImageDB:        nil,
	}
}

func TestNewPool(t *testing.T) {
	newPool := NewPool()

	if newPool.registerChan == nil {
		t.Error("registerChan not created")
	}

	if newPool.unregisterChan == nil {
		t.Error("unregisterChan not created")
	}

	if newPool.Map == nil {
		t.Error("map not created")
	}

	if newPool.ImageDB == nil {
		t.Error("new ImageDB object not created")
	}
}

func TestPool_List(t *testing.T) {
	customMap := make(map[uuid.UUID]types.API)
	uid, _ := uuid.NewUUID()

	newAPI := text.NewAPI(nil, nil, uid)
	customMap[uid] = newAPI

	customPool := NewFakePool()
	customPool.Map = customMap

	keys := customPool.List()

	if len(keys) != 1 {
		t.Error("keys is not the correct size")
	}

	if keys[0] != newAPI {
		t.Error("items not iterated over properly in List")
	}
}

func TestPool_GetImageDB(t *testing.T) {
	pool := NewFakePool()
	imageDB := types.ImageDB{}
	pool.ImageDB = &imageDB

	if pool.GetImageDB() != &imageDB {
		t.Error("not returning correct imageDB")
	}
}
