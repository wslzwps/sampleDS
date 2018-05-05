package datanode

import(
	"os"
	)

type ILogManager interface {
	Append()
	Compact()
	ApplySnapShot()
	SnapShot()
}

type LogManager struct {
	lf LogFile
}

func (lm *LogManager) Append(data []byte) {
	lm.lf.Append(data)

}
func (lm *LogManager) Compact() {

}
func (lm *LogManager) ApplySnapShot() {

}
func (lm *LogManager) CreateSnapShot() {

}
func (lm *LogManager) SnapShot() {

}

type Buffer []byte
type LogFile struct {
	path       string
	fd   *os.File
	prelogsize int
	buf   []Buffer
}

func newLogFile(){

}

func (lf *LogFile) create() (err error){
	rl.fd,err=os.OpenFile(path,os.O_RDWR|os.O_APPEND|os.O_CREATE,0777)
}

func (lf *LogFile) Append(log []byte) (err error){
	if len(log)!=lf.prelogsize{
		//return "unexcept log size please check it"
	}
	_,err=lf.fd.Write(log)
	lf.buf=append(lf.buf,log)
}
func (lf *LogFile) Read(index uint32) ([]byte,error){
	//将buf[0]转化为entry的结构
	//entry:=unmarsh(bufer[0])
	//需要判断index是否在first和last之间
	//offset:=(index-entry.index)*lf.prelogsize
	//b:=make([]byte,lf.prelogsize)
	//_,err:=lf.fd.ReadAt(b,offset)
	//return b,err
}
func (lf *LogFile) Truncate(index uint32) {
	//如果index大于last。将文件truncate为0
	//lf.buf=make([]Buffer,0)

	//如果Index小于first，什么也不做。

	//如果index位于first和last之间，则按偏移截取文件。
	//lf.buf原有数据清空。重新加载截取后的文件。
}
