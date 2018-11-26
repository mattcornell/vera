package vera
// config.go is the most rudimentary parser 

import (
	"os"
	//"fmt"
	"time"
	"bufio"
	"github.com/BurntSushi/toml"
)

var Cfg cfgType = cfgType{File:".config.toml"}

func (c cfgType) isfresh() bool {
	 if RefreshOpt  { return false }
 // if this is a passive request and data is fresh 
	 if  ( ! RefreshOpt && (Cmd.Do=="all"||Cmd.Do=="list"||Cmd.Do=="details") && ((time.Now().Unix()-Cfg.Lastpull) <  Cfg.Refresh)) {
		 return true
	 }
	 return false
}

type cfgType struct {
	//Host string
	//Port string
	//File string
	File string
	Cache string
	host string
	Host string
	Port string
	//LastPull  int64
	Lastpull int64
	Refresh  int64
}

func WriteCfg (){ 
	file,err := os.Create(Cfg.File) 
	if err != nil { panic(err) }
	f := bufio.NewWriter(file)
	_,err = f.WriteString(mkstr("Host=%q\nPort=%q\nRefresh=%v\nLastpull=%v\nCache=%q\n",
		Cfg.Host, Cfg.Port, Cfg.Refresh, Cfg.Lastpull,Cfg.Cache))
	if err != nil {
        ErrorExit(mkstr("failed writing to config file: %s", err),1)
    }
	f.Flush()
}

func ReadCfg () {
	if _,err := toml.DecodeFile(Cfg.File,&Cfg)
	err != nil { 
		ErrorExit(mkstr("Config file %v not found",Cfg.File),1)
	}
	DWhatsThis(Cfg.Host)
	DWhatsThis(Cfg.Port)
	DWhatsThis(Cfg.Refresh)
	DWhatsThis(Cfg.Lastpull)
	DWhatsThis(Cfg.Cache)
	dPause(mkstr("Found Cfg.Refresh=%v",Cfg.Refresh))

	//if empty(Cfg.Host) { Cfg.Host = Cfg.host }
	//if empty(Cfg.Port) { Cfg.Host = Cfg.port }
	//if empty(Cfg.LastPull) { Cfg.LastPull = 0 }
	//if empty(Cfg.refresh) { Cfg.refresh = 180 }
	if empty(Cfg.Cache) { Cfg.Cache =".lastpull" }
	if empty(Cfg.Lastpull) { Cfg.Lastpull = 0 }
	if empty(Cfg.Refresh) { Cfg.Refresh = 180 }
	return
}
