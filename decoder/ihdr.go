package decoder

import "encoding/binary"

// TODO: color type
type IHDRChunk struct {
    BaseChunk
    Width uint32
    Height uint32
    BitDepth byte
    ColorType byte
    Compression byte
    Filter byte
    Interlace byte
    Bpp byte
}

var bpp_map = map[byte]byte {
    0:1,
    2:3,
    3:1,
    4:2,
    6:4,
}

type ColorType int
const (
    RGB ColorType = iota
    RGBA
    G
    GA
    IDX
)

var color_type_map = map[byte]ColorType{
    0:G,
    2:RGB,
    3:IDX,
    4:GA,
    6:RGBA,
}

func (c *IHDRChunk) GetColorType() ColorType {
    return color_type_map[c.ColorType]
}

func parseIHDR(data []byte) *IHDRChunk {
    w := binary.BigEndian.Uint32(data[:4])
    h := binary.BigEndian.Uint32(data[4:8])
    return &IHDRChunk{
	BaseChunk: BaseChunk{
	    data: data,
	    length: 13,
	    type_: IHDR,
	},
	Width: w,
	Height: h,
	BitDepth: data[8],
	ColorType: data[9],
	Compression: data[10],
	Filter: data[11],
	Interlace: data[12],
	Bpp: bpp_map[data[9]],
    }
}
