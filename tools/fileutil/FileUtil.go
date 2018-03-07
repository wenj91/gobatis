package fileutil

import (
	"os"
	"log"
	"io"
)

func ReadFile(name string) ([]byte, error){
	file, err := os.Open(name)
	if nil != err {
		log.Printf("read file err:", err)
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	var res []byte
	for  {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		res = append(buf[:n])
	}

	return res, nil
}
