package socket

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/pisign/pisign-backend/types"
)

func TestSocket_CloseChan(t *testing.T) {
	closeChan := make(chan bool)
	socket := Socket{nil, closeChan, nil}

	if socket.CloseChan() != closeChan {
		t.Error("CloseChan does not return the proper object")
	}
}

func TestSocket_ConfigChan(t *testing.T) {
	configChan := make(chan types.ClientMessage)
	socket := Socket{nil, nil, configChan}

	if socket.ConfigChan() != configChan {
		t.Error("CloseChan does not return the proper object")
	}
}

func TestSocket_Close(t *testing.T) {
	socket := Socket{nil, nil, nil}

	err := socket.Close()
	if err == nil {
		t.Error("nil conn did not throw error when attempting to close from Socket type")
	}
}

func TestSocket_Create(t *testing.T) {
	configChan := make(chan types.ClientMessage)
	conn := new(websocket.Conn)

	createdSocket := Create(configChan, conn)
	if createdSocket.configChan != configChan {
		t.Error("config chans do not match for created socket")
	}

	if createdSocket.conn != conn {
		t.Error("conn objects do not match for created socket")
	}

	if createdSocket.closeChan == nil {
		t.Error("socket Create did not create a new closeChan")
	}
}
