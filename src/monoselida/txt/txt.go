package txt

import (
	"bytes"
)

type Txt struct {
	Title      string
	Paragraphs []string
}

func (txt *Txt) SetTitle(title string) {
	txt.Title = title
}
func (txt *Txt) SetAnnotation(annotation string) {
	txt.Paragraphs = append(txt.Paragraphs, annotation)
}

func (txt *Txt) AppendChapter(_, text string) {
	txt.Paragraphs = append(txt.Paragraphs, text)
}

func (txt Txt) Bytes() []byte {
	var buffer bytes.Buffer
	if txt.Title != "" {
		buffer.WriteString("# " + txt.Title + "\n\n")
	}
	for _, paragraph := range txt.Paragraphs {
		if paragraph != "" {
			buffer.WriteString(paragraph + "\n\n")
		}
	}
	return buffer.Bytes()
}

func (txt Txt) String() string {
	return string(txt.Bytes())
}

func Init(title, annotation string) *Txt {
	return &Txt{Title: title}
}
