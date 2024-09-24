package main

import (
	"fmt"
	"os"
	"os/signal"
	"proxytun/io_rw"
	"proxytun/tun"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

const (
	ExitSetupFailed = 1
)

const (
	EnvWgTunFd = "WG_TUN_FD"
)

const DefaultMTU = 1420

// Vars:
// N=number of goroutines
// S=size of payload
// opt=Operation, between read and write

func main() {
	interfaceName := "tun0"

	operation   := os.Args[1]
	payloadSize, errPayload := strconv.Atoi(os.Args[2])
	nGoroutines, errNGoroutines := strconv.Atoi(os.Args[3])

	if errPayload != nil {
		fmt.Println("Give a numerical payload size please...")
		os.Exit(-1)
	}

	if errNGoroutines != nil {
		fmt.Println("Give a numerical number of go routines please...")
		os.Exit(-1)
	}

	tdev, err := func() (tun.Device, error) {
		tunFdStr := os.Getenv(EnvWgTunFd)
		if tunFdStr == "" {
			return tun.CreateTUN(interfaceName, DefaultMTU)
		}

		// construct tun device from supplied fd

		fd, err := strconv.ParseUint(tunFdStr, 10, 32)
		if err != nil {
			return nil, err
		}

		err = unix.SetNonblock(int(fd), true)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}

		file := os.NewFile(uintptr(fd), "")
		return tun.CreateTUNFromFile(file, DefaultMTU)
	}()

	if err == nil {
		realInterfaceName, err2 := tdev.Name()
		if err2 == nil {
			interfaceName = realInterfaceName
		}
	}

	if err != nil {
		fmt.Printf("Failed to create TUN device: %v", err)
		os.Exit(ExitSetupFailed)
	}

	// file := tdev.File()

	// fmt.Println("TUN-DEV file:", file)
	fmt.Println("Duration")

	switch operation {
	case "read":
		for i := 0; i < nGoroutines; i++ {
			go io_rw.RoutineReadFromTun(tdev)
		}
	case "write":
		time.Sleep(5 * time.Second)
		for i := 0; i < nGoroutines; i++ {
			go io_rw.RoutineWriteToTun(tdev, payloadSize)
		}
	default:
		fmt.Println("Please give either read or write as operations")
		os.Exit(-1)
	}

	errs := make(chan error)
	term := make(chan os.Signal, 1)

	signal.Notify(term, unix.SIGTERM)
	// signal.Notify(term, os.Interrupt)

	select {
	case <-term:
	case <-errs:
	}

	fmt.Println("Exiting...")
}
