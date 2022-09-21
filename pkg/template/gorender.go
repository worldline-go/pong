package template

import (
	"bytes"
	"fmt"
	"io"

	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

var txtTemplate = textTemplate.New("txt")

type TemplateExecuter interface {
	Execute(wr io.Writer, data interface{}) error
}

func Ext(v map[string]interface{}, content string) ([]byte, error) {
	funct := txtTemplate

	funct2, err := funct.Clone()
	if err != nil {
		return nil, fmt.Errorf("clone text template error: %w", err)
	}

	funct2 = funct2.Funcs(sprig.TxtFuncMap())

	tmp, err := funct2.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parse text template error: %w", err)
	}

	return Execute(v, tmp)
}

// Execute executes the template with the given data.
func Execute(v map[string]interface{}, tmp TemplateExecuter) ([]byte, error) {
	var b bytes.Buffer

	// Execute the template and write the output to the buffer
	if err := tmp.Execute(&b, v); err != nil {
		return nil, fmt.Errorf("execute template error: %w", err)
	}

	return b.Bytes(), nil
}
