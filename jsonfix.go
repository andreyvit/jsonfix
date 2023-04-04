// Package jsonfix allows trailing commas and comments in JSON.
// Just run it before passing data to json.Unmarshal.
package jsonfix

var (
	isSpace [256]bool
)

func init() {
	isSpace[' '] = true
	isSpace['\t'] = true
	isSpace['\n'] = true
	isSpace['\r'] = true
}

// Bytes removes trailing commas and comments from JSON data.
// The result can be json.Unmarshal'ed normally.
// Assumes valid JSON on input, otherwise might produce invalid output
// Preserves line numbers and formatting.
func Bytes(data []byte) []byte {
	n := len(data)
	result := make([]byte, 0, n)
	const (
		normal = iota
		inString
		inStringAfterSlash
		inLineComment
	)
	var state int = normal
	var start int
	var comma int = -1 // index of undecided comma in result (could be followed by whitespace)
	for i, b := range data {
		switch state {
		case inStringAfterSlash:
			state = inString
		case inString:
			if b == '"' {
				state = normal
			} else if b == '\\' {
				state = inStringAfterSlash
			}
		case inLineComment:
			if b == '\n' || b == '\r' {
				start = i
				state = normal
			}
		case normal:
			if isSpace[b] {
				continue
			}
			if b == '/' && (i+1 < n) && data[i+1] == '/' {
				result = append(result, data[start:i]...)
				state = inLineComment
				continue
			}
			if comma >= 0 {
				if b == ']' || b == '}' {
					result = deleteChar(result, comma)
				}
				comma = -1
			}
			if b == '"' {
				state = inString
			} else if b == ',' {
				result = append(result, data[start:i+1]...)
				comma = len(result) - 1
				start = i + 1
			}
		}
	}
	if state != inLineComment {
		result = append(result, data[start:]...)
	}
	return result
}

func deleteChar(data []byte, at int) []byte {
	n := len(data)
	if at+1 < n {
		copy(data[at:], data[at+1:])
	}
	return data[:n-1]
}
