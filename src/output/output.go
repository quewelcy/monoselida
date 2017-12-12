package output

type OutputFormat interface {
	SetTitle(title string)
	SetAnnotation(annotation string)
	AppendChapter(title, text string)
	String() string
	Bytes() []byte
}
