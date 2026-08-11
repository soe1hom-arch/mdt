package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	reset     = "\033[0m"
	red       = "\033[31m"
	green     = "\033[32m"
	yellow    = "\033[33m"
	blue      = "\033[34m"
	magenta   = "\033[35m"
	cyan      = "\033[36m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	checkMark = "✓"
	crossMark = "✗"
)

func isTerminal() bool {
	stat, _ := os.Stdout.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func color(c, s string) string {
	if !isTerminal() {
		return s
	}
	return c + s + reset
}

func PrintTasks(tl *TaskList) {
	if len(tl.Tasks) == 0 {
		fmt.Println(color(yellow, "✨ No tasks found. Use 'mdt add <task>' to add one."))
		return
	}

	grouped := tl.TasksByCategory()
	sortedCats := sortedKeys(grouped)

	for i, cat := range sortedCats {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf(" %s\n", color(bold+blue, cat))
		fmt.Println(strings.Repeat("─", 40))

		tasks := grouped[cat]
		for _, t := range tasks {
			prefix := color(red, crossMark)
			desc := t.Description
			num := color(cyan, fmt.Sprintf("%d.", t.Number))

			if t.Done {
				prefix = color(green, checkMark)
				desc = color(dim, desc)
				num = color(dim, fmt.Sprintf("%d.", t.Number))
			}

			fmt.Printf("   %s %s %s\n", num, prefix, desc)
		}
	}

	total, done, pending := tl.GetStats()
	fmt.Println()
	fmt.Printf(" %s  %s%d  %s%d  %s%d\n",
		color(bold, "Summary:"),
		color(green, "Done: "), done,
		color(yellow, "Pending: "), pending,
		color(cyan, "Total: "), total,
	)
}

func PrintStats(tl *TaskList) {
	total, done, pending := tl.GetStats()
	fmt.Printf("\n %s\n", color(bold+blue, "📊 Task Statistics"))
	fmt.Println(strings.Repeat("─", 30))
	fmt.Printf("   %-12s %d\n", color(bold, "Total:"), total)
	fmt.Printf("   %-12s %d\n", color(green, "Done:"), done)
	fmt.Printf("   %-12s %d\n", color(yellow, "Pending:"), pending)

	if total > 0 {
		pct := float64(done) / float64(total) * 100
		barLen := 20
		filled := int(float64(barLen) * float64(done) / float64(total))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barLen-filled)

		fmt.Printf("   %-12s %s %3.0f%%\n", color(bold, "Progress:"), bar, pct)
	}

	grouped := tl.TasksByCategory()
	if len(grouped) > 0 {
		fmt.Println()
		fmt.Printf("   %s\n", color(bold, "By Category:"))
		for _, cat := range sortedKeys(grouped) {
			tasks := grouped[cat]
			catDone := 0
			for _, t := range tasks {
				if t.Done {
					catDone++
				}
			}
			fmt.Printf("   %s %s %d/%d\n", color(cyan, "▶"), cat, catDone, len(tasks))
		}
	}
	fmt.Println()
}

func PrintHelp() {
	help := `
 %s

  %s      List all tasks (default)
  %s      Add a new task
  %s      Mark task as done
  %s      Mark task as not done
  %s      Remove a task
  %s      Show statistics
  %s      Create a TASKS.md template
  %s      Show version
  %s      Show this help

 %s
  %s
  %s
  %s
  %s

`
	fmt.Printf(help,
		color(bold+blue, "mdt - MarkDown Task Manager"),
		color(cyan, "mdt"),
		color(cyan, "mdt add <description>"),
		color(cyan, "mdt done <number>"),
		color(cyan, "mdt undo <number>"),
		color(cyan, "mdt rm <number>"),
		color(cyan, "mdt stats"),
		color(cyan, "mdt init"),
		color(cyan, "mdt version"),
		color(cyan, "mdt help"),
		color(bold, "Examples:"),
		color(dim, "  mdt add \"Buy groceries\""),
		color(dim, "  mdt add -c Python \"Write unit tests\""),
		color(dim, "  mdt done 1"),
		color(dim, "  mdt stats"),
	)
}

func PrintVersion() {
	fmt.Printf("%s %s\n", color(bold+blue, "mdt"), version)
}

func sortedKeys(m map[string][]Task) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	if len(keys) > 1 {
		// put "General" first if it exists
		for i, k := range keys {
			if k == "General" {
				keys = append([]string{k}, append(keys[:i], keys[i+1:]...)...)
				break
			}
		}
	}
	return keys
}

func sortStrings(s []string) {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}
