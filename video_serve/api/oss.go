package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/wonderivan/logger"
	"testgin/api/config"
)

var client *oss.Client
var err error

func init() {
	if config.Set.OpenOssServe && config.Set.EndPoint != "" && config.Set.AccessKeyID != "" && config.Set.AccessKeySecret != "" {
		client, err = oss.New(config.Set.EndPoint, config.Set.AccessKeyID, config.Set.AccessKeySecret)
		if err != nil {
			logger.Emer(err)
		}
	} else {
		url := "https://common-buy.aliyun.com/?spm=5176.7933691.1309840.1.5c072a66XIyKNT&commodityCode=ossbag#/buy"
		logger.Warn("Oss Serve Not Open.if You Want Oss Please Open The Address:%v", url)
	}

}

func UploadToOss(filename string, path string, bucketName string) bool {

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logger.Emer("Bucket Connect Failed: %v", err)
		return false
	}

	err = bucket.UploadFile(filename, path, 1024*500, oss.Routines(3))
	//err = bucket.PutObjectFromFile(filename, "./videos/ddd")
	if err != nil {
		logger.Emer("Upload to Oss Failed: %v", err)
		return false
	}
	return true

}

func DeleteOssFile(filename string, bucketName string) bool {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logger.Emer("Bucket Connect Failed: %v", err)
		return false
	}
	err = bucket.DeleteObject(filename)
	if err != nil {
		logger.Emer("Deleting Oss Object Failed: %v", err)
		return false
	}
	return true
}
