package vera
import (
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
	"time"
	"strconv"
	"strings"
)

func Empty(object interface{}) bool {
    //First check normal definitions of empty
	//return true
    if object == nil {
        return true
    } else if object == "" {
        return true
    } else if object == false {
        return true
    }
    return false
}

const(
      // personal preference for date format
      dateTxt string = "2006-01-01_150405.00000"
	  dateNice string = "2006-01-01 15:04:05"
)

var mkstr = fmt.Sprintf 

var CatNames = []string{"", "NA Interface", "Dimmable_Light", "Switch", "Security_Sensor", "HVAC", "Camera", "Door_Lock", "Window_Convering", "Remote_Control", "IR_Transmitter", "Generic_I_O", "Generic_Sensor", "Serial_Port", "Scene_Controller", "AV", "Humidity_Sensor", "Temperature_Sensor", "Light_Sensor", "Z-Wave_Interface", "Insteon_Interface", "Power_Meter", "Alarm_Panel", "Alarm_Partition", "Siren", "Weather", "Philips_Controller", "Appliance", "UV_Sensor", "Mouse_Trap", "Doorbell", "Keypad"}

var Data VeraRoot
var Xml []byte

func Populate () {
	err := xml.Unmarshal(Xml, &Data)
	if err != nil { ErrorExit(mkstr("Error getting xml root data %v",err),1) }
	return
}

type VeraRoot struct {
	XMLName          xml.Name `xml:"root"`
	Timezone         string   `xml:"timezone,attr"`
	FirmwareVersion string    `xml:"firmware_version,attr"`
	CityDescription string    `xml:"City_description,attr"`
	Model            string   `xml:"model,attr"`
	DeviceList          []Device `xml:"devices>device"`
	Scenes           []Scene  `xml:"scenes>scene"`
	Users            []User   `xml:"users>user"`
	RoomList            []Room   `xml:"rooms>room"`
}

type Devices []Device
type Rooms []Room
type Users []User
type Scenes []Scene

func (l Rooms) Match(match string) (r Rooms)  { 
    if _ , err := strconv.Atoi(match); err == nil {
		// int room id, return a match
		for _, v := range l { 
			if v.Id == match {
				r = append(r,v)
			   return r 
			}
		}
		ErrorExit(mkstr("No matching room: %v",match),1) 
		return r 
	} else { 
		//return room with matching Name
		for _, v := range l { 
			if (strings.ToUpper(v.Name)==strings.ToUpper(match)) {
				r = append(r,v)
			   return r
			}
		}
		return  r 
}
}
func (l Devices) Matches ( match string ) (d Devices) { 
    if _, err := strconv.Atoi(match); err == nil {
		//look for Int Dev ID
	} else { 
		for _, this := range l { 
			if (strings.Contains(strings.ToUpper(this.Name),strings.ToUpper(match)) ){ 
				   d = append(d,this)
			}
		}
		return d
	}
	ErrorExit(mkstr("No matching device: %v",match),1) 
	return d
}

func (l Devices) Match(match string) (d Device) { 
    if _, err := strconv.Atoi(match); err == nil {
		// match int dev num
		for _, v := range l {
			c,_:= strconv.Atoi(match)
			if v.Id==c  {
				return v
			}
		} 
	} else { 
		//return first matching name 
		for _, v := range l { 
			if (strings.ToUpper(v.Name)==strings.ToUpper(match)) {
				return v
			} 
		}//end of range over l 
	} //end of if int
	return d
}

func (l VeraRoot) DevMatches(match string) (d Device) { 
    if _, err := strconv.Atoi(match); err == nil {
		//look for Int Dev ID
		return l.DevId(match)
	} else { 
		//look for String Dev ID
		return l.DevMatchesName(match)
	}
	ErrorExit(mkstr("No matching device: %v",match),1) 
	return d
}

func (l VeraRoot) DevMatchesName(match string) (d Device) { 
	for _, v := range Data.DeviceList { 
		if (strings.ToUpper(v.Name)==strings.ToUpper(match)) {
		   return v
		}
	}
	ErrorExit(mkstr("No matching device: %v",match),1) 
	return d
}

func (l VeraRoot) DevId(id string) (d Device) { 
	for _, v := range Data.DeviceList {
		c,_:= strconv.Atoi(id)
		if v.Id==c  {
			return v
		}
	}
	return d
}

func (l VeraRoot) DevsContainsName(match string) (r Devices) { 
	for _, v := range l.DeviceList { 
		if (strings.Contains(strings.ToUpper(v.Name),strings.ToUpper(match)) ){ 
	       r = append(r,v)
		}
	}
	return r
}

func (l VeraRoot) DevContainsName(match string) (d Device) { 
	for _, v := range l.DeviceList { 
		if (strings.Contains(strings.ToUpper(v.Name),strings.ToUpper(match)) ){ 
	       return v
		}
	}
	return  d 
}
func (l VeraRoot) DevsId(id string) ( d Devices ) { 
	for _, v := range Data.DeviceList {
		c,_:= strconv.Atoi(id)
		if v.Id==c  {
			d=append(d,v)
		}
	}
	return d
}

func (d Device) RoomName() string {
	for _, v := range Data.RoomList {
		if v.Id == d.RoomNum {
			return v.Name
		}
	}
	return ""
}

type Room struct {
	XMLName xml.Name `xml:"room"`
	/*id="1373941" Name="buster" Level="1" IsGuest="0"></user>*/
	Id      string `xml:"id,attr"`
	Section string `xml:"section,attr"`
	Name    string `xml:"name,attr"`
}

type User struct {
	XMLName xml.Name `xml:"user"`
	Id      string   `xml:"id,attr"`
	Name    string   `xml:"Name,attr"`
	Level   int      `xml:"Level,attr"`
}

