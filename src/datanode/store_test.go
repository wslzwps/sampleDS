package datanode

import (
	"testing"
	"github.com/coreos/etcd/raft/raftpb"
	"fmt"
)

func TestStore_Append(t *testing.T) {
	e1:=raftpb.Entry{Index:1,Term:2,Data:[]byte("aaaa")}
	e2:=raftpb.Entry{Index:2,Term:2,Data:[]byte("bbbb")}
	e3:=raftpb.Entry{Index:3,Term:3,Data:[]byte("cccc")}
	entries:=make([]raftpb.Entry,0)
	entries=append(entries,e1)
	entries=append(entries,e2)
	entries=append(entries,e3)

	store,err:=NewStore()
	fmt.Println(store,err)
	if err!=nil {
		return
	}

	entries1,err:=store.Entries(1,3,3)
	for _,e:=range entries1 {
		fmt.Println(string(e.Data))
	}
	fmt.Println(store.FirstIndex())
	fmt.Println(store.LastIndex())
	fmt.Println(store.Term(3))
	//store.Append(entries)
}
