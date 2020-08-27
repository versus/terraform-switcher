package main

/*** OPERATION WORKFLOW ***/
/*
* 1- Create /usr/local/terraform directory if does not exist
* 2- Download zip file from url to /usr/local/terraform
* 3- Unzip the file to /usr/local/terraform
* 4- Rename the file from `terraform` to `terraform_version`
* 5- Remove the downloaded zip file
* 6- Read the existing symlink for terraform (Check if it's a homebrew symlink)
* 7- Remove that symlink (Check if it's a homebrew symlink)
* 8- Create new symlink to binary  `terraform_version`
 */

import (
	"fmt"
	"github.com/pborman/getopt"
	"github.com/versus/terraform-switcher/cmd"
	"os"
)

const (
	defaultPath   = "/usr/local/bin/terraform" //default bin installation dir
	version = "terraform-switcher 0.9.1\n\n"
)

func main() {
	var path, tfversion string

	customBinPathFlag := getopt.StringLong("bin", 'b', "", "Custom binary path. For example: /Users/username/bin/terraform")
	listReleaseFlag := getopt.BoolLong("list", 'r', "List release versions of terraform")
	listAllFlag := getopt.BoolLong("list-all", 'l', "List all versions of terraform - including beta and rc")
	programVersionFlag := getopt.BoolLong("version", 'v', "Displays the version of tfswitch")
	latestVersionFlag := getopt.BoolLong("latest", 'u', "Switch to the latest terraform version")
	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")

	getopt.Parse()
	args := getopt.Args()

	fmt.Printf(version)
	if *programVersionFlag {
		os.Exit(0)
	}

	if *helpFlag {
		cmd.UsageMessage()
	}

	tfversion, path = cmd.GetConfigVariable()
	fmt.Println("tfversion=", tfversion)
	fmt.Println("path=",path)


	if *customBinPathFlag != "" {
		path = *customBinPathFlag
	}

	if path == "" {
		path = defaultPath
	}

	if *listReleaseFlag {
		cmd.Install(false, path)
	}

	if *listAllFlag {
		cmd.Install(true, path)
		os.Exit(0)
	}

	if *latestVersionFlag {
		cmd.InstallLatest(path)
	}

	if len(args) == 0 {
		if tfversion != ""  &&  path != ""{
			cmd.InstallSelectVersion(tfversion, path)
		} else {
			cmd.Install(false, path)
			os.Exit(0)
		}
	} else 	if len(args) == 1 {
		cmd.InstallSelectVersion(args[0], path)
	} else if len(args) > 1 {
		cmd.UsageMessage()
	}

}
