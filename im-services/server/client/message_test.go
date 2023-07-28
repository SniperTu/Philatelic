package client

import (
	"testing"
)

func TestMessage(t *testing.T) {
	var service GrpcMessageService

	service.SendGpcMessage("", "127.0.0.1:8002")
}
