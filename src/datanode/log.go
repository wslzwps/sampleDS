package datanode

import (
	"errors"
	"os"
)

type ILogManager interface {
	Index(i uint64) (b []byte, err error)
	First() (b []byte, err error)
	Last() (b []byte, err error)
	Append(data []byte) error
	Compact()
}

type LogManager struct {
	path string
	file LogFile
}

func (lm *LogManager) Append(data []byte) error{
	return lm.file.append(data)

}
func (lm *LogManager) Compact() {
	lm.file.truncate()
}
func (lm *LogManager) Index(i uint64) (b []byte, err error) {
	if len(lm.file.buf) == 0 {
		err = errors.New("Empty log buffer")
		return
	}

	if i > uint64(len(lm.file.buf)) {
		err = errors.New("Invalid index")
		return
	}

	b = lm.file.buf[i]
	return
}

func (lm *LogManager) First() (b []byte, err error) {
	return lm.Index(0)
}
func (lm *LogManager) Last() (b []byte, err error) {
	return lm.Index(uint64(len(lm.file.buf) - 1))
}



//**********************************************************************************
type Buffer []byte

type LogFile struct {
	fd         *os.File
	buf        []Buffer
	prelogsize uint32
}

func newLogFile(path string) (file *LogFile, err error) {
	file = new(LogFile)
	file.buf = make([]Buffer, 0)
	file.fd, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	return
}

func (file *LogFile) close() {
	file.fd.Close()
}

func (file *LogFile) append(b []byte) (err error) {
	if len(b) != int(file.prelogsize) {
		err = errors.New("unexcept log size please check that.")
		return
	}

	_, err = file.fd.Write(b)
	file.buf = append(file.buf, b)
	return
}

func (file *LogFile) reload()(err error){
	fi,err:=file.fd.Stat()
	if err!=nil {
		return
	}

	blockcount:=fi.Size()/int64(file.prelogsize)
	
	for i:=int64(0);i<blockcount;i++{
		offset:=blockcount*i
		if b,err:=file.read(offset,file.prelogsize);err!=nil{
			break
		}else{
			file.buf=append(file.buf,b)
		}

	}
	return

}

func (file *LogFile) read(offset int64, size uint32) ([]byte, error) {
	b := make([]byte, size)
	_, err := file.fd.ReadAt(b, offset)
	return b, err
}

func (file *LogFile) truncate() error {
	return file.fd.Truncate(0)
}
