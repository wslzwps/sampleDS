package datanode

import (
	"fmt"
	"os"
)

//Chunk is the abstraction of the container file.
//For user. chunk only provides read and write functions. and users do not need to pay attention
//to when to create a file .The file will be automatically created.
//There are two types of chunk. active[Read-write] and unactive[Read-only].
//DataManager has only one [active chunk] and multiple [unactive chunk]
type Chunk struct {
	Id   uint32
	Path string

	fp *os.File
	isActive bool
	mustSync bool
}

//Write success returns file id and offset.Failed to return error
func (chunk *Chunk) Write(data []byte, size uint32) (fid, offset uint32, err error) {
	//Panic may appear if not checked
	if uint32(len(data)) < size {
		err = UnMatchSizeErr
		return
	}

	chunk.loadfile()

	fid = chunk.Id
	if chunk.isActive == true {
		offset = uint32(chunk.size())
		_, err = chunk.fp.Write(data[:size])
	} else {
		err = UnActiveChunkWriteErr
	}
	return
}
//Write success returns data .Failed to return error
func (chunk *Chunk) Read(fid, offset, size uint32) (data []byte, err error) {
	chunk.loadfile()
	if fid != chunk.Id {
		err = MisUpChunkIdErr
		return
	}
	//To be optimized.this should be obtained from the memory pool
	data = make([]byte, size)
	_, err = chunk.fp.ReadAt(data, int64(offset))

	return
}
//UnActive turn active chunks into inactive.without return
func (chunk *Chunk) UnActive() {
	path := fmt.Sprintf("%s/%d", chunk.Path, chunk.Id)
	//before reopening the file, you need to close the previous fd
	if chunk.isActive == true && chunk.fp != nil {
		chunk.fp.Close()
		if fp, err := os.OpenFile(path, os.O_RDONLY, 0444); err == nil {
			chunk.fp = fp
			chunk.isActive = false
		} else {
			panic("chunk to be UnActive failed:" + err.Error())
		}
	}
	return
}

func (chunk *Chunk) Close() {
	if chunk.fp != nil {
		chunk.fp.Close()
	}
}

func (chunk *Chunk) size() int64 {
	if chunk.fp == nil {
		panic("chunk: fp is nil")
	}

	fi, err := chunk.fp.Stat()
	if err != nil {
		panic("chunk:" + err.Error())
	}
	return fi.Size()
}


//Active file open with read-write
//Inactive file read-only opens
func (chunk *Chunk) loadfile() {
	if chunk.fp != nil {
		return
	}

	var err error
	path := fmt.Sprintf("%s/%d", chunk.Path, chunk.Id)
	if chunk.isActive == true {
		//open file with  RDWR|APPEND mode. if not exsit create it。
		chunk.fp, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_EXCL, 0644)

	} else {
		//open file with readonyl mode.
		chunk.fp, err = os.OpenFile(path, os.O_RDONLY, 0444)
	}

	if err != nil {
		panic("Load chunk failed:" + err.Error())
	}
	return
}
