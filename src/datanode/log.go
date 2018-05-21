package datanode

import (
	"github.com/wslzwps/sampleDS/src/datanode/raftlogpb"
	"encoding/binary"
	"github.com/coreos/etcd/raft/raftpb"
	"github.com/gogo/protobuf/proto"
	"os"
)

type ILogManager interface {
	Append(entries []raftpb.Entry) (err error)
	Index(i uint64) (entry *raftpb.Entry)
	First() (entry *raftpb.Entry)
	Last() (entry *raftpb.Entry)
}

//LBO means: [Log Binary Object]
//We've wrapped a layer above the Entry
//Added offset and size.This way we can easily load the specified log by offset
//More importantly.If the log is very large.We use to load the log between snapshot and last index.
//so our inner snapshot needs to include LBO.
//type LBO struct {
//	raftlogpb.LBO
//}

type LogManager struct {
	path string
	file *LogFile
	lbos []*raftlogpb.LBO
}

func NewLogManager(path string) (*LogManager,error) {
	lm := new(LogManager)
	lm.path = path
	lm.lbos = make([]*raftlogpb.LBO, 0)
	//todo: error 需要处理下~
	lm.file, _ = newLogFile(path)
	lm.lbos,_=lm.file.read(0)
	return lm,nil
}

//需要多种情况---网络隔离带来的日志不一致。
func (lm *LogManager) Append(entries []raftpb.Entry) (err error) {
	if len(entries) == 0 {
		return
	}
	offset:=uint64(0)
	first:=uint64(0)
	last:=uint64(0)

	if len(lm.lbos)!=0 {
		first = lm.lbos[0].Entry.Index
		last  = entries[0].Index + uint64(len(entries)) - 1
	}


	if last < first {
		return
	}

	if entries[0].Index < first {
		entries = entries[first-entries[0].Index:]
	}

	if len(lm.lbos)!=0 {
		offset = entries[0].Index - lm.lbos[0].Entry.Index
	}

	switch {
	case uint64(len(lm.lbos)) > offset:
		if offset >= 0 {
			//将文件截断到指定位置.
			lm.file.truncate(*lm.lbos[offset].Offset)
		}else{
			//将文件截到上一次snapshot的位置
		}

		//然后将所有的entry转化成[]lbo,逐条追加到日志文件
		newlbos, err := lm.file.append(entries)
		if err != nil {
			return err
		}

		//最后修改内存的数据.
		lm.lbos = lm.lbos[:offset]
		lm.lbos = append(lm.lbos, newlbos...)
	case uint64(len(lm.lbos)) == offset:

		//将所有的entry转化成[]byte,逐条追加到日志文件
		newlbos, err := lm.file.append(entries)
		if err != nil {
			return err
		}
		//最后操作内存的数据.
		lm.lbos = append(lm.lbos, newlbos...)
	default:
		panic("missing log entry")
	}

	return

}

func (lm *LogManager) Index(i uint64) (entry *raftpb.Entry) {
	firstentry:= lm.First()
	lastentry:= lm.Last()
	if firstentry==nil||lastentry==nil{
		return
	}

	if i < firstentry.Index || i > lastentry.Index {
		return
	}

	offset := i - firstentry.Index
	entry = lm.lbos[offset].Entry

	return
}

func (lm *LogManager) First() (entry *raftpb.Entry) {
	if len(lm.lbos) != 0 {
		entry = lm.lbos[0].Entry
	}
	return
}
func (lm *LogManager) Last() (entry *raftpb.Entry) {
	if len(lm.lbos) != 0 {
		entry = lm.lbos[len(lm.lbos)-1].Entry
	}
	return
}

//**********************************************************************************
type LogFile struct {
	fd *os.File
}

func newLogFile(path string) (file *LogFile, err error) {
	file = new(LogFile)
	file.fd, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	return
}

//需要完善错误检查...
func (file *LogFile) append(entries []raftpb.Entry) (lbos []*raftlogpb.LBO,err error) {
	for _, e := range entries {
		offset, err := file.size()
		if err != nil {
			return nil, err
		}

		//参考下面的注释
		entry:=new(raftpb.Entry)
		entry.Type=e.Type
		entry.Index=e.Index
		entry.Term=e.Term
		entry.Data=append(entry.Data,e.Data...)

		//lbo这个局部变量在每次循环的时候，里面成员的entry一直地址是不变的，
		//但是lbo对象本身的地址是变化的,offset成员地址也变化的
		//可能是因为golang对于局部变量内部的结构体成员地址分配在另外的固定的地方
		//所以才需要对上边entry做一次深拷贝
		lbo := &raftlogpb.LBO{Offset: &offset, Entry: entry}

		b, err := proto.Marshal(lbo)
		if err != nil {
			return nil, err
		}
		//fmt.Println("append lbo entry:",lbo)
		lbos = append(lbos, lbo)


		//先将lbo序列化后的长度加在头部
		data := make([]byte, 4)
		binary.BigEndian.PutUint32(data, uint32(len(b)))
		///然后写数据
		data = append(data, b...)
		if _, err = file.fd.Write(data); err != nil {
			return nil, err
		}
	}

	return  lbos,nil
}

func (file *LogFile) read(offset int64) (lbos []*raftlogpb.LBO, err error) {
	fsize, _ := file.size()

	for offset < fsize {
		//需要先读出来头部的长度信息
		b := make([]byte, 4)
		_, err = file.fd.ReadAt(b, offset)
		size := binary.BigEndian.Uint32(b)

		data := make([]byte, size)
		//根据长度信息,读取剩余的数据
		offset += 4
		_, err = file.fd.ReadAt(data, offset)
		//将lbo对象数据序列化
		lbo := raftlogpb.LBO{}
		err=proto.Unmarshal(data, &lbo)
		lbos = append(lbos, &lbo)
		//fmt.Println("read offset:",offset)
		offset += int64(size)
	}
	return
}

func (file *LogFile) size() (size int64, err error) {
	fi, err := file.fd.Stat()
	if err != nil {
		return
	}
	size = fi.Size()
	return
}

func (file *LogFile) truncate(offset int64) {
	file.fd.Truncate(offset)
}

func (file *LogFile) close() {
	file.fd.Close()
}
