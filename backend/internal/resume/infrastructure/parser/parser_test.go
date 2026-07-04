package parser

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"

	"workspace-app/internal/resume/domain"
)

func makeDocx(t *testing.T, body string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	xml := `<?xml version="1.0"?><w:document><w:body>` + body + `</w:body></w:document>`
	if _, err := w.Write([]byte(xml)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractDocx(t *testing.T) {
	data := makeDocx(t, `<w:p><w:r><w:t>Senior Operations Manager</w:t></w:r></w:p><w:p><w:r><w:t>Led teams &amp; budgets</w:t></w:r></w:p>`)
	p := New()
	text, err := p.ExtractText("resume.docx", "", data)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !strings.Contains(text, "Senior Operations Manager") {
		t.Errorf("missing first paragraph: %q", text)
	}
	if !strings.Contains(text, "Led teams & budgets") {
		t.Errorf("entities not unescaped / second paragraph missing: %q", text)
	}
}

func TestExtractTxt(t *testing.T) {
	p := New()
	text, err := p.ExtractText("resume.txt", "text/plain", []byte("hello   world"))
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if text != "hello world" {
		t.Errorf("expected normalized whitespace, got %q", text)
	}
}

func TestExtractUnsupportedAndEmpty(t *testing.T) {
	p := New()
	if _, err := p.ExtractText("resume.xyz", "application/octet-stream", []byte("data")); err != domain.ErrUnsupported {
		t.Errorf("expected ErrUnsupported, got %v", err)
	}
	if _, err := p.ExtractText("resume.txt", "text/plain", nil); err != domain.ErrEmptyUpload {
		t.Errorf("expected ErrEmptyUpload, got %v", err)
	}
}
