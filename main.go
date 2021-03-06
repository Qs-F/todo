package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	init_            = iota + 1 // init todo on current directory
	showTodo_                   // no arg
	addTodo_                    // add todo
	all_                        // show all todo
	showArchive_                // show archive
	archive_                    // archives one todo
	help_                       // show help
	unexpectedErrMsg = "sorry, but something occured:("
)

var (
	currentDir string
	fileDir    string
	message    string
)

type Todo struct {
	Todo    []string
	Archive []string
}

func parse() (t *Todo, err error) {
	f, err := os.Open(fileDir)
	defer f.Close()
	t = &Todo{}
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(f)
	i := 0
SCANNER:
	for scanner.Scan() {
		if scanner.Text() == "[ archived ]" {
			i++
		} else {
			if i == 0 {
				t.Todo = append(t.Todo, scanner.Text())
			} else if i == 1 {
				t.Archive = append(t.Archive, scanner.Text())
			} else {
				break SCANNER
			}
		}
	}
	return
}

func (t *Todo) AddTodo(s string) {
	t.Todo = append(t.Todo, s)
}

func (t *Todo) ArchiveTodo() {
	// display all todo
	if len(t.Todo) == 0 {
		fmt.Println("There is no todo🍻 ")
		return
	}
	for i, v := range t.Todo {
		fmt.Printf(color.RedString("[ %d ] ")+"%s\n", i, fmt.Sprintf("%s", v))
	}
	todo := []string{}
CONFIRM:
	for { // choosing form
		fmt.Printf("choose one: ")
		s := ""
		fmt.Scan(&s)                                                  // input
		if i, err := strconv.Atoi(s); err == nil && len(t.Todo) > i { // input is number and valid item.
			for j, v := range t.Todo { // choose from all todo
				if j != i { // current loop number is not selected
					todo = append(todo, v)
				} else if j == i { // current loop number is selected
					t.Archive = append(t.Archive, v)
				}
			}
			break CONFIRM
		} else {
			fmt.Println("sorry, try again.")
			continue CONFIRM
		}
	}
	t.Todo = todo
}

func (t *Todo) ShowTodo() {
	if len(t.Todo) == 0 {
		fmt.Println("nothing is here. have a good day:)")
		return
	}
	fmt.Printf(color.YellowString("[ Active Todo ]\n"))
	for _, v := range t.Todo {
		fmt.Println(color.RedString("+ ") + v)
	}
}

func (t *Todo) ShowArchive() {
	if len(t.Archive) == 0 {
		fmt.Println("nothing is here.")
		return
	}
	fmt.Printf(color.YellowString("[ Archived Todo ]\n"))
	for _, v := range t.Archive {
		fmt.Println(color.RedString("+ ") + v)
	}
}

func (t *Todo) ShowAll() {
	fmt.Printf(color.YellowString("[ Active Todo ]\n"))
	if len(t.Todo) == 0 {
		fmt.Println("nothing is here.")
	}
	for _, v := range t.Todo {
		fmt.Println(color.RedString("+ ") + v)
	}
	fmt.Printf(color.YellowString("[ Archived Todo ]\n"))
	if len(t.Archive) == 0 {
		fmt.Println("nothing is here.")
	}
	for _, v := range t.Archive {
		fmt.Println(color.RedString("+ ") + v)
	}
}

func (t *Todo) Save() error {
	// fmt.Println("> " + fileDir)
	s := ""
	for _, v := range t.Todo {
		s += fmt.Sprintf("%s\n", formatMessage(v))
	}
	if len(t.Archive) >= 1 {
		s += "[ archived ]\n"
	}
	for _, v := range t.Archive {
		s += fmt.Sprintf("%s\n", formatMessage(v))
	}
	return ioutil.WriteFile(fileDir, []byte(s), 0644)
}

