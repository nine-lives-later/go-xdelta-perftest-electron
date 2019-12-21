package perftest

import (
	"bytes"
	"crypto/sha1"
	"github.com/dustin/go-humanize"
	"io"
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
		ExeFilePath:     "./xdelta-bin-gpl/xdelta3-3.1.0-x86_64.exe",
	}

	t.Run("CreatePatch", func(t *testing.T) { testExeRoundtrip_CreatePatch(t, ctx) })
	t.Run("DumpPatchInfo", func(t *testing.T) { testExeRoundtrip_DumpPatchInfo(t, ctx) })
	t.Run("ApplyPatch", func(t *testing.T) { testExeRoundtrip_ApplyPatch(t, ctx) })
	t.Run("CompareHash", func(t *testing.T) { testExeRoundtrip_CompareHash(t, ctx) })
}

func testExeRoundtrip_CreatePatch(t *testing.T, ctx *testExeRoundtrip_Context) {
	cmd := exec.Command(ctx.ExeFilePath, "-S", "none", "-B", "8589934592", "-W", "67108864", "-s", ctx.FromFilePath, ctx.ToFilePath, ctx.PatchFilePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed with %s", err)
	}
}

func testExeRoundtrip_DumpPatchInfo(t *testing.T, ctx *testExeRoundtrip_Context) {
	fromFileStat, err := os.Stat(ctx.FromFilePath)
	if err != nil {
		t.Fatalf("Failed to get FROM filesize: %v", err)
	}

	toFileStat, err := os.Stat(ctx.ToFilePath)
	if err != nil {
		t.Fatalf("Failed to get TO filesize: %v", err)
	}

	patchFileStat, err := os.Stat(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to get PATCH filesize: %v", err)
	}

	t.Logf("FROM  file size: %v (%v)", fromFileStat.Size(), humanize.Bytes(uint64(fromFileStat.Size())))
	t.Logf("TO    file size: %v (%v)", toFileStat.Size(), humanize.Bytes(uint64(toFileStat.Size())))
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

func testExeRoundtrip_CompareHash(t *testing.T, ctx *testExeRoundtrip_Context) {
	// open the files
	toFile, err := os.Open(ctx.ToFilePath)
	if err != nil {
		t.Fatalf("Failed to open TO file: %v", err)
	}
	defer toFile.Close()

	// calculate hash
	toHash := sha1.New()

	_, err = io.Copy(toHash, toFile)
	if err != nil {
		t.Fatalf("Failed to hash TO file: %v", err)
	}

	toFile.Close()

	toHashResult := toHash.Sum(nil)

	// open the files
	appliedFile, err := os.Open(ctx.AppliedFilePath)
	if err != nil {
		t.Fatalf("Failed to open APPLIED file: %v", err)
	}
	defer appliedFile.Close()

	// calculate hash
	appliedHash := sha1.New()

	_, err = io.Copy(appliedHash, appliedFile)
	if err != nil {
		t.Fatalf("Failed to hash APPLIED file: %v", err)
	}

	appliedFile.Close()

	appliedHashResult := appliedHash.Sum(nil)

	// compare
	t.Logf("APPLIED file hash: %x", appliedHashResult)

	if !bytes.Equal(toHashResult, appliedHashResult) {
		t.Fatalf("File hash of TO and APPLIED file are different")
	}
}
