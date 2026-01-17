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

type InvalidChunkLength struct {
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

func NewInvalidChunkLength(msg string) *InvalidChunkLength {
	return &InvalidChunkLength{&BaseDecoderError{Message:msg}}
}
