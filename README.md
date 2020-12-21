
# Terraform Switcher
Inspired by  [warrensbox/terraform-switcher](https://github.com/warrensbox/terraform-switcher)

<img style="text-allign:center" src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/smallerlogo.png" alt="drawing" width="120" height="130"/>

<!-- ![gopher](https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/logo.png =100x20) -->

The `tfswitch` command line tool lets you switch between different versions of [terraform](https://www.terraform.io/).
If you do not have a particular version of terraform installed, `tfswitch` will download the version you desire.
The installation is minimal and easy.
Once installed, simply select the version you require from the dropdown and start using terraform.

## Installation

`tfswitch` is available for MacOS and Linux based operating systems (Windows experemetal).

### Homebrew

Installation for MacOS/Linux is the easiest with Homebrew. [If you do not have homebrew installed, click here](https://brew.sh/).


```ruby
brew install versus/tap/tfswitch
```

### Linux

Installation for other linux operation systems.

```sh
curl -L https://raw.githubusercontent.com/versus/terraform-switcher/release/install.sh | bash
```

### Build and install SNAP package 


```bash
snap install snapcraft --classic

snap install multipass

snapcraft

snap install terraform-switcher_*.snap --devmode --dangerous

tfswitch -v

multipass stop snapcraft-tfswitch && multipass delete snapcraft-tfswitch && multipass purge

```


### Get binary releases or install from source

Alternatively, you can get releases or install the binary from source [here](https://github.com/versus/terraform-switcher/releases)

## How to use:
### Use dropdown menu to select version
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch.gif#1" alt="drawing" style="width: 370px;"/>

1.  You can switch between different versions of terraform by typing the command `tfswitch` on your terminal.
2.  Select the version of terraform you require by using the up and down arrow.
3.  Hit **Enter** to select the desired version.

The most recently selected versions are presented at the top of the dropdown.

### Supply version on command line
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v4.gif#1" alt="drawing" style="width: 370px;"/>

1. You can also supply the desired version as an argument on the command line.
2. For example, `tfswitch 0.10.5` for version 0.10.5 of terraform.
3. Hit **Enter** to switch.

### See all versions including beta, alpha and release candidates(rc)
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v5.gif#1" alt="drawing" style="width: 370px;"/>

1. Display all versions including beta, alpha and release candidates(rc). 
2. For example, `tfswitch -l` or `tfswitch --list-all` to see all versions.
3. Hit **Enter** to select the desired version.

### Use version.tf file  
If a .tf file with the terraform constrain is included in the current directory, it should automatically download or switch to that terraform version. For example, the following should automatically switch terraform to the lastest version:     
```ruby
terraform {
  required_version = ">= 0.12.9"

  required_providers {
    aws        = ">= 2.52.0"
    kubernetes = ">= 1.11.1"
  }
}
```
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/versiontf.gif#1" alt="drawing" style="width: 370px;"/>


### Use .tfswitch.toml file  (For non-admin - users with limited privilege on their computers)
This is similiar to using a .tfswitchrc file, but you can specify a custom binary path for your terraform installation

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v7.gif#1" alt="drawing" style="width: 370px;"/>     

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v8.gif#1" alt="drawing" style="width: 370px;"/>

1. Create a custom binary path. Ex: `mkdir /Users/warrenveerasingam/bin` (replace warrenveerasingam with your username)
2. Add the path to your PATH. Ex: `export PATH=$PATH:/Users/warrenveerasingam/bin` (add this to your bash profile or zsh profile)
3. Pass -b or --bin parameter with your custom path to install terraform. Ex: `tfswitch -b /Users/warrenveerasingam/bin/terraform 0.10.8 `
4. Optionally, you can create a `.tfswitch.toml` file in your home directory for global settings.
5. Your `.tfswitch.toml` file should look like this:
```ruby
bin = "/Users/versus/bin/terraform"
version = "0.11.3"
```
4. Run `tfswitch` and it should automatically install the required terraform version in the specified binary path

Alternatively, you can generate .tfswitch.toml in current directory just use `tfswitch --init ` or `tfswitch --init 0.13.2`


### Use .tfswitchrc file
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/tfswitch-v6.gif#1" alt="drawing" style="width: 370px;"/>

1. Create a `.tfswitchrc` file containing the desired version
2. For example, `echo "0.10.5" >> .tfswitchrc` for version 0.10.5 of terraform
3. Run the command `tfswitch` in the same directory as your `.tfswitchrc`

### Use environment variable TFSWITCH_PATH

1. Create a `TFSWITCH_PATH` environment variable with your custom path to install terraform Ex: `export  TFSWITCH_PATH=/Users/versus/bin/terraform`
2. Run the command `tfswitch` 

#### *Instead of a `.tfswitchrc` file, a `.terraform-version` file may be used for compatibility with [`tfenv`](https://github.com/tfutils/tfenv#terraform-version-file) and other tools which use it*

**Automatically switch with bash**

Add the following to the end of your `~/.bashrc` file:
(Use either `.tfswitchrc` or `.tfswitch.toml` or `.terraform-version`)

```sh
cdtfswitch(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tfswitchrc" ]; then
    tfswitch
  fi
}
alias cd='cdtfswitch'
```

**Automatically switch with zsh**

Add the following to the end of your `~/.zshrc` file:

```sh
load-tfswitch() {
  local tfswitchrc_path=".tfswitchrc"

  if [ -f "$tfswitchrc_path" ]; then
    tfswitch
  fi
}
add-zsh-hook chpwd load-tfswitch
load-tfswitch
```
> NOTE: if you see an error like this: `command not found: add-zsh-hook`, then you might be on an older version of zsh (see below), or you simply need to load `add-zsh-hook` by adding this to your `.zshrc`:
>    ```
>    autoload -U add-zsh-hook
>    ```

*older version of zsh*
```sh
cd(){
  builtin cd "$@";
  cdir=$PWD;
  if [ -e "$cdir/.tfswitchrc" ]; then
    tfswitch
  fi
}
```

### Jenkins setup
<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/jenkins_tfswitch.png" alt="drawing" style="width: 170px;"/>

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/versus/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

./bin-directory/tfswitch
```

If you have limited permission, try:

```sh
#!/bin/bash 

echo "Installing tfswitch locally"
wget https://raw.githubusercontent.com/versus/terraform-switcher/release/install.sh 
chmod 755 install.sh
./install.sh -b bin-directory

CUSTOMBIN=`pwd`/bin             #set custom bin path
mkdir $CUSTOMBIN                #create custom bin path
export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

./bin-directory/tfswitch -b $CUSTOMBIN/terraform 0.11.7

terraform -v                    #testing version
```

### Circle CI setup

<img src="https://s3.us-east-2.amazonaws.com/kepler-images/warrensbox/tfswitch/circleci_tfswitch.png" alt="drawing" style="width: 280px;"/>


Example config yaml
```yaml
version: 2
jobs:
  build:
    docker:
      - image: ubuntu

    working_directory: /go/src/github.com/versus/terraform-switcher

    steps:
      - checkout
      - run: 
          command: |    
            set +e   
            apt-get update 
            apt-get install -y wget 
            rm -rf /var/lib/apt/lists/*

            echo "Installing tfswitch locally"

            wget https://raw.githubusercontent.com/versus/terraform-switcher/release/install.sh 
            chmod 755 install.sh
            ./install.sh -b bin-directory

            CUSTOMBIN=`pwd`/bin             #set custom bin path
            mkdir $CUSTOMBIN                #create custom bin path
            export PATH=$PATH:$CUSTOMBIN    #Add custom bin path to PATH environment

            ./bin-directory/tfswitch -b $CUSTOMBIN/terraform 0.11.7

            terraform -v                    #testing version
```

### Adds two new flags:

 - no-symlink (n)
 Makes tfswitch detect and install a version only, without creating a permanent symlink to it. This allows tfswitch to be used as a helper when dynamically switching Terraform versions in environments where multiple tasks run with different versions in parallel, such as CI/CD agents.

 - quiet (q)
Makes tfswitch switch to a version if it's detected or specified as an argument, but prevents it from prompting the user for input if none is found or provided. This makes tfswitch more automation friendly.

#### Example use case on a CI/CD agent
An agent host is scheduled to run multiple Terraform Plan/Apply tasks simultaneously from different repositories, which each contain different Terraform configurations, using different versions of Terraform.

With a regular symlink approach, they would have to be queued and run synchronously, as a version switch would have to be made between each run.

However, if the terraform binary (or symlink) is replaced with a script (example below) that catches the arguments sent to terraform, performs a tfswitch using both the no-symlink and quiet flags, then catches the output from tfswitch and uses the desired binary to run the specific task, this can be done dynamically with each run, in parallel.

This utilizes tfswitch's brilliant ability to detect and install the required version, while allowing concurrency.

Example PoC script

```sh
#!/bin/bash

TF_DEFAULT_VERSION=0.13.2
TF_BIN_LOCATION=~/.terraform.versions

TF_SWITCH=`tfswitch-test -nq`
TF_DETECTED_VERSION=`echo "$TF_SWITCH" | awk -F '\"|\"' '{print $2}' | grep "\S"`

echo "$TF_SWITCH"

if [ "$TF_DETECTED_VERSION" == "" ]; then
    TF_DETECTED_VERSION=$TF_DEFAULT_VERSION
fi

TF_BIN="$TF_BIN_LOCATION/terraform_$TF_DETECTED_VERSION"

echo "Using version: $TF_DETECTED_VERSION"
echo ""

$TF_BIN "$@"
```
Suggestions and improvements are welcome.

-Sindre(c) https://github.com/sindrel

## Issues

Please open  *issues* here: [New Issue](https://github.com/versus/terraform-switcher/issues)
