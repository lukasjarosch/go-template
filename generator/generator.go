package generator

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	template2 "text/template"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"

	"bitbucket.org/jdbergmann/godin/pkg/template"
	"bitbucket.org/jdbergmann/godin/pkg/writer"
)

// Target defines the all required metadata in order to render a template into
// a specific location.
type Target struct {
	PathTemplate      string // path template into which the target (file) is written
	GoImports         bool   // If true, the generator will call 'goimports' after writing the target
	OverwriteExisting bool   // If true, existing targets will be overwritten (aka. 'file is managed by godin')
	TemplateName      string // Name of the target's template
	TemplatePath      string // path to the template to render
}

// Path is a convenience method to quickly render Target.PathTemplate.
// The method throws a panic() if the rendering fails. This will most likely be the case if you pass the wrong
// template context. It's up to the caller managing this.
func (t *Target) Path(ctx interface{}) string {
	 str, err := ParseString(ctx, t.PathTemplate)
	 if err != nil {
	 	panic(fmt.Errorf("RenderedPath failed: %s", err))
	 }
	 return str
}

// Generator is a convenience struct to abstract away the template and file-writing logic
// and provide an easy interface to render files.
type Generator struct {
	target        *Target
	name          string // name is used as template name
	tplFilesystem http.FileSystem
	funcMap       template2.FuncMap // funcMap to inject into the template
}

// NewGenerator returns a newly initialized Generator. The passed FileSystem must provide the templates
// with the File defined in the Target. The 'Masterminds/sprig' FuncMap is also added no matter whether an additional one was passed or not.
func NewGenerator(name string, target *Target, templateFilesystem http.FileSystem, funcMap template2.FuncMap) *Generator {
	return &Generator{
		target:        target,
		name:          name,
		tplFilesystem: templateFilesystem,
		funcMap:       funcMap,
	}
}

// ParseString is a convenient helper function to quickly parse simple string templates
func ParseString(context interface{}, template string) (string, error) {
	tmpl := template2.New("parse-string").Funcs(sprig.TxtFuncMap())

	tmpl, err := tmpl.Parse(template)
	if err != nil {
		return "", errors.Wrap(err, "unable to parse Target.File as template")
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", errors.Wrap(err, "unable to execute Target.File template")
	}
	return buf.String(), nil
}

// Run the generator on the target, passing in the given templateContext to be used within the template.
// If the 'appendMode' is enabled, the generation Target will be appended to the file specified in 'File'.
// The 'forceOverwrite' flag allows to overwrite the Target's configuration during runtime. It's used
// in combination with the 'force' flag of the cli.
//
// If the Target's 'GoImports' flag is set, then the template renderer will 'gofmt' it first. After that 'goimports'
// is also executed to make sure the imports are fine. Although this will eventually re-format as well.
func (g *Generator) Run(templateContext interface{}, appendMode bool, forceOverwrite bool) error {
	if forceOverwrite {
		g.target.OverwriteExisting = true
	}

	path, err := ParseString(templateContext, g.target.PathTemplate)
	if err != nil {
		return err
	}
	g.target.PathTemplate = path

	tmpl := template2.New("target-path").Funcs(sprig.TxtFuncMap())

	if g.funcMap != nil {
		tmpl = tmpl.Funcs(g.funcMap)
	}

	tmpl, err = tmpl.Parse(g.target.PathTemplate)
	if err != nil {
		return errors.Wrap(err, "unable to parse Target.File as template")
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateContext); err != nil {
		return errors.Wrap(err, "unable to execute Target.File template")
	}
	g.target.PathTemplate = buf.String()

	tpl := template.NewFullTemplate(
		template.Name(g.target.TemplateName),
		template.Path(g.target.TemplatePath),
		template.GoSource(g.target.GoImports),
		template.UseFilesystem(g.tplFilesystem),
		template.FuncMap(g.funcMap),
	)
	rendered, err := tpl.Render(templateContext)
	if err != nil {
		return errors.Wrap(err, "unable to GenerateFull")
	}

	fw := writer.NewFileWriter(
		g.target.PathTemplate,
		writer.Overwrite(g.target.OverwriteExisting),
		writer.Append(appendMode),
	)

	// write rendered target to file, ignoring ErrNoOp which only indicates that the writer didn't needed to be run.
	// This happens for example if the file exists and the overwrite and append flags are 'false'.
	if err := fw.WriteFile(rendered); err != nil {
		switch err {
		case writer.ErrNoOp:
			return nil
		default:
			return errors.Wrap(err, "unable to WriteFile")
		}
	}

	if g.target.GoImports {
		if err := g.goImports(g.target.PathTemplate); err != nil {
			return errors.Wrap(err, "goimports failed")
		}
	}

	return nil
}

// WithFuncMap returns a generator with the passed FuncMap attached.
func (g *Generator) WithFuncMap(funcMap template2.FuncMap) *Generator {
	g.funcMap = funcMap
	return g
}

// goImports executes 'goimports -w <path>' on the host system
func (g *Generator) goImports(path string) error {
	modCmd := exec.Command("goimports", "-w", path)
	err := modCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
