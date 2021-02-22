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

type Templater struct {
	RepositoryGitURL string
	RepositoryBranch string
	RepositoryFsPath string
	RepositoryPullInterval time.Duration
	TheTemplates []markdowntemplates.Template
}

var theTemplater *Templater

func NewTemplater(RepositoryGitURL, RepositoryBranch, RepositoryFsPath string, RepositoryPullInterval time.Duration) error {

	templater := Templater{
		RepositoryGitURL: RepositoryGitURL,
		RepositoryBranch: RepositoryBranch,
		RepositoryFsPath: RepositoryFsPath,
		RepositoryPullInterval: RepositoryPullInterval,
	}

	theTemplater = &templater

	err := templater.CloneTemplatesFromGitHub()
	if err != nil {
		return err
	}

	// Start routine to fetch latest templates
	go templater.FetchLatestTemplatesFromGithub()
	return nil
}

func (templater *Templater) FetchLatestTemplatesFromGithub() {
	ticker := time.NewTicker(templater.RepositoryPullInterval)

	for range ticker.C {
		cmd := exec.Command("git", "pull", "origin", templater.RepositoryBranch)
		cmd.Dir = templater.RepositoryFsPath
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing %v\n", cmd.Args)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error pulling git repository %v\n", err)
		}
		err = templater.slurpTemplates(templater.RepositoryFsPath)
		if err != nil {
			fmt.Printf("Error reading templates %v\n", err)
		}
	}
}

func (templater *Templater) CloneTemplatesFromGitHub() error {
	// Check if templates folder exists
	if _, err := os.Stat(templater.RepositoryFsPath); os.IsNotExist(err) {
		fmt.Println("Templates do not exists; Go and clone repository")
		cmd := exec.Command("git", "clone", "-b", templater.RepositoryBranch, templater.RepositoryGitURL, templater.RepositoryFsPath)
		cmd.Stderr = os.Stderr
		fmt.Printf("Executing %v\n", cmd.Args)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return templater.slurpTemplates(templater.RepositoryFsPath)
}

func (templater *Templater) slurpTemplates(templatesRepositoryFsPath string) error {

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

	if len(templates) == 0 {
		return fmt.Errorf("No templates imported from template repository")
	}

	templater.TheTemplates = templates

	return nil
}

// GetTemplate returns the template matching the templateName and the language
func GetTemplate(templateName, language string) (markdowntemplates.Template, error) {
	if theTemplater == nil {
		err := fmt.Errorf("Templater is not available")
		if err != nil {
			return markdowntemplates.Template{}, err
		}
	}

	// check for direct or prefix (de_DE -> de) match
	for _, template := range theTemplater.TheTemplates {
		if template.Name == templateName && (template.Language == language || strings.HasPrefix(language, template.Language) == true) {
			return template, nil
		}
	}

	return markdowntemplates.Template{}, fmt.Errorf("Template not available (template: '%s', language: '%s')", templateName, language)
}
