package workdir

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tmp, _ := TempDir()
	if !tmp.IsEmpty("") {
		t.Error("Expected new temporary directory to be empty")
	}
	tmp.TouchAll("k")
	if tmp.IsEmpty("") {
		t.Error("Dir with contents should not be empty")
	}

	tmp.RemoveAll()
	if tmp.IsEmpty("") {
		t.Error("IsEmpty should be false for non existing")
	}

}

func TestRemoveAll(t *testing.T) {
	err := WorkDir("/").RemoveAll() // :-)
	if err == nil {
		t.Fatal("Well we've probably erased the entire disk")
	}
}

func Test_color(t *testing.T) {
	cases := []string{"?? ", "M  ", " M "}
	for _, input := range cases {
		got := color(input)
		if !strings.Contains(got, "\x1b[0m") || len(got) <= len(input) {
			t.Errorf("Failed to colorize %q, got %q", input, got)
		}
	}
}

func TestTempDir_error(t *testing.T) {
	os.Setenv("TMPDIR", "/_no_such_dir")
	defer os.Setenv("TMPDIR", "/tmp")
	_, err := TempDir()
	if err == nil {
		t.Fail()
	}
}

func ExampleWorkDir_Ls() {
	wd, _ := setup()
	defer wd.RemoveAll()
	wd.Ls(nil)
	// output:
	// A
	// B
	// empty/
	// ex/
	// newdir/
	// sub/
}

func Test_Ls_error(t *testing.T) {
	wd, _ := TempDir()
	wd.RemoveAll()
	err := wd.Ls(nil)
	if err == nil {
		t.Fail()
	}
}

func Test_Ls(t *testing.T) {
	wd, _ := setup()
	defer wd.RemoveAll()

	w := bytes.NewBufferString("")
	wd.Ls(w)
	got := w.String()
	exp := `A
B
empty/
ex/
newdir/
sub/
`
	if exp != got {
		t.Errorf("Expected \n'%s'\ngot \n'%s'", exp, got)
	}
}

func TestNew(t *testing.T) {
	got := New(".")
	if got != "." {
		t.Fail()
	}
}

func Test_String(t *testing.T) {
	got := WorkDir(".").String()
	exp := "."
	if exp != got {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}

func TestTouch(t *testing.T) {
	wd, _ := TempDir()
	defer wd.RemoveAll()
	wd.MkdirAll("x")
	os.Chmod(wd.Join("x"), 0000)
	_, err := wd.TouchAll("x")
	if err == nil {
		t.Error("Expected to fail")
	}
	// Cleanup
	os.Chmod(wd.Join("x"), 0644)
}

func TestMkdirAll(t *testing.T) {
	wd, _ := TempDir()
	defer wd.RemoveAll()
	os.Chmod(string(wd), 0000)
	err := wd.MkdirAll("hepp")
	if err == nil {
		t.Error("Expected to fail")
	}
	os.Chmod(wd.Join("x"), 0644)
}
