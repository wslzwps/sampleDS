package datanode

import (
	"fmt"
	"testing"
	"os"
)

func TestCreateChunkFailed(t *testing.T) {

}

//先调通，后整理case
func TestCreateChunkSuccess(t *testing.T) {
	return
	path := "/home/cintell/newds/src/datanode"
	chunk := &Chunk{Id: 1, Path: path, isActive: true}
	chunk.load()
	b1 := []byte("21232fdsafdsagsagfdfas")
	b2 := []byte("wslzwps")
	b3 := []byte("lizhi195")
	//size := len(b)

	fmt.Println(chunk.Write(b1, uint32(len(b1))))
	fid1, offset1, err := chunk.Write(b2, uint32(len(b2)))
	fmt.Println(fid1, offset1, err)
	fid2, offset2, err := chunk.Write(b3, uint32(len(b3)))
	fmt.Println(fid2, offset2, err)

	chunk.UnActive()

	fmt.Println(chunk.Write(b3, uint32(len(b3))))

	data, err := chunk.Read(fid1, offset1, uint32(len(b2)))
	fmt.Println(string(data), err)

	data, err = chunk.Read(fid2, offset2, uint32(len(b3)))
	fmt.Println(string(data), err)
	chunk.Close()
	os.Remove(chunk.Path)
}

func TestChunkWrite(t *testing.T) {

}

func TestChunkRead(t *testing.T) {

}

func TesChunkUnActive(t *testing.T) {

}

func TestSwitchChunk(t *testing.T) {

}
