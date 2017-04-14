package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"preprocess/modules/mconfig"
	"preprocess/modules/mlog"
	"preprocess/modules/pushkafka"

	"github.com/julienschmidt/httprouter"
)

var routerA *httprouter.Router

func wafServInit() {
	route()
}

func route() {
	routerA = httprouter.New()
	routerA.POST("/alert/waf", wafAlertWatch)
}

func RunWafServer() {
	mlog.Debug("WAF-ALERT server start running...")
	port, _ := mconfig.Conf.String("server", "HttpPort")
	addr := ":" + port
	mlog.Debug("http listen on ", port, "...")
	http.ListenAndServe(addr, routerA)
	return
}

func wafAlertWatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	content, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	mlog.Debug("get waf alert:", string(content))

	//check form
	//TODO parse json

	//push kafka
	topicname, _ := mconfig.Conf.String("kafka", "WafAlertTopicName")
	data := &DataType{
		topicName: topicname,
		handlePre: func(data []byte) ([]byte, error) {
			return data, nil
		},
		origiData: content,
		partition: 0,
	}
	if err := pushkafka.PushKafka(data); err != nil {
		mlog.Error("waf push kafka error:", err.Error())
	}

	//response
	mlog.Debug("waf alert push kafka success!")
	Write(w, ErrOkErr, 10000)
}

var ErrOkErr = errors.New("success")

func Write(w http.ResponseWriter, err error, code int) {
	rsp := RespData{
		Code:    int32(code),
		Message: err.Error(),
	}
	writeContent, _ := json.Marshal(rsp)
	mlog.Debug(string(writeContent))
	w.Write(writeContent)
}

type RespData struct {
	Code    int32  `json:"code"`
	Message string `json:"msg"`
}
