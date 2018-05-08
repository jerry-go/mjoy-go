package defaults

import (
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

var (
	AppName 					= getAppName()

	DefaultDataDir        		= getAppCurrentDir()
	DefaultTOMLConfigPath 		= getAppCurrentDir() + "/" + AppName + ".conf"
	DefaultLogPath        		= getAppCurrentDir() + "/log/" + AppName + ".log"
	DefaultLogLevel       		= "info"
	DefaultNodeName       		= AppName
	DefaultKeystore       		= getAppCurrentDir() + "/" + "keystore"
	DefaultNodePort       		= 36180

	//Rpc
	DefaultHttpModules    		= "mjoy,personal,txpool"
	DefaultHttpHost       		= "localhost"
	DefaultHttpPort       		= 8989
	//Miner
	DefaultBlockproducerStart	= false
	//Net
	DefaultWorkingNet			="alpha"
)

func getAppName() string {
	name := os.Args[0]
	if strings.HasSuffix(name, ".exe") {
		name = strings.TrimSuffix(name, ".exe")
		if name == "" {
			panic("empty executable name")
		}
	}
	name = strings.Replace(name, "\\", "/", -1)
    v := strings.SplitAfterN(name, "/", -1)
	name = v[len(v) - 1]
	return name
}

func getAppCurrentDir() string {
	// discard error !!
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	//fmt.Println("CurrentDir:",dir)

	return strings.Replace(dir, "\\", "/", -1)
}

func PrintAllDefalts(){
	fmt.Printf("DefaultDataDir:%s\n" , DefaultDataDir)
	fmt.Printf("DefaultKeystore:%s\n", DefaultKeystore)
}