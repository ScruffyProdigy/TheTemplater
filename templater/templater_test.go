package templater

import (
	"os"
	"log"
)

func Example_HelloWorld() {
	out := os.Stdout
	vars := map[string]string{"Name": "World"}

	layouts := NewAndLogErrors("layouts_test",log.New(out,"Templater - ",0))
	layouts.Render("hello",out,vars)
	//output: Hello World
}