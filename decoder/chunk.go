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

// pHYs: Non-square pixels can be represented (see the pHYs chunk), but viewers are not required to account for them; a viewer can present any PNG file as though its pixels are square.
type chunk_type int
const (
        IHDR chunk_type = iota
		pHYs
		tEXt
		iTXt
		IDAT
		IEND
		PLTE
)

//TODO: all types
var chunk_type_map = map[string]chunk_type{
	"IHDR": IHDR,
	"PLTE": PLTE,
	"pHYs": pHYs,
	"tEXt": tEXt,
	"iTXt": iTXt,
	"IDAT": IDAT,
	"IEND": IEND,
}

type PLTEChunk struct {

}

type IDATChunk struct {

}

type IENDChunk struct {

}

