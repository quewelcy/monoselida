package md

import (
	"bytes"
)

type Chapter struct {
	Title string
	Text  string
}

type Md struct {
	TitleInfo   string
	Description string
	Chapters    []Chapter
}

func (md *Md) SetTitle(title string) {
	md.TitleInfo = title
}

func (md *Md) SetAnnotation(annotation string) {
	md.Description = annotation
}

func (md *Md) AppendChapter(title, text string) {
	md.Chapters = append(md.Chapters, Chapter{
		Title: title,
		Text:  text,
	})
}

func (md Md) Bytes() []byte {
	var buffer bytes.Buffer
	if md.TitleInfo != "" {
		buffer.WriteString("# " + md.TitleInfo + "\n\n")
	}
	if md.Description != "" {
		buffer.WriteString(md.Description + "\n\n")
	}
	for _, chapter := range md.Chapters {
		if chapter.Title != "" {
			buffer.WriteString("## " + chapter.Title + "\n\n")
		}
		if chapter.Text != "" {
			buffer.WriteString(chapter.Text + "\n\n")
		}
	}
	return buffer.Bytes()
}

func (md Md) String() string {
	return string(md.Bytes())
}

func Init(title, annotation string) *Md {
	return &Md{
		TitleInfo:   title,
		Description: annotation,
	}
}
