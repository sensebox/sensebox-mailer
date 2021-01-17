package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	markdowntemplates "github.com/sensebox/sensebox-mailer-templates"
)

const (
	templatesRepositoryGitURL = "https://github.com/sensebox/sensebox-mailer-templates.git"
	templatesRepositoryBranch = "main"
	templatesRepositoryFsPath = "./mailer-templates"
)

var theTemplates []markdowntemplates.Template

func FetchLatestTemplatesFromGithub() {
	ticker := time.NewTicker(60000 * time.Millisecond)

	for range ticker.C {
		cmd := exec.Command("git", "pull", "origin", templatesRepositoryBranch)
		cmd.Dir = templatesRepositoryFsPath
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing %v", cmd.Args)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error pulling git repository %v", err)
		}
		err = slurpTemplates()
		if err != nil {
			fmt.Printf("Error reading templates %v", err)
		}
	}
}

func CloneTemplatesFromGitHub() error {
	// Check if templates folder exists
	if _, err := os.Stat(templatesRepositoryFsPath); os.IsNotExist(err) {
		fmt.Println("Templates do not exists; Go and clone repository")
		cmd := exec.Command("git", "clone", "-b", templatesRepositoryBranch, templatesRepositoryGitURL, templatesRepositoryFsPath)
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing %v", cmd.Args)
		return cmd.Run()
	}
	return slurpTemplates()
}

func slurpTemplates() error {

	templates := []markdowntemplates.Template{}

	err := filepath.Walk(templatesRepositoryFsPath+"/templates", func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// check if it is a regular file (not dir)
		if info.Mode().IsRegular() && strings.HasSuffix(path, ".md") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			tpls, err := markdowntemplates.Slurp(file)
			if err != nil {
				return err
			}
			templates = append(templates, tpls...)

		}

		return nil
	})

	if err != nil {
		return err
	}

	theTemplates = templates

	return nil
}

// GetTemplate returns the template matching the templateName and the language
func GetTemplate(templateName, language string) (markdowntemplates.Template, error) {
	if theTemplates == nil {
		err := slurpTemplates()
		if err != nil {
			return markdowntemplates.Template{}, err
		}
	}

	// check for direct or prefix (de_DE -> de) match
	for _, template := range theTemplates {
		if template.Name == templateName && (template.Language == language || strings.HasPrefix(language, template.Language) == true) {
			return template, nil
		}
	}

	return markdowntemplates.Template{}, fmt.Errorf("Template not available (template: '%s', language: '%s')", templateName, language)
}
