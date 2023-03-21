package base

// Used to transform the driver cursor representation (can be any type, most likely a literal, array or map)
// into a "string" (that will be converted to a portable format through the Encoder).
// Keep this as simple as possible
type Marshaller interface {
	// When input is nil, must return an empty string
	Marshal(input interface{}) ([]byte, error)

	// When encoded is an empty string, return value must be nil
	Unmarshal(encoded []byte) (map[string]interface{}, error)
}

// Used to manipulate the output of the Marshaller (convert to base64 for portability or encrypt the cursor)
type Encoder interface {
	// When input is nil, must return an empty string
	Encode(input []byte) (string, error)

	// When encoded is an empty string, return value must be nil
	Decode(encoded string) ([]byte, error)
}
