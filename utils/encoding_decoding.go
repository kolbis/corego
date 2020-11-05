package utils

import (
	"github.com/mitchellh/mapstructure"
)

// Decoder perform decoding operations
type Decoder struct {
}

// NewDecoder creates a new instance of the decoder
func NewDecoder() Decoder {
	return Decoder{}
}

// MapDecode perform decode from an input into an output
// the output must be reference type. Example:
// var output outerStruct
// decoder := core.NewDecoder()
// err := decoder.MapDecode(input, &output)
func (d Decoder) MapDecode(input interface{}, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &output,
		TagName:          "json",
		WeaklyTypedInput: true,
		ZeroFields:       false,
	}

	decoder, _ := mapstructure.NewDecoder(config)
	return decoder.Decode(input)
}
