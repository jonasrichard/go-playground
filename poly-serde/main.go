package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

type Command interface {
	Handle() error
}

type Request struct {
	Originator string
	Command    Command
}

type FirstCommand struct {
	Name string
}

func (f FirstCommand) Handle() error {
	return nil
}

func SerializeToJSON(req Request) ([]byte, error) {
	cmdType := reflect.TypeOf(req.Command).Name()

	var cmdData interface{}
	switch req.Command.(type) {
	case FirstCommand:
		cmdData = req.Command.(FirstCommand)
		// Handle other concrete types here as needed...
	}

	rawReq := struct {
		Originator  string      `json:"originator"`
		CommandType string      `json:"command_type"`
		CommandData interface{} `json:"command_data"`
	}{
		Originator:  req.Originator,
		CommandType: cmdType,
		CommandData: cmdData,
	}

	return json.Marshal(rawReq)
}

func DeserializeFromJSON(data []byte) (Request, error) {
	var rawReq struct {
		Originator  string          `json:"originator"`
		CommandType string          `json:"command_type"`
		CommandData json.RawMessage `json:"command_data"`
	}

	if err := json.Unmarshal(data, &rawReq); err != nil {
		return Request{}, err
	}

	var cmd Command
	switch rawReq.CommandType {
	case "FirstCommand":
		var fc FirstCommand
		if err := json.Unmarshal(rawReq.CommandData, &fc); err != nil {
			return Request{}, err
		}
		cmd = fc
	// Handle other concrete types here as needed...

	default:
		return Request{}, fmt.Errorf("unknown CommandType: %s", rawReq.CommandType)
	}

	return Request{
		Originator: rawReq.Originator,
		Command:    cmd,
	}, nil
}

func SerializeToMessagePack(req Request) ([]byte, error) {
	cmdType := reflect.TypeOf(req.Command).Name()

	var cmdData []byte
	switch req.Command.(type) {
	case FirstCommand:
		fc := req.Command.(FirstCommand)
		var err error
		cmdData, err = msgpack.Marshal(fc)
		if err != nil {
			return nil, err
		}
		// Handle other concrete types here as needed...
	}

	rawReq := struct {
		Originator  string `msgpack:"originator"`
		CommandType string `msgpack:"command_type"`
		CommandData []byte `msgpack:"command_data"`
	}{
		Originator:  req.Originator,
		CommandType: cmdType,
		CommandData: cmdData,
	}

	return msgpack.Marshal(rawReq)
}

func DeserializeFromMessagePack(data []byte) (Request, error) {
	var rawReq struct {
		Originator  string `msgpack:"originator"`
		CommandType string `msgpack:"command_type"`
		CommandData []byte `msgpack:"command_data"`
	}

	if err := msgpack.Unmarshal(data, &rawReq); err != nil {
		return Request{}, err
	}

	var cmd Command
	switch rawReq.CommandType {
	case "FirstCommand":
		var fc FirstCommand
		if err := msgpack.Unmarshal(rawReq.CommandData, &fc); err != nil {
			return Request{}, err
		}
		cmd = fc
		// Handle other concrete types here as needed...
	}

	return Request{
		Originator: rawReq.Originator,
		Command:    cmd,
	}, nil
}

func main() {
	// Create a Request with a polymorphic Command field
	request := Request{
		Originator: "user123",
		Command: FirstCommand{
			Name: "example",
		},
	}

	// Serialize to JSON
	jsonData, err := SerializeToJSON(request)
	if err != nil {
		panic(err)
	}
	fmt.Println("JSON Serialized:", string(jsonData))

	// Deserialize from JSON
	deserializedJSON, err := DeserializeFromJSON(jsonData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("JSON Deserialized: %+v\n", deserializedJSON)

	// Serialize to MessagePack
	msgpackData, err := SerializeToMessagePack(request)
	if err != nil {
		panic(err)
	}
	fmt.Println("MessagePack Serialized:", msgpackData)

	// Deserialize from MessagePack
	deserializedMsgPack, err := DeserializeFromMessagePack(msgpackData)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MessagePack Deserialized: %+v\n", deserializedMsgPack)
}
