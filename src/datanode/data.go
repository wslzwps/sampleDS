package datanode

import (
	"io/ioutil"
	"strconv"
	"sync"
	"fmt"
	"errors"
)

type IDataManager interface {
	Write(data []byte, size uint32) (fid, offset uint32, err error)
	Read(fid, offset, size uint32) (data []byte, err error)
}

type DataManager struct {
	workChunk      *Chunk
	workChunkLock  sync.Mutex
	switchLock     sync.RWMutex
	readOnlyChunks map[uint32]*Chunk

	maxId uint32
	dir   string
}

func NewDataManager(path string) (dm *DataManager, err error) {
	dm = new(DataManager)
	dm.dir=path
	dm.readOnlyChunks = make(map[uint32]*Chunk)

	//todo
	//这个地方需要重构。
	//不应该在启动的时候加载所有的文件
	//每次读请求过触发文件的加载。
	//当文件fd个数超过阈值，可以采取关闭一半的方法,使其衰减。
	if err = dm.loadChunks();err==nil{
		dm.workChunk=&Chunk{Id:dm.maxId,Path:dm.dir,isActive:true}
	}

	fmt.Println("NewDataManager maxid",dm.maxId)

	return
}

func (dm *DataManager) loadChunks() (err error) {
	dir, err := ioutil.ReadDir(dm.dir)
	if err != nil {
		return
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		} else {
			id, err := strconv.ParseUint(fi.Name(), 10, 64)
			if err != nil {
				return err
			}

			dm.readOnlyChunks[uint32(id)] = &Chunk{Id: uint32(id), Path: dm.dir}

			if dm.maxId < uint32(id) {
				dm.maxId = uint32(id)
			}
		}
	}

	dm.maxId=dm.maxId+1

	return
}

//todo:切换的workChunk的时候与读请求互斥。（可优化）。
func (dm *DataManager) triggerSwitchWorkChunk() {
	dm.switchLock.Lock()
	defer dm.switchLock.Unlock()


	//1将workChunk变为普通chunk
	dm.workChunk.UnActive()
	//2将workChunk转移到readOnlyChunks中
	dm.readOnlyChunks[dm.workChunk.Id] = dm.workChunk
	//3创建新的workChunk
	dm.maxId=dm.maxId+1
	fmt.Println("triggerSwitchWorkChunk maxid:",dm.maxId)
	dm.workChunk = &Chunk{Id: dm.maxId, Path: dm.dir, isActive: true}

}

//implenment interface IDataManager
func (dm *DataManager) Write(data []byte, size uint32) (fid, offset uint32, err error) {
	dm.workChunkLock.Lock()
	defer dm.workChunkLock.Unlock()

	fid, offset, err = dm.workChunk.Write(data, size)
	if offset > 64*1024*1024 {
		dm.triggerSwitchWorkChunk()
	}

	return
}

func (dm *DataManager) Read(fid, offset, size uint32) (data []byte, err error) {
	dm.switchLock.RLock()
	defer dm.switchLock.RUnlock()

	var chunk *Chunk
	if chunk = dm.findChunk(fid); chunk == nil {
		return nil,errors.New("chunk not found")
	}

	return chunk.Read(fid, offset, size)
}

//todo:如果内存里面没有，就去磁盘上找
func (dm *DataManager) findChunk(fid uint32) (chunk *Chunk) {

	if v, ok := dm.readOnlyChunks[fid]; ok {
		chunk = v
	} else if dm.workChunk.Id == fid {
		chunk = dm.workChunk
	}
	return
}
