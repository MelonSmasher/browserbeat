package cmd

import (
	"github.com/MelonSmasher/browserbeat/beater"
	version "github.com/MelonSmasher/browserbeat/version"

	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/libbeat/cmd/instance"
)

// Name of this beat
var Name = "browserbeat"
var Version = version.AppVersion

//var Version = "0.0.1-alpha1"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name, Version: Version})
