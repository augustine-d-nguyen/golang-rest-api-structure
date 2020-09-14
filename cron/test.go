package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()
	time.Sleep(2 * time.Second)
	second := int(time.Since(startTime).Seconds())
	fmt.Println(second)
}
