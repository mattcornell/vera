package vera
import (
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
	"time"
)

var mkstr = fmt.Sprintf 
var CatNames = []string{"", "NA Interface", "Dimmable_Light", "Switch", "Security_Sensor", "HVAC", "Camera", "Door_Lock", "Window_Convering", "Remote_Control", "IR_Transmitter", "Generic_I_O", "Generic_Sensor", "Serial_Port", "Scene_Controller", "AV", "Humidity_Sensor", "Temperature_Sensor", "Light_Sensor", "Z-Wave_Interface", "Insteon_Interface", "Power_Meter", "Alarm_Panel", "Alarm_Partition", "Siren", "Weather", "Philips_Controller", "Appliance", "UV_Sensor", "Mouse_Trap", "Doorbell", "Keypad"}

type VeraRoot struct {
	XMLName          xml.Name `xml:"root"`
	Timezone         string   `xml:"timezone,attr"`
	Firmware_version string   `xml:"firmware_version,attr"`
	City_description string   `xml:"City_description,attr"`
	Model            string   `xml:"model,attr"`
	Devices          []Device `xml:"devices>device"`
	Scenes           []Scene  `xml:"scenes>scene"`
	Users            []User   `xml:"users>user"`
	Rooms            []Room   `xml:"rooms>room"`
}

func (v VeraRoot) RoomName(id string) string {
	for _, this := range v.Rooms {
		if this.Id == id {
			return this.Name
		}
	}
	return ""
}

type Room struct {
	XMLName xml.Name `xml:"room"`
	/*id="1517941" Name="mariel" Level="1" IsGuest="0"></user>*/
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
	Timestamps        int64    `xml:"Timestamp,attr"`
	Name              string   `xml:"name,attr"`
	Room              string   `xml:"room,attr"`
	Triggers_operator string   `xml:"triggers_operator,attr"`
	users             string   `xml:"users,attr"`
	Paused            string   `xml:"paused,attr"`
	ModeStatus        int      `xml:"modeStatus,attr"`
	Id                int      `xml:"id,attr"`
	Last_Run          int64    `xml:"last_run,attr"`
	Trigger           Trigger  `xml:"triggers>trigger"`
}
type Trigger struct {
	XMLName  xml.Name `xml:trigger`
	Name     string   `xml:"name,attr"`
	Enabled  int      `xml:"enabled,attr"`
	Device   int      `xml:"device,attr"`
	Last_Run int64    `xml:"last_run"`
}

//method
func (d Device) Category() string {
	return CatNames[d.Category_num]
}

/* func (s State) debug() {
	debugPrint(mkstr("status.value: %v s.Variable: %v", s.Value, s.Variable))
	return

}*/
/*func (d Device) debug() {
	debugPrint(mkstr("device.id: %v d.Name: %v d.CatNum: %v", d.Id, d.Name, d.Category_num))
	return

}*/
func (d Device) value(V string) (r string) {
	for _, this := range d.States {
		if this.Variable == V {
			return this.Value
		}
	}
	return ""
}

func (d Device) StatusTxt() (r string) {
	r = ""
	switch d.Category_num {
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
		r = mkstr("controls scene %v", d.value("Scenes"))
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
	Device_type     string   `xml:"device_type,attr"`
	Room            string   `xml:"room,attr"`
	Device_file     string   `xml:"device_file,attr"`
	Category_num    int      `xml:"category_num,attr"`
	Subcategory_num string   `xml:"subcategory_num,attr"`
	Time_created    string   `xml:"time_created,attr"`
	Invisible       string   `xml:"invisible,attr"`
	Local_udn       string   `xml:"local_udn,attr"`
	States          []State  `xml:"states>state"`
}

type State struct {
	XMLName  xml.Name `xml:"state"`
	Service  string   `xml:"service,attr"`
	Variable string   `xml:"variable,attr"`
	Value    string   `xml:"value,attr"`
	Id       string   `xml:"value,id"`
}

//func stdErrOut (msg string ){
//fmt.Fprintf(os.Stderr, "%v", msg)
//}

func getVera(byteIn []byte) (x VeraRoot, err error) {
	err = xml.Unmarshal(byteIn, &x)
	return
}
// end of matt library type functions
func OpenFile (f string) (b []byte, err error) {
    defer DTook(time.Now(),"open verafile")
    x,err := os.Open(f)
    // open vera xml data file
    if err != nil {
        errorExit(mkstr("Error opening data file: %s",err),1)
    }
    DTime("Opened Vera xml file\n")
	b,err = ioutil.ReadAll(x)
    return
}





