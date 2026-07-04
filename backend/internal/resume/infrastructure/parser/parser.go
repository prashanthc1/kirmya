// Package parser extracts plain text from resume uploads. DOCX and TXT are
// fully supported (stdlib only); PDF is best-effort (decodes Flate streams and
// pulls text operators). Implements resume/domain.Parser.
package parser

import (
	"archive/zip"
	"bytes"
	"compress/zlib"
	"io"
	"regexp"
	"strings"

	"workspace-app/internal/resume/domain"
)

type Parser struct{}

func New() *Parser { return &Parser{} }

func (p *Parser) ExtractText(filename, contentType string, data []byte) (string, error) {
	if len(data) == 0 {
		return "", domain.ErrEmptyUpload
	}
	switch detectKind(filename, contentType) {
	case "txt":
		return normalize(string(data)), nil
	case "docx":
		return extractDocx(data)
	case "pdf":
		return extractPDF(data), nil // best-effort; empty text is tolerated downstream
	default:
		return "", domain.ErrUnsupported
	}
}

func detectKind(filename, contentType string) string {
	name := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(name, ".docx") || strings.Contains(contentType, "wordprocessingml"):
		return "docx"
	case strings.HasSuffix(name, ".pdf") || strings.Contains(contentType, "pdf"):
		return "pdf"
	case strings.HasSuffix(name, ".txt") || strings.HasPrefix(contentType, "text/"):
		return "txt"
	default:
		return ""
	}
}

var (
	reTag        = regexp.MustCompile(`<[^>]+>`)
	reWhitespace = regexp.MustCompile(`[ \t]+`)
	reBlankLines = regexp.MustCompile(`\n{3,}`)
	reParen      = regexp.MustCompile(`\(((?:[^()\\]|\\.)*)\)`)
)

func extractDocx(data []byte) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", domain.ErrUnsupported
	}
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		raw, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return "", err
		}
		xml := string(raw)
		// Paragraph and tab breaks become whitespace before stripping tags.
		xml = strings.ReplaceAll(xml, "</w:p>", "\n")
		xml = strings.ReplaceAll(xml, "<w:tab/>", " ")
		xml = strings.ReplaceAll(xml, "<w:br/>", "\n")
		text := reTag.ReplaceAllString(xml, "")
		return normalize(unescapeXML(text)), nil
	}
	return "", domain.ErrUnsupported
}

// extractPDF is a best-effort extractor: it inflates FlateDecode streams and
// collects text drawn with the (string) Tj/TJ operators. Returns "" if nothing
// readable is found (the scorer then flags it).
func extractPDF(data []byte) string {
	var sb strings.Builder
	const marker = "stream"
	rest := data
	for {
		i := bytes.Index(rest, []byte(marker))
		if i < 0 {
			break
		}
		after := rest[i+len(marker):]
		// Skip CRLF/LF following the stream keyword.
		after = bytes.TrimLeft(after, "\r\n")
		end := bytes.Index(after, []byte("endstream"))
		if end < 0 {
			break
		}
		chunk := after[:end]
		if decoded, ok := inflate(chunk); ok {
			collectText(&sb, decoded)
		} else {
			collectText(&sb, chunk)
		}
		rest = after[end+len("endstream"):]
	}
	return normalize(sb.String())
}

func inflate(b []byte) ([]byte, bool) {
	zr, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, false
	}
	defer func() { _ = zr.Close() }()
	out, err := io.ReadAll(zr)
	if err != nil || len(out) == 0 {
		return nil, false
	}
	return out, true
}

func collectText(sb *strings.Builder, content []byte) {
	for _, m := range reParen.FindAllStringSubmatch(string(content), -1) {
		s := strings.NewReplacer(`\(`, "(", `\)`, ")", `\\`, `\`).Replace(m[1])
		if strings.TrimSpace(s) != "" {
			sb.WriteString(s)
			sb.WriteByte(' ')
		}
	}
}

func unescapeXML(s string) string {
	return strings.NewReplacer(
		"&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&apos;", "'",
	).Replace(s)
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = reWhitespace.ReplaceAllString(s, " ")
	s = reBlankLines.ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}
