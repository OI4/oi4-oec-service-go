package utils

func DNPEncode(input string) string {
	var encoded = []rune{}
	for _, char := range input {
		switch char {
		case '#':
			encoded = append(encoded, ',', '2', '3')
		case '/':
			encoded = append(encoded, ',', '2', 'F')
		case ':':
			encoded = append(encoded, ',', '3', 'A')
		case '?':
			encoded = append(encoded, ',', '3', 'F')
		case '@':
			encoded = append(encoded, ',', '4', '0')
		case '[':
			encoded = append(encoded, ',', '5', 'B')
		case ']':
			encoded = append(encoded, ',', '5', 'D')
		case '!':
			encoded = append(encoded, ',', '2', '1')
		case '$':
			encoded = append(encoded, ',', '2', '4')
		case '&':
			encoded = append(encoded, ',', '2', '6')
		case '‘':
			encoded = append(encoded, ',', '2', '7')
		case '(':
			encoded = append(encoded, ',', '2', '8')
		case ')':
			encoded = append(encoded, ',', '2', '9')
		case '*':
			encoded = append(encoded, ',', '2', 'A')
		case '+':
			encoded = append(encoded, ',', '2', 'B')
		case ',':
			encoded = append(encoded, ',', '2', 'C')
		case ';':
			encoded = append(encoded, ',', '3', 'B')
		case '=':
			encoded = append(encoded, ',', '3', 'D')
		case ' ':
			encoded = append(encoded, ',', '2', '0')
		case '"':
			encoded = append(encoded, ',', '2', '2')
		case '%':
			encoded = append(encoded, ',', '2', '5')
		case '<':
			encoded = append(encoded, ',', '3', 'C')
		case '>':
			encoded = append(encoded, ',', '3', 'E')
		case '\\':
			encoded = append(encoded, ',', '5', 'C')
		case '^':
			encoded = append(encoded, ',', '5', 'E')
		case '`':
			encoded = append(encoded, ',', '6', '0')
		case '{':
			encoded = append(encoded, ',', '7', 'B')
		case '|':
			encoded = append(encoded, ',', '7', 'C')
		case '}':
			encoded = append(encoded, ',', '7', 'D')
		default:
			encoded = append(encoded, char)
		}

	}
	return string(encoded)
}
