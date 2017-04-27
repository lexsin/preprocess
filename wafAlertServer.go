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
	if err := http.ListenAndServe(addr, routerA); err != nil {
		panic("http listenAndServe " + addr + " error!")
	}
	return
}

type WafAlert struct {
	Client    string       `json:"Client"`
	Rev       string       `json:"Rev"`
	Msg       string       `json:"Msg"`
	Attack    string       `json:"Attack"`
	Severity  int32        `json:"Severity"`
	Maturity  int32        `json:"Maturityc"`
	Accuracy  int32        `json:"Accuracy"`
	Hostname  string       `json:"Hostname"`
	Uri       string       `json:"Uri"`
	Unique_id string       `json:"Unique_id"`
	Ref       string       `json:"Ref"`
	Tags      []string     `json:"Tags"`
	Rule      WafAlertRule `json:"Rule"`
	Version   string       `json:"Version"`
}

type WafAlertRule struct {
	File string `json:"File"`
	Line uint64 `json:"Line"`
	Id   uint64 `json:"Id"`
	Data string `json:"Data"`
	Ver  string `json:"Ver"`
}

type WafAlertObj struct {
	BackendObj
	Alert WafAlert `json:"Alert"`
}

func wafAlertWatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	content, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	mlog.Debug("get waf alert:", string(content))

	//check form
	//TODO parse json
	var wafAlert []WafAlertObj
	if err := json.Unmarshal([]byte(content), &wafAlert); err != nil {
		mlog.Error("waf alert form err!")
		Write(w, ErrOkErr, 10000)
		return
	}

	for _, alert := range wafAlert {
		altJson, err := json.Marshal(alert)
		if err != nil {
			continue
		}
		//push kafka
		data := &DataType{
			topicName: WafAlertTopic,
			handlePre: func(data []byte) ([]byte, error) {
				return data, nil
			},
			origiData: altJson,
			partition: 0,
		}
		if err := pushkafka.PushKafka(data); err != nil {
			mlog.Error("waf push kafka error:", err.Error())
		}
	}

	//response
	mlog.Debug("waf alert push kafka success!")
	Write(w, ErrOkErr, 10000)
	return
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
