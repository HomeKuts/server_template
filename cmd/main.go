package main 

import (
	srv "server_template"
)

const versionMajor = "0.1"

var (
	version string
)

func main() {
	srv.Start(versionMajor, version);
}

