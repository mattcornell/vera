package main 
// vera.go: cli command for vera micasa home controller
/* date: 2018-11-09_100434
 * by: matt@teamcornell.com
 * https://teamcornell.com/code/vera/
 * --------------------
 */
import ( 
"fmt"
"os"
//"io/ioutil"
//"net/http"
//"path/filepath"
//"crypto/tls"
//"strings"
"time"
//"regexp"
//"flag"
"text/tabwriter"
//"github.com/spf13/viper"
//"io"
//"regexp"
//"bytes"
//"gopkg.in/xmlpath.v2"
//"encoding/xml"
v "code.teamcornell.com/vera"
)

var (
	//rename oft used functions for fun
	mkstr      = fmt.Sprintf
	print      = fmt.Printf
)

const(
	// personal preference for date format
	dateFormat string = "2006-01-01_150405.00000"
)

//var c v.Cmd

func empty(object interface{}) bool {
	//First check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}
	return false
}

var t = new(tabwriter.Writer)
var Root v.VeraRoot

func main () { 
	t.Init(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug|tabwriter.AlignRight)
	v.GetOptions() 
	v.ReadCfg()
	v.DMsg(mkstr("Now Time: %v\n",time.Now().Unix()))


	/* v.DMsg(mkstr("Cmd.do: %v %v\n ", v.Cmd.Do, v.Cmd.Next))
	v.DMsg(mkstr("Cfg.Uri: %v\n ", v.Cfg.Uri))
	v.DMsg(mkstr("Cfg.Host: %v\n ", v.Cfg.Host))
	v.DMsg(mkstr("Cfg.Port: %v\n ", v.Cfg.Port))
	v.DMsg(mkstr("HelpOpt: %v\n ",v.HelpOpt))
	v.DMsg(mkstr("InfoOpt: %v\n ",v.InfoOpt))
	v.DMsg(mkstr("UpdateOpt: %v\n ",v.UpdateOpt))
	v.DMsg(mkstr("DebugOpt: %v\n ",v.DebugOpt))
	v.DMsg(mkstr("Cfg.File: %v\n ",v.Cfg.File))
	*/

	data,err:=v.Cmd.Execute()
	if err != nil {
		v.ErrorExit(mkstr("Error fetching data %v",err),1)
	}
	err = v.Cmd.WriteTemp(data)
	if err != nil {
		v.ErrorExit(mkstr("Error writing temp data %v",err),1)
	}

	v.Populate(data)
	v.WriteCfg()

	listDevs(v.Data.Devices)
	listRooms(v.Data.Rooms)
	listScene(v.Data.Scenes)
}


func listDevs (r v.DeviceList) {
	fmt.Fprintf(t,"ID:\tName\tRoom\tType\n")
	fmt.Fprintf(t,"---\t----\t----\t----\n")
	for _,this:= range r {
		fmt.Fprintf(t,"%v:\t%v\t%v\t%v\n",this.Id, this.Name, this.RoomName(),this.Category())
	}
	fmt.Fprintln(t,"\n")
	t.Flush()
}

func listRooms (r v.RoomList) {
	fmt.Fprintf(t,"ID:\tRoom Name\n")
	fmt.Fprintf(t,"---\t---------\n")
	for _,this:= range r {
		fmt.Fprintf(t,"%v:\t%v\n",this.Id, this.Name)
	}
	fmt.Fprintln(t,"\n")
	t.Flush()
}

func listUsers (r v.Users) {
	fmt.Fprintf(t,"ID:\tUser\tType\n")
	fmt.Fprintf(t,"---\t----\t----\n")
	for _,this:= range r {
		fmt.Fprintf(t,"%v:\t%v\t%v\n",this.Id, this.Name,this.Level)
	}
	t.Flush()
}

func listScene (r v.Scenes) {
	fmt.Fprintf(t,"ID:\tName\tLast Run\n")
	fmt.Fprintf(t,"---\t----\t--------\n")
	for _,this:= range r {
		fmt.Fprintf(t,"%v:\t%v\t%v\n",this.Id, this.Name, this.LastRun())
	}
	t.Flush()
}
