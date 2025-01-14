package util

import (
	"io"
	"net/http"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func OssGetObject(key string) (io.ReadCloser, error) {
	client, err := oss.New(os.Getenv("ALI_OSS_REGION")+".aliyuncs.com", os.Getenv("ALI_ACCESS_KEY_ID"), os.Getenv("ALI_ACCESS_KEY_SECRET"))
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(os.Getenv("ALI_OSS_ORIGIN"))
	if err != nil {
		return nil, err
	}
	return bucket.GetObject(key)
}

func OssPutObject(key string, f io.Reader) error {
	client, err := oss.New(os.Getenv("ALI_OSS_REGION")+".aliyuncs.com", os.Getenv("ALI_ACCESS_KEY_ID"), os.Getenv("ALI_ACCESS_KEY_SECRET"))
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(os.Getenv("ALI_OSS_ORIGIN"))
	if err != nil {
		return err
	}
	bucket.PutObject(key, f)
	return nil
}

func OssGetUrl(key string, options ...oss.Option) (string, error) {
	client, err := oss.New(os.Getenv("ALI_OSS_REGION")+".aliyuncs.com", os.Getenv("ALI_ACCESS_KEY_ID"), os.Getenv("ALI_ACCESS_KEY_SECRET"))
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(os.Getenv("ALI_OSS_ORIGIN"))
	if err != nil {
		return "", err
	}
	return bucket.SignURL(key, oss.HTTPGet, 3600, options...)
}

func OssDeleteObject(key string) error {
	client, err := oss.New(os.Getenv("ALI_OSS_REGION")+".aliyuncs.com", os.Getenv("ALI_ACCESS_KEY_ID"), os.Getenv("ALI_ACCESS_KEY_SECRET"))
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(os.Getenv("ALI_OSS_ORIGIN"))
	if err != nil {
		return err
	}
	return bucket.DeleteObject(key)
}

func OssGetObjectMeta(key string) (http.Header, error) {
	client, err := oss.New(os.Getenv("ALI_OSS_REGION")+".aliyuncs.com", os.Getenv("ALI_ACCESS_KEY_ID"), os.Getenv("ALI_ACCESS_KEY_SECRET"))
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(os.Getenv("ALI_OSS_ORIGIN"))
	if err != nil {
		return nil, err
	}
	return bucket.GetObjectMeta(key)
}
