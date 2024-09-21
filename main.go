package main

import (
	"fmt"
	"os"
	"os/signal"
	"proxytun/io_rw"
	"proxytun/tun"
	"strconv"

	"golang.org/x/sys/unix"
)

const (
	ExitSetupSuccess = 0
	ExitSetupFailed  = 1
)

const (
	ENV_WG_TUN_FD             = "WG_TUN_FD"
	ENV_WG_UAPI_FD            = "WG_UAPI_FD"
	ENV_WG_PROCESS_FOREGROUND = "WG_PROCESS_FOREGROUND"
)

const DefaultMTU = 1420

func main() {
	interfaceName := "tun0"

	tdev, err := func() (tun.Device, error) {
		tunFdStr := os.Getenv(ENV_WG_TUN_FD)
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

	file := tdev.File()

	fmt.Println("TUNDEV file:", file)
	fmt.Println("Everything should be fine...TUN device should exist somewhere")

	go io_rw.RoutineReadFromTun(tdev)

	errs := make(chan error)
	term := make(chan os.Signal, 1)


	signal.Notify(term, unix.SIGTERM)
	signal.Notify(term, os.Interrupt)

	select {
		case <- term:
		case <- errs:
	}

	fmt.Println("Exiting...")
}