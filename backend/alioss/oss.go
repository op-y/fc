package alioss

import (
    "io"
    "log"

    "github.com/op-y/fc/config"

    "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
    client *oss.Client
    bucket *oss.Bucket
}

var Client *OssClient

func init() {
    Client = new(OssClient)
}

func (oc *OssClient) GetClient() error {
    if oc.client != nil {
        return nil
    }
    c, err := oss.New(config.Cfg.Oss.Endpoint, config.Cfg.Oss.AccessKeyId, config.Cfg.Oss.AccessKeySecret)
    if err != nil {
        return err
    }
    oc.client = c
    return nil
}

func (oc *OssClient) GetBucket() error {
    b, err := oc.client.Bucket(config.Cfg.Oss.Bucket)
    if err != nil {
        log.Printf("failed to get bucket %s: %s", config.Cfg.Oss.Bucket, err.Error())
        return err
    }
    oc.bucket = b
    return nil
}

func (oc *OssClient) PutObject(objectKey string, objectFile io.Reader) error {
    if err := oc.GetClient(); err != nil {
        log.Printf("failed to get cilent")
        return err
    }
    if err := oc.GetBucket(); err != nil {
        log.Printf("failed to get bucket")
        return err
    }
    if err := oc.bucket.PutObject(objectKey, objectFile); err != nil {
        log.Printf("failed to put object %s: %s", objectKey, err.Error())
        return err
    }
    log.Printf("upload %s to Aliyun OSS", objectKey)
    return nil
}
