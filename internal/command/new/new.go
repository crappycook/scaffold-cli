package new

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/crappycook/scaffold-cli/internal/config"
	"github.com/crappycook/scaffold-cli/internal/utils"
	"github.com/spf13/cobra"
)

type Project struct {
	ModuleName string `survey:"name"`
	Dirname    string
}

var CmdNew = &cobra.Command{
	Use:     "new",
	Example: "scaffold-cli new demo",
	Short:   "Create a new project.",
	Long:    `Create a new project from your template repo layout.`,
	Run:     run,
}

var (
	repoURL string
)

func init() {
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
}

func NewProject() *Project {
	return &Project{}
}

func run(cmd *cobra.Command, args []string) {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "Input your module name: ",
			Help:    "module name like github.com/gin-gonic/gin",
			Suggest: nil,
		}, &p.ModuleName, survey.WithValidator(survey.Required))
		if err != nil {
			return
		}
	} else {
		p.ModuleName = args[0]
	}

	strArr := strings.Split(p.ModuleName, "/")
	p.Dirname = strArr[len(strArr)-1]

	// clone repo
	yes, err := p.cloneTemplate()
	if err != nil || !yes {
		return
	}

	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	err = p.modTidy()
	if err != nil || !yes {
		return
	}

	p.rmTemplateGit()

	p.projectGitInit()

	fmt.Printf("🎉 Project \u001B[36m%s\u001B[0m created successfully!\n\n", p.ModuleName)
}

func (p *Project) cloneTemplate() (bool, error) {
	stat, _ := os.Stat(p.Dirname)
	if stat != nil {
		var overwrite = false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Folder %s already exists, do you want to overwrite it?", p.Dirname),
			Help:    "Remove old project and create new project.",
		}
		err := survey.AskOne(prompt, &overwrite)
		if err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
		err = os.RemoveAll(p.Dirname)
		if err != nil {
			fmt.Println("remove old project error: ", err)
			return false, err
		}
	}

	repo := repoURL
	if len(repoURL) == 0 {
		repo = config.DefaultRepoURL
	}

	fmt.Printf("git clone %s\n", repo)
	cmd := exec.Command("git", "clone", repo, p.Dirname)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", repo, err)
		return false, err
	}
	return true, nil
}

func (p *Project) replacePackageName() error {
	packageName := utils.GetProjectName(p.Dirname)

	err := p.replaceFiles(packageName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ModuleName)
	cmd.Dir = p.Dirname
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod edit error: ", err)
		return err
	}
	return nil
}

func (p *Project) modTidy() error {
	fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.Dirname
	if err := cmd.Run(); err != nil {
		fmt.Println("go mod tidy error: ", err)
		return err
	}
	return nil
}

func (p *Project) rmTemplateGit() {
	os.RemoveAll(p.Dirname + "/.git")
}

func (p *Project) projectGitInit() {
	cmd := exec.Command("git", "init")
	cmd.Dir = p.Dirname
	if err := cmd.Run(); err != nil {
		fmt.Println("git init error: ", err)
	}
}

func (p *Project) replaceFiles(packageName string) error {
	err := filepath.Walk(p.Dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newData := bytes.ReplaceAll(data, []byte(packageName), []byte(p.ModuleName))
		if err := os.WriteFile(path, newData, 0644); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("walk file error: ", err)
		return err
	}
	return nil
}
