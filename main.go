package main

/*
#include <linux/bpf.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/cilium/cilium/pkg/bpf"
	"github.com/jessfraz/bpfps/version"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = ` _            __
| |__  _ __  / _|_ __  ___
| '_ \| '_ \| |_| '_ \/ __|
| |_) | |_) |  _| |_) \__ \
|_.__/| .__/|_| | .__/|___/
      |_|       |_|

 A tool to list and diagnose bpf programs.  (Who watchs the watchers..? :)
 Version: %s
`
)

var (
	debug bool
	vrsn  bool
)

func init() {
	// parse flags
	flag.BoolVar(&vrsn, "version", false, "print version and exit")
	flag.BoolVar(&vrsn, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&debug, "d", false, "run in debug mode")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if vrsn {
		fmt.Printf("bpfps version %s, build %s", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	// set log level
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Mount the bpf filesystem.
	bpf.MountFS()
}

func main() {
	// On ^C, or SIGTERM handle exit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			logrus.Infof("Received %s, exiting.", sig.String())
			os.Exit(0)
		}
	}()

	// Get the next id.
	var id, nextID uint32
	if err := getNextID(unsafe.Pointer(&id), unsafe.Pointer(&nextID)); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("id: %#v", id)
	logrus.Infof("next id: %#v", nextID)

	// Get the file descriptor
	fd, err := getFDByID(nextID)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("fd: %d", fd)

	// Get the object's info.
	var info bpfProgInfo
	if err := getObjInfo(fd, unsafe.Pointer(&info), unsafe.Sizeof(info)); err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("info: %#v", info)
}

type bpfProgGetNextID struct {
	startID   uint32
	progID    uint32
	nextID    uint32
	openFlags uint32
}

type bpfProgInfo struct {
	progType     bpf.ProgType
	id           uint32
	mapIDs       uint64
	createdByUID uint32
	loadTime     uint64 /* ns since boottime */
	name         string
}

type bpfAttrObj struct {
	bpfFD   uint32
	infoLen uint32
	info    uint64
}

func getNextID(start, next unsafe.Pointer) error {
	attr := bpfProgGetNextID{
		startID: uint32(uintptr(start)),
		nextID:  uint32(uintptr(next)),
	}

	ret, _, err := unix.Syscall(unix.SYS_BPF, C.BPF_PROG_GET_NEXT_ID, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if ret != 0 || err != 0 {
		return fmt.Errorf("Unable to get next id: %v", err)
	}

	return nil
}

func getFDByID(id uint32) (int, error) {
	attr := bpfProgGetNextID{
		progID: uint32(uintptr(id)),
	}

	fd, _, err := unix.Syscall(unix.SYS_BPF, C.BPF_PROG_GET_FD_BY_ID, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if fd < 0 || err != 0 {
		return int(fd), fmt.Errorf("Unable to get fd for program id %d: %v", id, err)
	}

	return int(fd), nil
}

func getObjInfo(fd int, info unsafe.Pointer, size uintptr) error {
	attr := bpfAttrObj{
		bpfFD:   uint32(fd),
		infoLen: uint32(size),
		info:    uint64(uintptr(info)),
	}

	ret, _, err := unix.Syscall(unix.SYS_BPF, C.BPF_OBJ_GET_INFO_BY_FD, uintptr(unsafe.Pointer(&attr)), unsafe.Sizeof(attr))
	if ret != 0 || err != 0 {
		return fmt.Errorf("Unable to get object info: %v", err)
	}

	return nil
}

func usageAndExit(message string, exitCode int) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(exitCode)
}
