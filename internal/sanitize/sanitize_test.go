package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text without tags",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "text with paragraph tags",
			input:    "<p>Hello</p><p>World</p>",
			expected: "Hello\nWorld\n",
		},
		{
			name:     "text with br tags",
			input:    "Line 1<br>Line 2<br/>Line 3<br />Line 4",
			expected: "Line 1\nLine 2\nLine 3\nLine 4",
		},
		{
			name:     "text with various tags",
			input:    "<div><strong>Bold</strong> <em>Italic</em></div>",
			expected: "Bold Italic",
		},
		{
			name:     "text with script tags",
			input:    "<script>alert('xss')</script>Hello",
			expected: "alert('xss')Hello", // script tag is removed but content remains
		},
		{
			name:     "text with common entities",
			input:    "&#8216;quote&#8217; &#8220;double&#8221; &nbsp; &quot;test&quot; &apos;apos&apos;",
			expected: "'quote' \"double\"   \"test\" 'apos'",
		},
		{
			name:     "text with nested tags",
			input:    "<div><p>Nested <span>content</span> here</p></div>",
			expected: "Nested content here\n",
		},
		{
			name:     "text with anchor tags",
			input:    "<a href='http://example.com'>Link</a>",
			expected: "Link",
		},
		{
			name:     "text with multiple newlines",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3", // newlines are preserved in plain text
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only tags no content",
			input:    "<div><span></span></div>",
			expected: "",
		},
		{
			name:     "text with html entities",
			input:    "&lt;script&gt;alert('xss')&lt;/script&gt;",
			expected: "&lt;script&gt;alert('xss')&lt;/script&gt;", // entities are escaped for safety
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HTML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple path",
			input:    "hello/world",
			expected: "hello/world",
		},
		{
			name:     "path with spaces",
			input:    "hello world/test path",
			expected: "hello world/test path", // Path() is restrictive, spaces are allowed
		},
		{
			name:     "path with uppercase",
			input:    "Hello/World",
			expected: "hello/world",
		},
		{
			name:     "path with special characters",
			input:    "hello!@#$%^&*()world",
			expected: "hello-@-$%^-*()world", // Path() uses illegalPath regex which is very restrictive
		},
		{
			name:     "path with accents",
			input:    "café/naïve",
			expected: "cafe/naive",
		},
		{
			name:     "path with double dots",
			input:    "../../../etc/passwd",
			expected: "/etc/passwd", // path.Clean() processes this to /etc/passwd
		},
		{
			name:     "path with dots",
			input:    "./test/../file.txt",
			expected: "test/file.txt", // path.Clean() preserves relative structure
		},
		{
			name:     "path with tilde",
			input:    "~/documents/file",
			expected: "~/documents/file",
		},
		{
			name:     "path with dash",
			input:    "my-file/my-folder",
			expected: "my-file/my-folder",
		},
		{
			name:     "empty string",
			input:    "",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Path(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "file.txt",
			expected: "file-txt",
		},
		{
			name:     "filename with spaces",
			input:    "my file name.mp3",
			expected: "my file name-mp3", // Only dots/slashes replaced, not spaces
		},
		{
			name:     "filename with path",
			input:    "/path/to/file.txt",
			expected: "-path-to-file-txt", // Slashes replaced with dashes before basename
		},
		{
			name:     "filename with special characters",
			input:    "file!@#$%name.txt",
			expected: "file-@-$%name-txt", // Some special chars preserved
		},
		{
			name:     "filename with accents",
			input:    "café.mp3",
			expected: "cafe-mp3",
		},
		{
			name:     "filename with multiple dots",
			input:    "my.file.name.txt",
			expected: "my-file-name-txt",
		},
		{
			name:     "filename with slashes",
			input:    "my/file/name.txt",
			expected: "my-file-name-txt", // Slashes replaced before basename
		},
		{
			name:     "empty string",
			input:    "",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Name(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBaseName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple basename",
			input:    "filename",
			expected: "filename",
		},
		{
			name:     "basename with dots",
			input:    "file.name.txt",
			expected: "file-name-txt",
		},
		{
			name:     "basename with slashes",
			input:    "path/to/file",
			expected: "path-to-file",
		},
		{
			name:     "basename with special characters",
			input:    "file!@#name",
			expected: "file-@-name", // Some special chars preserved
		},
		{
			name:     "basename with spaces",
			input:    "my file name",
			expected: "my file name", // Spaces not replaced
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BaseName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAccents(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no accents",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "french accents",
			input:    "café",
			expected: "cafe",
		},
		{
			name:     "german umlauts",
			input:    "Müller",
			expected: "Mueller",
		},
		{
			name:     "spanish characters",
			input:    "niño",
			expected: "nino",
		},
		{
			name:     "multiple accents",
			input:    "àáâãäå",
			expected: "aaaaaeaa", // å = aa, ä = ae
		},
		{
			name:     "uppercase accents",
			input:    "ÀÁÂÃÄÅ",
			expected: "AAAAAAA",
		},
		{
			name:     "mixed content",
			input:    "Café München naïve",
			expected: "Cafe Muenchen naive",
		},
		{
			name:     "polish characters",
			input:    "Łódź",
			expected: "Lodź", // Capital Ł mapped but ź not in map
		},
		{
			name:     "nordic characters",
			input:    "Øresund",
			expected: "OEresund",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Accents(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHTMLAllowing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		args     [][]string
		expected string
	}{
		{
			name:     "allow default tags",
			input:    "<p>Hello <strong>World</strong></p>",
			expected: "<p>Hello <strong>World</strong></p>",
		},
		{
			name:     "strip script tags",
			input:    "<p>Hello</p><script>alert('xss')</script>",
			expected: "<p>Hello</p>",
		},
		{
			name:     "strip style tags",
			input:    "<p>Hello</p><style>.bad{color:red}</style>",
			expected: "<p>Hello</p>",
		},
		{
			name:     "allow links",
			input:    "<a href='http://example.com'>Link</a>",
			expected: "<a href=\"http://example.com\">Link</a>",
		},
		{
			name:     "allow images",
			input:    "<img src='image.jpg' alt='test'>",
			expected: "<img src=\"image.jpg\" alt=\"test\">", // No self-closing slash
		},
		{
			name:     "allow custom tags",
			input:    "<custom>content</custom><p>para</p>",
			args:     [][]string{{"custom"}, {}},
			expected: "<custom>content</custom>para", // <p> is stripped but content remains
		},
		{
			name:     "strip iframe",
			input:    "<p>Hello</p><iframe src='evil.com'></iframe>",
			expected: "<p>Hello</p>",
		},
		{
			name:     "nested allowed tags",
			input:    "<div><span><strong>Nested</strong></span></div>",
			expected: "<div><span><strong>Nested</strong></span></div>",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			var err error
			if len(tt.args) > 0 {
				result, err = HTMLAllowing(tt.input, tt.args...)
			} else {
				result, err = HTMLAllowing(tt.input)
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIncludes(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		search   string
		expected bool
	}{
		{
			name:     "found in slice",
			slice:    []string{"a", "b", "c"},
			search:   "b",
			expected: true,
		},
		{
			name:     "not found in slice",
			slice:    []string{"a", "b", "c"},
			search:   "d",
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			search:   "a",
			expected: false,
		},
		{
			name:     "empty search string",
			slice:    []string{"a", "b", "c"},
			search:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := includes(tt.slice, tt.search)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello world",
			expected: "hello world", // cleanString doesn't replace spaces by default
		},
		{
			name:     "string with separators",
			input:    "hello!world&test",
			expected: "hello-world-test",
		},
		{
			name:     "string with multiple dashes",
			input:    "hello---world",
			expected: "hello-world",
		},
		{
			name:     "string with accents",
			input:    "café münchën",
			expected: "cafe muenchen", // Accents removed but spaces preserved
		},
		{
			name:     "string with trailing spaces",
			input:    "  hello world  ",
			expected: "hello world", // Trailing spaces trimmed
		},
		{
			name:     "string with special chars",
			input:    "test#123+456:789",
			expected: "test-123-456-789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanString(tt.input, illegalName)
			assert.Equal(t, tt.expected, result)
		})
	}
}
