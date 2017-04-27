package main

import "github.com/fsufitch/slickrlink-transfer-node/server"
import "github.com/fsufitch/slickrlink-transfer-node/protobufs"
import "github.com/golang/protobuf/proto"

func main() {
	x := &protobufs.TransferNodeToClientMessage{}
	proto.Marshal(x)
	server.StartServer(8888)
}
