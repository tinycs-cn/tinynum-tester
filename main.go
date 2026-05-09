package main

import (
	"os"

	"github.com/tinycs-cn/tinynum-tester/internal/stages"
	tester_utils "github.com/tinycs-cn/tester-utils"
)

func main() {
	definition := stages.GetDefinition()
	os.Exit(tester_utils.Run(os.Args[1:], definition))
}
