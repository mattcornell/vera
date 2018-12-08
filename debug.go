package vera 

import ( 
"fmt"
"io"
"io/ioutil"
"time"
"runtime"
)

var DebugOpt bool
var DOut io.Writer = ioutil.Discard
var	DebugTime   = time.Now()

var DPause = dPause

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
func DMemPause(msg string) { 
	if DebugOpt {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("------------------\n")
        fmt.Printf("%v\n------------------\n", msg)
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
		var answer string
		fmt.Scanf("%s", &answer)
	}
}


func dPause(this string) {
	if DebugOpt {
		fmt.Fprintf(DOut, "%v", this)
		var answer string
		fmt.Scanf("%s", &answer)
	}
}

func DWhatsThis(this interface{}) {
	if DebugOpt {
		fmt.Fprintf(DOut, "DWhatsThis %T: %v\n", this, this)
	}
}

func DTook(t time.Time, msg string) {
    if DebugOpt {
        fmt.Fprintf(DOut, "took: %v: %v\n", msg, time.Since(t))
    }
}

func DTime(msg string) {
	if DebugOpt {
		fmt.Fprintf(DOut, "time: %v		 %v", time.Since(DebugTime), msg)
	}
}
func DMsg(msg string) {
	if DebugOpt {
		fmt.Fprintf(DOut, "%v", msg)
	}
}

