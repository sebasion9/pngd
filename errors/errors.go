package errors

import "fmt"

type BaseDecoderError struct {
	Message string
}

func (e *BaseDecoderError) Error() string {
	return fmt.Sprintf("[!] %s\n", e.Message)
}

type InvalidSignatureError struct {
	*BaseDecoderError
}

type InvalidChunkError struct {
	*BaseDecoderError
}

type InvalidCRCError struct {
	*BaseDecoderError
}

type InvalidIHDRChunk struct {
	*BaseDecoderError
}

type InvalidtEXtChunk struct {
	*BaseDecoderError
}

func NewInvalidSignatureError(msg string) *InvalidSignatureError {
	return &InvalidSignatureError{&BaseDecoderError{Message: msg}}
}

func NewInvalidChunkError(msg string) *InvalidChunkError{
	return &InvalidChunkError{&BaseDecoderError{Message: msg}}
}

func NewInvalidCRCError(msg string) *InvalidCRCError {
	return &InvalidCRCError{&BaseDecoderError{Message:msg}}
}

func NewInvalidIHDRChunk(msg string) *InvalidIHDRChunk {
	return &InvalidIHDRChunk{&BaseDecoderError{Message:msg}}
}

func NewInvalidtEXtChunk(msg string) *InvalidtEXtChunk {
	return &InvalidtEXtChunk{&BaseDecoderError{Message:msg}}
}
