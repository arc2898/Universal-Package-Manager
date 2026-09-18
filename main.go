package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arc2898/Universal-Package-Manager/internal/core"
	"github.com/arc2898/Universal-Package-Manager/internal/logger"
)

const version = "0.2.0"

func main() {
	u := core.NewUPM()
	u.Init()

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	args := os.Args[1:]
	var err error
	switch args[0] {
	case "install":
		err = requirePackage(args, "install", u.Install)
	case "remove":
		err = requirePackage(args, "remove", u.Remove)
	case "search":
		err = requirePackage(args, "search", u.Search)
	case "-b", "--bulk":
		err = handleBulk(u, args[1:])
	case "-u", "--update":
		err = u.UpdateAll()
	case "refresh":
		u.Init()
		u.ListManagers()
	case "-l", "--log":
		err = logger.ReadLogs()
	case "-r", "--remove":
		err = u.UninstallSelf()
	case "-v", "--version":
		fmt.Println("UPM version " + version)
	case "-h", "--help":
		printHelp()
	default:
		err = fmt.Errorf("unknown command: %s", args[0])
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		logger.Log("command failed: " + err.Error())
		os.Exit(1)
	}
}

func requirePackage(args []string, command string, operation func(string) error) error {
	if len(args) != 2 || strings.TrimSpace(args[1]) == "" {
		return fmt.Errorf("usage: upm %s <pkg_name>", command)
	}
	return operation(args[1])
}

func handleBulk(u *core.UPM, args []string) error {
	if len(args) < 3 {
		return errors.New("usage: upm -b [install|remove] [manager] [pkg1,pkg2] ...")
	}
	items, err := parseBulkArgs(args)
	if err != nil {
		return err
	}
	if args[0] == "install" {
		return u.BulkInstall(items)
	}
	return u.BulkRemove(items)
}

func parseBulkArgs(args []string) (map[string][]string, error) {
	if len(args) < 3 {
		return nil, errors.New("bulk action requires an action and at least one manager/package pair")
	}
	if args[0] != "install" && args[0] != "remove" {
		return nil, fmt.Errorf("unknown bulk action %q", args[0])
	}
	pairs := args[1:]
	if len(pairs)%2 != 0 {
		return nil, errors.New("bulk arguments must contain complete manager/package pairs")
	}

	items := make(map[string][]string)
	for i := 0; i < len(pairs); i += 2 {
		managerName := strings.TrimSpace(pairs[i])
		if managerName == "" {
			return nil, errors.New("bulk manager name cannot be empty")
		}
		var packages []string
		for _, pkgName := range strings.Split(pairs[i+1], ",") {
			pkgName = strings.TrimSpace(pkgName)
			if pkgName == "" {
				return nil, fmt.Errorf("empty package in manager %q", managerName)
			}
			packages = append(packages, pkgName)
		}
		items[managerName] = append(items[managerName], packages...)
	}
	return items, nil
}

func printHelp() {
	fmt.Println("UPM - Universal Package Manager")
	fmt.Println("Usage: upm <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  install <pkg>       Install a package using the native manager")
	fmt.Println("  remove <pkg>        Remove a package using the native manager")
	fmt.Println("  search <pkg>        Search for a package across detected managers")
	fmt.Println("  -b, --bulk <action> Bulk install/remove")
	fmt.Println("                      Example: upm -b install apt git,vim snap spotify")
	fmt.Println("  -u, --update        Update all detected package managers")
	fmt.Println("  refresh             Re-detect known package managers")
	fmt.Println("  -l, --log           View operation logs")
	fmt.Println("  -r, --remove        Remove the installed UPM binary")
	fmt.Println("  -v, --version       Show version")
	fmt.Println("  -h, --help          Show help")
}
