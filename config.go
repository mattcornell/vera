package vera
// config.go is the most rudimentary parser 

import (
	"os"
	//"fmt"
	//"time"
	"bufio"
	"github.com/BurntSushi/toml"
)

var Cfg cfgType = cfgType{File:".config.toml"}


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
	_,err = f.WriteString(mkstr("host=%q\nport=%q\nrefresh=%v\nlastpull=%v\ncache=%q\n",
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
	if Empty(Cfg.Cache) { Cfg.Cache =".lastpull" }
	if Empty(Cfg.Lastpull) { Cfg.Lastpull = 0 }
	if Empty(Cfg.Refresh) { Cfg.Refresh = 180 }
	return
}
