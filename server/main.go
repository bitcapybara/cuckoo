package main

import (
	"flag"
	"github.com/bitcapybara/raft"
	"log"
	"strings"
)

const NoneOption = ""

func main() {
	// 命令行参数定义
	var me string
	flag.StringVar(&me, "me", NoneOption, "当前节点nodeId")
	var peerStr string
	flag.StringVar(&peerStr, "peers", NoneOption, "指定所有节点地址，nodeId@nodeAddr，多个地址使用逗号间隔")
	var role string
	flag.StringVar(&role, "role", NoneOption, "当前节点角色")
	flag.Parse()

	// 命令行参数解析
	if me == "" {
		log.Fatal("未指定当前节点id！")
	}

	if peerStr == "" {
		log.Fatalln("未指定集群节点")
	}

	if role == "" {
		log.Fatalln("未指定节点角色")
	}

	peerSplit := strings.Split(peerStr, ",")
	peers := make(map[raft.NodeId]raft.NodeAddr, len(peerSplit))
	for _, peerInfo := range peerSplit {
		idAndAddr := strings.Split(peerInfo, "@")
		peers[raft.NodeId(idAndAddr[0])] = raft.NodeAddr(idAndAddr[1])
	}

	// 启动 server
	startServer(raft.RoleFromString(role), raft.NodeId(me), peers)
}
