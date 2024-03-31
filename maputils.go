package transfig

// Helper functions that work with maps
func MapDeepCopy(m map[KeyString]interface{}) map[KeyString]interface{} {
	newMap := make(map[KeyString]interface{})
	for k, v := range m {
		if vAsMap, ok := v.(map[KeyString]interface{}); ok {
			newMap[k] = MapDeepCopy(vAsMap)
		} else {
			newMap[k] = v
		}
	}
	return newMap
}

// MapSetNested sets a nested key in a map
func MapSetNested(m map[KeyString]interface{}, keys []KeyString, value interface{}) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		m[keys[0]] = value
		return
	}
	topValue := m[keys[0]]
	topValueAsMap, topValueIsMap := topValue.(map[KeyString]interface{})
	if !topValueIsMap {
		topValueAsMap = make(map[KeyString]interface{})
	}
	MapSetNested(topValueAsMap, keys[1:], value)
	m[keys[0]] = topValueAsMap
}

// MapGetNested gets a nested key in a map
func MapGetNested(m map[KeyString]interface{}, keys []KeyString) (value interface{}, found bool) {
	if len(keys) == 0 {
		return nil, false
	}
	if len(keys) == 1 {
		v, f := m[keys[0]]
		return v, f
	}
	topValue := m[keys[0]]
	if topValueAsMap, ok := topValue.(map[KeyString]interface{}); ok {
		return MapGetNested(topValueAsMap, keys[1:])
	}
	return nil, false
}
