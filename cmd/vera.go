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

	v.DMsg(mkstr("Now Time: %v\n", time.Now().Unix()))

	v.DoRefresh()

	v.Populate()

	v.WriteCfg()

	switch v.Cmd.Do {
	case "all", "list":
		listDev(v.Data.Devices)
	case "on", "off","switch","toggle":
		switchDev()
	case "room", "rooms":
		listRoom(v.Data.Rooms)
	case "users", "user":
		listUser(v.Data.Users)
	case "scene", "scenes":
		listScene(v.Data.Scenes)
	case "value", "status":
		listDev(v.Data.Devices)
	}
}

func header(s string) {
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

func listDev(r v.DeviceList) {
	if len(v.Cmd.Next) > 1 {
		if !(v.Empty(v.Cmd.Next[2])) {
			if isInt(v.Cmd.Next[2]) {
				r = v.Data.DevId(v.Cmd.Next[2])
			} else {
				r = v.Data.DevContainsName(v.Cmd.Next[2])
			}
		} //end of ! Empty(v.Cmd.Next)
	}

	if len(r) > 0 {
		header("ID\tName\tRoom\tType\tStatus\n---\t----\t----\t----\t------\n")
	}
	for _, this := range r {
		fmt.Fprintf(t, "%v:\t%v\t%v\t%v\t%v\n", this.Id, this.Name, this.RoomName(), this.Category(), this.StatusTxt())
	} //end range over d devices
	t.Flush()
}

func listRoom(r v.RoomList) {
	header("Room\tRoom Name\n-----\t---------\n")
	for _, this := range r {
		fmt.Fprintf(t, "%v:\t%v\n", this.Id, this.Name)
	}
	fmt.Fprintln(t, "\n")
	t.Flush()
}

func listUser(r v.Users) {
	header("ID\tUser\tType\n---\t----\t----\n")
	for _, this := range r {
		fmt.Fprintf(t, "%v\t%v\t%v\n", this.Id, this.Name, this.Level)
	}
	t.Flush()
}

func listScene(r v.Scenes) {
	header("Scene\tName\tLast Run\n-----\t----\t--------\n")
	for _, this := range r {
		fmt.Fprintf(t, "%v:\t%v\t%v\n", this.Id, this.Name, this.LastRun())
	}
	t.Flush()
}
