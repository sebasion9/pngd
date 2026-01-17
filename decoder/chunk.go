package decoder

type Chunk interface {
	Data() []byte
	Type() chunk_type
	Length() uint64
}

type BaseChunk struct {
	type_ chunk_type
	data []byte
	length uint64
}

func(c *BaseChunk) Data() []byte { return c.data}
func(c *BaseChunk) Type() chunk_type { return c.type_}
func(c *BaseChunk) Length() uint64 { return c.length}

type chunk_type int
const (
        IHDR chunk_type = iota
		PLTE
		IDAT
		IEND

		tEXt
)

//TODO: all types
var chunk_type_map = map[string]chunk_type{
	"IHDR": IHDR,
	"PLTE": PLTE,
	"IDAT": IDAT,
	"IEND": IEND,

	"tEXt": tEXt,
}

type PLTEChunk struct {

}

