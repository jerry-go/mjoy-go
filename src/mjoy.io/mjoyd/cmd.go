package main

import (
	"fmt"
	"mjoy.io/mjoyd/utils"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"
)

// TODO: To be complete

func version(ctx *cli.Context) error {
	fmt.Println("Version:", utils.Version)
	if gitCommit != "" {
		fmt.Println("Git Commit:", gitCommit)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	//fmt.Println("Protocol Versions:", mjoy.ProtocolVersions)
	//fmt.Println("Network Id:", mjoy.DefaultConfig.NetworkId)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}

var (
	versionCommand = cli.Command{
		Action:      version,
		Name:        "version",
		Usage:       "Print version numbers",
		ArgsUsage:   " ",
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The output of this command is supposed to be machine-readable.`,
	}
)
