package main

// vera.go: cli command for vera micasa home controller
/* date: 2019-03-11_100434
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
	//	"github.com/pkg/profile"
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
	//defer profile.Start().Stop()
	//	defer profile.Start(profile.MemProfile).Stop()
	t.Init(os.Stdout, 0, 0, 3, ' ', 0)

	v.GetOptions()
	v.ReadCfg()
	v.DoRefresh()
	v.Populate()
	v.DMemPause("Populate")
	v.WriteCfg()

	switch v.Cmd.Do {
	case "all", "list":
		var d *v.Devices = &v.Data.DeviceList
		listDev(d.Matches(v.SecondArg()))
	case "lock", "unlock":
		this := v.Cmd.MakeUri()
		var d *v.Devices = &v.Data.DeviceList
		if d.Match(v.Cmd.Dev).CategoryNum == 7 {
			this.Fetch()
		} else {
			v.ErrorExit(mkstr("Device %v does not appear to be a lock", d.Match(v.Cmd.Dev).Name), 1)
		}
		v.WriteCfg()
	case "on", "off", "switch", "toggle":
		switchDev()
		time.Sleep(time.Second * 2)
		v.RefreshAfterCommand()
		var d *v.Devices = &v.Data.DeviceList
		print(mkstr("%v\n", d.Match(v.Cmd.Dev).Value()))
		r := v.Devices{d.Match(v.Cmd.Dev)}
		listDev(r)
	case "room", "rooms":
		var r *v.Rooms = &v.Data.RoomList
		//var d *v.Devices = &v.Data.DeviceList
		if !v.Empty(v.SecondArg()) {
			//listRooms(d.Matches(r.Matches(v.SecondArg())))
			if (len(r.Matches(v.SecondArg()))==1) { 
				listDevsFromRoom(r.Matches(v.SecondArg())[0],v.Data.DeviceList)
			} else  { 
				listRooms(r.Matches(v.SecondArg()))
			}
		} else {
			listRooms(*r)
		}
	case "users", "user":
		listUser(v.Data.Users)
		v.DMsg("users")
	case "scene", "scenes":
		v.DMsg("scenes")
		listScene(v.Data.Scenes)
	case "value", "status":
		this := v.Cmd.MakeUri()
		this.Fetch()
		if v.BareOpt {
			var d v.Devices = v.Data.DeviceList
			print(mkstr("%v\n", d.Match(v.Cmd.Dev).Value()))

		} else {
			var d v.Devices = v.Data.DeviceList
			print(mkstr("%v\n", d.Match(v.Cmd.Dev).StatusTxt()))
		}
	}
}

func printChoice(n int, s string) {
	if (n > 0) && !v.BareOpt {
		fmt.Fprintf(t, s)
	}
}

func isInt(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func switchDev() {
	v.Cmd.MakeUri().Fetch()
	return 
}

func listDev(r v.Devices) {
	printChoice(len(r), "ID\tName\tRoom\tType\tStatus\n---\t----\t----\t----\t------\n")
	v.DMemPause("list Dev ")
	for _, val := range r {
		printChoice(len(r), mkstr("%v:\t", val.Id))
		fmt.Fprintf(t, "%v\t", val.Name)
		printChoice(len(r), mkstr("%v\t", val.RoomName()))
		printChoice(len(r), mkstr("%v\t", *val.Category()))
		fmt.Fprintf(t, "%v\n", val.StatusTxt())
	} //end range over d devices
	t.Flush()
}

func listRooms(r v.Rooms) {
	printChoice(len(r), "Room\tRoom Name\n-----\t---------\n")
	for _, v := range r {
		printChoice(len(r), mkstr("%v:\t", v.Id))
		fmt.Fprintf(t, "%v\n", v.Name)
	}
	t.Flush()
}

func listDevsFromRoom(room v.Room,d v.Devices) {
	printChoice(1, mkstr("\nRoom: %v\n\n",room.Name)) 
		printChoice(len(room.Name), "ID\tName\tRoom\tType\tStatus\n---\t----\t----\t----\t------\n")
		v.DMemPause("list Dev From Room")
		for _, val := range d {
			//if (val.Id == thisroom.Id) { 
			if room.Match(val.RoomNum){
				printChoice(len(room.Name), mkstr("%v:\t", val.Id))
				fmt.Fprintf(t, "%v\t", val.Name)
				printChoice(len(room.Name), mkstr("%v\t", val.RoomName()))
				printChoice(len(room.Name), mkstr("%v\t", *val.Category()))
				fmt.Fprintf(t, "%v\n", val.StatusTxt())
			}
		} //end range over d devices
		t.Flush()
} //end of listDevsFromRoom

func listUser(r v.Users) {
	printChoice(len(r), "ID\tUser\tType\n---\t----\t----\n")
	for _, v := range r {
		printChoice(len(r), mkstr("%v", v.Id))
		fmt.Fprintf(t, "%v\t%v\n", v.Name, v.Level)
	}
	t.Flush()
}

func listScene(r v.Scenes) {
	printChoice(len(r), "Scene\tName\tStatus\tLast Run\n-----\t----\t----\t--------\n")
	for _, v := range r {
		printChoice(len(r), mkstr("%v:\t", v.Id))
		fmt.Fprintf(t, "%v\t%v\t%v\n", v.Name, v.PauseState(), v.LastRun())

	}
	t.Flush()
}
