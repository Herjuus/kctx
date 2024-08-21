package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

func getKubeContexts() ([]string, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	contexts := strings.Split(strings.TrimSpace(string(out)), "\n")
	return contexts, nil
}

func switchContext(context string) error {
	cmd := exec.Command("kubectl", "config", "use-context", context)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	contexts, err := getKubeContexts()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching contexts: %v\n", err)
		os.Exit(1)
	}

	prompt := promptui.Select{
		Label: "Select context",
		Items: contexts,
	}

	_, selectedContext, err := prompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting context: %v\n", err)
		os.Exit(1)
	}

	err = switchContext(selectedContext)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error switching context: %v\n", err)
	}

	fmt.Printf("Set current context to: %s\n", selectedContext)
}
