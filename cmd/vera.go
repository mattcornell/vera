package main
// vera.go: cli command for vera micasa home controller
/* date: 2018-11-09_100434
 * by: matt@teamcornell.com
 * https://teamcornell.com/code/vera/
 * --------------------
 */
import (
	v "code.teamcornell.com/vera"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

var (
	//rename oft used functions for fun
	mkstr = fmt.Sprintf
	print = fmt.Printf
)

const (
	// personal preference for date format
	dateFormat string = "2006-01-01_150405.00000"
)

var t = new(tabwriter.Writer)
//var Root v.VeraRoot

func main() {
	padding := 3
	t.Init(os.Stdout, 0, 0, padding, ' ', 0)

	v.GetOptions()
	v.ReadCfg()
	v.DoRefresh()
	v.Populate()
	v.WriteCfg()

	switch v.Cmd.Do {
	case "all", "list":
	    var d v.Devices = v.Data.DeviceList 
	    listDev(d.Matches(v.SecondArg()))
	case "lock","unlock":
		this:=v.Cmd.MakeUri()
	    var d v.Devices = v.Data.DeviceList 
		if (d.Match(v.Cmd.Dev).CategoryNum == 7) { 
			this.Fetch()
		} else { 
			v.ErrorExit(mkstr("Device %v does not appear to be a lock",v.Data.DevId(v.Cmd.Dev).Name),1)
		}
		/* matt
		if (v.Data.DevMatches(v.Cmd.Dev).CategoryNum == 7) { 
			this.Fetch()
		} else { 
			v.ErrorExit(mkstr("Device %v does not appear to be a lock",v.Data.DevId(v.Cmd.Dev).Name),1)
		} matt */
	case "on", "off","switch","toggle":
		switchDev()
		time.Sleep(time.Second * 4) 
		v.RefreshAfterCommand()
		var d v.Devices = v.Data.DeviceList 
		print(mkstr("%v\n",d.Match(v.Cmd.Dev).Value()))
		r := v.Devices { v.Data.DevId(v.Cmd.Dev), }
		listDev(r)
	case "room", "rooms":
		var r v.Rooms = v.Data.RoomList
		listRoom(r.Match(v.SecondArg()))
	case "users", "user":
		listUser(v.Data.Users)
	case "scene", "scenes":
		listScene(v.Data.Scenes)
	case "value", "status":
		this:=v.Cmd.MakeUri()
		this.Fetch()
		if v.BareOpt {
			var d v.Devices = v.Data.DeviceList 
			print(mkstr("%v\n",d.Match(v.Cmd.Dev).Value()))

		} else {
			var d v.Devices = v.Data.DeviceList 
			print(mkstr("%v\n",d.Match(v.Cmd.Dev).StatusTxt()))
		}
	}
}

func printChoice(s string) {
	if !v.BareOpt {
		fmt.Fprintf(t, s)
	}
}

func isInt(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func switchDev() error { 
	v.Cmd.MakeUri().Fetch()
	return nil
}

func listDev(r v.Devices) {
	//r := v.Data.Devices
	if len(r) > 0 {
		printChoice("ID\tName\tRoom\tType\tStatus\n---\t----\t----\t----\t------\n")
	}
	for _, v := range r {
		printChoice(mkstr("%v:\t", v.Id))
		fmt.Fprintf(t, "%v\t",v.Name)
		printChoice(mkstr("%v\t",v.RoomName()))
		printChoice(mkstr("%v\t", v.Category()))
		fmt.Fprintf(t,"%v\n", v.StatusTxt())
	} //end range over d devices
	t.Flush()
}

func listRoom(r v.Rooms) {
	printChoice("Room\tRoom Name\n-----\t---------\n")
	for _, v := range r {
		printChoice(mkstr("%v:\t", v.Id))
		fmt.Fprintf(t, "%v\n", v.Name)
	}
	t.Flush()
}

func listUser(r v.Users) {
	printChoice("ID\tUser\tType\n---\t----\t----\n")
	for _, v := range r {
		printChoice(mkstr("%v", v.Id))
		fmt.Fprintf(t, "%v\t%v\n", v.Name, v.Level)
	}
	t.Flush()
}

func listScene(r v.Scenes) {
	printChoice("Scene\tName\tLast Run\n-----\t----\t--------\n")
	for _, v := range r {
		printChoice(mkstr("%v:\t", v.Id))
		fmt.Fprintf(t, "%v\t%v\n", v.Name, v.LastRun())
	}
	t.Flush()
}
