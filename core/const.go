package core

import (
	"fmt"
)

const serverName = "magic_engine"

func traceInfo(info string) {
	fmt.Printf("[%s] %s\n", serverName, info)
}

func panicInfo(info string) {
	msg := fmt.Sprintf("[%s] %s\n", serverName, info)
	panic(msg)
}
