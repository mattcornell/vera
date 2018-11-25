package vera 
/*
###h### supporting package for vera command line client for Vera Micasa
###h### home controller
*/

import (
	"regexp"
	"fmt"
	"os"
)

const (
helpMsg = `
###i### -------------------------------------------------
###v### date: 2018-11-09_100434
###i### https://teamcornell.com/code/vera/
###v### by: matt@teamcornell.com
###i### vera
###i### this script is for controlling a veraPro Zwave home controller 
###i### 
###i### More notes:
###i### Naming devices with no spaces help xmllint pick up devices by name
###i### 
###i### source: https://mmcis.com/matt/projects/vera/
###i### -----------------------------------------------------------------
###h### vera [-qrvVhi] command [command arg]
###i### Vera stores a CONFIGFILE=./vera_control and TMPFILE=./vera_tmpfile
###h###    
###h### commands:
###h### help, info, list, room, scene, refresh, status, detail, on, off
###h### dim, lock, name, unlock, tripped, toggle, watch 
###h### [info for more...]
###h###    
###i###	help - this output
###i###	info - this output (but more)
###i###	list [device_nubmer|name|grep string] - list device(s)
###i###	room[s] [roomnumber] - list rooms 
###i###	scene[s] [scenenumber] - list scenes 
###i###	refresh - poll the zwave network for changes 
###i###	status dev_number - give the status value for device 
###i###	detail[s] device_num - show everything about the device
###i###	on dev_number - turn device on 
###i###	off dev_number - turn device off
###i###	dim dev_number level - change a dimmable light to 'level'
###i###	lock dev_number - lock a door
###i###	name dev_number - return the name of dev_number
###i###	unlock dev_number - unlock a door
###i###	tripped dev_number - toggle device state 
###i###	toggle dev_number - toggle device state 
###i###	watch dev_number [seconds] - watch a security device for "trip" 
###i###    
###h###    vera options:
###h###        -q  = quiet 	run quiet except for errors
###h###        -r  = force refresh of vera status
###h###        -u  = check version for update
###h###        -h  = print help
###h###        -i  = print more info with help
###i###        
###i### Configuration options (stored in the .config.toml): 
###i### host="vera.teamcornell.com"
###i### port="3480"
###i### lastpull="1543096934"
###i### freshfile="2"
###i### FRESHFILE=2 # <- minutes before forces a reload of vera settings
###i###  (the script will force a refresh when changing any state)
###i###        
###i###  examples: 
###i###   ./vera list   # show all devices 
###i###   ./vera list garage  # show all device with garage in the name
###i###   ./vera on 10 # turn on switch that is device num 10
###i###   ./vera dim 33 10  # Set dimmable light device 33 to 10 
###i###`
)
func configQuit(m string) {                                                                            
	if ! empty(m) { errOut(mkstr("\n vera: %v",m)) }
    reg, err := regexp.Compile(`###h?i###(.*)\r?\n`)                                                          
    if err != nil {                                                                                           
        ErrorExit("bad regexp compile",1)                                                                     
    }                                                                                                         
    for _,i := range reg.FindAllString(helpMsg,-1) {                                                          
        fmt.Printf("%v",reg.ReplaceAllString(i,"$1\n"))                                                            
    }                                                                                                         
    os.Exit(0)                                                                                                
}                                                                                                             
func InfoQuit(m string) {                                                                                             
	if ! empty(m) { errOut(mkstr("\nvera: %v",m)) }
    reg,err := regexp.Compile(`###h?i###(.*)\r?\n`)                                                           
    if err != nil {                                                                                           
        ErrorExit("bad regexp compile",1)                                                                     
    }                                                                                                         
    for _,i := range reg.FindAllString(helpMsg,-1) {                                                          
        fmt.Printf("%v",reg.ReplaceAllString(i,"$1\n"))                                                            
    }                                                                                                         
    os.Exit(0)                                                                                                
}                                                                                                             
                                                                                                              
func HelpQuit(m string) {                                                                                             
	if ! empty(m) { errOut(mkstr("\n vera: %v",m)) }
    reg,err := regexp.Compile(`###h###(.*)\r?\n`)                                                             
    if err != nil {                                                                                           
        ErrorExit("bad regexp compile",1)                                                                     
    }                                                                                                         
    for _,i := range reg.FindAllString(helpMsg,-1) {                                                          
        fmt.Printf("%v",reg.ReplaceAllString(i,"$1\n"))                                                            
    }                                                                                                         
    os.Exit(0)                                                                                                
}                                             

func errOut (m string ){
    fmt.Fprintf(os.Stderr, "%v", m)
}

func ErrorExit(m string, err int) {
    errOut(mkstr("error: %v\n", m))
    os.Exit(err)
}


