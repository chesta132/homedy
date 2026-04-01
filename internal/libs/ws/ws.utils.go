package ws

func GetSubprotocolValue(protocols []string, key string) (value string) {
	for i, p := range protocols {
		if p == key && i+1 < len(protocols) {
			value = protocols[i+1]
			break
		}
	}
	return
}
