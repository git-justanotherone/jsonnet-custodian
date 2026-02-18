package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	patchDir = "./patches"

	files = map[string]string{
		// Download go-jsonnet source files from the specified version.
		"https://raw.githubusercontent.com/google/go-jsonnet/refs/tags/v0.21.0/cmd/jsonnet/cmd.go":        "go-jsonnet/cmd.go",
		"https://raw.githubusercontent.com/google/go-jsonnet/refs/tags/v0.21.0/cmd/internal/cmd/utils.go": "go-jsonnet/cmd/utils.go",
	}
)

func main() {
	must(copyFiles())
	must(applyPatches())
	must(gofmt())
	fmt.Println("✔ upstream synced from third_party successfully")
}

func copyFiles() error {
	for src, dst := range files {
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	var in io.ReadCloser
	var err error

	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		resp, err := http.Get(src)
		if err != nil {
			return err
		}
		in = resp.Body
	} else {
		in, err = os.Open(src)
		if err != nil {
			return err
		}
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func applyPatches() error {
	patches, err := filepath.Glob(filepath.Clean(filepath.Join(patchDir, "*.patch")))
	if err != nil {
		return err
	}
	if len(patches) == 0 {
		return nil
	}

	cmd := exec.Command(
		"git",
		append([]string{"apply", "--no-index", "--whitespace=nowarn"}, patches...)...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gofmt() error {
	cmd := exec.Command("gofmt", "-w", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
