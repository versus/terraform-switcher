package cmd

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/versus/terraform-switcher/lib"
	"github.com/kiranjthomas/terraform-config-inspect/tfconfig"
	"io/ioutil"
	"log"
	"os"
	"github.com/spf13/viper"
	"os/user"
	"sort"
	"strings"
)

const (
	tfvFilename  = ".terraform-version"
	rcFilename   = ".tfswitchrc"
	tomlFilename = ".tfswitch.toml"
	defaultBin   = "/usr/local/bin/terraform" //default bin installation dir
)

func GetConfigVariable() (string, string) {
	exist, tfversion, path := findConfig()
	if exist == true {
		if path == "" {
			if exist, _, p := checkHomeDirToml(); exist == true {
				 return tfversion, p
			}
			return tfversion, defaultBin
		}
		return tfversion, path
	}
	return  "" , defaultBin
}

func findConfig() (bool, string, string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("Failed to get current directory %v\n", err)
		os.Exit(1)
	}

	var tfversion, path string
	var exist bool

	if exist, tfversion, path = checkToml(dir); exist == true {
		return exist, tfversion, path
	}

	if exist, tfversion = checkTFswitchrc(dir); exist == true {
		return exist, tfversion, ""
	}

	if exist, tfversion = checkTFvfile(dir); exist == true {
		return exist, tfversion, ""
	}

	if exist, tfversion = checkTFVAR(dir); exist == true {
		return exist, tfversion, ""
	}

	if exist, tfversion, path = checkHomeDirToml(); exist == true {
		return exist, tfversion, path
	}

	return false, "", ""
}

func checkToml(dir string) (bool, string, string) {
	configfile := dir + fmt.Sprintf("/%s", tomlFilename) //settings for .tfswitch.toml file in current directory (option to specify bin directory)
	if _, err := os.Stat(configfile); err == nil {
		fmt.Printf("Reading configuration from %s\n", configfile)

		path :=  ""                        //takes the default bin (defaultBin) if user does not specify bin path
		configfileName := lib.GetFileName(tomlFilename) //get the config file
		viper.SetConfigType("toml")
		viper.SetConfigName(configfileName)
		viper.AddConfigPath(dir)

		errs := viper.ReadInConfig() // Find and read the config file
		if errs != nil {
			fmt.Printf("Unable to read %s provided\n", tomlFilename) // Handle errors reading the config file
			fmt.Println(err)
			os.Exit(1) // exit immediately if config file provided but it is unable to read it
		}

		bin := viper.Get("bin")                  // read custom binary location
		path = os.ExpandEnv(bin.(string))
		tfversion := viper.Get("version") //attempt to get the version if it's provided in the toml
		return true, tfversion.(string), path
	}

	return false, "", ""
}

func checkHomeDirToml() (bool, string, string) {
	usr, errCurr := user.Current()
	if errCurr != nil {
		return false, "", ""
	}
	return checkToml(usr.HomeDir)
}

func checkTFswitchrc(dir string) (bool, string) {
	rcfile := dir + fmt.Sprintf("/%s", rcFilename)               //settings for .tfswitchrc file in current directory (backward compatible purpose)
	if _, err := os.Stat(rcfile); err == nil {
		fmt.Printf("Reading required terraform version %s \n", rcFilename)
		fileContents, err := ioutil.ReadFile(rcfile)
		if err != nil {
			fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/terraform-switcher/blob/master/README.md\n", rcFilename)
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		tfversion := strings.TrimSuffix(string(fileContents), "\n")
		return true, tfversion
	}
	return false, ""
}

func checkTFvfile(dir string) (bool, string) {
	tfvfile := dir + fmt.Sprintf("/%s", tfvFilename)     //settings for .terraform-version file in current directory (tfenv compatible)

	if _, err := os.Stat(tfvfile); err == nil  { //if there is a .terraform-version file, and no command line arguments
		fmt.Printf("Reading required terraform version %s \n", tfvFilename)

		fileContents, err := ioutil.ReadFile(tfvfile)
		if err != nil {
			fmt.Printf("Failed to read %s file. Follow the README.md instructions for setup. https://github.com/warrensbox/terraform-switcher/blob/master/README.md\n", tfvFilename)
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}
		tfversion := strings.TrimSuffix(string(fileContents), "\n")
		return true, tfversion
	}
	return false, ""
}

func checkTFVAR(dir string)  (bool, string) {
	if module, _ := tfconfig.LoadModule(dir); len(module.RequiredCore) >= 1  { //if there is a version.tf file, and no commmand line arguments
		tfversion := ""
		tfconstraint := module.RequiredCore[0]        //we skip duplicated definitions and use only first one
		listAll := true                               //set list all true - all versions including beta and rc will be displayed
		tflist, _ := lib.GetTFList(hashiURL, listAll) //get list of versions
		fmt.Printf("Reading required version from terraform file, constraint: %s\n", tfconstraint)
		constrains, err := semver.NewConstraint(tfconstraint) //NewConstraint returns a Constraints instance that a Version instance can be checked against
		if err != nil {
			return false, ""
		}
		versions := make([]*semver.Version, len(tflist))
		for i, tfvals := range tflist {
			version, err := semver.NewVersion(tfvals) //NewVersion parses a given version and returns an instance of Version or an error if unable to parse the version.
			if err != nil {
				return false, ""
			}
			versions[i] = version
		}
		sort.Sort(sort.Reverse(semver.Collection(versions)))
		for _, element := range versions {
			if constrains.Check(element) { // Validate a version against a constraint
				tfversion = element.String()
				fmt.Printf("Matched version: %s\n", tfversion)
				return true, tfversion
			}
		}
	}
	return false, ""
}