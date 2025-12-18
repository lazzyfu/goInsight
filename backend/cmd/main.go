package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lazzyfu/goinsight/internal/app"
)

// Build information injected by -ldflags during build
var (
	Version   string
	BuildTime string
	GitCommit string
	GitBranch string
)

func main() {
	// Command-line flags
	var configFile string
	var showVersion bool

	// Read local config file
	flag.StringVar(&configFile, "config", "config.yaml", "config file path (required)")
	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.Parse()

	if showVersion {
		printVersion()
		return
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Printf("Config file %s does not exist. Use -config to specify the config file path.\n", configFile)
		os.Exit(1)
	}

	app.Run(configFile)
}

func printVersion() {
	fmt.Println("GoInsight")
	if Version != "" {
		fmt.Printf("Version:    %s\n", Version)
	} else {
		fmt.Println("Version:    dev")
	}
	if BuildTime != "" {
		fmt.Printf("BuildTime:  %s\n", BuildTime)
	}
	if GitCommit != "" {
		fmt.Printf("GitCommit:  %s\n", GitCommit)
	}
	if GitBranch != "" {
		fmt.Printf("GitBranch:  %s\n", GitBranch)
	}
}
