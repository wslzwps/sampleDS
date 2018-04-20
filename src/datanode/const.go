package datanode

import (
	"errors"
)

//有时间需要好好组织下描述
var (
	UnActiveChunkWriteErr = errors.New("UnActive chunk can not write")
	UnMatchSizeErr        = errors.New("size is unmatch with buffer")
	MisUpChunkIdErr       = errors.New("mix up fid")
)
