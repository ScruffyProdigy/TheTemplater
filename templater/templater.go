/*

A simple bulk template loader

It is not uncommon for a web project to have a folder filled with templates to help simplify the "view" section of a MVC project.
Often times these templates are separated into folders based on which controller needs to call them.
And often times these templates need to access one another in order to DRY up the template and create partial views.

TheTemplater makes it as simple as possible to load all of these views in as few commands as possible.
Specifically, one command to load the root folder of the templates.
And one command each time you wish to render one of the templates.

*/
package templater

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//eventually, I would like to allow different extension types to refer to different types of templates (ie, something like HAML)
//To my knowledge, something like that doesn't yet exist, so for now, only basic templates
var translators map[string]func(string)(string,error) = map[string]func(string)(string,error){
	"tmpl":func(s string)(string,error) {
		return s,nil
	},
}

type Group struct {
	t *template.Template
	l *log.Logger
}

/*
	New() will search a directory, and recursively search all subdirectories for template files.
	It then stores all template files into one location to make for easy lookup.
	It also associates all of the templates found, so any template loaded can reference any other loaded templated
	Currently, it looks for all files with a .tmpl extension, and then strips the extension
*/
func New(directory string) *Group {
	return NewAndLogErrors(directory,nil)
}

/*
	NewAndLogErrors() does the same thing as New, except it takes a logger, and will log any errors it encounters for debugging later
*/
func NewAndLogErrors(directory string,logger *log.Logger) *Group {
	this := new(Group)
	
	this.l = logger
	
	this.t = template.New("")
	this.t.Parse("")
	
	this.loadFolder(directory)
	
	return this
}

func (this *Group) create(name, text string) {
	a := this.t.New(name)
	_,err := a.Parse(text)
	if err != nil && this.l != nil {
		this.l.Println("**Warning** "+err.Error())
	}
}

func getTranslator(info os.FileInfo) func(string)(string,error) {
	extension := strings.TrimLeft(filepath.Ext(info.Name()), ".")
	return translators[extension]
}

func getFileName(base,path string) string {
	filename,err := filepath.Rel(base, path)
	if err != nil {
		return ""
	}
	return filename
}

func stripExtension(file string) string {
	//we used the extension to determine the template type, now we strip the extension from the file name
	index := strings.LastIndex(file, ".")
	if index == -1 {
		//no "." in the filename (therefore, not a .tmpl file), skip
		return ""
	}
	return file[0:index]
}

func (this *Group) loadFolder(dir string) {
	base := filepath.Clean(dir)

	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			return nil
		}
		
		translator := getTranslator(info)
		if translator == nil {
			return nil
		}

		fileName := getFileName(base,path)
		if fileName == "" {
			return nil
		}

//we used the extension to determine which template type we were using, 
//we strip it for the template name
		templateName := stripExtension(fileName)
		if templateName == "" {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			//couldn't load the file, skip
			return nil
		}
		
		s,err := translator(string(b))
		if err != nil {
			return nil
		}

		this.create(templateName, s)

		return nil
	})
}


/*
	Get() will get a specified template from the previously loaded templates
*/
func (this *Group) Get(name string) *template.Template {
	return this.t.Lookup(name)
}

/*
	Render() will get the named template from the previously loaded templates
	it will then execute it using data for the pipeline, and outputting to a writer
*/
func (this *Group) Render(name string, out io.Writer, data interface{}) {
	t := this.Get(name)
	if t == nil {
		this.l.Println("**Warning** could not find template:"+name)
		return
	}
	err := t.Execute(out,data)
	if err != nil {
		this.l.Println("**Warning** "+err.Error())
	}
}
