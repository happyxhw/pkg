package fsutil

import (
	"crypto/md5" //nolint:gosec
	"encoding/hex"
	"io"
	"os"

	"github.com/happyxhw/pkg/util"
)

type File struct {
	Path  string
	MD5   string
	Size  int64
	Parts []*Part
}

type Part struct {
	MD5    string
	Size   int64
	Offset int64
}

func Split(in string, splitSize int64) (*File, error) {
	fi, err := os.Stat(in)
	if err != nil {
		return nil, err
	}
	var file File
	file.Path = in
	file.Size = fi.Size()
	if fi.Size() <= splitSize {
		file.MD5, err = GetFileMD5(in)
		if err != nil {
			return nil, err
		}
		file.Parts = append(file.Parts, &Part{
			MD5:  file.MD5,
			Size: file.Size,
		})
		return &file, nil
	}

	inFile, err := os.Open(in)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()
	var offset int64
	md5File := md5.New() //nolint:gosec
	for {
		buf := make([]byte, splitSize)
		n, err := inFile.Read(buf)
		if err == io.EOF {
			break
		}
		var p Part
		p.Offset = offset
		offset += int64(n)
		if n < int(splitSize) {
			buf = buf[:n]
		}
		p.Size = int64(len(buf))
		p.MD5 = util.Md5FromBytes(buf)
		file.Parts = append(file.Parts, &p)
		_, _ = md5File.Write(buf)
	}
	file.MD5 = hex.EncodeToString(md5File.Sum(nil))
	return &file, nil
}
