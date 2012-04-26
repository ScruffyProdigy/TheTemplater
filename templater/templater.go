package templater

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var t *template.Template

/*
	This error gets returned if you request access to a template that doesn't exist
*/
type TemplateNotFoundError struct {
	tmpl string
}

func (this TemplateNotFoundError) Error() string {
	return "Unable to find template \"" + this.tmpl + "\""
}

/*
	Get() will get a specified template from the previously loaded templates
*/
func Get(name string) (result *template.Template, err error) {
	result = t.Lookup(name)
	if result == nil {
		err = TemplateNotFoundError{name}
	}
	return
}

/*
	Available() will let you know whether or not a specified template exists
*/
func Available(name string) bool {
	result := t.Lookup(name)
	if result == nil {
		return false
	}
	return true
}

/*
	Create will add a template with the specified text under the specified name for later use with Get() or Available()
*/
func Create(name, text string) {
	a := t.New(name)
	a.Parse(text)
}

/*
	LoadFromFiles() loads all of the templates within a specified file directory
	It dig into subdirectories recursively
*/
func LoadFromFiles(dir string, logger *log.Logger) {
	base := filepath.Clean(dir)

	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.TrimLeft(filepath.Ext(info.Name()), ".")
		if ext != "tmpl" {
			//eventually, I would like to allow different extension types to refer to different types of templates, but for now, we just skip them
			if logger != nil {
				logger.Print("Warning - Unknown Template Type: ." + ext + " - Skipping")
			}
			return nil
		}

		file, err := filepath.Rel(base, path)
		if err != nil {
			if logger != nil {
				logger.Print("Unable to find file at " + path + " - Skipping")
			}
			return nil
		}

		index := strings.Index(file, ".")
		if index == -1 {
			if logger != nil {
				logger.Print("Warning - No file type in " + file + " - Skipping")
			}
			return nil
		}
		name := file[0:index]

		b, err := ioutil.ReadFile(path)
		if err != nil {
			if logger != nil {
				logger.Print("Error - Could not load template file:" + path + " - Skipping")
			}
			return nil
		}
		s := string(b)

		//once I have other types of templates, I would need to look up a translator in a extension-translator map right here, and translate s into a "*/template" friendly format

		Create(name, s)

		return nil
	})
}

func init() {
	t = template.New("")
	t.Parse("")
}
