package main

import (
	"flag"
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

func getCurrentContextCursorPos(contexts []string) int {
	cmd := exec.Command("kubectl", "config", "current-context")
	out, _ := cmd.Output()

	for i, context := range contexts {
		if strings.TrimSpace(string(out)) == context {
			return i
		}
	}

	return 0
}

func switchContext(context string) error {
	cmd := exec.Command("kubectl", "config", "use-context", context)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func exitIfError(err error, msg string, a ...any) {
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s: %s\n", msg, err.Error()), a...)
		os.Exit(1)
	}
}

func main() {
	var nContexts = flag.Int("n", 10, "Number of contexts to display at once")
	flag.Parse()

	contexts, err := getKubeContexts()
	exitIfError(err, "Error fetching contexts")

	selectedContext := flag.Arg(0)
	if selectedContext == "" {
		prompt := promptui.Select{
			Label: "Select context",
			Items: contexts,
			Size:  *nContexts,
			Searcher: func(input string, index int) bool {
				return strings.Contains(strings.ToLower(contexts[index]), strings.ToLower(input))
			},
		}

		_, selectedContext, err = prompt.RunCursorAt(getCurrentContextCursorPos(contexts), 0)
		exitIfError(err, "Error selecting context '%s'", selectedContext)
	}

	err = switchContext(selectedContext)
	exitIfError(err, "Error switching context to '%s'", selectedContext)

	fmt.Printf("Set current context to '%s'\n", selectedContext)
}
