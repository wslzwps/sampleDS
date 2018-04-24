package datanode

import (
	"sync"
)

type IDataManager interface {
	Write(data []byte, size uint32) (fid, offset uint32, err error)
	Read(fid, offset, size uint32) (data []byte, err error)
}

type DataManager struct {
	workChunk     *Chunk
	workChunkLock sync.Mutex

	switchLock sync.RWMutex

	readOnlyChunks map[uint32]*Chunk
	snapshot       SnapShot

	maxId uint32
	dir   string
}

func (dm *DataManager) loadChunks() {

}

//切换的workChunk的时候与读请求互斥。（可优化）。
func (dm *DataManager) triggerSwitchWorkChunk() {
	dm.switchLock.Lock()
	defer dm.switchLock.Unlock()

	//1将workChunk变为普通chunk
	dm.workChunk.UnActive()
	//2将workChunk转移到readOnlyChunks中
	dm.readOnlyChunks[dm.workChunk.Id] = dm.workChunk
	//3创建新的workChunk
	dm.workChunk = &Chunk{Id: dm.maxId + 1, Path: dm.dir, isActive: true}
}

//implenment interface IDataManager
func (dm *DataManager) Write(data []byte, size uint32) (fid, offset uint32, err error) {
	dm.workChunkLock.Lock()
	defer dm.workChunkLock.Unlock()

	fid, offset, err = dm.workChunk.Write(data, size)
	if offset > 512*1024*1024 {
		dm.triggerSwitchWorkChunk()
	}

	return
}

func (dm *DataManager) Read(fid, offset, size uint32) (data []byte, err error) {
	dm.switchLock.RLock()
	defer dm.switchLock.RUnlock()

	var chunk *Chunk
	if chunk = dm.findChunk(fid); chunk == nil {
		//return nil,"chunk not found"
	}

	return chunk.Read(fid, offset, size)
}

func (dm *DataManager) findChunk(fid uint32) (chunk *Chunk) {

	if v, ok := dm.readOnlyChunks[fid]; ok {
		chunk = v
	} else if dm.workChunk.Id == fid {
		chunk = dm.workChunk
	}
	return
}
