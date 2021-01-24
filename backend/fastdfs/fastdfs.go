package fastdfs

import (
	"errors"
	"io"
	"log"
	"strings"

    "github.com/op-y/fc/config"
    "github.com/op-y/fc/utils"

	"github.com/op-y/weilaihui/fdfs_client"
)

// interface
type Sizer interface {
	Size() int64
}

type FdfsClient struct {
	client *fdfs_client.FdfsClient
}

var Client *FdfsClient

func init() {
	Client = new(FdfsClient)
}

func (fdfsc *FdfsClient) GetClient() error {
	if fdfsc.client != nil {
		return nil
	}
	c, err := fdfs_client.NewFdfsClient(config.Cfg.FastDfs.FdfsConf)
	if err != nil {
		return err
	}
	fdfsc.client = c
	return nil
}

func (fdfsc *FdfsClient) PutObject(objectKey string, objectFile io.Reader) error {
	seg := strings.Split(objectKey, "/")
	filename := seg[len(seg)-1]
	ext := utils.GetExt(filename)
	fileSizer, ok := objectFile.(Sizer)
	if !ok {
		log.Printf("objectFile object miss Size() method")
		return errors.New("No Size method")
	}
	buffer := make([]byte, fileSizer.Size())
	_, err := objectFile.Read(buffer)
	if err != nil {
		log.Printf("failed to read file: %s", err.Error())
		return err
	}

	// ignore the RemoteFileID
	_, err = fdfsc.Upload(buffer, ext)
	if err != nil {
		return err
	}
	return nil
}

func (fdfsc *FdfsClient) Upload(buf []byte, ext string) (string, error) {
	if err := fdfsc.GetClient(); err != nil {
		log.Printf("fail to get Fast DFS client: %s", err.Error())
		return "", err
	}
	response, err := fdfsc.client.UploadByBuffer(buf, ext)
	if err != nil {
		log.Printf("fail to upload file to Fast DFS: %s", err.Error())
		return "", err
	}
	id := response.RemoteFileId
	return id, nil
}

func (fdfsc *FdfsClient) Delete(id string) error {
	if err := fdfsc.GetClient(); err != nil {
		log.Printf("fail to get Fast DFS client: %s", err.Error())
		return err
	}
	if err := fdfsc.client.DeleteFile(id); err != nil {
		log.Printf("fail to delete file from Fast DFS: %s", err.Error())
		return err
	}
	return nil
}
