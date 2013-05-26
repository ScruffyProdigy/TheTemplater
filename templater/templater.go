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
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//eventually, I would like to allow different extension types to refer to different types of templates (ie, something like HAML)
//To my knowledge, something like that doesn't yet exist, so for now, only basic templates
var translators map[string]func(string) (string, error) = map[string]func(string) (string, error){
	"tmpl": func(s string) (string, error) {
		return s, nil
	},
}

type Group struct {
	t *template.Template
}

/*
	New() will search a directory, and recursively search all subdirectories for template files.
	It then stores all template files into one location to make for easy lookup.
	It also associates all of the templates found, so any template loaded can reference any other loaded templated
	Currently, it looks for all files with a .tmpl extension, and then strips the extension
*/
func New(directory string) (*Group, []error) {
	this := new(Group)

	this.t = template.New("")
	this.t.Parse("")

	errs := this.loadFolder(directory)

	return this, errs
}

func (this *Group) create(name, text string) error {
	a := this.t.New(name)
	_, err := a.Parse(text)
	return err
}

func getTranslator(info os.FileInfo) (func(string) (string, error), error) {
	extension := strings.TrimLeft(filepath.Ext(info.Name()), ".")
	translator := translators[extension]
	if translator == nil {
		return nil, errors.New("No translator found for " + extension + " file")
	}
	return translator, nil
}

func getFileName(base, path string) (string, error) {
	filename, err := filepath.Rel(base, path)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func stripExtension(file string) (string, error) {
	//we used the extension to determine the template type, now we strip the extension from the file name
	index := strings.LastIndex(file, ".")
	if index == -1 {
		//no "." in the filename (therefore, not a .tmpl file), skip
		return "", errors.New("No extension found")
	}
	return file[0:index], nil
}

func (this *Group) loadFolder(dir string) []error {
	base := filepath.Clean(dir)
	errs := []error{}

	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errs = append(errs, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		translator, err := getTranslator(info)
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		fileName, err := getFileName(base, path)
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		//we used the extension to determine which template type we were using, 
		//we strip it for the template name
		templateName, err := stripExtension(fileName)
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			//couldn't load the file, skip
			errs = append(errs, err)
			return nil
		}

		s, err := translator(string(b))
		if err != nil {
			errs = append(errs, err)
			return nil
		}

		err = this.create(templateName, s)
		if err != nil {
			errs = append(errs, err)
		}

		return nil
	})
	return errs
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
func (this *Group) Render(name string, out io.Writer, data interface{}) error {
	t := this.Get(name)
	if t == nil {
		return errors.New("could not find template:" + name)
	}
	err := t.Execute(out, data)
	if err != nil {
		return err
	}
	return nil
}
