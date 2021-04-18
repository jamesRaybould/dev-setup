# Dev setup

A tool that allows for the running of a set of scripts that are encoded within Json files, designed to help get developers up and working faster.

Usages include:

- Installing brew
- Installing OSX Dev tools
- Setting up the /etc/hosts
- Setting up git correctly - connect github to SSO, create ssh key, update git settings to include right username
- Setting up node & node tools
- Cloning your popular repos
- Installing all the useful apps

Uses a JSON config format for easy configuring and allows for chaining of files so you can drill into finer detail

N.B. It requires that a `./config/core.json` is *required*, there is a convenience example included. 


## JSON config

```JSON
{
  "name": "Name of the json file",
  "message": "What to display to the user for this block of questions",
  "type": "[select|multiSelect]", //select allows for a single selection from a list, multiselect lets the user select all, some, or none
  "options":[
    {
      "name": "Displayed to the user to aid selection",
      "command": "echo 'hello world'", // Command to run, must be bash compatible, runs in `bash -c`
      "check": "command -v myCommand", // OPTIONAL. A check to see if the command has already been run/completed
      "wait": true, // OPTIONAL. Bool, if true then will display a message to the user asking them to press [enter] when the command has finished
      "input": "Message to the user asking for input", // OPTIONAL. When input is used a `%s` can be used in the command property - it is NOT SANITISED
      "location": "./configs/my_config.json" // OPTIONAL. If used then it will nest the config
    }
  ] // The actual selections the user can make
}
```

## Example

``` JSON
{
  "name": "Core",
  "message": "What would you like to do today,",
  "type": "multiSelect",
  "options": [
    {
      "name": "Create ~/source directory",
      "command": "mkdir ~/source",
      "check": "ls ~/source"
    },
    {
      "name": "Install brew",
      "command": "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
      "check": "command -v brew"
    },
    {
      "name": "Install OSX dev tools",
      "command": "xcode-select --install",
      "wait": true,
      "check": "xcode-select --install 2>&1 | grep installed"
    },
    {
      "name": "Install node",
      "location": "./config/node.json",
      "check": "command -v node"
    },
    {
      "name": "Setup hosts file",
      "command": "echo '\n127.0.0.1 sub.my-hostname.co.uk' | sudo tee -a /etc/hosts",
      "check": "cat /etc/hosts | grep 'sub.my-hostname.co.uk'"
    },
    {
      "name": "Clone repos",
      "location": "./config/repos.json"
    }
  ]
}

```