func main() {
	var err error
	v, err := filepath.Abs(".") // get user current directory
	if err != nil {
		fmt.Println(unexpectedErrMsg)
		return
	}
	err = os.Chdir(v) // move to user current directory in go app
	if err != nil {
		fmt.Println(unexpectedErrMsg)
		return
	}
	currentDir, err = os.Getwd() // get current working directory (in this time, it must equal current user directory and go app working directory.
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	flag.StringVar(&message, "m", "", "set your message")
	flag.Parse()
	command := flagManage() // get what subcommand did user use
	switch command {
	case 0: // unexpected error
		fmt.Println("sorry, but something wrong")
		return
	case init_: // init current directory
		err = initTodo()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		return
	case help_:
		fmt.Println(`
todo               show active todo
todo YOURMESSAGE   add todo
todo archive       show archived todo
todo archive add   you can archive one todo
todo all           show all todo(active and archived)
todo help          open this help
`)
	default: // other command
		// START: LABEL: SEEK
		var info os.FileInfo
	SEEK:
		for {
			info, err = os.Stat("todo") // check working directory's todo file
			switch {
			case err != nil || (err == nil && !info.Mode().IsRegular()): // if there isn't todo file or it is directory.
				err = chdir() // prevent unexpected error(change directory is not working in if transaction.
				if err != nil {
					return // break when chdir error occured
				}
				if v, err = os.Getwd(); v == "/" || err != nil { // get working directory, and check whether it is root directory(it must be finished from infinite loop)
					err = errors.New("no todo file is found. you can add todo by `todo init`")
					break SEEK // break for when file not found
				}
			case err == nil && info.Mode().IsRegular(): // there is a file or directory, and check it is a file.
				fileDir, err = os.Getwd()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fileDir += "/todo"
				break SEEK // break for when file being found
			default:
				fmt.Println("sorry, but something wrong.")
				return // unexpected error
			}
		} // LABEL: SEEK
		// END: LABEL: SEEK
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		t, err := parse()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		switch command {
		case showTodo_:
			t.ShowTodo()
		case addTodo_:
			t.AddTodo(message)
		case all_:
			t.ShowAll()
		case showArchive_:
			t.ShowArchive()
		case archive_:
			t.ArchiveTodo()
		}
		err = t.Save()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func formatMessage(s string) string {
	s = strings.Replace(s, "\r", " ", -1)
	s = strings.Replace(s, "\r\n", " ", -1)
	return strings.Replace(s, "\n", " ", -1)
}

func chdir() error {
	return os.Chdir("..")
}

func initTodo() error {
	if currentDir == "/" {
	CONFIRM:
		for {
			fmt.Printf("current directory is root directory. this is not recommended. Are you sure? [y/N]: ")
			s := ""
			fmt.Scan(&s)
			if s == "y" || s == "Y" {
				break CONFIRM
			} else {
				fmt.Println("Aborted")
				return errors.New("abort")
			}
		}
	}
	if _, err := os.Stat(currentDir + "/todo"); err == nil { // existing a todo file or directory and it is fiel
		return errors.New("todo file is already exist.")
	}
	err := ioutil.WriteFile(currentDir+"/todo", []byte{}, 0644)
	if err != nil {
		fmt.Println("Permission denied.")
		return err
	}
	fmt.Println("Success! 🍻")
	return nil
}

func flagManage() int {
	if len(os.Args) == 1 {
		return showTodo_ // show todo
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "all":
			return all_ // show all todo
		case "help":
			return help_
		case "archive":
			if len(os.Args) > 2 {
				if os.Args[2] == "add" {
					return archive_ // archive one todo
				} else {
					return 0
				}
			} else {
				return showArchive_ // show archives
			}
		case "init":
			return init_ // init todo
		default:
			m := ""
			for _, v := range os.Args[1:] {
				m += v + " "
			}
			if message == "" {
				message = formatMessage(m)
			}
			return addTodo_ // add todo
		}
	}
	return 0
}
