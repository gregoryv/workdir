package workdir

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type nopWriter struct{}

func (nw *nopWriter) Write(p []byte) (n int, err error) {
	return
}

type WorkDir struct {
	Root        string
	Writer      io.Writer
	Skip        func(*WorkDir, string, os.FileInfo) bool
	Format      func(*WorkDir, string, os.FileInfo) string
	NewWalkFunc Visitor
}

type Filter func(string) string

func Hidden(wd *WorkDir, path string, f os.FileInfo) bool {
	if strings.Index(f.Name(), ".") == 0 {
		return true
	}
	return false
}

func unfiltered(path string) string                           { return path }
func nameOnly(wd *WorkDir, path string, f os.FileInfo) string { return f.Name() }
func fullPath(wd *WorkDir, path string, f os.FileInfo) string {
	line := string(path[len(wd.Root)+1:])
	if f.IsDir() {
		line += "/"
	}
	return line
}

func New() *WorkDir {
	return &WorkDir{
		Root:        ".",
		Writer:      os.Stdout,
		Skip:        Hidden,
		Format:      fullPath,
		NewWalkFunc: showVisible,
	}
}

// Returns a new temporary working directory.
func TempDir() (wd *WorkDir, err error) {
	tmpPath, err := ioutil.TempDir("", "workdir")
	if err != nil {
		return
	}
	wd = New()
	wd.Root = tmpPath
	wd.Writer = &nopWriter{}
	return
}

func (wd *WorkDir) WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(wd.Join(file), data, 0644)
}

func (wd *WorkDir) MkdirAll(subDirs ...string) error {
	for _, sub := range subDirs {
		err := os.MkdirAll(filepath.Join(wd.Root, sub), 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wd *WorkDir) Command(cmd string, args ...string) *exec.Cmd {
	os.Chdir(wd.Root)
	return exec.Command(cmd, args...)
}

func (wd *WorkDir) String() string {
	return wd.Root
}

func (wd *WorkDir) Ls() {
	ls(wd, showVisible, unfiltered)
}

func ls(wd *WorkDir, visit Visitor, filter Filter) {
	out := make(chan string)
	go func() {
		_ = filepath.Walk(wd.Root, visit(wd, out))
		close(out)
	}()
	for path := range out {
		line := filter(path)
		if line != "" {
			fmt.Fprintf(wd.Writer, "%s\n", line)
		}
	}
}

type Visitor func(wd *WorkDir, out chan string) filepath.WalkFunc

func showVisible(wd *WorkDir, out chan string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if wd.Skip(wd, path, f) {
			if f.IsDir() && path != "." {
				return filepath.SkipDir
			}
			return nil
		}
		if f.Name() != filepath.Base(wd.Root) {
			out <- wd.Format(wd, path, f)
		}
		return nil
	}
}

func (wd *WorkDir) TouchAll(filenames ...string) ([]string, error) {
	files := make([]string, len(filenames))
	for i, name := range filenames {
		name, err := wd.Touch(name)
		if err != nil {
			return files, err
		}
		files[i] = name
	}
	return files, nil
}

func (wd *WorkDir) Touch(filename string) (string, error) {
	fh, err := os.Create(path.Join(wd.Root, filename))
	if err != nil {
		return filename, err
	}
	return filename, fh.Close()
}

func (wd *WorkDir) Join(filename string) string {
	return filepath.Join(wd.Root, filename)
}

func (wd *WorkDir) RemoveAll() error {
	return os.RemoveAll(wd.Root)
}