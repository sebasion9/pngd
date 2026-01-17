package decoder


type IDATChunk struct {
    BaseChunk
}

func parseIDAT(data []byte, length uint64) (*IDATChunk, error) {

    return &IDATChunk{
	BaseChunk: BaseChunk{
	    data: data,
	    length: length,
	    type_: IDAT,
	},
    }, nil
}
