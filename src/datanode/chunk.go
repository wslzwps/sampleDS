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

func (chunk *Chunk) Write(data []byte, size uint32) (fid, offset uint32, err error) {
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

func (chunk *Chunk) Read(fid, offset, size uint32) (data []byte, err error) {
	chunk.loadfile()
	if fid != chunk.Id {
		err = MisUpChunkIdErr
		return
	}
	//should get from pool
	data = make([]byte, size)
	_, err = chunk.fp.ReadAt(data, int64(offset))
	return
}

func (chunk *Chunk) UnActive() {

	path := fmt.Sprintf("%s/%d", chunk.Path, chunk.Id)
	if chunk.isActive == true && chunk.fp != nil {
		//reopen to be readonly
		chunk.fp.Close()

		if fp, err := os.OpenFile(path, os.O_RDONLY, 0444); err == nil {
			chunk.fp = fp
			chunk.isActive = false
		} else {
			panic("set chunk to be UnActive failed:" + err.Error())
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
