package datanode

import (
	"os"
	"datanode/raftlogpb"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/gogo/protobuf/proto"
)


type ILogManager interface {
	Append(data []byte) error
	Index(i uint64) (b []byte, err error)
	First() (b []byte, err error)
	Last() (b []byte, err error)
}

//LBO means: [Log Binary Object]
//We've wrapped a layer above the Entry
//Added offset and size.This way we can easily load the specified log by offset
//More importantly.If the log is very large.We use to load the log between snapshot and last index.
//so our inner snapshot needs to include LBO.
type LBO struct {
	raftlogpb.LBO
}

type LogManager struct {
	path string
	file LogFile
	lbos []LBO
}

//需要多种情况。包含网络隔离带来的日志不一致。
func (lm *LogManager)MaybeAppend(){

}

func (lm *LogManager) append(entry *raftpb.Entry) (err error){
	offset,_:=lm.file.size()
	size:=uint32(12+len(entry.Data))

	lbo:=raftlogpb.LBO{Offset:&offset,Size:&size,Entry:entry}
	data,err:=proto.Marshal(&lbo)
	if err!=nil{
		return err
	}

	if err=lm.file.append(data);err==nil {
		lm.lbos=append(lm.lbos, lbo)
	}
	return
}

func (lm *LogManager) Index(i uint64) (b []byte, err error) {

	return
}

func (lm *LogManager) First() (b []byte, err error) {
	return lm.Index(0)
}
func (lm *LogManager) Last() (b []byte, err error) {
	return
}



//**********************************************************************************
type LogFile struct {
	fd         *os.File
}

func newLogFile(path string) (file *LogFile, err error) {
	file = new(LogFile)
	file.fd, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	return
}



func (file *LogFile) append(b []byte) (err error) {
	_, err = file.fd.Write(b)
	return
}


func (file *LogFile) read(offset int64, size uint32) ([]byte, error) {
	b := make([]byte, size)
	_, err := file.fd.ReadAt(b, offset)
	return b, err
}

func (file *LogFile)size()(size int64,err error){
	fi,err:=file.fd.Stat()
	if err!=nil{
		return
	}
	size=fi.Size()
	return
}


func (file *LogFile)truncate(lo,hi int64){

}

func (file *LogFile) close() {
	file.fd.Close()
}