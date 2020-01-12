package template

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
)

// Renderer defines the crude interface for a template implementation
type Renderer interface {
	Render(templateContext interface{}) ([]byte, error)
}

type FileTemplate struct {
	baseTemplate
}

func NewFileTemplate(options ...Option) *FileTemplate {
	opts := newOptions(options...)
	return &FileTemplate{baseTemplate{opts:opts}}
}

// Render implements the Renderer interface for the FileTemplate struct.
// The FileTemplate always handles full-file templates and does not support append-writing to
// existing files.
func (tpl *FileTemplate) Render(templateContext interface{}) ([]byte, error) {
	templateData, err := tpl.readTemplate()
	if err != nil {
		return nil, NewFileReadError(err)
	}

	tmpl, err := tpl.createTemplate(string(templateData))
	if err != nil {
		return nil, NewCreateTemplateError(err)
	}

	// execute template
	out := bytes.Buffer{}
	if err := tmpl.Execute(&out, templateContext); err != nil {
		return nil, NewExecutionError(err)
	}

	// format as Go source if required
	if tpl.opts.GoSource {
		formatted, err := tpl.formatGoSource(out.Bytes())
		if err != nil {
			return nil, NewGoFormatError(err)
		}
		out = *formatted
	}

	return out.Bytes(), nil
}

type PartialFileTemplate struct {
	baseTemplate
}

func NewPartialFileTemplate(options ...Option) *PartialFileTemplate {
	opts := newOptions(options...)
	return &PartialFileTemplate{baseTemplate{opts:opts}}
}


func (tpl *PartialFileTemplate) Render(templateContext interface{}) ([]byte, error) {
	templateData, err := tpl.readTemplate()
	if err != nil {
		return nil, NewFileReadError(err)
	}

	tmpl, err := tpl.createTemplate(string(templateData))
	if err != nil {
		return nil, NewCreateTemplateError(err)
	}

	// in order to use the partial template definitions, we need to include them somewhere
	// thus we just add the template to use itself
	templateData = append(templateData, []byte(fmt.Sprintf("{{ template \"%s\" . }}", tpl.opts.Name))...)

	out := bytes.Buffer{}
	if err := tmpl.Execute(&out, templateContext); err != nil {
		return nil, NewExecutionError(err)
	}

	// format as Go source if required
	if tpl.opts.GoSource {
		formatted, err := tpl.formatGoSource(out.Bytes())
		if err != nil {
			return nil, NewGoFormatError(err)
		}
		out = *formatted
	}

	return out.Bytes(), nil
}

// baseTemplate defines default behaviour for all templates
type baseTemplate struct {
	opts Options
}

// templateFileExists returns 'true' if the source template file exists
// The given path must be pointing to a file. If the path points to a directory, it's invalid.
func (tpl *baseTemplate) templateFileExists() bool {
	var info os.FileInfo
	var err error

	if tpl.opts.Path == "" {
		return false
	}

	if info, err = os.Stat(tpl.opts.Path); os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// openTemplate decides from which filesystem the file should be read and provides an io.Reader
func (tpl *baseTemplate) openTemplate() (reader io.Reader, err error) {
	if tpl.opts.Filesystem != nil {
		f, err := tpl.opts.Filesystem.Open(tpl.opts.Path)
		if err != nil {
			return nil, NewFileOpenError(err)
		}
		reader = f

	} else {
		if !tpl.templateFileExists() {
			return nil, NewFileOpenError(err)
		}

		f, err := os.Open(tpl.opts.Path)
		if err != nil {
			return nil, NewFileOpenError(err)
		}
		reader = f
	}
	return reader, nil
}

// readTemplate is a convenience function which calls openTemplate. If the call is successful, it sluprs the Reader empty
// and returns the read bytes-slice.
func (tpl *baseTemplate) readTemplate() ([]byte, error) {
	reader, err := tpl.openTemplate()
	if err != nil {
		return nil, NewFileOpenError(err)
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, NewFileReadError(err)
	}

	return buf, nil
}

// formatGoSource formats go source files
func (tpl *baseTemplate) formatGoSource(data []byte) (out *bytes.Buffer, err error) {
		formatted, err := format.Source(data)
		if err != nil {
			return nil, NewGoFormatError(err)
		}
		out = bytes.NewBuffer(formatted)
	return out, nil
}

// createTemplate creates and parses the template with attached FuncMaps
func (tpl *baseTemplate) createTemplate(templateData string) (*template.Template, error){
	tmpl := template.New(tpl.opts.Name).Funcs(sprig.TxtFuncMap()).Funcs(tpl.opts.FuncMap)
	tmpl, err := tmpl.Parse(templateData)
	if err != nil {
		return nil, NewParseError(err)
	}
	return tmpl, nil
}

