package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func cmdModUsage(o io.Writer) {
	fmt.Fprintln(o, "Custodian mod provides module management for jsonnet.")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Usage:")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "    custodian mod <command> [arguments]")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "The commands are:")
	fmt.Fprintln(o, "    init    Initialize a new module")
	fmt.Fprintln(o, "    get     Download modules to the local module cache")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Use \"custodian mod <command> -h\" for more information about a command.")
}

func cmdModMain(o io.Writer, args []string) error {
	mod := flag.NewFlagSet("mod", flag.ExitOnError)

	mod.Usage = func() {
		cmdModUsage(o)
	}

	mod.Parse(args)
	nargs := mod.Args()
	if len(nargs) == 0 {
		cmdModUsage(o)
		os.Exit(1)
	}

	//subArgs := nargs[1:] // Drop program name and command.
	switch nargs[0] {
	case "init":
		return cmdModInitMain(o, nargs[1:])
	case "get":
		return cmdModGetMain(o, nargs[1:])
	default:
		cmdModUsage(o)
		fmt.Printf("error: unknown command - %q\n", nargs[0])
		os.Exit(1)
	}
	return nil
}
