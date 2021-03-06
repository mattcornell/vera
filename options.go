package vera 

import(
	"flag"
	"os"
	"strings"
	"strconv"
)

var Cmd CmdType

func (c CmdType) new() (r CmdType)  { 
	if Empty(c.Do) {r.Do="" } else { r.Do = c.Do}
	if Empty(c.Value) { r.Value="" }else { r.Value = c.Value}
	if Empty(c.Dev) { r.Dev="" }else { r.Dev = c.Dev}
	if Empty(c.Uri) { r.Uri="" } else { r.Uri = c.Uri}
	if len(c.Next)>0 { r.Next=c.Next }else {  r.Next=append(r.Next,"") } 
	return
}

type CmdType struct {
	Do		string
	Next  []string
	Value	string
	Dev     string
	Uri		string 
	Body	[]byte
}

func SecondArg() (r string) {
    for i, next := range Cmd.Next {
        if (i == 2) { return next }
    }
    return r
}

var (
	HelpOpt, BareOpt, InfoOpt, OfflineOpt, RefreshOpt, UpdateOpt bool 
)

func isInt(s string) bool {
    if _, err := strconv.Atoi(s); err == nil {
        return true
    }
    return false
}


func GetOptions() {
	Cmd = Cmd.new() 
    DTime("Getting options\n")
    flag.BoolVar(&HelpOpt, "help", false, "Help info")
    flag.BoolVar(&HelpOpt, "h", false, "Help info")
    flag.BoolVar(&InfoOpt, "info", false, "Help, but more so")
    flag.BoolVar(&InfoOpt, "i", false, "Help, but more so")
    flag.BoolVar(&UpdateOpt, "U", false, "Check for an updated version ")
    flag.BoolVar(&UpdateOpt, "u", false, "Check for an updated version ")
    flag.BoolVar(&UpdateOpt, "update", false, "Check for an updated version")
    flag.BoolVar(&OfflineOpt, "O", false, "Offline, use the cache file ")
    flag.BoolVar(&OfflineOpt, "o", false, "Offline, use the cache file")
    flag.BoolVar(&OfflineOpt, "offline", false, "Offline, use the cache fil")
    flag.BoolVar(&RefreshOpt, "R", false, "Force a refresh of cached data")
    flag.BoolVar(&RefreshOpt, "r", false, "Force a refresh of cached data")
    flag.BoolVar(&RefreshOpt, "refresh", false, "Force a refresh of cached data")
    flag.BoolVar(&BareOpt, "B", false, "least output")
    flag.BoolVar(&BareOpt, "b", false, "least output")
    flag.BoolVar(&BareOpt, "bare", false, "least output")
    flag.BoolVar(&DebugOpt, "D", false, "Debug messages")
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
		if ! Empty(Cmd.Next[1]) { 
			switch Cmd.Next[1]  { 
			case "on":
				DWhatsThis(len(Cmd.Next))
				if (len(Cmd.Next)<3) || Empty(Cmd.Next[2]){ 
					HelpQuit("on command needs a device number\n")
				}
				Cmd.Dev = Cmd.Next[2]
				Cmd.Value = "1"
			case "off":
				if (len(flag.Args())<2) || Empty(Cmd.Next[2] ){ 
					HelpQuit("off command needs a device number")
				}
				Cmd.Dev = Cmd.Next[2]
				Cmd.Value = "off"
			case "dim":
				if (len(flag.Args())<2) || Empty (Cmd.Next[2]) || Empty(Cmd.Next[3]){ 
					HelpQuit("dim command needs a device number and intesity value\n")
				}
					HelpQuit("dim command needs a device number and intesity value\n")
				Cmd.Dev = Cmd.Next[2]
				Cmd.Value = Cmd.Next[3]
			case "toggle","switch", "lock", "unlock","status","value":
			   if len(Cmd.Next) > 2 {
					if Empty ( Cmd.Next[2] ){ 
						HelpQuit(mkstr("command %v needs a device name or number\n",Cmd.Do))
					}
				}
				Cmd.Dev = Cmd.Next[2]
			case "scene","scenes":
			}
		}
    } else {
       Cmd.Do = "list"
   }
    if HelpOpt || Cmd.Do=="help" {
        HelpQuit("")
    }
    if InfoOpt || Cmd.Do=="info"{
        InfoQuit("")
    }
    DMsg(mkstr("Args: %v \n", strings.Join(os.Args, " ")))
} 
