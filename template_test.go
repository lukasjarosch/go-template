package template

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestBaseTemplate_OpenTemplate(t *testing.T) {
	g := NewGomegaWithT(t)
	fixtures := []struct{
		input *baseTemplate
		errExpected error
	}{
		{
			input:       &baseTemplate{Options{
				Name:       "test-template-not-existing",
				Path:       "/some/where/over/the/rainbow",
				GoSource:   false,
				FuncMap:    nil,
				Filesystem: nil,
			}},
			errExpected: NewFileOpenError(nil),
		},
	}

	// TODO: more tests

	for _, tt := range fixtures {
		_, err := tt.input.openTemplate()

		if tt.errExpected != nil {
			g.Expect(err).To(HaveOccurred())
			g.Expect(err).To(MatchError(tt.errExpected))
		}
	}
}

func TestBaseTemplate_TemplateFileExists(t *testing.T) {
	g := NewGomegaWithT(t)
	fixtures := []struct{
		input *baseTemplate
		expected bool
	}{
		{
			input:    &baseTemplate{Options{
				Name:       "test-template-not-existing",
				Path:       "/some/where/over/the/rainbow",
				GoSource:   false,
				FuncMap:    nil,
				Filesystem: nil,
			}},
			expected: false,
		},
		{
			input:    &baseTemplate{Options{
				Name:       "test-template-existing",
				Path:       "/tmp/test.go",
				GoSource:   false,
				FuncMap:    nil,
				Filesystem: nil,
			}},
			expected: true,
		},
		{
			input:    &baseTemplate{Options{
				Name:       "test-template-dir",
				Path:       "/tmp",
				GoSource:   false,
				FuncMap:    nil,
				Filesystem: nil,
			}},
			expected: false,
		},
	}

	for _, tt := range fixtures {
		if tt.expected {
			 _, _ = os.Create(tt.input.opts.Path)
		}

		out := tt.input.templateFileExists()
		g.Expect(out).To(Equal(tt.expected))

		if tt.expected {
			_, _ = os.Create(tt.input.opts.Path)
		}
	}
}

