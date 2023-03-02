package main

import (
	"os"
    "os/exec"
	"path/filepath"

	"github.com/caseymrm/menuet"
	"github.com/sqweek/dialog"
)

func open(project string) {
    warning(project)
    os.Chdir(project)
    cmd := exec.Command("neovide")
    cmd.Dir = project
    _ = cmd.Start()
}

const projectKey = "projectKey"

func openProject() menuet.MenuItem {
    var projects []string
    menuet.Defaults().Unmarshal(projectKey, &projects)

    var items []menuet.MenuItem
    for _, p := range projects {
        items = append(items, menuet.MenuItem{
            Text: filepath.Base(p),
            Clicked: func() { open(p) },
        })
    }

    if len(projects) > 0 {
        items = append(items, menuet.MenuItem{ Type: menuet.Separator })
    }

    items = append(items, addProjectMenu())
    return menuet.MenuItem{
        Text: "Open Neovide Project",
        Clicked: nil,
        Children: func() []menuet.MenuItem {
            return items
        },
    }
}

func addProjectMenu() menuet.MenuItem {
    return menuet.MenuItem{
        Text: "Choose new Project",
        Clicked: func() {
            homeDir, _ := os.UserHomeDir()
            project, err := dialog.Directory().SetStartDir(homeDir).Browse()
            if err != nil {
                if err.Error() == "Cancelled" {
                   return 
                }

                warning("Can not open project: " + err.Error())
                return
            }

            var projects []string
            menuet.Defaults().Unmarshal(projectKey, &projects)
            var exists bool
            for _, p := range projects {
                if p == project {
                    exists = true
                }
            }

            if !exists {
                projects = append(projects, project)
            }
            menuet.Defaults().Marshal(projectKey, projects)
            open(project)
        },
    }
}