type Scene struct {
	XMLName           xml.Name `xml:"scene"`
	Id                int      `xml:"id,attr"`
	Timestamps        int64    `xml:"Timestamp,attr"`
	Name              string   `xml:"name,attr"`
	Room              string   `xml:"room,attr"`
	TriggersOperator string   `xml:"triggers_operator,attr"`
	users             string   `xml:"users,attr"`
	Paused            string   `xml:"paused,attr"`
	ModeStatus        int      `xml:"modeStatus,attr"`
	RawLastRun        int64    `xml:"last_run,attr"`
	Trigger           Trigger  `xml:"triggers>trigger"`
}

func epochDate(s int64) string { 
	t:= time.Unix(s,0).Format(dateNice)
   return t
}

func (s Scene) LastRun () string{ 
	//return  mkstr("%s",time.Unix(s.RawLastRun,0))
	return  epochDate(s.RawLastRun)
}

type Trigger struct {
	XMLName  xml.Name `xml:trigger`
	Name     string   `xml:"name,attr"`
	Enabled  int      `xml:"enabled,attr"`
	Device   int      `xml:"device,attr"`
	LastRun int64    `xml:"last_run"`
}

//method
func (d Device) Category() string {
	return CatNames[d.CategoryNum]
}

func (d Device) value (V string) (r string) {
	for _, v := range d.States {
		if v.Variable == V {
			return v.Value
		}
	}
	return ""
}

func (d Device) Value() (string) { 
	switch d.CategoryNum {
	case 2: //dimmable light
		return d.value("LoadLevelStatus")
	case 3: //lights, switches
		return d.value("Target") 
	case 4: // security device (motion, window sensor)
		return d.value("Armed")
	case 5: //HVAC
		return d.value("BatteryLevel")
	case 6: //cameras
		return "" //just null out the status
	case 7: // door lock
		return d.value("Status") 
	case 11: // Generic IO
		return d.value("ArmedTripped") 
	case 14: // remote controller
		return ""  //TODO, return battery level? Last seen?
	case 16: // humidity level
		return d.value("CurrentLevel")
	case 17: // degree temperature
		return d.value("CurrentTemperature")
	case 18: // lux level
		return d.value("CurrentLevel")
	case 19: // ZigBee Network
		return  ""
	case 28: // UV sensor
		return d.value("CurrentLevel")
	default:
		return ""
	}
	return ""
}

func (d Device) StatusTxt() (r string) {
	//r = ""
	switch d.CategoryNum {
	case 2: //dimmable light
		if d.value("LoadLevelStatus") == "0" {
			r = "[Off]"
		} else {
			r = mkstr("[On] %v", d.value("LoadLevelStatus"))
		}
	case 3: //lights, switches
		if d.value("Target") == "1" {
			r = "[On]"
		} else {
			r = "[Off]"
		}
	case 4: // security device (motion, window sensor)
		if d.value("Armed") == "0" {
			r = "[Off]"
		} else {
			r = "[Armed]"
			if d.value("Tripped") == "1" {
				r = mkstr("%v [tripped]", r)
			}
		}
	case 5: //HVAC
		r = mkstr("state:%v %vF bat:%v", d.value("ModeStatus"), d.value("CurrentTemperature"), d.value("BatteryLevel"))
	case 6: //cameras
		r = " " //just null out the status
	case 7: // door lock
		if d.value("Status") == "0" {
			r = "[Open]"
		} else {
			r = "[Locked]"
		}
	case 11: // Generic IO
		if d.value("ArmedTripped") == "0" {
			r = "[Fine]"
		} else {
			r = "[Alarm]"
		}
	case 14: // remote controller
		if len(d.value("Scenes"))>0 {  r = mkstr("controls %v", d.value("Scenes"))}
	case 16: // humidity level
		r = mkstr("%v%v", d.value("CurrentLevel"), `%`)
	case 17: // degree temperature
		r = mkstr("%vF", d.value("CurrentTemperature"))
	case 18: // lux level
		r = mkstr("%v lux", d.value("CurrentLevel"))
	case 19: // ZigBee Network
		r = " "
	case 28: // UV sensor
		r = mkstr("%v ", d.value("CurrentLevel"))
	default:
		r = "n/a"
	}
	return r
}

type Device struct {
	XMLName         xml.Name `xml:"device"`
	Id              int      `xml:"id,attr"`
	Name            string   `xml:"name,attr"`
	DeviceType     string   `xml:"device_type,attr"`
	RoomNum            string   `xml:"room,attr"`
	DeviceFile     string   `xml:"device_file,attr"`
	CategoryNum    int      `xml:"category_num,attr"`
	SubcategoryNum string   `xml:"subcategory_num,attr"`
	Time_created    string   `xml:"time_created,attr"`
	Invisible       string   `xml:"invisible,attr"`
	LocalUdn       string   `xml:"local_udn,attr"`
	States          []State  `xml:"states>state"`
}

type State struct {
	XMLName  xml.Name `xml:"state"`
	Service  string   `xml:"service,attr"`
	Variable string   `xml:"variable,attr"`
	Value    string   `xml:"value,attr"`
	Id       string   `xml:"value,id"`
}



func GetRoot(byteIn []byte) (x VeraRoot, err error) {
	err = xml.Unmarshal(byteIn, &x)
	return
}
// end of matt library type functions
func OpenFile (f string) (b []byte, err error) {
    defer DTook(time.Now(),"open verafile")
    x,err := os.Open(f)
    // open vera xml data file
    if err != nil {
        ErrorExit(mkstr("Error opening data file: %s",err),1)
    }
    DTime("Opened Vera xml file\n")
	b,err = ioutil.ReadAll(x)
    return
}
