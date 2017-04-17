package main

import (
	"bufio"
	//"crypto/md5"
	"errors"
	"fmt"
	"os"
	"preprocess/modules/mlog"
	"preprocess/modules/xdrParse"
	"time"
)

func saveToCephPerXdr(obj *xdrParse.DpiXdr) error {
	jtype := obj.CheckType()
	mlog.Debug("jtype=", jtype)
	switch jtype {
	case xdrParse.XdrType:
		//xdrTypeToCeph()
	case xdrParse.XdrHttpType:
		if err := xdrHttpTypeToCeph(obj); err != nil {
			mlog.Error("httpxdr Write ceph error:", err.Error())
			return err
		}
	case xdrParse.XdrFileType:
		if err := xdrFileTypeToCeph(obj); err != nil {
			mlog.Error("filexdr Write ceph error:", err.Error())
			return err
		}
	default:
		mlog.Error("CheckType error! return ", jtype)
		return errors.New(fmt.Sprintf("CheckType error! return %d", jtype))
	}
	mlog.Debug("save to ceph success!")
	return nil
}

func saveToCeph(objlist []*xdrParse.DpiXdr) error {
	for _, obj := range objlist {
		saveToCephPerXdr(obj)
	}
	mlog.Debug("save to ceph success!")
	return nil
}

func xdrHttpTypeToCeph(data *xdrParse.DpiXdr) error {
	//write file
	httpResp := data.HttpRespInfo
	httpReq := data.HttpReqInfo
	respFileName := createFilenameByMd5(httpResp)
	reqFileName := createFilenameByMd5(httpReq)
	rootPath := "/cephfs/DPI/http"
	path := createPathByTime()
	fullPath := rootPath + "/" + string(path)
	if exist, _ := IsExist(string(fullPath)); !exist {
		//create dir
		mlog.Debug(string(fullPath), " not exist")
		if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
			mlog.Error("MkdirAll dir:", fullPath, err.Error())
		}
		mlog.Debug(string(fullPath), " created!")
	}
	httprespFile := fullPath + "/" + string(respFileName)
	httpreqFile := fullPath + "/" + string(reqFileName)
	if err := wirteFile(httprespFile, httpResp); err != nil {
		return err
	}
	if err := wirteFile(httpreqFile, httpReq); err != nil {
		return err
	}

	//modify object
	data.HttpRespInfo = []byte(httprespFile)
	data.HttpReqInfo = []byte(httpreqFile)

	//del other unnormal big data
	if len(data.FileContent) != 0 {
		mlog.Warning("len(FileContent)=", len(data.FileContent))
		mlog.Warning("HttpRespInfo HttpReqInfo FileContent both have data")
		data.FileContent = nil
	}

	return nil
}

func xdrFileTypeToCeph(data *xdrParse.DpiXdr) error {
	//write file
	mlog.Alert("cephtime1=", time.Now().Unix())
	content := data.FileContent
	fileName := createFilenameByMd5(content)
	mlog.Alert("cephtime2=", time.Now().Unix())
	rootPath := "/cephfs/DPI/file"
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
	mlog.Alert("cephtime3=", time.Now().Unix())
	//modify obj
	data.FileContent = []byte(fullFile)

	//del other big data
	if len(data.HttpReqInfo) != 0 {
		mlog.Warning("len(HttpReqInfo)=", len(data.HttpReqInfo))
		mlog.Warning("HttpRespInfo HttpReqInfo FileContent both have data")
		data.HttpReqInfo = nil
	}
	if len(data.HttpRespInfo) != 0 {
		mlog.Warning("len(HttpRespInfo)=", len(data.HttpRespInfo))
		mlog.Warning("HttpRespInfo HttpReqInfo FileContent both have data")
		data.HttpRespInfo = nil
	}
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

func createFilenameByMd5(data []byte) string {
	sum := Md5Sum(data)
	return fmt.Sprintf("%x", sum)
}
