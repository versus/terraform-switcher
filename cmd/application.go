package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/versus/terraform-switcher/lib"
	"log"
	"os"
	"strings"

	"github.com/pborman/getopt"

)

const (
	hashiURL     = "https://releases.hashicorp.com/terraform/"
	invalidVersion = "The provided terraform version does not exist. Try `tfswitch -r ot tfswitch -l` to see all available versions.")

/* installOption : displays & installs tf version */
/* listAll = true - all versions including beta and rc will be displayed */
/* listAll = false - only official stable release are displayed */
func Install(listAll bool, path string) {

	tflist, _ := lib.GetTFList(hashiURL, listAll) //get list of versions
	recentVersions, _ := lib.GetRecentVersions()  //get recent versions from RECENT file
	tflist = append(recentVersions, tflist...)    //append recent versions to the top of the list
	tflist = lib.RemoveDuplicateVersions(tflist)  //remove duplicate version

	/* prompt user to select version of terraform */
	prompt := promptui.Select{
		Label: "Select Terraform version",
		Items: tflist,
	}

	_, tfversion, errPrompt := prompt.Run()
	tfversion = strings.Trim(tfversion, " *recent") //trim versions with the string " *recent" appended

	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	lib.Install(tfversion, path)
	os.Exit(0)
}

func InstallLatest(path string)  {
	tfversion, err := lib.GetTFLatest(hashiURL)
	if err != nil {
		fmt.Println("Error get latest version: ", err)
		os.Exit(1)
	}
	if lib.ValidVersionFormat(tfversion) { //check if version is correct
		lib.Install(string(tfversion), path)
	} else {
		fmt.Println(invalidVersion)
		os.Exit(1)
	}
}

func InstallSelectVersion(tfversion string, path string)  {
		if lib.ValidVersionFormat(tfversion) { //check if version is correct
			lib.Install(string(tfversion), path)
		} else {
			fmt.Println(invalidVersion)
			os.Exit(1)
		}
}

func UsageMessage() {
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu\n")
	os.Exit(0)
}
