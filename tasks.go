package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Task struct {
	Index      int
	Number     int
	Description string
	Done       bool
	Category   string
}

type TaskList struct {
	Tasks     []Task
	Categories []string
	Filename  string
}

func (tl *TaskList) Add(desc, category string) error {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	if category == "" {
		category = "General"
	}
	num := tl.nextNumber()
	tl.Tasks = append(tl.Tasks, Task{
		Number:      num,
		Description: desc,
		Done:        false,
		Category:    category,
	})
	tl.addCategory(category)
	return nil
}

func (tl *TaskList) Toggle(num int, done bool) error {
	for i := range tl.Tasks {
		if tl.Tasks[i].Number == num {
			tl.Tasks[i].Done = done
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", num)
}

func (tl *TaskList) Remove(num int) error {
	for i := range tl.Tasks {
		if tl.Tasks[i].Number == num {
			tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)
			tl.rebuildCategories()
			return nil
		}
	}
	return fmt.Errorf("task #%d not found", num)
}

func (tl *TaskList) GetStats() (total, done, pending int) {
	total = len(tl.Tasks)
	for _, t := range tl.Tasks {
		if t.Done {
			done++
		}
	}
	pending = total - done
	return
}

func (tl *TaskList) nextNumber() int {
	max := 0
	for _, t := range tl.Tasks {
		if t.Number > max {
			max = t.Number
		}
	}
	return max + 1
}

func (tl *TaskList) addCategory(cat string) {
	for _, c := range tl.Categories {
		if c == cat {
			return
		}
	}
	tl.Categories = append(tl.Categories, cat)
	sort.Strings(tl.Categories)
}

func (tl *TaskList) rebuildCategories() {
	seen := make(map[string]bool)
	for _, t := range tl.Tasks {
		seen[t.Category] = true
	}
	tl.Categories = nil
	for cat := range seen {
		tl.Categories = append(tl.Categories, cat)
	}
	sort.Strings(tl.Categories)
}

func (tl *TaskList) TasksByCategory() map[string][]Task {
	grouped := make(map[string][]Task)
	for _, t := range tl.Tasks {
		cat := t.Category
		if cat == "" {
			cat = "General"
		}
		grouped[cat] = append(grouped[cat], t)
	}
	return grouped
}

func FindTaskNum(args []string) (int, error) {
	if len(args) < 3 {
		return 0, fmt.Errorf("task number required")
	}
	num, err := strconv.Atoi(args[2])
	if err != nil {
		return 0, fmt.Errorf("invalid task number: %s", args[2])
	}
	return num, nil
}
