package decoder

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"pngd/errors"
	"pngd/internal"
)


type Decoder struct {
	Filter
	source []byte
	pos uint64

	raw []byte
	chunks []Chunk

	textChunks []TEXTChunk
	idatChunks []IDATChunk
	ihdr IHDRChunk
	iend IENDChunk

	// PLTE PLTEChunk

	warnings []string
}

func (d *Decoder) SetSrc(src []byte) {
	d.source = src
	d.pos = 0
	d.chunks = d.chunks[:0]
	d.textChunks = d.textChunks[:0]
	d.idatChunks = d.idatChunks[:0]
	d.raw = d.raw[:0]
	d.ihdr = IHDRChunk{}
	d.iend = IENDChunk{}
	d.Filter = Filter{}
}

func (d *Decoder) Chunks() []Chunk {
	return d.chunks
}

func (d *Decoder) TextChunks() []TEXTChunk {
	return d.textChunks
}

func (d *Decoder) IdatChunks() []IDATChunk {
	return d.idatChunks
}

func (d *Decoder) IHDR() *IHDRChunk {
	return &d.ihdr
}

func (d *Decoder) IEND() *IENDChunk {
	return &d.iend
}

func (d *Decoder) Info() []string {
	info := make([]string, 4)
	info[0] = "[i] IHDR dump start\n"
	info[1] = fmt.Sprintf("[i] interlace:  %d\n", d.ihdr.Interlace)
	info[2] = fmt.Sprintf("[i] color type: %v\n", d.ihdr.ColorType)
	info[3] = "[i] IHDR dump end\n"
	return info
}

func (d *Decoder) Warnings() []string {
	return d.warnings
}

func NewDecoder(source []byte) *Decoder {
	return &Decoder{ source: source }
}

func (d *Decoder) ValidateSignature() error {
	if len(d.source) < 8 {
		return errors.NewInvalidSignatureError("Invalid signature (less than 8 bytes)")
	}

	var signature uint64
	signature = 0x89504E470D0A1A0A

	packed := uint64(d.source[7]) |
	(uint64(d.source[6]) << 8)	|
	(uint64(d.source[5]) << 16)	|
	(uint64(d.source[4]) << 24)	|
	(uint64(d.source[3]) << 32)	|
	(uint64(d.source[2]) << 40)	|
	(uint64(d.source[1]) << 48)	|
	(uint64(d.source[0]) << 56);

	if packed != signature {
		return errors.NewInvalidSignatureError(
			fmt.Sprintf("Invalid signature:\n%s",
			internal.GotExpectedFmt(signature, packed)),
		)
	}

	d.pos += 8

	return nil
}

// uint4 chunk_len
// 4 bytes chunk type
// chunk_len bytes of data
// 4 byte crc - check for data integrity

func (d *Decoder) Decode() ([]byte, error) {
	for d.pos < uint64(len(d.source)) {
		if err := d.ReadChunk(); err != nil {
			return nil, err
		}
		if d.chunks[len(d.chunks)-1].Type() == IEND {
			break
		}
	}

	var buf bytes.Buffer
	for _, idat := range d.idatChunks {
		buf.Write(idat.data)
	}


	zr, err := zlib.NewReader(&buf)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	decomp, err := io.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	d.raw = decomp


	row_size := int(1 + d.ihdr.Width * uint32(d.ihdr.Bpp))
	d.Filter.Scanlines = make([][]byte, d.ihdr.Height)
	for y := 0; y < int(d.ihdr.Height); y++ {
		start := y * row_size
		end := start + row_size

		d.Filter.Scanlines[y] = d.raw[start:end]
	}

	recons, err := d.Filter.Reconstruct(d.ihdr.Bpp)
	if err != nil {
		return nil, err
	}


	return recons, nil
}

func (d *Decoder) ReadChunk() error {
	var chunk_len []byte
	var chunk_type []byte
	var chunk_data []byte
	var chunk_crc []byte

	if d.pos < 8 {
		return errors.NewInvalidSignatureError("Did not validate signature")
	}

	if d.pos + 4 > uint64(len(d.source)) {
		return errors.NewInvalidChunkError("Failed to read chunk len")
	}

	if d.pos + 8 > uint64(len(d.source)) {
		return errors.NewInvalidChunkError("Failed to read chunk type")
	}

	if d.pos + 12 > uint64(len(d.source)) {
		return errors.NewInvalidChunkError("Failed to read chunk crc")
	}

	chunk_len = d.source[d.pos:d.pos+4]
	chunk_type = d.source[d.pos+4:d.pos+8]

	chunk_len_uint := uint64(binary.BigEndian.Uint32(chunk_len))
	chunk_data = d.source[d.pos+8:d.pos+8+chunk_len_uint]

	chunk_crc = d.source[d.pos+8+chunk_len_uint:d.pos+12+chunk_len_uint]

	if !d.checkCRC(chunk_type, chunk_data, chunk_crc) {
		return errors.NewInvalidCRCError("CRC Hash compare failed, image file could be corrupted")
	}

	ct, ok := chunk_type_map[string(chunk_type)];
	if !ok {
		d.warnings = append(d.warnings, fmt.Sprintf("[W] \"%s\" is not supported", string(chunk_type)))
		d.pos += 12 + chunk_len_uint
		return nil
	}

	switch ct {
	case IHDR:
		if chunk_len_uint != 13 {
			return errors.NewInvalidIHDRChunk(
				fmt.Sprintf("IHDR chunk lenght is incorrect:\n%s",
				internal.GotExpectedFmt(chunk_len_uint,"13"),
			))
		}
		if len(d.chunks) > 0 {
			return errors.NewInvalidIHDRChunk(
				fmt.Sprintf("IHDR chunk should be FIRST, chunks:\n%s",
				internal.GotExpectedFmt(len(d.chunks), "0"),
			))
		}

		d.ihdr = *parseIHDR(chunk_data)

		if d.ihdr.Interlace == 1 {
			return errors.NewNotImplementedError("Interlacing not implemented")
		}

		d.chunks = append(d.chunks, &d.ihdr)
	case tEXt:
		tEXt_chunk, err := parseTEXT(chunk_data, chunk_len_uint)
		if err != nil {
			return err
		}

		d.textChunks = append(d.textChunks, *tEXt_chunk)
		d.chunks = append(d.chunks, tEXt_chunk)
	case IDAT:
		idat_chunk, err := parseIDAT(chunk_data, chunk_len_uint)
		if err != nil {
			return err
		}

		d.idatChunks = append(d.idatChunks, *idat_chunk)
		d.chunks = append(d.chunks, idat_chunk)
	// case PLTE:
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
		d.iend = *iend_chunk
		d.chunks = append(d.chunks, iend_chunk)
	default:
	}


	d.pos += 12 + chunk_len_uint
	return nil 
}


func (d *Decoder) checkCRC(chunk_type []byte, chunk_data []byte, chunk_crc []byte) bool {
	crc := crc32.ChecksumIEEE(append(chunk_type, chunk_data...))
	expected := binary.BigEndian.Uint32(chunk_crc)
	if crc != expected {
		return false
	}
	return true
}


