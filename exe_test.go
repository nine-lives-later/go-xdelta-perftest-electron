package perftest

import (
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)


type testExeRoundtrip_Context struct {
	FromFilePath    string
	ToFilePath      string
	PatchFilePath   string
	AppliedFilePath string
	ExeFilePath     string
}

func TestExeRoundtrip(t *testing.T) {
	// get temporary directory
	tempDir, err := ioutil.TempDir("", "go-xdelta")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tempDir)

	ctx := &testExeRoundtrip_Context{
		FromFilePath:    "./test-data/electron-v2.0.17-win32-x64/electron.exe",
		ToFilePath:      "./test-data/electron-v5.0.12-win32-x64/electron.exe",
		PatchFilePath:   filepath.Join(tempDir, "patch"),
		AppliedFilePath: filepath.Join(tempDir, "to_applied"),
		ExeFilePath:      "./xdelta-bin-gpl/xdelta3-3.1.0-x86_64.exe",
	}

	t.Run("CreatePatch", func(t *testing.T) { testExeRoundtrip_CreatePatch(t, ctx) })
	t.Run("DumpPatchInfo", func(t *testing.T) { testExeRoundtrip_DumpPatchInfo(t, ctx) })
	t.Run("ApplyPatch", func(t *testing.T) { testExeRoundtrip_ApplyPatch(t, ctx) })
}

func testExeRoundtrip_CreatePatch(t *testing.T, ctx *testExeRoundtrip_Context) {
	cmd := exec.Command(ctx.ExeFilePath, "-0", "-B", "8589934592", "-W", "67108864", "-s", ctx.FromFilePath, ctx.ToFilePath, ctx.PatchFilePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed with %s", err)
	}
}

func testExeRoundtrip_DumpPatchInfo(t *testing.T, ctx *testExeRoundtrip_Context) {
	patchFileStat, err := os.Stat(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to get patch filesize: %v", err)
	}

	t.Logf("PATCH file size: %v (%v)", patchFileStat.Size(), humanize.Bytes(uint64(patchFileStat.Size())))
}

func testExeRoundtrip_ApplyPatch(t *testing.T, ctx *testExeRoundtrip_Context) {
	cmd := exec.Command(ctx.ExeFilePath, "-d", "-B", "8589934592", "-W", "67108864", "-s", ctx.FromFilePath, ctx.PatchFilePath, ctx.AppliedFilePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed with %s", err)
	}
}
