package vera 

import(
	"flag"
	"os"
	"strings"
)


var Cmd CmdType 

type CmdType struct {
	Do		string
	Next  []string
	Value	string
	Dev     string
	Uri		string 
	Body	[]byte
}

var (
	HelpOpt, InfoOpt, UpdateOpt bool 
)

func GetOptions() {
    DTime("Getting options\n")
    flag.BoolVar(&HelpOpt, "help", false, "Help info")
    flag.BoolVar(&HelpOpt, "h", false, "Help info")
    flag.BoolVar(&InfoOpt, "info", false, "Help, but more so")
    flag.BoolVar(&InfoOpt, "i", false, "Help, but more so")
    flag.BoolVar(&UpdateOpt, "u", false, "Check for an updated version ")
    flag.BoolVar(&UpdateOpt, "update", false, "Check for an updated version")
    flag.BoolVar(&DebugOpt, "d", false, "Debug messages")
    flag.BoolVar(&DebugOpt, "debug", false, "Debug messages")
    flag.StringVar(&Cfg.File, "c", Cfg.File, "specify config file location ")
    flag.StringVar(&Cfg.File, "config", Cfg.File, "specify config file location")
    flag.Parse()

    if DebugOpt {
        DOut = os.Stderr
        DMsg("Debug Enabled\n")
    }

    if (len(os.Args) > 1) && (len(flag.Args()) > 0) && (len(flag.Args()[0]) > 0) {
		Cmd.Do = flag.Args()[0]
		for _,v := range flag.Args() { 
			Cmd.Next = append(Cmd.Next, v)
		}

    } else {
       Cmd.Do = "list"
   }
    if Cmd.Do=="help" {
        HelpQuit("")
    }
    if HelpOpt || Cmd.Do=="help" {
        HelpQuit("")
    }
    if InfoOpt || Cmd.Do=="info"{
        InfoQuit("")
    }
    DMsg(mkstr("Args: %v \n", strings.Join(os.Args, " ")))
} 
