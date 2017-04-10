package main

import (
	"bufio"
	//"crypto/md5"
	"fmt"
	"os"
	"preprocess/modules/mlog"
	"preprocess/modules/xdrParse"
	"time"
)

func saveToCeph(objlist []*xdrParse.DpiXdr) error {
	mlog.Debug("111")
	for _, obj := range objlist {
		jtype := obj.CheckType()
		switch jtype {
		case xdrParse.XdrType:
			//xdrTypeToCeph()
		case xdrParse.XdrHttpType:
			xdrHttpTypeToCeph(obj)
		case xdrParse.XdrFileType:
			xdrFileTypeToCeph(obj)
		default:
			mlog.Error("CheckType error! return ", jtype)
		}
	}
	mlog.Debug("save to ceph success!")
	return nil
}

func xdrHttpTypeToCeph(data *xdrParse.DpiXdr) error {
	//write file
	httpResp := data.HttpRespInfo
	httpReq := data.HttpReqInfo
	respFileName := Md5Sum(httpResp)
	reqFileName := Md5Sum(httpReq)
	rootPath := "DPI/http"
	path := createPathByTime()
	fullPath := rootPath + "/" + string(path)
	if exist, _ := IsExist(string(fullPath)); !exist {
		//create dir
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			mlog.Error("MkdirAll dir:", fullPath, err.Error())
		}
	}
	httprespFile := fullPath + "/" + string(respFileName)
	httpreqFile := fullPath + "/" + string(reqFileName)
	wirteFile(httprespFile, httpResp)
	wirteFile(httpreqFile, httpReq)

	//modify obj
	data.HttpRespInfo = []byte(httprespFile)
	data.HttpReqInfo = []byte(httpreqFile)
	return nil
}

func xdrFileTypeToCeph(data *xdrParse.DpiXdr) error {
	//write file
	content := data.FileContent
	fileName := Md5Sum(content)
	rootPath := "DPI/file"
	path := createPathByTime()
	fullPath := rootPath + "/" + string(path)
	if exist, _ := IsExist(string(fullPath)); !exist {
		//create dir
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			mlog.Error("MkdirAll dir:", fullPath, err.Error())
		}
	}
	fullFile := fullPath + "/" + string(fileName)
	wirteFile(fullFile, content)

	//modify obj
	data.FileContent = []byte(fullFile)
	return nil
}

func wirteFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		mlog.Error("create file", file, err.Error())
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		mlog.Error("write file", filename, err.Error())
		return err
	}
	writer.Flush()
	return nil
}

/**
 * check if file or dir exist
 */
func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createPathByTime() string {
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	path := fmt.Sprintf("%d/%d/%d", year, month, day)
	return path
}
