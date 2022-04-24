package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"time"
)

var zkaddr = []string{"127.0.0.1:2184", "127.0.0.1:2182", "127.0.0.1:2183"}

func main() {
	get()
}

func update()  {
	conn, _, err := zk.Connect(zkaddr, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	path := "/hello"
	_, state, _ := conn.Get(path)
	state, err = conn.Set(path, []byte("girl"), state.Version)
	if err != nil {
		panic(err)
	}
	fmt.Println("state ->")
	fmt.Printf("cZxid=%d\nctime=%d\nmZxid=%d\nmtime=%d\npZxid=%d\ncversion=%d\ndataVersion=%d\naclVersion=%d\nephemeralOwner=%v\ndataLength=%d\nnumChildren=%d\n", state.Czxid, state.Ctime, state.Mzxid, state.Mtime, state.Pzxid, state.Cversion, state.Version, state.Aversion, state.EphemeralOwner, state.DataLength, state.NumChildren)

	data, _, err := conn.Get(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nnew value: ", string(data))
}

func delete()  {
	conn, _, err := zk.Connect(zkaddr, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	path := "/hello"
	exists, state, err := conn.Exists(path)
	fmt.Printf("\npath[%s] exists: %v\n", path, exists)

	err = conn.Delete(path, state.Version)
	if err != nil {
		panic(err)
	}
	fmt.Printf("path[%s] is deleted.", path)

	exists, _, err = conn.Exists(path)
	fmt.Printf("\npath[%s] exists: %v\n", path, exists)
}

func get()  {

	conn, _, err := zk.Connect(zkaddr, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	result, state, err := conn.Get("/hello")
	if err != nil {
		panic(err)
	}
	fmt.Println("result: ", string(result))
	fmt.Println("state ->")
	fmt.Printf("cZxid=%d\n" +
		"ctime=%d\n" +
		"mZxid=%d\n" +
		"mtime=%d\n" +
		"pZxid=%d\n" +
		"cversion=%d\n" +
		"dataVersion=%d\n" +
		"aclVersion=%d\n" +
		"ephemeralOwner=%v\n" +
		"dataLength=%d\nn" +
		"umChildren=%d\n", state.Czxid,
		state.Ctime,
		state.Mzxid,
		state.Mtime,
		state.Pzxid,
		state.Cversion,
		state.Version,
		state.Aversion,
		state.EphemeralOwner,
		state.DataLength, state.NumChildren)
}

func insert()  {
	conn, _, err := zk.Connect(zkaddr, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 创建持久节点
	path, err := conn.Create("/hello", []byte("world"), 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
	println("Created", path)

	// 创建临时节点，创建此节点的会话结束后立即清除此节点
	ephemeral, err := conn.Create("/ephemeral", []byte("1"), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
	println("Ephemeral node created:", ephemeral)

	// 创建持久时序节点
	sequence, err := conn.Create("/sequence", []byte("1"), zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
	println("Sequence node created:", sequence)

	// 创建临时时序节点，创建此节点的会话结束后立即清除此节点
	ephemeralSequence, err := conn.Create("/ephemeralSequence", []byte("1"), zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		panic(err)
	}
	println("Ephemeral-Sequence node created:", ephemeralSequence)
}