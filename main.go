package main

import (
	"fmt"
	"os"
	"strings"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintTasks(tl)
		return
	}

	if os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h" {
		PrintHelp()
		return
	}

	if os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v" {
		PrintVersion()
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "init":
		filename := DefaultFilename
		if len(os.Args) > 2 {
			filename = os.Args[2]
		}
		if err := InitTasksFile(filename); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created %s\n", filename)
		fmt.Println("Run 'mdt' to view tasks.")

	case "list", "ls":
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintTasks(tl)

	case "add":
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		category := ""
		descParts := os.Args[2:]

		if len(descParts) >= 2 && descParts[0] == "-c" {
			category = descParts[1]
			descParts = descParts[2:]
		}

		desc := strings.Join(descParts, " ")
		if desc == "" {
			fmt.Fprintln(os.Stderr, "Error: task description required")
			fmt.Fprintln(os.Stderr, "Usage: mdt add [-c category] <description>")
			os.Exit(1)
		}

		if err := tl.Add(desc, category); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := WriteTasksFile(tl); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added task #%d\n", tl.nextNumber()-1)

	case "done":
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		num, err := FindTaskNum(os.Args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := tl.Toggle(num, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := WriteTasksFile(tl); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d marked as done %s\n", num, color(green, checkMark))

	case "undo":
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		num, err := FindTaskNum(os.Args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := tl.Toggle(num, false); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := WriteTasksFile(tl); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d marked as pending\n", num)

	case "rm", "remove":
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		num, err := FindTaskNum(os.Args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := tl.Remove(num); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := WriteTasksFile(tl); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing tasks: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task #%d removed\n", num)

	case "stats", "stat":
		tl, err := EnsureTasksFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		PrintStats(tl)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		fmt.Fprintf(os.Stderr, "Run 'mdt help' for usage.\n")
		os.Exit(1)
	}
}
