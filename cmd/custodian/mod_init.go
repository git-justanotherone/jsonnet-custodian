package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/git-justanotherone/jsonnet-custodian/pkg/modules"
)

func modInit(o io.Writer, name string) error {
	if _, err := os.Stat(modules.ModuleFileName); err == nil {
		return fmt.Errorf("module file '%s' already exists", modules.ModuleFileName)
	}

	fmt.Fprintf(o, "Initializing module: %s\n", name)
	modFile := &modules.ModuleFile{
		Module:  name,
		Require: make(map[string]string),
	}

	data, err := modules.SerializeModuleFile(modFile)
	if err != nil {
		return err
	}
	err = os.WriteFile(modules.ModuleFileName, data, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(o, "Module '%s' initialized successfully.\n", name)
	return nil
}

func cmdModInitUsage(o io.Writer) {
	fmt.Fprintln(o, "Custodian mod init initializes a new module in the current directory.")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Usage:")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "    custodian mod init [name]")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Arguments:")
	fmt.Fprintln(o, "    name    Name of the module to initialize")
	fmt.Fprintln(o)
}

func cmdModInitMain(o io.Writer, args []string) error {
	init := flag.NewFlagSet("init", flag.ExitOnError)

	init.Usage = func() {
		cmdModInitUsage(o)
	}

	init.Parse(args)
	nargs := init.Args()
	if len(nargs) == 0 {
		init.Usage()
		os.Exit(1)
	}

	moduleName := nargs[0]
	return modInit(o, moduleName)
}
