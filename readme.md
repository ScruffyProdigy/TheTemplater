# Template

A simple bulk template loader

run `go get github.com/HairyMezican/TheTemplater/templater` to install it

## Example

**main.go:**

	package main

	import (
		"github.com/HairyMezican/TheTemplater/templater"
		"os"
	)

	func main() {
		out := os.Stdout
		vars := map[string]string{"Name": "World"}

		templater.LoadFromFiles("templates", nil)
		t, _ := templater.Get("hello")
		t.Execute(out, vars)
	}
  
  
**templates/hello.tmpl:**
	
	Hello {{.Name}}
	

When run, this program should output "Hello World"