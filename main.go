package main

import (
	"context"
	"errors"
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
	"github.com/genuinetools/pkg/cli"
	"github.com/sirupsen/logrus"
)

var (
	debug bool
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "bpfps"
	p.Description = "A tool to list and diagnose bpf programs. (Who watches the watchers..? :)"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
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

		// Mount the bpf filesystem.
		bpf.CheckOrMountFS("/sys/fs/bpf")

		infos := getProgInfos()

		if len(infos) <= 0 {
			return errors.New("no bpf programs found")
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

		return nil
	}

	// Run our program.
	p.Run()
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
