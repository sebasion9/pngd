package decoder

import (
	"bytes"
	"pngd/errors"
)

type TEXTChunk struct {
    BaseChunk
    Keyword string
    Text string
}

var keyword_map = map[string]Keyword{
    "Title": Title,
    "Author": Author,
    "Description": Description,
    "Copyright": Copyright,
    "Creation": Creation,
    "Software": Software,
    "Disclaimer": Disclaimer,
    "Warning": Warning,
    "Source": Source,
    "Comment": Comment,
}

type Keyword int
const (
    Title = iota
    Author
    Description
    Copyright
    Creation
    Software
    Disclaimer
    Warning
    Source
    Comment
)

func parseTEXT(data []byte, length uint64) (*TEXTChunk, error) {
    idx := bytes.IndexByte(data, 0)
    if idx <= 0 || idx > 79 {
	return nil, errors.NewInvalidtEXtChunk("Null separator either doesn't exit or is too far")
    }
    keyword := string(data[:idx])
    text := string(data[idx+1:])

    return &TEXTChunk{
	BaseChunk: BaseChunk{
	    data: data,
	    length: length,
	    type_: tEXt,
	},
	Keyword: keyword,
	Text: text,
    }, nil
}
