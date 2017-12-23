package md

import (
	"bytes"
)

//PTitle describes body title
type PTitle struct {
	Text string
}

//PText describes body text
type PText struct {
	Text string
}

//Md describes whole Markdown file
type Md struct {
	TitleInfo   string
	Description string
	Paragraphs  []interface{}
}

//SetBookTitle sets main document title
func (md *Md) SetBookTitle(title string) {
	md.TitleInfo = title
}

//SetAnnotation sets document annotation
func (md *Md) SetAnnotation(annotation string) {
	md.Description = annotation
}

//AppendTitle appends title to document body
func (md *Md) AppendTitle(title string) {
	md.Paragraphs = append(md.Paragraphs, PTitle{title})
}

//AppendText appends text to document body
func (md *Md) AppendText(text string) {
	md.Paragraphs = append(md.Paragraphs, PText{text})
}

//Bytes marshals whole markdown file to byte array
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

//String marshals whole markdown file to string
func (md Md) String() string {
	return string(md.Bytes())
}

//InitWithTitle inits markdown file with title and annotation
func InitWithTitle(title, annotation string) *Md {
	return &Md{
		TitleInfo:   title,
		Description: annotation,
	}
}

//Init inits empty markdown file
func Init() *Md {
	return &Md{}
}
