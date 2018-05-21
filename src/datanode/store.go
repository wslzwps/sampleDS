package datanode

import (
	"github.com/coreos/etcd/raft/raftpb"
	"errors"
	"github.com/coreos/etcd/raftsnap"
	"encoding/json"
	"os"
)

type Store struct {
	dm IDataManager
	lm ILogManager
	hardstate raftpb.HardState
	snapshot  raftpb.Snapshot

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


func NewStore(basePath string) (store *Store,err error){
	store=new(Store)
	store.dm,err=NewDataManager(basePath+"/data")
	if err!=nil{
		return nil,err
	}
	store.lm,err=NewLogManager(basePath+"/log/raft.log")
	if err!=nil{
		return nil,err
	}
	store.snapshotter=raftsnap.New(nil,basePath+"/snapshot")
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

	if hi>last+1 {
		panic("hi is out if bound lastindex")
	}

	if hi-lo>=maxsize{
		hi=lo+maxsize-1
	}

	for i:=lo;i<=hi;i++{
		if e:=s.lm.Index(i);e!=nil{
			entries=append(entries,*e)
		}
	}
	return
}

func (s *Store) Term(i uint64) (t uint64,err error){
	entry:=s.lm.Index(i)
	if entry==nil{
		return
	}
	t=entry.Term
	return
}
func (s *Store) LastIndex()(index uint64,err error) {
	entry:=s.lm.Last()
	if entry==nil{
		return
	}
	index=entry.Index
	return
}
func (s *Store) FirstIndex()(index uint64,err error) {
	entry:=s.lm.First()
	if entry==nil{
		return 1,nil
	}
	index=entry.Index
	return
}

func (s *Store) Snapshot() (raftpb.Snapshot,error){
	snap,err:=s.snapshotter.Load()
	if err==nil{
		s.snapshot=*snap
	}

	return s.snapshot,err
}


//保存snapshot
func (s *Store)SaveSnap(snap raftpb.Snapshot) error{
	if s.snapshot.Metadata.Index >snap.Metadata.Index {
		return errors.New("snapshot out fo date")
	}

	s.snapshot=snap
	return s.snapshotter.SaveSnap(snap)
}
//保存hardstate,暂时先存在内存
//todo:需要持久化到文件里面
func (s *Store)SaveHardState(hs raftpb.HardState) error{
	s.hardstate=hs
	return nil
}

// ApplySnapshot overwrites the contents of this Storage object with
// those of the given snapshot.
func (s *Store) ApplySnapshot(snap raftpb.Snapshot) error {

	return nil
}

// CreateSnapshot makes a snapshot which can be retrieved with Snapshot() and
// can be used to reconstruct the state at that point.
// If any configuration changes have been made since the last compaction,
// the result of the last ApplyConfChange must be passed in.
func (s *Store) CreateSnapshot(i uint64, cs *raftpb.ConfState, data []byte) (raftpb.Snapshot, error) {

	return s.snapshot, nil
}

// Compact discards all log entries prior to compactIndex.
// It is the application's responsibility to not attempt to compact an index
// greater than raftLog.applied.
func (s *Store) Compact(compactIndex uint64) error {

	return nil
}

func (s *Store)Append(entries []raftpb.Entry) error{
	for i,_:=range entries{
		size:=uint32(len(entries[i].Data))
		if size==0{
			break
		}
		fid,offset,err:=s.dm.Write(entries[i].Data,size)
		if err!=nil {
			panic(err.Error())
		}
		//暂时将mate以json的方式存在log中。
		entries[i].Type=raftpb.EntryNormalPersistent
		entries[i].Data,_=json.Marshal(Mate{Size:size,Fid:fid,Offset:offset,Opcode:opWrite})
	}

	return s.lm.Append(entries)
}
func (s *Store)Read(fid, offset, size uint32)([]byte,error){
	return s.dm.Read(fid,offset,size)
}

func (s *Store)Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
