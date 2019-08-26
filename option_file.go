package gobatis

import (
	"io/ioutil"
	"os"
)

type FileOption struct {
	path string
}

var _ IOption = &FileOption{}

// NewFileOption db config file path, default: db.yml
func NewFileOption(pt ...string) *FileOption {
	path := "db.yml"
	if len(pt) > 0 {
		path = pt[0]
	}
	return &FileOption{
		path: path,
	}
}

func (f *FileOption) Type() OptionType {
	return OptionTypeFile
}

func (f *FileOption) ToDBConf() *DBConfig {
	file, err := os.Open(f.path)
	if nil != err {
		panic("Open db conf err:" + err.Error())
	}

	r, err := ioutil.ReadAll(file)
	if nil != err {
		panic("Read db conf err:" + err.Error())
	}

	dbConf := buildDbConfig(string(r))
	return dbConf
}
