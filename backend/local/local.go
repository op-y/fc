package local

import (
	"io"
	"log"
	"os"
	"strings"
)

type LocalClient struct {
	name string
}

var Client *LocalClient

func init() {
	Client = &LocalClient{name: "local"}
}

func (lc *LocalClient) MkDir(objectKey string) error {
	seg := strings.Split(objectKey, "/")
	path := strings.Join(seg[0:len(seg)-1], "/")
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("fail to ensure path %s: %s", path, err.Error())
		return err
	}
	return nil
}

func (lc *LocalClient) PutObject(objectKey string, objectFile io.Reader) error {
	if err := lc.MkDir(objectKey); err != nil {
		log.Printf("failed to mkdir: %s", err.Error())
		return err
	}
	f, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("failed to open file: %s", err.Error())
		return err
	}
	defer f.Close()
	length, err := io.Copy(f, objectFile)
	if err != nil {
		log.Printf("failed to copy file: %s", err.Error())
		return err
	}
	log.Printf("upload file %s(%d bytes) to local storage", objectKey, length)
	return nil
}
