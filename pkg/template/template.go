package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/rasteiro11/PogCore/pkg/logger"
)

var global Template
var once sync.Once

var validExtensions = []string{".tmpl", ".tpl", ".html"}

type (
	Template interface {
		Parse(name string, payload any) (string, error)
	}

	templates struct {
		data map[string]*template.Template
	}
)

type options struct {
	templates []string
}

type Option func(*options)

func WithRequiredTemplate(template string) Option {
	return func(o *options) {
		o.templates = append(o.templates, template)
	}
}

func New(opts ...Option) error {
	opt := options{
		templates: []string{},
	}

	templates := lazyLoadingTemplate()

	for _, o := range opts {
		o(&opt)
	}

	for _, template := range opt.templates {
		_, ok := templates.data[template]
		if !ok {
			return fmt.Errorf("template is required: %s", template)
		}
	}

	return nil
}

func Instance() Template {
	once.Do(func() {
		global = lazyLoadingTemplate()
	})

	return global
}

func (t *templates) Parse(filename string, data any) (string, error) {
	tmpl, ok := t.data[filename]
	if !ok {
		return "", fmt.Errorf("template not found: %+v", filename)
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func lazyLoadingTemplate() (t *templates) {
	t = &templates{
		data: make(map[string]*template.Template),
	}

	exe, err := os.Executable()
	if err != nil {
		logger.Global().Errorf("[template.newTemplate] os.Executable() returned error: %+v\n", err)
		return
	}

	callback := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			logger.Global().Errorf("[template.lazyLoadingTemplate] returned error: %+v\n", err)
			return err
		}

		if ok := keepWalking(path, info); ok {
			return nil
		}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			logger.Global().Errorf("[template.lazyLoadingTemplate] template.ParseFiles() returned error: %+v\n", err)
			return err
		}

		logger.Global().Debugf("Loaded template: %+s\n", info.Name())

		t.data[info.Name()] = tmpl

		return nil
	}

	if err := filepath.WalkDir(filepath.Dir(exe), callback); err != nil {
		logger.Global().Errorf("[template.newTemplate] filepath.WalkDir() returned error: %+v\n", err)
	}

	return
}

func keepWalking(path string, fi fs.DirEntry) bool {
	if fi.IsDir() {
		return true
	}

	for _, extension := range validExtensions {
		if ok := strings.HasSuffix(fi.Name(), extension); ok {
			return !ok
		}
	}

	return true
}
