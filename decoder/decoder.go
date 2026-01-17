package decoder

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"pngd/errors"
	"pngd/util"
)


type Decoder struct {
	Warnings []string
	Source []byte
	pos uint64
	Chunks []Chunk
	TEXTChunks []TEXTChunk
	IDATChunks []IDATChunk
	IHDR IHDRChunk
	// PLTE PLTEChunk
	IDAT []IDATChunk
	IEND IENDChunk
}

func NewDecoder(source []byte) *Decoder {
	return &Decoder{ Source: source }
}

func (d *Decoder) ValidateSignature() error {
	if len(d.Source) < 8 {
		return errors.NewInvalidSignatureError("Invalid signature (less than 8 bytes)")
	}

	var signature uint64
	signature = 0x89504E470D0A1A0A

	packed := uint64(d.Source[7]) |
	(uint64(d.Source[6]) << 8)	|
	(uint64(d.Source[5]) << 16)	|
	(uint64(d.Source[4]) << 24)	|
	(uint64(d.Source[3]) << 32)	|
	(uint64(d.Source[2]) << 40)	|
	(uint64(d.Source[1]) << 48)	|
	(uint64(d.Source[0]) << 56);

	if packed != signature {
		return errors.NewInvalidSignatureError(
			fmt.Sprintf("Invalid signature:\n%s",
			util.GotExpectedFmt(signature, packed)),
		)
	}

	d.pos += 8

	return nil
}

// uint4 chunk_len
// 4 bytes chunk type
// chunk_len bytes of data
// 4 byte crc - check for data integrity

func (d *Decoder) Decode() error {
	for d.pos < uint64(len(d.Source)) {
		if err := d.ReadChunk(); err != nil {
			return err
		}
		if d.Chunks[len(d.Chunks)-1].Type() == IEND {
			return nil
		}
	}
	return nil
}

func (d *Decoder) ReadChunk() error {
	var chunk_len []byte
	var chunk_type []byte
	var chunk_data []byte
	var chunk_crc []byte
	if d.pos + 4 > uint64(len(d.Source)) {
		return errors.NewInvalidChunkError("Failed to read chunk len")
	}

	if d.pos + 8 > uint64(len(d.Source)) {
		return errors.NewInvalidChunkError("Failed to read chunk type")
	}

	if d.pos + 12 > uint64(len(d.Source)) {
		return errors.NewInvalidChunkError("Failed to read chunk crc")
	}

	chunk_len = d.Source[d.pos:d.pos+4]
	chunk_type = d.Source[d.pos+4:d.pos+8]

	chunk_len_uint := uint64(binary.BigEndian.Uint32(chunk_len))
	chunk_data = d.Source[d.pos+8:d.pos+8+chunk_len_uint]

	chunk_crc = d.Source[d.pos+8+chunk_len_uint:d.pos+12+chunk_len_uint]

	if !d.checkCRC(chunk_type, chunk_data, chunk_crc) {
		return errors.NewInvalidCRCError("CRC Hash compare failed, image file could be corrupted")
	}

	ct, ok := chunk_type_map[string(chunk_type)];
	if !ok {
		d.Warnings = append(d.Warnings, fmt.Sprintf("[W] \"%s\" is not supported", string(chunk_type)))
		d.pos += 12 + chunk_len_uint
		return nil
	}

	switch ct {
	case IHDR:
		if chunk_len_uint != 13 {
			return errors.NewInvalidIHDRChunk(
				fmt.Sprintf("IHDR chunk lenght is incorrect:\n%s",
				util.GotExpectedFmt(chunk_len_uint,"13"),
			))
		}
		if len(d.Chunks) > 0 {
			return errors.NewInvalidIHDRChunk(
				fmt.Sprintf("IHDR chunk should be FIRST, chunks:\n%s",
				util.GotExpectedFmt(len(d.Chunks), "0"),
			))
		}

		d.IHDR = *parseIHDR(chunk_data)
		d.Chunks = append(d.Chunks, &d.IHDR)
	case tEXt:
		tEXt_chunk, err := parseTEXT(chunk_data, chunk_len_uint)
		if err != nil {
			return err
		}

		d.TEXTChunks = append(d.TEXTChunks, *tEXt_chunk)
		d.Chunks = append(d.Chunks, tEXt_chunk)
	case IDAT:
		idat_chunk, err := parseIDAT(chunk_data, chunk_len_uint)
		if err != nil {
			return err
		}

		d.IDATChunks = append(d.IDATChunks, *idat_chunk)
		d.Chunks = append(d.Chunks, idat_chunk)
	case PLTE:
		//TODO:
	case IEND:
		if chunk_len_uint != 0 {
			return errors.NewInvalidIENDChunk("IEND chunk should be 0 bytes long")
		}
		iend_chunk := &IENDChunk{
			BaseChunk: BaseChunk {
				data: chunk_data, length: chunk_len_uint, type_: IEND,
			},
		}
		d.IEND = *iend_chunk
		d.Chunks = append(d.Chunks, iend_chunk)
	default:
	}


	fmt.Println("[L] ", chunk_len_uint)
	fmt.Println("[T] " + string(chunk_type[:]))
	fmt.Printf("[D] %x\n", chunk_data)

	d.pos += 12 + chunk_len_uint
	return nil 
}


func (d *Decoder) checkCRC(chunk_type []byte, chunk_data []byte, chunk_crc []byte) bool {
	crc := crc32.ChecksumIEEE(append(chunk_type, chunk_data...))
	expected := binary.BigEndian.Uint32(chunk_crc)
	if crc != expected {
		fmt.Printf("calculated:\t%x\n", crc)
		fmt.Printf("chunk_crc:\t%x\n", expected)
		return false
	}
	return true
}


