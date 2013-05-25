package templater

import (
	"fmt"
	"os"
)

func Example_HelloWorld() {
	out := os.Stdout
	vars := map[string]string{"Name": "World"}

	layouts, errs := New("layouts_test")
	for _, err := range errs {
		fmt.Fprintln(out, "Layout - "+err.Error())
	}

	err := layouts.Render("hello", out, vars)
	if err != nil {
		fmt.Fprintln(out, err.Error())
	}
	//output: Hello World
}

func Example_Errors() {
	out := os.Stdout
	layouts, errs := New("layouts_with_errors_test")
	err := layouts.Render("errors", out, errs)
	if err != nil {
		fmt.Fprintln(out, err.Error())
	}
	fmt.Fprintln(out, len(errs), "errors")

	/*output:
	Error - No translator found for err file
	Error - template: b:1: unexpected unclosed action in command
	2 errors
	*/
}
