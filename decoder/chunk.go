package decoder

import "encoding/binary"

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
func chunkTypeFromBytes(type_ []byte) chunk_type {
    switch string(type_) {
    case "IHDR":
        return IHDR
    case "PLTE":
        return PLTE
    case "pHYs":
        return pHYs
    case "tEXt":
        return tEXt
    case "iTXt":
        return iTXt
    case "IDAT":
        return IDAT
    case "IEND":
        return IEND
    default:
        return -1 // unknown
    }
}

// TODO: color type
type IHDRChunk struct {
	Width           uint32
    Height          uint32
    BitDepth        byte
    ColorType       byte
    Compression     byte
    Filter          byte
    Interlace       byte
}

func parseIHDR(data []byte) *IHDRChunk {
	w := binary.BigEndian.Uint32(data[:4])
	h := binary.BigEndian.Uint32(data[4:8])
	return &IHDRChunk{
		Width: w,
		Height: h,
		BitDepth: data[8],
		ColorType: data[9],
		Compression: data[10],
		Filter: data[11],
		Interlace: data[12],
	}
}

type PLTEChunk struct {

}

type IDATChunk struct {

}

type IENDChunk struct {

}
