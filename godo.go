package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var Config struct {
	taskfile string
}

func init() {
	usr := getUser()
	default_taskfile := usr.HomeDir + "/.godo.taskpaper"
	flag.StringVar(&Config.taskfile, "taskfile", default_taskfile, "File for storing godo tasks (taskpaper format)")
}

func getUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

func main() {
	args := os.Args
	if len(args) == 1 {
		cmd := exec.Command("nvim", Config.taskfile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}

	} else if args[1][0] != '-' {
		task_arr := args[1:]
		task := strings.Join(task_arr, " ")

		f, err := os.OpenFile(Config.taskfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err = f.WriteString("- " + task + "\n"); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Added:", task)
	} else {
		search := flag.String("s", "", "Search for term in tasks")
		flag.Parse()

		cmd := exec.Command("nvim", "-c", "lvimgrep "+*search+" "+Config.taskfile, "-c", "lopen")
		// cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}

	}
}
