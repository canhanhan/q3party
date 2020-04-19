package protocol

const maxStringTokens = 1024

func TokenizeString(data []byte) []string {
	tokens := []string{}
	var search string
	for i := 0; i < len(data); i++ {
		if search != "" {
			if string(data[i:i+len(search)]) == search {
				i += len(search)
				search = ""
			}

			continue
		}

		var j int
		for j = i; j < len(data); j++ {
			// Skip whitespace
			if data[j] > ' ' {
				break
			}
		}

		i = j
		if i >= len(data) {
			break
		}

		// Skip // comments
		if data[i] == '/' && data[i+1] == '/' {
			break
		}

		// Skip /* */ comments
		if data[i] == '/' && data[i+1] == '*' {
			search = "*/"
			continue
		}

		if search != "" {
			continue
		}

		var token string
		if data[i] == '"' {
			start := i + 1

			var j int
			for j = i + 1; j < len(data); j++ {
				if data[j] == '"' && data[j-1] != '\\' {
					break
				}
			}

			i = j
			token = string(data[start:i])
			tokens = append(tokens, token)
		} else {
			// skip until whitespace, quote, or command
			var j int
			for j = i + 1; j < len(data); j++ {
				if data[j] <= ' ' {
					break
				}

				if data[j] == '"' {
					break
				}

				// skip // comments
				if data[j] == '/' && data[j+1] == '/' {
					break
				}

				// skip /* */ comments
				if data[j] == '/' && data[j+1] == '*' {
					break
				}
			}

			token = string(data[i:j])
			tokens = append(tokens, token)
			i = j
		}
	}

	return tokens
}

func SliceEqual(data []byte, value string) bool {
	l := len(value)
	if len(data) < l {
		return false
	}

	return string(data[0:l]) == value
}
