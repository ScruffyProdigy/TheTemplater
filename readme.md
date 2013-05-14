# TheTemplater

A simple bulk template loader

It is not uncommon for a web project to have a folder filled with templates to help simplify the "view" section of a MVC project.
Often times these templates are separated into folders based on which controller needs to call them.
And often times these templates need to access one another in order to DRY up the template and create partial views.

TheTemplater makes it as simple as possible to load all of these views in as few commands as possible.
Specifically, one command to load the root folder of the templates.
And one command each time you wish to render one of the templates.

## Usage

1) Create a folder for your templates
2) Put a bunch of templates in the folder
3) Feel free to organize your templates with subfolders
4) Import this Library
5) Call `t := templater.New()` and pass it the name of the folder with the templates
6) Later, render a template by calling t.Render()

## Example

**main.go:**

	package main

	import (
		"github.com/ScruffyProdigy/TheTemplater/templater"
		"os"
	)

	func main() {
		out := os.Stdout
		vars := map[string]string{"Name": "World"}

		layouts := templater.New("layouts")
		layouts.Render("hello",out,vars)
	}
  
  
**layouts/hello.tmpl:**
	
	Hello {{template "world" .}}
	
**layouts/world.tmpl

	{{.Name}}

When run, this program should output "Hello World"

## Documentation

http://godoc.org/github.com/ScruffyProdigy/TheTemplater/templater