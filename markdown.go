package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const DefaultFilename = "TASKS.md"

func detectTasksFile() string {
	candidates := []string{"TASKS.md", "tasks.md", "TODO.md", "todo.md"}
	for _, name := range candidates {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	return DefaultFilename
}

func ParseTasksFile(filename string) (*TaskList, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	tl := &TaskList{Filename: filename}
	content := string(data)
	lines := strings.Split(content, "\n")

	currentCategory := "General"
	taskRegex := regexp.MustCompile(`^\s*-\s+\[([ xX])\]\s+(.+)$`)
	headerRegex := regexp.MustCompile(`^##\s+(.+)$`)

	for _, line := range lines {
		if m := headerRegex.FindStringSubmatch(line); m != nil {
			currentCategory = strings.TrimSpace(m[1])
			continue
		}

		if m := taskRegex.FindStringSubmatch(line); m != nil {
			done := m[1] == "x" || m[1] == "X"
			desc := strings.TrimSpace(m[2])
			num := tl.nextNumber()
			tl.Tasks = append(tl.Tasks, Task{
				Number:      num,
				Description: desc,
				Done:        done,
				Category:    currentCategory,
			})
			tl.addCategory(currentCategory)
		}
	}

	return tl, nil
}

func WriteTasksFile(tl *TaskList) error {
	grouped := tl.TasksByCategory()

	var cats []string
	for cat := range grouped {
		cats = append(cats, cat)
	}
	sort.Strings(cats)

	var sb strings.Builder
	sb.WriteString("# Tasks\n\n")

	for _, cat := range cats {
		sb.WriteString(fmt.Sprintf("## %s\n\n", cat))
		tasks := grouped[cat]
		for _, t := range tasks {
			status := " "
			if t.Done {
				status = "x"
			}
			sb.WriteString(fmt.Sprintf("- [%s] %s\n", status, t.Description))
		}
		sb.WriteString("\n")
	}

	return os.WriteFile(tl.Filename, []byte(sb.String()), 0644)
}

func InitTasksFile(filename string) error {
	template := `# Tasks

## General

- [ ] Task 1
- [ ] Task 2

## Ideas

- [ ] Idea 1

`
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("%s already exists", filename)
	}
	return os.WriteFile(filename, []byte(template), 0644)
}

func EnsureTasksFile() (*TaskList, error) {
	filename := detectTasksFile()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &TaskList{Filename: filename}, nil
	}
	return ParseTasksFile(filename)
}

func FindProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "TASKS.md")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "TODO.md")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}
