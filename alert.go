package main

type VdsFullAlert struct {
	BackendObj
	Alert VdsAlert `json:"Alert"`
}

type VdsAlert struct {
	Threatname       string `json:threatname`
	Subfile          string `json:subfile`
	Local_threatname string `json:local_threatname`
	Local_vtype      string `json:local_vtype`
	Local_platfrom   string `json:local_platfrom`
	Local_vname      string `json:local_vname`
	Local_extent     string `json:local_extent`
	Local_enginetype string `json:local_enginetype`
	Local_logtype    string `json:local_logtype`
	Local_engineip   string `json:local_engineip`
	Log_time         uint64 `json:log_time`
}

type IdsAlert struct {
	Time        uint64 `json:time`
	Src_ip      string `json:src_ip`
	Src_port    uint32 `json:src_port`
	Dest_ip     string `json:dest_ip`
	Dest_port   uint32 `json:dest_port`
	Proto       uint32 `json:proto`
	Attack_type string `json:attack_type`
	Details     string `json:details`
	Severity    uint32 `json:severity`
	Engine      string `json:engine`
}
