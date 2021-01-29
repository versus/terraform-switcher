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
	version     = "terraform-switcher 0.21.4\n\n"
)

//var version string

func main() {
	var path, tfversion string

	customBinPathFlag := getopt.StringLong("bin", 'b', "", "Custom binary path. For example: /Users/username/bin/terraform")
	listReleaseFlag := getopt.BoolLong("list", 'r', "List release versions of terraform")
	listAllFlag := getopt.BoolLong("list-all", 'l', "List all versions of terraform - including beta and rc")
	programVersionFlag := getopt.BoolLong("version", 'v', "Displays the version of tfswitch")
	latestVersionFlag := getopt.BoolLong("latest", 'u', "Switch to the latest terraform version")
	preVersionFlag := getopt.BoolLong("pre", 'p', "Latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest)\"")
	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")
	initFlag := getopt.BoolLong("save", 's', "Generate .tfswitch.toml in current directory with current version terraform")
	removeFlag := getopt.BoolLong("delete", 'd', "Remove terraform version from filesystem")
	noSymlinkFlag := getopt.BoolLong("no-symlink", 'n', "Skip symlink creation")
	noPromptFlag := getopt.BoolLong("quiet", 'q', "Only switch if version is detected or specified")

	getopt.Parse()
	args := getopt.Args()

	if len(args) > 1 || *helpFlag {
		cmd.UsageMessage()
	}

	fmt.Printf(version)
	if *programVersionFlag {
		os.Exit(0)
	}

	tfversion, path = cmd.GetConfigVariable()

	envPath := os.Getenv("TFSWITCH_PATH")
	if envPath != "" {
		fmt.Println("TFSWITCH_PATH detected")
		path = envPath
	}

	envVersion := os.Getenv("TFSWITCH_VERSION")
	if envVersion != "" {
		fmt.Println("TFSWITCH_VERSION  detected")
		tfversion = envVersion
	}

	if *customBinPathFlag != "" {
		path = *customBinPathFlag
		if tfversion == "" {
			v, err := cmd.GetInstalledVersion(*customBinPathFlag)
			if err != nil {
				tfversion = ""
			}
			tfversion = v
		}
	}

	if path == "" {
		path = defaultPath
	}

	createSymlink := true
	if *noSymlinkFlag {
		createSymlink = false
	}

	//fmt.Println("tfversion=", tfversion)
	//fmt.Println("path=", path)

	if *removeFlag {
		if len(args) == 1 {
			cmd.RemoveSelectVersion(args[0], path)
			os.Exit(0)
		} else {
			cmd.Remove(path)
		}
	}

	if *initFlag {
		if len(args) == 1 {
			cmd.InitConfigVersion(args[0], path, createSymlink)
		} else if tfversion != "" {
			cmd.InitConfigVersion(tfversion, path, createSymlink)
		} else {
			if *latestVersionFlag {
				cmd.InitConfigLatestVersion(path, createSymlink)
			}
			if *listAllFlag {
				cmd.InitConfig(true, path, createSymlink)
			} else {
				cmd.InitConfig(false, path, createSymlink)
			}
		}
	}

	if *listReleaseFlag {
		cmd.Install(false, path, createSymlink)
	}

	if *listAllFlag {
		cmd.Install(true, path, createSymlink)
	}

	if *latestVersionFlag {
		cmd.InstallLatest(path, createSymlink)
	}

	if *preVersionFlag {
		cmd.InstallPreReleaseVersion(args[0], path, createSymlink)
	}

	if len(args) == 0 {
		if tfversion != "" && path != "" {
			cmd.InstallSelectVersion(tfversion, path, createSymlink)
		} else if *noPromptFlag {
			fmt.Println("No terraform version detected")
			os.Exit(0)
		} else {
			cmd.Install(false, path, createSymlink)
		}
	}

	if len(args) == 1 {
		cmd.InstallSelectVersion(args[0], path, createSymlink)
	}

}
