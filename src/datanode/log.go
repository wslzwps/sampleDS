package datanode

type ILogManager interface {
	Append()
	Compact()
	ApplySnapShot()
	SnapShot()
}

type LogManager struct {
	dm IDataManager
	rl RaftLog
}

func (lm *LogManager) Append() {
	//1:将数据从entry中提取出来,获取
	//	data，size
	//2:将数据持久化到dm
	//	fid,offset,err:=lm.dm.Write(data,size)
	//3:将fid,offset,size。填充到entry中。
	//	 entry.Data=[]byte{fid,offset,size}
	//4:将填充后的entry持久化到磁盘的log文件里

}
func (lm *LogManager) Compact() {

}
func (lm *LogManager) ApplySnapShot() {

}
func (lm *LogManager) CreateSnapShot() {

}
func (lm *LogManager) SnapShot() {

}

type RaftLog struct {
	path       string
	prelogsize uint32
}

func (rl *RaftLog) Append() {

}
func (rl *RaftLog) Write(index uint32) {

}
func (rl *RaftLog) Read(index uint32) {

}
func (rl *RaftLog) Truncate(index uint32) {

}
