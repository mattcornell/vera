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
	print("Cmd.do: %v %v\n", vera.Cmd.Do, vera.Cmd.Next)
	print("Cfg.File: %v\n", vera.Cfg.File)
	print("HelpOpt: %v\n",vera.HelpOpt)
	print("InfoOpt: %v\n",vera.InfoOpt)
	print("UpdateOpt: %v\n",vera.UpdateOpt)
	print("DebugOpt: %v\n",vera.DebugOpt)
	print("Cfg.File: %v\n",vera.Cfg.File)
}


