package vera
// config.go is the most rudimentary parser 

import (
	"github.com/BurntSushi/toml"
)

var Cfg cfgType = cfgType{File:"config.toml"}

type cfgType struct {
	/*uri="http://vera.teamcornell.com:3480/data_request?id=user_data&output_format=xml&ns=1"
	host="vera.teamcornell.com"
	port="3480"*/
	Host string
	Port string
	Uri  string
	File string
	host string
	port string
	uri  string
	lastPull  int64

}


func ReadCfg () { 
	if _,err := toml.DecodeFile(Cfg.File,&Cfg)
	err != nil {
		panic(err)
	}
	if empty(Cfg.Host) { Cfg.Host = Cfg.host }
	if empty(Cfg.Port) { Cfg.Host = Cfg.port }
	if empty(Cfg.Uri) { Cfg.Host = Cfg.uri }

	//Cfg.Host=Cfg.host
	//Cfg.Port=Cfg.port
	//Cfg.Uri=Cfg.uri
	return

}
