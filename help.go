package cmd

import (
	"fmt"
)

func UnknownCommand() {
	fmt.Println(`Unknown command, usage :
    grest new app        = Create a new grest app in the current folder
    grest new service    = Create a new service of grest app in the current folder
    grest version        = Check grest-cli version
    grest fmt            = Formatting struct tag
    `)
}
