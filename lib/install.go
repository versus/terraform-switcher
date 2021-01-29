package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const (
	hashiURL       = "https://releases.hashicorp.com/terraform/"
	installFile    = "terraform"
	installVersion = "terraform_"
	installPath    = "/.terraform.versions/"
	recentFile     = "RECENT"
)

var (
	installLocation = "/tmp"
)

// initialize : removes existing symlink to terraform binary
func initialize(createSymlink bool) {

	/* Step 1 */
	/* initilize default binary path for terraform */
	/* assumes that terraform is installed here */
	/* we will find the terraform path instalation later and replace this variable with the correct installed bin path */
	installedBinPath := "/usr/local/bin/terraform"

	/* find terraform binary location if terraform is already installed*/
	cmd := NewCommand("terraform")
	next := cmd.Find()

	/* overrride installation default binary path if terraform is already installed */
	/* find the last bin path */
	for path := next(); len(path) > 0; path = next() {
		installedBinPath = path
	}

	/* remove current symlink if exist*/
	if CheckSymlink(installedBinPath) && createSymlink {
		RemoveSymlink(installedBinPath)
	}

}

// getInstallLocation : get location where the terraform binary will be installed,
// will create a directory in the home location if it does not exist
func getInstallLocation() string {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	/* set installation location */
	installLocation = usr.HomeDir + installPath

	/* Create local installation directory if it does not exist */
	CreateDirIfNotExist(installLocation)

	return installLocation

}

func RemoveVersion(tfversion string, binPath string) bool {
	if runtime.GOOS == "windows" {
		tfversion = tfversion+".exe"
	}
	if !ValidVersionFormat(tfversion) {
		fmt.Printf("The provided terraform version format does not exist - %s. Try `tfswitch -l` to see all available versions.\n", tfversion)
		os.Exit(1)
	}
 	//initialize path
	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	if CheckFileExist(installLocation + installVersion + tfversion) {

		/* remove current symlink if exist*/
		if CheckSymlink(binPath) {
			RemoveSymlink(binPath)
		}

		RemoveFiles(installLocation+installVersion+tfversion)
		fmt.Println("File ", installLocation+installVersion+tfversion, " deleted")
		return true
	}

	fmt.Println("File ", installLocation+installVersion+tfversion, " not found")
	return false
}

//Install : Install the provided version in the argument
func Install(tfversion string, binPath string, createSymlink bool) {

	if runtime.GOOS == "windows" {
		tfversion = tfversion+".exe"
	}

	if !createSymlink {
		fmt.Printf("Skipping symlink creation\n")
	}

	if !ValidVersionFormat(tfversion) {
		fmt.Printf("The provided terraform version format does not exist - %s. Try `tfswitch -l` to see all available versions.\n", tfversion)
		os.Exit(1)
	}

	pathDir := Path(binPath)              //get path directory from binary path
	binDirExist := CheckDirExist(pathDir) //check bin path exist

	if !binDirExist {
		fmt.Printf("Error - Binary path does not exist: %s\n", pathDir)
		fmt.Printf("Create binary path: %s for terraform installation\n", pathDir)
		os.Exit(1)
	}

	initialize(createSymlink)                           //initialize path
	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS


	/* if selected version already exist, */

	if CheckFileExist(installLocation + installVersion + tfversion) {

		if !createSymlink {
			fmt.Printf("Terraform version %q already exists\n", tfversion)
			os.Exit(0)
		}

		/* remove current symlink if exist*/
		if CheckSymlink(binPath) {
			RemoveSymlink(binPath)
		}

		/* set symlink to desired version */
		CreateSymlink(installLocation+installVersion+tfversion, binPath)
		fmt.Printf("Switched terraform to version %q \n", tfversion)
		AddRecent(tfversion) //add to recent file for faster lookup
		os.Exit(0)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := hashiURL + tfversion + "/" + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip"
	zipFile, errDownload := DownloadFromURL(installLocation, url)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		fmt.Println(errDownload)
		os.Exit(1)
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("Unable to unzip downloaded zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

	/* rename unzipped file to terraform version name - terraform_x.x.x */
	RenameFile(installLocation+installFile, installLocation+installVersion+tfversion)

	/* remove zipped file to clear clutter */
	RemoveFiles(installLocation + installVersion + tfversion + "_" + goos + "_" + goarch + ".zip")

	/* if no symlink is to be created, stop here */
	if !createSymlink {
		fmt.Printf("Terraform version %q installed\n", tfversion)
		os.Exit(0)
	}

	/* remove current symlink if exist*/

	if CheckSymlink(binPath) {
		RemoveSymlink(binPath)
	}

	/* set symlink to desired version */
	CreateSymlink(installLocation+installVersion+tfversion, binPath)
	fmt.Printf("Switched terraform to version %q \n", tfversion)
	AddRecent(tfversion) //add to recent file for faster lookup
	os.Exit(0)
}

// AddRecent : add to recent file
func AddRecent(requestedVersion string) {

	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		lines, errRead := ReadLines(installLocation + recentFile)

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return
		}

		for _, line := range lines {
			if !ValidVersionFormat(line) {
				fmt.Println("File dirty. Recreating cache file.")
				RemoveFiles(installLocation + recentFile)
				CreateRecentFile(requestedVersion)
				return
			}
		}

		versionExist := VersionExist(requestedVersion, lines)

		if !versionExist {
			if len(lines) >= 3 {
				_, lines = lines[len(lines)-1], lines[:len(lines)-1]

				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			} else {
				lines = append([]string{requestedVersion}, lines...)
				WriteLines(lines, installLocation+recentFile)
			}
		}

	} else {
		CreateRecentFile(requestedVersion)
	}
}

func GetInstalledVersions() ([]string, error) {
	var versions []string
	installLocation = getInstallLocation()
	files, err := ioutil.ReadDir(installLocation)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name() == "RECENT" {
			continue
		}
		r, err := regexp.MatchString(".zip", file.Name())
		if err == nil && r {
			continue
		}
		versions = append(versions, strings.Trim(file.Name(), "terraform_"))
	}
	var str  sort.StringSlice = versions
	sort.Sort(sort.Reverse(str[:]))
	return str, nil
}


// GetRecentVersions : get recent version from file
func GetRecentVersions() ([]string, error) {

	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	fileExist := CheckFileExist(installLocation + recentFile)
	if fileExist {
		lines, errRead := ReadLines(installLocation + recentFile)
		outputRecent := []string{}

		if errRead != nil {
			fmt.Printf("Error: %s\n", errRead)
			return nil, errRead
		}

		for _, line := range lines {
			/* 	checks if versions in the recent file are valid.
			If any version is invalid, it will be consider dirty
			and the recent file will be removed
			*/
			if !ValidVersionFormat(line) {
				RemoveFiles(installLocation + recentFile)
				return nil, errRead
			}

			/* 	output can be confusing since it displays the 3 most recent used terraform version
			append the string *recent to the output to make it more user friendly
			*/
			outputRecent = append(outputRecent, fmt.Sprintf("%s *recent", line))
		}

		return outputRecent, nil
	}

	return nil, nil
}

//CreateRecentFile : create a recent file
func CreateRecentFile(requestedVersion string) {

	installLocation = getInstallLocation() //get installation location -  this is where we will put our terraform binary file

	WriteLines([]string{requestedVersion}, installLocation+recentFile)
}
