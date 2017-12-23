package md

import (
	"bytes"
)

type PTitle struct {
	Text string
}

type PText struct {
	Text string
}

type Md struct {
	TitleInfo   string
	Description string
	Paragraphs  []interface{}
}

func (md *Md) SetTitle(title string) {
	md.TitleInfo = title
}

func (md *Md) SetAnnotation(annotation string) {
	md.Description = annotation
}

func (md *Md) AppendTitle(title string) {
	md.Paragraphs = append(md.Paragraphs, PTitle{title})
}

func (md *Md) AppendText(text string) {
	md.Paragraphs = append(md.Paragraphs, PText{text})
}

func (md Md) Bytes() []byte {
	var buffer bytes.Buffer
	if md.TitleInfo != "" {
		buffer.WriteString("# " + md.TitleInfo + "\n\n")
	}
	if md.Description != "" {
		buffer.WriteString(md.Description + "\n\n")
	}
	for _, paragraph := range md.Paragraphs {
		if v, ok := paragraph.(PTitle); ok {
			buffer.WriteString("## " + v.Text + "\n\n")
		}
		if v, ok := paragraph.(PText); ok {
			buffer.WriteString(v.Text + "\n\n")
		}
	}
	return buffer.Bytes()
}

func (md Md) String() string {
	return string(md.Bytes())
}

func InitWithTitle(title, annotation string) *Md {
	return &Md{
		TitleInfo:   title,
		Description: annotation,
	}
}

func Init() *Md {
	return &Md{}
}
