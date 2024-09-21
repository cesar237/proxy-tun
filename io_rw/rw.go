package io_rw

import (
	"fmt"

	// "time"

	"proxytun/tun"
)

func RoutineReadFromTun(tdev tun.Device) {
	var (
		batchSize = 128
		elemBuf   = make([]byte, (1<<16)-1)
		readErr   error
		bufs      = make([][]byte, batchSize)
		count     = 0
		sizes     = make([]int, batchSize)
		offset    = 16
		rxBytes   = 0
		rxCount   = 0
	)

	for i := 0; i < batchSize; i++ {
		bufs[i] = elemBuf[:]
	}

	for {
		// read packets
		count, readErr = tdev.Read(bufs, sizes, offset)
		rxCount += count

		for i := 0; i < count; i++ {
			if sizes[i] < 1 {
				continue
			}
			rxBytes += sizes[i]
		}

		fmt.Println("Count: ", rxCount, "rxBytes: ", rxBytes)

		if readErr != nil {
			fmt.Println(readErr)
			return
		}
	}
}

func RoutineWriteToTun(tdev tun.Device) {

}
