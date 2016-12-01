package main

import (
	"github.com/atotto/clipboard"
	"fmt"
	"os"
)




// fatalPrint is a convenience function that will quit the program with the specified Exit Code, print some friendly
// context and a colon, and the error message from golang. Pass a nil error in to avoid printing the error string.
func fatalPrint(exitCode int, friendlyContext string, err error) {
	if err == nil {
		stdErr(fmt.Sprintf("Error: %s", friendlyContext))
	} else {
		stdErr(fmt.Sprintf("Error: %s: %s", friendlyContext, err))
	}
	os.Exit(exitCode)
}

// stdErr is a convenience function to send output to stderr. Automatically adds newline.
func stdErr(m string) {
	fmt.Fprintf(os.Stderr, "%s\n", m)
}


func getClipboardString() (out string, err error) {
	if clipboard.Unsupported {
		return out, fmt.Errorf("Clipboard functionality unsupported.")
	}
	return clipboard.ReadAll()
}

func deDuplicateStringSlice(sliceToDeDupe *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *sliceToDeDupe {
		if !found[x] {
			found[x] = true
			(*sliceToDeDupe)[j] = (*sliceToDeDupe)[i]
			j++
		}
	}
	*sliceToDeDupe = (*sliceToDeDupe)[:j]
}
