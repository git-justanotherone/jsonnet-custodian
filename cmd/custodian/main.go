package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	gojsonnet "github.com/git-justanotherone/jsonnet-custodian/cmd/internal/upstream/go-jsonnet"
)

func cmdMainUsage(o io.Writer) {
	fmt.Fprintln(o, "Custodian is a tool for managing jsonnet extended source code.")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Usage:")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "    custodian <command> [arguments]")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "The commands are:")
	fmt.Fprintln(o, "    mod        Module management commands")
	fmt.Fprintln(o, "    jsonnet    Run the jsonnet-extended interpreter. (like jsonnet but with extensions)")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Use \"custodian <command> -h\" for more information about a command.")
}

func cmdMain(o io.Writer, args []string) error {
	global := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	global.Usage = func() {
		cmdMainUsage(o)
	}
	// parse apenas flags globais
	global.Parse(args)
	nargs := global.Args()
	if len(nargs) == 0 {
		cmdMainUsage(o)
		os.Exit(1)
	}

	subArgs := nargs[1:] // Drop the command.
	switch nargs[0] {
	case "mod":
		err := cmdModMain(o, subArgs)
		if err != nil {
			panic(err)
		}
	case "jsonnet":
		err := gojsonnet.CmdJsonnetMain(subArgs)
		if err != nil {
			panic(err)
		}
	default:
		cmdMainUsage(o)
		fmt.Printf("error: unknown command - %q\n", nargs[0])
		os.Exit(1)
	}
	return nil
}

func main() {
	err := cmdMain(os.Stdout, os.Args[1:])
	if err != nil {
		panic(err)
	}
}
