package vera
// config.go is the most rudimentary parser 

import (
	"os"
	//"fmt"
	"bufio"
	"github.com/BurntSushi/toml"
)

var Cfg cfgType = cfgType{File:".config.toml",Cache:".lastpull"} 

type cfgType struct {
	/*uri="http://vera.teamcornell.com:3480/data_request?id=user_data&output_format=xml&ns=1"
	host="vera.teamcornell.com"
	port="3480"*/
	Host string
	Port string
	Uri  string
	File string
	Cache string
	host string
	port string
	uri  string
	LastPull  int64
}

func WriteCfg (){ 
	//file,err := os.Open(Cfg.File) 
	file,err := os.Open("myoutput") 
	if err != nil { panic(err) }
	f := bufio.NewWriter(file)
		defer file.Close()
		f.WriteString(mkstr("host=%q\n",Cfg.Host))
		f.WriteString(mkstr("port=%q\n",Cfg.Port))
		f.WriteString(mkstr("uri=%q\n",Cfg.Uri))
		f.WriteString(mkstr("lastPull=%v\n",Cfg.LastPull))
	dPause(mkstr("I tried to write the file %v\n",Cfg.File))
		f.Flush()
}

func ReadCfg () { 
	if _,err := toml.DecodeFile(Cfg.File,&Cfg)
	err != nil {
		panic(err)
	}
	if empty(Cfg.Host) { Cfg.Host = Cfg.host }
	if empty(Cfg.Port) { Cfg.Host = Cfg.port }
	if empty(Cfg.Uri) { Cfg.Host = Cfg.uri }
	if empty(Cfg.LastPull) { Cfg.LastPull = 0 }

	//Cfg.Host=Cfg.host
	//Cfg.Port=Cfg.port
	//Cfg.Uri=Cfg.uri
	return

}
