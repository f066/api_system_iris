package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"time"
)

type File struct {
	Name 	string
	Body	[]byte
	Size	int64
	Mode	int64
	ModTime	time.Time
}

type TarFile struct {
	files	[]File
	tw		*tar.Writer
	buff	bytes.Buffer
}

type GzFile struct {
	body 	[]byte
	modtime	time.Time
	buff 	bytes.Buffer
}

func (this *File) CalcSize() {
	this.Size = int64(len(this.Body))
}

func (this *TarFile) AddFile(file File) {
	if file.Size == 0 {
		file.CalcSize()
	}
	this.files = append(this.files,file)
	hdr := &tar.Header{
		Name:file.Name,
		Mode:file.Mode,
		ModTime:file.ModTime,
		Size:file.Size,
	}
	this.tw.WriteHeader(hdr)
	this.tw.Write(file.Body)
}

func (this *TarFile) AddFiles(files []File) {
	for _,file := range files{
		this.AddFile(file)
	}
}

func (this *TarFile) GetFile() []byte {
	this.tw.Close()
	return this.buff.Bytes()
}

func (this *GzFile) Set(body []byte) {
	this.body = body
}

func (this *GzFile) GetFile() []byte {
	gw := gzip.NewWriter(&this.buff)
	gw.ModTime = time.Now()
	gw.Write(this.body)
	gw.Close()
	return this.buff.Bytes()
}


func NewFile(name string,body []byte) File {
	var file File
	file.Name,file.Body,file.Mode,file.ModTime = name,body,0644,time.Now()
	file.CalcSize()
	return file
}

func NewTar() *TarFile {
	tf := new(TarFile)
	tf.tw = tar.NewWriter(&tf.buff)
	return tf
}