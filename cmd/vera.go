 // vera.go: cli command for vera micasa home controller
package main 
/* date: 2018-11-09_100434
 * by: matt@teamcornell.com
 * https://teamcornell.com/code/vera/
 * --------------------
 */
import ( 
"fmt"
//"os"
//"io/ioutil"
//"net/http"
//"path/filepath"
//"crypto/tls"
//"strings"
//"time"
//"regexp"
//"flag"
//"text/tabwriter"
//"github.com/spf13/viper"
//"io"
//"regexp"
//"bytes"
//"gopkg.in/xmlpath.v2"
//"encoding/xml"
"code.teamcornell.com/vera"
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

//var c vera.Cmd

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

func main () { 
	vera.GetOptions() 
	vera.DMsg(mkstr("Cfg.File: %v\n", vera.Cfg.File))

	vera.ReadCfg()
	vera.DMsg(mkstr("Cmd.do: %v %v\n ", vera.Cmd.Do, vera.Cmd.Next))
	vera.DMsg(mkstr("Cfg.Uri: %v\n ", vera.Cfg.Uri))
	vera.DMsg(mkstr("Cfg.Host: %v\n ", vera.Cfg.Host))
	vera.DMsg(mkstr("Cfg.Port: %v\n ", vera.Cfg.Port))
	vera.DMsg(mkstr("HelpOpt: %v\n ",vera.HelpOpt))
	vera.DMsg(mkstr("InfoOpt: %v\n ",vera.InfoOpt))
	vera.DMsg(mkstr("UpdateOpt: %v\n ",vera.UpdateOpt))
	vera.DMsg(mkstr("DebugOpt: %v\n ",vera.DebugOpt))
	vera.DMsg(mkstr("Cfg.File: %v\n ",vera.Cfg.File))
	vera.Cmd.Uri=vera.MakeUrl(vera.Cmd)
	if empty(vera.Cmd.Uri) { 
		vera.HelpQuit(mkstr("Command %v not found\n",vera.Cmd.Do))
	}

}


