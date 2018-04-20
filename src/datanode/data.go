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

	readOnlyChunks map[uint32]*Chunk
	snapshot       SnapShot

	maxId uint32
	dir   string
}

func (dm *DataManager) loadChunks() {

}
func (dm *DataManager) allocateChunk() {

}
func (dm *DataManager) triggerSwitchWorkChunk() {
	dm.workChunkLock.Lock()
	defer dm.workChunkLock.Unlock()

	dm.workChunk = &Chunk{Id: dm.maxId + 1}
}

//implenment interface IDataManager
func (dm *DataManager) Write(data []byte, size uint32) (fid, offset uint32, err error) {
	dm.workChunkLock.Lock()
	defer dm.workChunkLock.Unlock()

	return dm.workChunk.Write(data, size)
}
func (dm *DataManager) Read(fid, offset, size uint32) (data []byte, err error) {
	return dm.workChunk.Read(fid, offset, size)
}
