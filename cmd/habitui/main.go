//nolint:forbidigo //prints for command line client are not debug statements
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/bazko1/habitui/client"
	"github.com/bazko1/habitui/habit"
	"github.com/bazko1/habitui/tui"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultFile string = ".habitui.json"

// getIOFiles returns filenames for reading and writing habits data based on
// tasksFile provided by user and files found in system.
// If not tasksFile is provided there is check for .habitui.json
// firstly in $(pwd) equivalent and then at $HOME.
// If no file is found at pwd or home location then
// new file will be written to user home.
func getIOFiles(tasksFile string) (string, string) {
	if tasksFile != "" {
		return tasksFile, tasksFile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return defaultFile, defaultFile
	}

	if _, err := os.Stat(defaultFile); err == nil {
		return defaultFile, defaultFile
	}

	homef := home + "/" + defaultFile
	if _, err := os.Stat(homef); err == nil {
		return homef, homef
	}

	return "", homef
}

func main() {
	tasksFile := flag.String("data", "", "file name for loading/saving tasks data")
	disableDebug := flag.Bool("no-debug", false, "do not log debug data to file")

	remoteAddress := flag.String("remote-server", "localhost:3000", "address of remote server for loading saving tasks data")
	remoteUser := flag.String("remote-user", "", "username for remote login")
	remotePassword := flag.String("remote-password", "", "password for remote login")
	enableRemote := flag.Bool("enable-remote", false, "enable storing data into remote location")
	flag.Parse()

	var tasks habit.TaskList

	var outputFile string

	if !*enableRemote {
		var inputFile string
		inputFile, outputFile = getIOFiles(*tasksFile)
		file, err := os.ReadFile(inputFile)

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("failed to open tasks file '%s': %v\n", inputFile, err)
			os.Exit(1)
		}

		if !errors.Is(err, os.ErrNotExist) {
			tasks, err = habit.JSONLoadTasks(file)
			if err != nil {
				fmt.Println("failed to load tasks:", err)
				os.Exit(1)
			}
		}
	}

	if *enableRemote {
		if *remoteUser == "" || *remotePassword == "" {
			fmt.Println("Username and password must be provided for remote connection.")
			os.Exit(1)
		}

		remoteClient := client.HTTPClient{
			Address:  *remoteAddress,
			Username: *remoteUser,
			Password: *remotePassword,
		}

		var err error

		tasks, err = remoteClient.LoadTasksOrCreateUser()
		if err != nil {
			fmt.Println("Failed to connect to remote error: ", err)
			os.Exit(1)
		}
	}

	if !*disableDebug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		log.SetOutput(io.Discard)
	}

	logger := log.Default()
	logger.Println("starting tui program")

	model := tui.NewTuiModel(tasks)
	prog := tea.NewProgram(model)

	out, err := prog.Run()
	if err != nil {
		logger.Printf("Running tui error: %v", err)
	}

	model, _ = out.(tui.Model)

	defer func() {
		if *enableRemote {
			// TODO: Refactor to reuse previous remote client and
			// it might be better as func arguments instead of
			// methods on strucs - think about that.
			remoteClient := client.HTTPClient{
				Address:  *remoteAddress,
				Username: *remoteUser,
				Password: *remotePassword,
			}

			if err := remoteClient.SaveUserTasks(model.Tasks()); err != nil {
				logger.Printf("failed to save tasks to remote: %v", err)
				os.Exit(1)
			}
		} else {
			err := habit.JSONSaveTasks(outputFile, model.Tasks())
			if err != nil {
				logger.Printf("failed to save tasks: %v", err)
				os.Exit(1)
			}
		}

		logger.Println("saved state closing")
		os.Exit(0)
	}()
}
