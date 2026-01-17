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
    }
}
