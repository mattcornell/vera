package vera 
import (
	"net/http"
	"crypto/tls"
	"time"
	"io/ioutil"
	"bytes"
)

func fetchUrl(c CmdType) (byteResults []byte, err error) {
	var req *http.Request 
	if (len(c.Body)>0) { 
		req, err = http.NewRequest("POST", c.Uri, bytes.NewBuffer([]byte(c.Body)))
	} else { 
		req, err = http.NewRequest("POST", c.Uri, nil)
	}
    if err != nil {
        DMsg(mkstr("fetchUrl: there was an error creating request: %v!\n", err))
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
        DMsg(mkstr("netClient.Do: there was an error: %v!\n", err))
        return nil, err
    }

    switch response.StatusCode {
        case 200:
         DMsg("netclient Response was a code 200\n")
    default:
        response.Body.Close()
        errorExit(mkstr("Error with netclient connection\n%v",err),1)
        return
    }

    var body []byte
    if response.StatusCode == http.StatusOK {
        body, _ = ioutil.ReadAll(response.Body)
    }
    defer response.Body.Close()
    return body, err
}
