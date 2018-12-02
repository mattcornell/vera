package vera 

import ( 
"fmt"
"io"
"io/ioutil"
"time"
)

var DebugOpt bool
var DOut io.Writer = ioutil.Discard
var	DebugTime   = time.Now()

var DPause = dPause

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

