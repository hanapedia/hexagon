package generate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func RenderAndSave(dir string, name string, tmpl string, data interface{}) error {
	manifest, err := renderTemplate(name, tmpl, data)
	if err != nil {
		return err
	}
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, fmt.Sprintf("%s.yaml", name))
	err = ioutil.WriteFile(path, []byte(manifest), 0644)
	if err != nil {
		return err
	}
	return nil
}

func renderTemplate(name string, tmpl string, data interface{}) (string, error) {
	parser, err := template.New(name).Funcs(template.FuncMap{"indent": indent}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = parser.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func indent(input string, spaces int) string {
	lines := strings.Split(input, "\n")
	for i, line := range lines {
		if len(line) > 0 {
			lines[i] = strings.Repeat(" ", spaces) + line
		}
	}
	return strings.Join(lines, "\n")
}
