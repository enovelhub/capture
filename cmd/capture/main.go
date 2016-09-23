package main

import (
	"fmt"
	"github.com/enovelhub/capture/cmd/capture/fetch"
	"github.com/enovelhub/capture/cmd/capture/gen"
	"io"
	"os"
	"sync"
	"text/template"
)

var (
	cmdsMx = &sync.Mutex{}
	cmds   = make(map[string]Command)
)

func Register(cmd Command) {
	cmdsMx.Lock()
	defer cmdsMx.Unlock()

	cmds[cmd.Name()] = cmd
}

func Find(name string) (cmd Command, exist bool) {
	cmdsMx.Lock()
	defer cmdsMx.Unlock()

	cmd, exist = cmds[name]
	return
}

func List() []Command {
	cmdsMx.Lock()
	defer cmdsMx.Unlock()
	var keys []string
	for k, _ := range cmds {
		keys = append(keys, k)
	}

	retcmds := make([]Command, len(keys))
	for i, k := range keys {
		retcmds[i] = cmds[k]
	}

	return retcmds
}

type Command interface {
	Name() string
	Desc() string
	Run([]string) error
}

func init() {
	Register(gen.New())
	Register(fetch.New())

}

func main() {
	argc := len(os.Args)

	if argc < 2 {
		Usage(os.Stderr)
		os.Exit(1)
	}

	subName := os.Args[1]
	if subcmd, exist := Find(subName); exist {
		err := subcmd.Run(os.Args[1:])
		if err != nil {
			errorExit(err.Error())
		}
		return
	}

	Usage(os.Stderr)
	os.Exit(1)
}

const usageTmpl = `Usage of {{ .self }} cmd [args]:
	{{ range .cmds }}{{ .Name }}:	{{ .Desc }}
	{{ end }}
`

func Usage(w io.Writer) {
	fmt.Println(List())
	tmpl := template.Must(
		template.New("usage").Parse(usageTmpl))
	tmpl.Execute(w, map[string]interface{}{
		"self": os.Args[0],
		"cmds": List(),
	})
}

func errorExit(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
