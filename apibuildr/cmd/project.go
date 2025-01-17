package cmd

import (
	"fmt"
	"github.com/focks/apibuildr/apibuildr/cmd/tpl"
	"os"

	"text/template"
)

type Project struct {
	Name         string
	AbsolutePath string
	Wd           string
	PackageName  string
}

func (p *Project) Create() error {
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, 0754); err != nil {
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create cmd/server.go
	if _, err = os.Stat(fmt.Sprintf("%s/cmd", p.AbsolutePath)); os.IsNotExist(err) {
		CheckError(os.Mkdir(fmt.Sprintf("%s/cmd", p.AbsolutePath), 0751))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/cmd/server.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.ServerTemplate())))
	err = rootTemplate.Execute(rootFile, p)
	if err != nil {
		return err
	}

	return nil
}
