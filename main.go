package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/cilium/cilium/pkg/bpf"
	"github.com/genuinetools/bpfps/version"
	"github.com/sirupsen/logrus"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = ` _            __
| |__  _ __  / _|_ __  ___
| '_ \| '_ \| |_| '_ \/ __|
| |_) | |_) |  _| |_) \__ \
|_.__/| .__/|_| | .__/|___/
      |_|       |_|

 A tool to list and diagnose bpf programs.  (Who watches the watchers..? :)
 Version: %s
 Build: %s

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
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, version.VERSION, version.GITCOMMIT))
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

	infos := getProgInfos()

	if len(infos) <= 0 {
		logrus.Fatal("no bpf programs found")
	}

	// Print the table.
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// print header
	fmt.Fprintln(w, "BID\tNAME\tTYPE\tUID\tMAPS\tLOADTIME")

	for _, info := range infos {
		dur, _ := time.ParseDuration(fmt.Sprintf("%dns", info.LoadTime))
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%#v\t%s\n", info.ID, info.Name, bpf.ProgType(int(info.ProgType)).String(), info.CreatedByUID, info.MapIDs, dur)
	}

	w.Flush()
}

func getProgInfos() []bpf.ProgInfo {
	var (
		bid   uint32
		err   error
		infos = []bpf.ProgInfo{}
	)
	for {
		bid, err = bpf.GetProgNextID(bid)
		if err != nil {
			if !strings.Contains(err.Error(), "no such file or directory") {
				logrus.Warn(err)
			}
			// break on error here
			break
		}

		// Get the file descriptor
		fd, err := bpf.GetProgFDByID(bid)
		if err != nil {
			logrus.Warn(err)
			continue
		}

		// Get the object's info.
		info, err := bpf.GetProgInfoByFD(fd)
		if err != nil {
			logrus.Warn(err)
			continue
		}

		if &info == nil {
			logrus.Warnf("No program info returned for fd %d, bid %d", fd, bid)
		}
		infos = append(infos, info)
	}

	return infos
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
