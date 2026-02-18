package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/git-justanotherone/jsonnet-custodian/cmd/internal/utils"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/modules"
	"github.com/git-justanotherone/jsonnet-custodian/pkg/resolvers"
)

const (
	MODULE_CACHE_DIR = "/tmp/jnetx/modules"
	LOCK_FILE_NAME   = "module.lock"
)

func cmdModGetUsage(o io.Writer) {
	fmt.Fprintln(o, "Custodian mod get resolves its arguments to modules at specific versions,")
	fmt.Fprintln(o, "adds them as dependencies in the custodian.json file, and downloads all")
	fmt.Fprintln(o, "and their dependencies into the local module cache.")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Usage:")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "    custodian mod get [module]")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "Arguments:")
	fmt.Fprintln(o, "    module    Module to download (if omitted, downloads all modules in the module file)")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "To add a dependency or upgrade it to its latest version:")
	fmt.Fprintln(o, "    custodian mod get <module>")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "To add a dependency, upgrade or downgrade it to a specific version:")
	fmt.Fprintln(o, "    custodian mod get <module>@<version>")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "To remove a dependency:")
	fmt.Fprintln(o, "    custodian mod get <module>@none")
	fmt.Fprintln(o)
	fmt.Fprintln(o, "To download all dependencies in the module.json file:")
	fmt.Fprintln(o, "    custodian mod get")
}

func cmdModGetMain(o io.Writer, args []string) error {
	get := flag.NewFlagSet("get", flag.ExitOnError)

	get.Usage = func() {
		cmdModGetUsage(o)
	}

	get.Parse(args)
	nargs := get.Args()

	// Open and parse the module file
	moduleFile, err := os.Open(modules.ModuleFileName)
	if err != nil {
		return fmt.Errorf("failed to open module file: %w", err)
	}
	moduleData, err := modules.ParseModuleFile(moduleFile)
	moduleFile.Close()
	if err != nil {
		return fmt.Errorf("failed to parse module file: %w", err)
	}

	if len(nargs) == 0 {
		for moduleName, moduleIdentifier := range moduleData.Require {
			resolvedIdentifier, err := utils.GetModule(moduleIdentifier)
			if err != nil {
				return err
			}
			moduleData.Require[moduleName] = resolvedIdentifier
		}
	}

	for _, moduleIdentifier := range nargs {
		resolvedIdentifier, err := utils.GetModule(moduleIdentifier)
		if err != nil {
			return err
		}
		mId := resolvers.GitModuleIdentifier(moduleIdentifier)
		moduleData.Require[mId.Repo()] = resolvedIdentifier
	}

	// Serialize and write back the updated module file
	data, err := modules.SerializeModuleFile(moduleData)
	if err != nil {
		return err
	}
	err = os.WriteFile(modules.ModuleFileName, data, 0644)
	if err != nil {
		return err
	}

	// Get the full dependency tree and write the lock file
	dt, err := utils.GetDependencyTree()
	if err != nil {
		return err
	}
	lockFileData := dt.GenerateLockFile()
	err = os.WriteFile(LOCK_FILE_NAME, lockFileData, 0644)
	if err != nil {
		return err
	}
	return nil
}
