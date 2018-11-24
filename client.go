package vera 
import (
	"net/http"
	"crypto/tls"
	"time"
	"io/ioutil"
	"bytes"
)


func (c CmdType) MakeUri() { 
	switch Cmd.Do { 
		case "all", "list": //list all devices
			Cmd.Uri="http://vera.teamcornell.com:3480/data_request?id=user_data&output_format=xml&ns=1"
		case  "switch":  //toggle a device 
			Cmd.Uri=mkstr("http://%v:%v/data_request?id=action&output_format=xml&DeviceNum=%v&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&newTargetValue=%v", Cfg.Host,Cfg.Port,Cmd.Dev,Cmd.Value )
		// http://${HOST}:${PORT}/data_request?id=action&output_format=xml&DeviceNum=${1}&serviceId=urn:upnp-org:serviceId:SwitchPower1&action=SetTarget&newTargetValue=$NTVAL"`
		default: 
			DMsg(mkstr("invalid or unused command %v, so nothing to do.\n",c))

	}
		return 
}

func (c CmdType) Execute() (byteResults []byte, err error) {
	if empty(Cmd.Uri) { 
		Cmd.MakeUri()
	}
	var req *http.Request 
	if (len(Cmd.Body)>0) { 
		req, err = http.NewRequest("POST", Cmd.Uri, bytes.NewBuffer([]byte(Cmd.Body)))
	} else { 
		req, err = http.NewRequest("POST", Cmd.Uri, nil)
	}
    if err != nil {
        DMsg(mkstr("Cmd.Execute: Error creating request: %v!\n", err))
        return nil, err
    }

    //req.Header.Set("Accept", "application/json")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    var netClient = &http.Client{
        Timeout:   time.Second * 5,
        Transport: tr,
    }

    DTime("before netClient\n")
    response, err := netClient.Do(req)
    DTime("after netClient\n")
    if err != nil {
        DMsg(mkstr("netClient.Do: Error with netclient: %v!\n", err))
        return nil, err
    }

    switch response.StatusCode {
        case 200:
         DMsg("netclient Response was a code 200\n")
    default:
        response.Body.Close()
        ErrorExit(mkstr("Error with netclient connection\n%v",err),1)
        return
    }

    var body []byte
    if response.StatusCode == http.StatusOK {
        body, _ = ioutil.ReadAll(response.Body)
    }
    defer response.Body.Close()
    return body, err
}
