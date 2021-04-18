package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	tm "github.com/buger/goterm"
	"github.com/mgutz/ansi"
)

// Set the correct context for the executable
// so that it can find the config files
func init() {
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(executable)
	if strings.Contains(exPath, "dev-setup") {
		os.Chdir(exPath)
	}
}

type Option struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	Check    string `json:"check"`
	Wait     bool   `json:"wait"`
	Location string `json:"location"`
	Default  bool   `json:"default"`
	Input    string `json:"input"`
}

type Config struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Default string
	Type    string   `json:"type"`
	Options []Option `json:"options"`
}

var phosphorize func(string) string

func parseConfigFile(configFile string) Config {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		fmt.Println("OPEN ERROR", err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("READING ERR", err)
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func isCommandAvailable(command string) bool {
	cmd := exec.Command("/bin/sh", "-c", command)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Creates a question from the passed in config struct
// and then returns a survey question, and a map of the options
func createQuestion(config Config) (survey.Question, map[string]Option) {
	options := make([]string, len(config.Options))
	optionsMap := make(map[string]Option)
	for index, data := range config.Options {

		if data.Check != "" && isCommandAvailable((data.Check)) {
			data.Name = data.Name + " " + phosphorize("[installed]")
		}

		options[index] = data.Name
		optionsMap[data.Name] = data
		if data.Default {
			config.Default = data.Name
		}
	}
	var prompt survey.Prompt
	switch config.Type {
	case "select":
		prompt = &survey.Select{
			Message:  config.Message,
			Options:  options,
			Default:  config.Default,
			PageSize: 20,
		}
	case "multiSelect":
		prompt = &survey.MultiSelect{
			Message:  config.Message,
			Options:  options,
			Default:  config.Default,
			PageSize: 20,
		}
	}
	return survey.Question{Name: config.Name, Prompt: prompt}, optionsMap
}

func printError(error string) {
	errorStyle := ansi.ColorFunc("red")
	fmt.Println(errorStyle(error))
}

func multiSelectQuestion(question survey.Question, optionsMap map[string]Option) []Option {
	var selections []string
	err := survey.AskOne(question.Prompt, &selections)
	if err != nil {
		printError(err.Error())
		return nil
	}

	var commands []Option
	for _, v := range selections {
		commands = append(commands, optionsMap[v])
	}

	return commands
}

func selectQuestion(question survey.Question, optionsMap map[string]Option) Option {
	selection := ""
	err := survey.AskOne(question.Prompt, &selection)
	if err != nil {
		printError(err.Error())
		return Option{}
	}

	return optionsMap[selection]
}

func runCommand(option Option) {
	var command string
	if option.Input != "" {
		fmt.Printf("\n%s, and press [enter] when done\n", phosphorize(option.Input))
		var input string
		fmt.Scanln(&input)

		command = fmt.Sprintf(option.Command, input)
	} else {
		command = option.Command
	}
	cmd := exec.Command("bash", "-c", command)
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Stdin = os.Stdin

	cmd.Run()

	if option.Wait {
		fmt.Print("Press [enter] when the external has finished")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
	}
}

func Question(config Config) {
	question, optionsMap := createQuestion(config)

	switch config.Type {
	case "select":
		selected := selectQuestion(question, optionsMap)
		runCommand(selected)

	case "multiSelect":
		selected := multiSelectQuestion(question, optionsMap)
		for _, option := range selected {
			// TODO: Tidy this
			fmt.Printf("\n%s\n", fmt.Sprintf(phosphorize("#### %s"), option.Name))
			if option.Location != "" {
				selectConfig := parseConfigFile(option.Location)
				Question(selectConfig)
			}

			runCommand(option)
		}
	}
}

func main() {
	tm.Clear()
	tm.MoveCursor(0, 1)
	tm.Flush()
	phosphorize = ansi.ColorFunc("green+h")

	// Start the app proper
	fmt.Print(phosphorize("############## Dev setup ##############\n\n"))

	config := parseConfigFile("./config/core.json")
	Question(config)

	fmt.Println(phosphorize("\n\n👋 Good bye 👋\n\n"))

	tm.Clear()
}
