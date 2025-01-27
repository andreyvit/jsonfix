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

type container byte

const (
	object container = 0
	array  container = 1
)

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
		inBareObjectKey
	)
	var state int = normal
	var start int
	var comma int = -1 // index of undecided comma in result (could be followed by whitespace)
	var isStartOfKey bool
	var stackBuf [16]container
	stack := stackBuf[:0]
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
		case inBareObjectKey:
			if !(isSpace[b] || b == ':' || b == '}' || b == '/') {
				continue
			}
			result = append(result, data[start:i]...)
			result = append(result, '"')
			start = i
			state = normal
			fallthrough
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
			switch b {
			case '"':
				state = inString
				isStartOfKey = false
			case ',':
				result = append(result, data[start:i+1]...)
				comma = len(result) - 1
				start = i + 1
				isStartOfKey = (len(stack) > 0) && (stack[len(stack)-1] == object)
			case '{':
				stack = append(stack, object)
				isStartOfKey = true
			case '[':
				stack = append(stack, array)
				isStartOfKey = false
			case ']':
				if len(stack) > 0 && stack[len(stack)-1] == array {
					stack = stack[:len(stack)-1]
				}
				isStartOfKey = false
			case '}':
				if len(stack) > 0 && stack[len(stack)-1] == object {
					stack = stack[:len(stack)-1]
				}
				isStartOfKey = false
			default:
				if isStartOfKey {
					result = append(result, data[start:i]...)
					result = append(result, '"')
					start = i
					state = inBareObjectKey
				}
				isStartOfKey = false
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
