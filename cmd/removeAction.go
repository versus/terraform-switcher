package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/versus/terraform-switcher/lib"
	"log"
	"os"
)

func RemoveSelectVersion(tfversion string, path string)  {
	lib.RemoveVersion(tfversion, path)
}

func Remove(path string)  {

	tflist, err := lib.GetInstalledVersions()
	if err != nil {
			log.Printf("Can't get versions: %v\n", err)
	}
	tflist = lib.RemoveDuplicateVersions(tflist)  //remove duplicate version

	/* prompt user to select version of terraform */
	prompt := promptui.Select{
		Label: "Select Terraform version for remove",
		Items: tflist,
	}

	_, tfversion, errPrompt := prompt.Run()


	if errPrompt != nil {
		log.Printf("Prompt failed %v\n", errPrompt)
		os.Exit(1)
	}

	RemoveSelectVersion(tfversion, path)
	os.Exit(0)
}