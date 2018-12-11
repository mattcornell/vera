package vera
// get data from cache file or netClient  
import (
	"net/http"
	"crypto/tls"
	"time"
	"io/ioutil"
	"bytes"
	"os"
	"strconv"
)

var err error

func (c cfgType) NeedRefresh() bool {
    if OfflineOpt  {
		Xml,err=ReadCache()
		if err != nil { return true  }
	}
    if RefreshOpt  {
		return true
	}
	Xml,err=ReadCache()
	if err != nil { return true  }
    if  ( ! RefreshOpt &&
		 /*(Cmd.Do=="all"||
		 Cmd.Do=="list"||
		 Cmd.Do=="details"||
		 Cmd.Do=="status") && */
		 ((time.Now().Unix()-Cfg.Lastpull) <  Cfg.Refresh)) {
         return false
     }
     return true
}
func RefreshAfterCommand() { 
		RefreshOpt = true
		DoRefresh()
		Populate()
		WriteCfg()
}

func DoRefresh() {
	if Cfg.NeedRefresh() {
		var c CmdType
			c.Uri=mkstr("http://%v:%v/data_request?id=user_data&output_format=xml&ns=1",Cfg.Host,Cfg.Port)
			c.Do="list"
			c.Fetch()
	} else { 
		DMsg( "Using cache file\n" )

	}
	err := Cmd.WriteTemp(Xml)
    if err != nil {
        ErrorExit(mkstr("Error writing temp data %v", err), 1)
    }
}

func (c CmdType) MakeUri() CmdType {
	switch c.Do {
		case "all","list","rooms","room","scene","scenes","user","users": //list all devices
			c.Uri=mkstr("http://%v:%v/data_request?id=user_data&output_format=xml&ns=1",Cfg.Host,Cfg.Port)
		case  "off":  //turn on device 
			if Empty(c.Dev) { ErrorExit("Missing device number",1) }
			if Empty(c.Value) { ErrorExit("Missing device value",1) }
			var d Devices = Data.DeviceList
			dev := d.Match(SecondArg())
			c.Dev=strconv.Itoa(dev.Id)
			if (dev.Value() == "0") {
				 ErrorExit(mkstr("device %v is already off",dev.Name),0)
			}
			c.Value="0"
			c.Uri=mkstr("http://%v:%v/data_request?id=action&output_format=xml&DeviceNum=%v&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&newTargetValue=%v", Cfg.Host,Cfg.Port,c.Dev,c.Value )
		case  "on":  //turn on device 
			if Empty(c.Dev) { ErrorExit("Missing device number",1) }
			if Empty(c.Value) { ErrorExit("Missing device value",1) }
			var d Devices = Data.DeviceList
			dev := d.Match(SecondArg())
			if (dev.Value() == "1") { 
				ErrorExit(mkstr("device %v is already on",dev.Name),0) }
			c.Value="1"
			c.Dev=strconv.Itoa(dev.Id)
			c.Uri=mkstr("http://%v:%v/data_request?id=action&output_format=xml&DeviceNum=%v&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&newTargetValue=%v", Cfg.Host,Cfg.Port,c.Dev,c.Value )
		case  "lock", "unlock":  //toggle a device 
			if Empty(c.Dev) { ErrorExit("Missing device number",1) }
			var d Devices = Data.DeviceList
			dev := d.Match(SecondArg())
			c.Dev=strconv.Itoa(dev.Id)
			var v string
			if ((dev.Value() == "0")&&(c.Do =="unlock")) {
				ErrorExit(mkstr("device %v is already open",dev.Name),0) }
			if ((dev.Value() == "1")&&(c.Do =="lock")) {
				ErrorExit(mkstr("device %v is already locked",dev.Name),0) }
			if (c.Do == "lock") { v="1" } else { v="0" }
			c.Uri=mkstr("http://%v:%v/data_request?id=action&DeviceNum=%v&serviceId=urn:micasaverde-com:serviceId:DoorLock1&action=SetTarget&newTargetValue=%v",Cfg.Host,Cfg.Port,c.Dev,v)
		case "status","value": //get device value
			if Empty(c.Dev) { ErrorExit("Missing device ",1) }
				var d Devices = Data.DeviceList
				c.Dev=strconv.Itoa(	d.Match(SecondArg()).Id)
				c.Uri=mkstr("http://%v:%v/data_request?id=user_data&output_format=xml&ns=1",Cfg.Host,Cfg.Port)
		case  "switch","toggle":  //toggle a device 
			if Empty(c.Dev) { ErrorExit("Missing device ",1) }
				var d Devices = Data.DeviceList
				c.Dev=strconv.Itoa(	d.Match(SecondArg()).Id)
			c.Uri=mkstr("http://%v:%v/data_request?id=action&DeviceNum=%v&serviceId=urn:micasaverde-com:serviceId:HaDevice1&action=ToggleState", Cfg.Host,Cfg.Port,c.Dev)
		default:
			ErrorExit(mkstr("invalid or unused command %q, so nothing to do.\n",c.Do),1)
	}
		return c
}

func errorChk (e error) {
	if e !=nil  {
		panic(e)
	}
}

func (c CmdType) WriteTemp(b []byte) (err error)  {
		f, err := os.Create(Cfg.Cache)
		errorChk(err)
		_, err = f.Write(b)
		errorChk(err)
		Cfg.Lastpull = time.Now().Unix()
		f.Sync()
		f.Close()
		return
}

var cacheBuff []byte  = nil
func ReadCache() ([]byte,error){
		return ioutil.ReadFile(Cfg.Cache)
}
func (c CmdType) Fetch() (err error) {
	if Empty(c.Uri) {
		c=Cmd
		c.MakeUri()
	}
	var req *http.Request
	if (len(c.Body)>0) {
		req, err = http.NewRequest("POST", c.Uri, bytes.NewBuffer([]byte(c.Body)))
	} else { 
		req, err = http.NewRequest("POST", c.Uri, nil)
	}
	if  err != nil { 
		ErrorExit(mkstr("Err posting to c.Uri %v",c.Uri),1) 
	}
    errorChk(err)

    //req.Header.Set("Accept", "application/json")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    var netClient = &http.Client{
        Timeout: time.Second * 5,
        Transport: tr,
    }

    DTime(mkstr("\nstart netClient command %v %v\n",c.Do,c.Uri))
    response, err := netClient.Do(req)
    DTime("done netClient\n")
    errorChk(err)

    switch response.StatusCode {
        case 200:
         DMsg("netclient Response was a code 200\n")
    default:
        response.Body.Close()
        ErrorExit(mkstr("Error with netclient connection\n%v",err),1)
        return
    }

    //var body []byte
    if response.StatusCode == http.StatusOK {
        Xml, _ = ioutil.ReadAll(response.Body)
    }
    defer response.Body.Close()
    //return body, err
    return err
}
