package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPacketFromJson(t *testing.T) {
	a := `{"computerid":"23","packetid":"42","packetType":"test","arguments":{"arg0":"value0"},"response":{"foo":"bar"}}`
	var packet Packet
	err := json.Unmarshal([]byte(a), &packet)
	if err != nil {
		t.Errorf("Could not parse packet test 1: %s", err)
	}
	if packet.Response["foo"] != "bar" {
		t.Errorf("Could not parse packet test 2: %s", err)
	}
}

func TestPacketToJson(t *testing.T) {
	arguments := make(PacketArgument)
	arguments["arg0"] = "value0"
	response := make(PacketResponse)
	c := NewPacket("test", "23", "42", arguments, response)

	reference := `{"computerid":"23","packetid":"42","packetType":"test","arguments":{"arg0":"value0"},"response":{}}`
	u, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	s := string(u)
	fmt.Println(s)
	if s != reference {
		t.Errorf("Error jsonify: " + s)
	}
}

func TestResponseA(t *testing.T) {
	response := make(PacketResponse, 0)
	response["test"] = "test"
	response["k0"] = "v0"
	response["k1"] = "v1"
	response["k2"] = "v2"

	arr := ResponseToArray("k", response)
	if arr[0] != "v0" {
		t.Error("Error")
	}
	if arr[1] != "v1" {
		t.Error("Error")
	}
	if arr[2] != "v2" {
		t.Error("Error")
	}
}

func TestResponseB(t *testing.T) {
	arr := []string{"v0", "v1", "v2"}

	response := make(PacketResponse, 0)

	AddArrayToResponse("k", arr, response)
	if response["k0"] != "v0" {
		t.Error("Error")
	}
	if response["k1"] != "v1" {
		t.Error("Error")
	}
	if response["k2"] != "v2" {
		t.Error("Error")
	}
}
