package datanode

import (
	"testing"
	"os"
	"fmt"
	"hash/crc32"
	"encoding/binary"
)

func TestDataManager_Write(t *testing.T) {
	dm,err:=NewDataManager("/home/cintell/newds/src/datanode/data")
	//fmt.Println(dm.dir)
	if err!=nil{
		t.Error(err.Error())
		return
	}

	out_fp,err:=os.OpenFile("key.txt",os.O_RDWR|os.O_APPEND,0666)
	if err!=nil{
		t.Error(err.Error())
		return
	}
	defer  out_fp.Close()


	b:=make([]byte,4*1024*1024)
	for i:=0;i<1024*10;i++{
		binary.BigEndian.PutUint32(b[0:4],uint32(i))
		fid,offset,_:=dm.Write(b,uint32(len(b)))
		//fid_offset_size_crc
		key:=fmt.Sprintf("%d_%d_%d_%d\n",fid,offset,len(b),crc32.ChecksumIEEE(b))
		out_fp.WriteString(key)
	}
}