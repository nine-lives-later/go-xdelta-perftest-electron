package perftest

import (
	"bytes"
	"context"
	"crypto/sha1"
	"github.com/dustin/go-humanize"
	xd "github.com/konsorten/go-xdelta"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

type testFullRoundtrip_Context struct {
	FromFilePath    string
	ToFilePath      string
	PatchFilePath   string
	AppliedFilePath string

	Header     []byte // set during seeding
	ToFileHash []byte
}

func TestFullRoundtrip(t *testing.T) {
	// get temporary directory
	tempDir, err := ioutil.TempDir("", "go-xdelta")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	defer os.RemoveAll(tempDir)

	ctx := &testFullRoundtrip_Context{
		FromFilePath:    "./test-data/electron-v2.0.17-win32-x64/electron.exe",
		ToFilePath:      "./test-data/electron-v5.0.12-win32-x64/electron.exe",
		PatchFilePath:   filepath.Join(tempDir, "patch"),
		AppliedFilePath: filepath.Join(tempDir, "to_applied"),
	}

	t.Run("CreatePatch", func(t *testing.T) { testFullRoundtrip_CreatePatch(t, ctx) })
	t.Run("DumpPatchInfo", func(t *testing.T) { testFullRoundtrip_DumpPatchInfo(t, ctx) })
	t.Run("ApplyPatch", func(t *testing.T) { testFullRoundtrip_ApplyPatch(t, ctx) })
	t.Run("CompareHash", func(t *testing.T) { testFullRoundtrip_CompareHash(t, ctx) })
}

func testFullRoundtrip_CreatePatch(t *testing.T, ctx *testFullRoundtrip_Context) {
	// open the files
	fromFile, err := os.Open(ctx.FromFilePath)
	if err != nil {
		t.Fatalf("Failed to open FROM file: %v", err)
	}
	defer fromFile.Close()

	toFile, err := os.Open(ctx.ToFilePath)
	if err != nil {
		t.Fatalf("Failed to open TO file: %v", err)
	}
	defer toFile.Close()

	patchFile, err := os.Create(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to open PATCH file: %v", err)
	}
	defer patchFile.Close()

	// prepare encoder
	options := xd.EncoderOptions{
		FileID:    "TestFullRoundtrip",
		FromFile:  fromFile,
		ToFile:    toFile,
		PatchFile: patchFile,
		Header:    ctx.Header,
	}

	enc, err := xd.NewEncoder(options)
	if err != nil {
		t.Fatalf("Failed to create encoder: %v", err)
	}
	defer enc.Close()

	// create the patch
	err = enc.Process(context.TODO())
	if err != nil {
		t.Fatalf("Failed to create patch: %v", err)
	}
}

func testFullRoundtrip_DumpPatchInfo(t *testing.T, ctx *testFullRoundtrip_Context) {
	patchFileStat, err := os.Stat(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to get patch filesize: %v", err)
	}

	t.Logf("PATCH file size: %v (%v)", patchFileStat.Size(), humanize.Bytes(uint64(patchFileStat.Size())))
}

func testFullRoundtrip_ApplyPatch(t *testing.T, ctx *testFullRoundtrip_Context) {
	// open the files
	fromFile, err := os.Open(ctx.FromFilePath)
	if err != nil {
		t.Fatalf("Failed to open FROM file: %v", err)
	}
	defer fromFile.Close()

	appliedFile, err := os.Create(ctx.AppliedFilePath)
	if err != nil {
		t.Fatalf("Failed to open APPLIED file: %v", err)
	}
	defer appliedFile.Close()

	patchFile, err := os.Open(ctx.PatchFilePath)
	if err != nil {
		t.Fatalf("Failed to open PATCH file: %v", err)
	}
	defer patchFile.Close()

	// prepare decoder
	options := xd.DecoderOptions{
		FileID:    "TestFullRoundtrip",
		FromFile:  fromFile,
		ToFile:    appliedFile,
		PatchFile: patchFile,
	}

	dec, err := xd.NewDecoder(options)
	if err != nil {
		t.Fatalf("Failed to create decoder: %v", err)
	}
	defer dec.Close()

	// retrieve header
	headerChannel := make(chan []byte, 1)
	dec.Header = headerChannel

	// apply the patch
	err = dec.Process(context.TODO())
	if err != nil {
		t.Fatalf("Failed to apply patch: %v", err)
	}

	// compare the header
	readHeader := <-headerChannel

	if !bytes.Equal(ctx.Header, readHeader) {
		t.Fatalf("Header of PATCH file does not match")
	}
}

func testFullRoundtrip_CompareHash(t *testing.T, ctx *testFullRoundtrip_Context) {
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

	if !bytes.Equal(ctx.ToFileHash, appliedHashResult) {
		t.Fatalf("File hash of TO and APPLIED file are different")
	}
}
