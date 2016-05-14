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
)

const (
	init_            = iota + 1 // init todo on current directory
	showTodo_                   // no arg
	addTodo_                    // add todo
	all_                        // show all todo
	showArchive_                // show archive
	archive_                    // archives one todo
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
	for i, v := range t.Todo {
		fmt.Printf("[ %d ] %s\n", i, fmt.Sprintf("%s", v))
	}
	todo := []string{}
	// archive := []string{}
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
	for _, v := range t.Todo {
		fmt.Println(v)
	}
}

func (t *Todo) ShowArchive() {
	for _, v := range t.Archive {
		fmt.Println(v)
	}
}

func (t *Todo) ShowAll() {
	fmt.Printf("[ ACTIVE TODO ]\n")
	for _, v := range t.Todo {
		fmt.Println(v)
	}
	fmt.Printf("\n[ ARCHIVED TODO ]\n")
	for _, v := range t.Archive {
		fmt.Println(v)
	}
	fmt.Println("")
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
	return ioutil.WriteFile(fileDir, []byte(s), 0755)
}

func main() {
	var err error
	v, err := filepath.Abs(".")
	if err != nil {
		fmt.Println(unexpectedErrMsg)
		return
	}
	err = os.Chdir(v)
	if err != nil {
		fmt.Println(unexpectedErrMsg)
		return
	}
	currentDir, err = os.Getwd()
	flag.StringVar(&message, "m", "", "set or archive one todo following this option")
	flag.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	command := flagManage()
	switch command {
	case 0: // unexpected error
		fmt.Println("sorry, but something wrong")
		return
	case 1: // init current directory
		err = initTodo()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		return
	default: // other command
	SEEK:
		for {
			_, err = os.Stat("todo")
			switch {
			case err != nil:
				err = chdir() // prevent unexpected error(change directory is not working in if transaction.
				if err != nil {
					return // break when chdir error occured
				}
				if v, err := os.Getwd(); v == "/" || err != nil {
					err = errors.New("no todo file is found. you can add todo by `todo init`")
					break SEEK // break for when file not found
				}
			case err == nil:
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
				fmt.Println("sorry, try again.")
				continue CONFIRM
			}
		}
	}
	if _, err := os.Stat(currentDir + "/todo"); err != nil {
		errors.New("todo file is already exist.")
		return
	}
	err := ioutil.WriteFile(currentDir+"/todo", []byte{}, 0755)
	if err != nil {
		fmt.Println("Permission denied.")
	}
	return err
}

func flagManage() int {
	if len(os.Args) == 1 {
		return showTodo_ // show todo
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "all":
			return all_ // show all todo
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
