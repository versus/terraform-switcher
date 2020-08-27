package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/versus/terraform-switcher/lib"
	"log"
	"os"
	"strings"
)

func InitConfig(listAll bool, path string)  {
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

	InitConfigVersion(tfversion, path)
	lib.Install(tfversion, path)
	os.Exit(0)
}

func InitConfigLatestVersion(path string) {
	tfversion, err := lib.GetTFLatest(hashiURL)
	if err != nil {
		fmt.Println("Error get latest version: ", err)
		os.Exit(1)
	}
	if lib.ValidVersionFormat(tfversion) { //check if version is correct
		InitConfigVersion(tfversion, path)
		lib.Install(tfversion, path)
		os.Exit(0)
	} else {
		fmt.Println(invalidVersion)
		os.Exit(1)
	}
}

func InitConfigVersion(tfversion string, path string) {
	// generate file in current directory
	val := "bin = \"" + path + "\"\n"
	val += "version = \"" + tfversion  +"\"\n"
	f, err := os.Create(".tfswitch.toml")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer f.Close()
	_, err = f.WriteString(val)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println("Writting configuration to  " + pwd + "/.tfswitch.toml\n")
	lib.Install(tfversion, path)
	os.Exit(0)
}
