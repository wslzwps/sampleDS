package datanode

import (
	"github.com/coreos/etcd/raft/raftpb"
	"errors"
	"github.com/coreos/etcd/raftsnap"
	"encoding/json"
)

type Store struct {
	dm IDataManager
	lm ILogManager
	hardstate raftpb.HardState
	snapshot *raftpb.Snapshot

	snapshotter *raftsnap.Snapshotter
}

type Mate struct {
	Oid     uint64
	Fid 	uint32
	Offset 	uint32
	Size 	uint32
	Opcode 	uint8
}
const (
	opWrite=0xF1
	opDelete=0xF2
)


func NewStore() (store *Store,err error){
	store=new(Store)
	store.dm,err=NewDataManager("/home/cintell/data")
	if err!=nil{
		return nil,err
	}
	store.lm,err=NewLogManager("/home/cintell/log/raft1111.log")
	if err!=nil{
		return nil,err
	}
	store.snapshotter=raftsnap.New(nil,"/snapshot/")
	return store,err
}


func (s *Store) InitialState() (raftpb.HardState,raftpb.ConfState,error){
	return s.hardstate,s.snapshot.Metadata.ConfState,nil
}

func (s *Store) Entries(lo,hi,maxsize uint64) (entries []raftpb.Entry,err error){
	if lo>hi {
		err=errors.New("lo比hi大")
		return
	}

	last,err:=s.LastIndex()
	if err!=nil{
		err=errors.New("获取lastindex失败:"+err.Error())
		return
	}

	if hi>last {
		panic("hi is out if bound lastindex")
	}

	if hi-lo>=maxsize{
		hi=lo+maxsize-1
	}

	for i:=lo;i<=hi;i++{
		if e,err:=s.lm.Index(i);err!=nil{
			return nil,err
		}else{
			entries=append(entries,*e)
		}
	}
	return
}

func (s *Store) Term(i uint64) (t uint64,err error){
	entry,err:=s.lm.Index(i)
	if err!=nil{
		return
	}
	t=entry.Term
	return
}
func (s *Store) LastIndex()(index uint64,err error) {
	entry,err:=s.lm.Last()
	if err!=nil{
		return
	}
	index=entry.Index
	return
}
func (s *Store) FirstIndex()(index uint64,err error) {
	entry,err:=s.lm.First()
	if err!=nil{
		return
	}
	index=entry.Index
	return
}

func (s *Store) Snapshot() (raftpb.Snapshot,error){
	if s.snapshot!=nil{
		return *s.snapshot,nil
	}
	snap,err:=s.snapshotter.Load()
	return *snap,err
}


//保存snapshot
func (s *Store)SaveSnap(snap raftpb.Snapshot) error{
	if s.snapshot.Metadata.Index >snap.Metadata.Index {
		return errors.New("snapshot out fo date")
	}

	s.snapshot=&snap
	return s.snapshotter.SaveSnap(snap)
}
//保存hardstate,暂时先存在内存
//todo:需要持久化到文件里面
func (s *Store)SaveHardState(hs raftpb.HardState) error{
	s.hardstate=hs
	return nil
}

func (s *Store)Append(entries []raftpb.Entry) error{
	for i,_:=range entries{
		size:=uint32(len(entries[i].Data))
		fid,offset,err:=s.dm.Write(entries[i].Data,size)
		if err!=nil {
			panic(err.Error())
		}
		//暂时将mate以json的方式存在log中。
		entries[i].Data,_=json.Marshal(Mate{Size:size,Fid:fid,Offset:offset,Opcode:opWrite})
	}

	return s.lm.Append(entries)
}
