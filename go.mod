module github.com/konsorten/go-xdelta-perftest-electron

require (
	github.com/dustin/go-humanize v1.0.0
	github.com/konsorten/go-xdelta v0.0.0-00010101000000-000000000000
)

// Take the source code from local directory.
// Remove the following line to take the remote repository.
// You might need to provide the library (go-xdelta-lib.dll) manually
// when using CGO_DISABLED=1 on Windows.
replace github.com/konsorten/go-xdelta => ../go-xdelta

go 1.13
