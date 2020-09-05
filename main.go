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
	"os"

	"github.com/pborman/getopt"
	"github.com/versus/terraform-switcher/cmd"
)

const (
	defaultPath = "/usr/local/bin/terraform" //default bin installation dir
	version     = "terraform-switcher 0.10.1\n\n"
)

func main() {
	var path, tfversion string

	customBinPathFlag := getopt.StringLong("bin", 'b', "", "Custom binary path. For example: /Users/username/bin/terraform")
	listReleaseFlag := getopt.BoolLong("list", 'r', "List release versions of terraform")
	listAllFlag := getopt.BoolLong("list-all", 'l', "List all versions of terraform - including beta and rc")
	programVersionFlag := getopt.BoolLong("version", 'v', "Displays the version of tfswitch")
	latestVersionFlag := getopt.BoolLong("latest", 'u', "Switch to the latest terraform version")
	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")
	initFlag := getopt.BoolLong("init", 'i', "Generate .tfswitch.toml in current directory with current version terraform")
	removeFlag := getopt.BoolLong("delete", 'd', "Remove terraform version from filesystem")

	getopt.Parse()
	args := getopt.Args()

	if len(args) > 1 || *helpFlag{
		cmd.UsageMessage()
	}

	fmt.Printf(version)
	if *programVersionFlag {
		os.Exit(0)
	}

	tfversion, path = cmd.GetConfigVariable()

	envPath := os.Getenv("TFSWITCH_PATH")
	if envPath != "" {
		path = envPath
	}

	if *customBinPathFlag != "" {
		path = *customBinPathFlag
	}

	if path == "" {
		path = defaultPath
	}

	//fmt.Println("tfversion=", tfversion)
	//fmt.Println("path=",path)


	if *removeFlag {
		if len(args) == 1{
			cmd.RemoveSelectVersion(args[0], path)
			os.Exit(0)
		} else {
			cmd.Remove(path)
		}
	}

	if *initFlag {
		if len(args) == 1 {
			cmd.InitConfigVersion(args[0], path)
		} else {
			if *latestVersionFlag {
				cmd.InitConfigLatestVersion(path)
			}
			if *listAllFlag {
				cmd.InitConfig(true, path)
			} else {
				cmd.InitConfig(false, path)
			}
		}
	}

	if *listReleaseFlag {
		cmd.Install(false, path)
	}

	if *listAllFlag {
		cmd.Install(true, path)
	}

	if *latestVersionFlag {
		cmd.InstallLatest(path)
	}

	if len(args) == 0 {
		if tfversion != "" && path != "" {
			cmd.InstallSelectVersion(tfversion, path)
		} else {
			cmd.Install(false, path)
		}
	}

	if len(args) == 1 {
		cmd.InstallSelectVersion(args[0], path)
	}

}
