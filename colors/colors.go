package colors

import (
	"fmt"
	"runtime"
)

var Reset = "\033[0m"
var Cyan = "\033[36m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Cyan = ""
	}
}

func ColorizeOutput(label, value string) {
	fmt.Println(Cyan+label+": "+Reset, value)
}
