package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"

func checkTitles(lines []string) bool {
	checker := 0
	for _, line := range lines {
		if line == "### Instructions" {
			checker++
		}
		if line == "### Usage" && strings.HasPrefix(line, "### ") {
			checker++
		}
	}
	return checker == 2
}

func getDescription(lines []string) []string {
	description := []string{}
	for i := 0; i < len(lines); i++ {
		if lines[i] == "### Instructions" {
			for j := i + 1; j < len(lines) && (lines[j] != "### Expected function" && lines[j] != "### Usage"); j++ {
				description = append(description, lines[j])
			}
			break
		}
	}
	return description
}

func isUpper(line string) bool {
	if line == "" {
		return true
	}
	i := 0
	for ; i < len(line); i++ {
		if line[i] == ' ' || line[i] == '-' {
			continue
		} else if line[i] == '\'' || line[i] == '"' {
			return true
		}
		break
	}
	if i != len(line) && line[i] >= 'A' && line[i] <= 'Z' {
		return true
	}
	return false
}

func FunctionCheck(lines []string) bool {
	checker := 0
	desc := 0
	for _, line := range lines {
		if line == "### Expected function" {
			checker++
		} else if line == "Here is a possible program to test your function :" {
			desc += 1
		} else if line == "And its output :" {
			desc += 2
		}
	}
	if checker == 1 && desc == 2 {
		fmt.Println(Red, "Usage [F] : add this line after Usage title :", Yellow, "[Here is a possible program to test your function :]", Reset)
	} else if checker == 1 && desc != 3 {
		fmt.Println(Red, "Usage [F] : add this line after Usage title :", Yellow, " [Here is a possible program to test your function :] ", Reset, "and this line after the main func ", Yellow, "[And its output :]", Reset)
	}
	return checker == 1
}

func checkDescription(description []string) (bool, string) {
	_type := "Nand"
	for _, line := range description {
		if strings.HasPrefix(line, "		") {
			fmt.Println(Red, "Description [F] : should have only one level of indentation", Reset)
			return false, line
		}
		if !isUpper(line) {
			fmt.Println(Red, "Description [F] : should begin with a upper case letter", Reset)
			return false, line
		}
		if strings.Contains(line, "function") {
			_type = "F"
		} else if strings.Contains(line, "program") {
			_type = "P"
		}
	}
	return true, _type
}

func checkDirName(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] >= 'a' && str[i] <= 'z' || (str[i] >= '0' && str[i] <= '9') {
			continue
		}
		return false
	}
	return true
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println(Red, "Usage : go run . <Readme file>  <Test dir>", Reset)
		return
	}
	file, err := os.ReadFile(os.Getenv("PWD") + "/" + os.Args[1])
	dir_name := strings.Split(os.Getenv("PWD"), "/")[len(strings.Split(os.Getenv("PWD"), "/"))-1]
	_, err_sol := os.ReadFile(os.Getenv("PWD") + "/" + os.Args[2] + "/solutions/" + dir_name + ".go")
	_, err_tes := os.ReadFile(os.Getenv("PWD") + "/" + os.Args[2] + "/tests/" + dir_name + "_test/main.go")
	if err != nil {
		fmt.Println(Red, err, Reset)
		return
	}
	if err_sol != nil {
		fmt.Println(Red, "Solution [F] : add this file in the solutions folder", err_sol, Reset)
		return
	}
	if err_tes != nil {
		fmt.Println(Red, "Test [F] : add this file in the test folder", err_tes, Reset)
		return
	}

	// run command of gofmt
	cmd := exec.Command("gofmt", "-w", os.Getenv("PWD")+"/"+os.Args[2]+"/solutions/"+dir_name+".go")
	cmd.Run()
	cmd = exec.Command("gofmt", "-w", os.Getenv("PWD")+"/"+os.Args[2]+"/tests/"+dir_name+"_test/main.go")
	cmd.Run()
	fmt.Println(Green, "gofmt [T] : gofmt has been run on the files", Reset)

	if checkDirName(dir_name) == false {
		fmt.Println(Red, "Directory [F] : should be in lower case and without in spliter", Reset)
	}
	lines := strings.Split(string(file), "\n")
	if strings.HasPrefix(lines[0], "## ") || strings.HasPrefix(lines[1], "## ") {
		fmt.Println(Green, "Title [T]", Reset)
	} else {
		fmt.Println(Red, "Title [F] : should have ## and titile in line 0", Reset)
	}

	if checkTitles(lines) == false {
		fmt.Println(Red, "SubTitles [F] : should add subtitle for instruction and usage the subtitle should have ###", Reset)
	} else {
		fmt.Println(Green, "SubTitles [T]", Reset)
	}

	_descrition := getDescription(lines)
	check_desc, _type := checkDescription(_descrition)
	if check_desc {
		fmt.Println(Green, "Description [T]", Reset)
	}
	if _type == "Nand" {
		fmt.Println(Green, "Description [F] : should have a type of function or program", Reset)
	}
	if _type == "F" {
		FunctionCheck(lines)
	}
}
