package main

import (
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"os"
	"runtime"
	"sort"

	"mjoy.io/log"
	"mjoy.io/mjoyd/defaults"
	"mjoy.io/mjoyd/limits"
	"mjoy.io/mjoyd/utils"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the "+ defaults.AppName + " command line interface")
	// basic flags
	basicFlags = []cli.Flag{
		utils.DataDirFlag,
		utils.ConfigFileFlag,
		utils.LogFileFlag,
		utils.LogLevelFlag,
		utils.NodeNameFlag,
		utils.ListenPortFlag,
		utils.BootNodeUrlFlag,
		utils.HttpPortFlag,
		utils.HttpHostFlag,
		utils.HttpModulesFlag,
		utils.StartBlockproducerFlag,
		utils.MetricsEnabledFlag,
		utils.WorkingNetFlag,
		utils.ResyncBlockFlag,
	}

	logTag = "mjoyd.main"
	logger log.Logger
)

func init() {
	// get a logger
	logger = log.GetLogger(logTag)
	if logger == nil {
		fmt.Errorf("Can not get logger(%s)\n", logTag)
		os.Exit(1)
	}

	// Initialize
	app.Action = mjoyd
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018 The " + defaults.AppName + " Authors"

	// add commands
	app.Commands = []cli.Command{
		versionCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	// add flags
	app.Flags = append(app.Flags, basicFlags...)

	// set before action
	app.Before = func(ctx *cli.Context) error {
		// TODO:

		// Use all processor cores.
		runtime.GOMAXPROCS(runtime.NumCPU())

		// Up some limits.
		if err := limits.SetLimits(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
			return err
		}

		return nil
	}

	// set after action
	app.After = func(ctx *cli.Context) error {
		// TODO:
		return nil
	}
}

// mjoyd is the real main entry pointclear
func mjoyd(ctx *cli.Context) error {
	// log instance init
	err := log.InitInstance(ctx.GlobalString(utils.LogFileFlag.Name), ctx.GlobalString(utils.LogLevelFlag.Name))
	if err != nil {
		os.Exit(1)
	}
	defer log.CloseInstance()
	logger.Infof("")
	logger.Infof("===============================")
	logger.Infof("Hi, %s is starting ...", defaults.AppName)
	logger.Infof("===============================")

	//// TODO:
	//node := createNode(ctx)
	//if node == nil {
	//	logger.Critical("Create node failed.")
	//	os.Exit(1)
	//}
	//registerService(node)
	//startNode(node)

	node := createMjoyNode(ctx)
	if node == nil {
		logger.Critical("Create node failed.")
		os.Exit(1)
	}

	startMjoyNode(node)
	logger.Infof("%s is shutdown.", defaults.AppName)
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
