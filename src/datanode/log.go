package datanode

import (
	"os"
	"github.com/coreos/etcd/raft/raftpb"
	"datanode/raftlogpb"
)

var _=raftlogpb.LBO{}

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
	offset int64
	size uint32
	entry raftpb.Entry
}

type LogManager struct {
	path string
	file LogFile
}

func (lm *LogManager) Append(data []byte) error{
	return lm.file.append(data)

}

func (lm *LogManager) Index(i uint64) (b []byte, err error) {

	return
}

func (lm *LogManager) First() (b []byte, err error) {
	return lm.Index(0)
}
func (lm *LogManager) Last() (b []byte, err error) {
	//return lm.Index(uint64(len(lm.file.buf) - 1))
	return
}



//**********************************************************************************
type Buffer []byte

type LogFile struct {
	fd         *os.File
	//buf        []Buffer
	prelogsize uint32
}

func newLogFile(path string) (file *LogFile, err error) {
	file = new(LogFile)
	//file.buf = make([]Buffer, 0)
	file.fd, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	return
}



func (file *LogFile) append(b []byte) (err error) {
	_, err = file.fd.Write(b)
	if err!=nil{
		return
	}

	//file.buf = append(file.buf, b)
	return
}


func (file *LogFile) read(offset int64, size uint32) ([]byte, error) {
	b := make([]byte, size)
	_, err := file.fd.ReadAt(b, offset)
	return b, err
}


func (file *LogFile) close() {
	file.fd.Close()
}