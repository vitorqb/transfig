package transfig

// Helper functions that work with maps
func mapDeepCopy(m map[KeyString]interface{}) map[KeyString]interface{} {
	newMap := make(map[KeyString]interface{})
	for k, v := range m {
		if vAsMap, ok := v.(map[KeyString]interface{}); ok {
			newMap[k] = mapDeepCopy(vAsMap)
		} else {
			newMap[k] = v
		}
	}
	return newMap
}

// mapSetNested sets a nested key in a map
func mapSetNested(m map[KeyString]interface{}, keys []KeyString, value interface{}) {
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
	mapSetNested(topValueAsMap, keys[1:], value)
	m[keys[0]] = topValueAsMap
}

// mapGetNested gets a nested key in a map
func mapGetNested(m map[KeyString]interface{}, keys []KeyString) (value interface{}, found bool) {
	if len(keys) == 0 {
		return nil, false
	}
	if len(keys) == 1 {
		v, f := m[keys[0]]
		return v, f
	}
	topValue := m[keys[0]]
	if topValueAsMap, ok := topValue.(map[KeyString]interface{}); ok {
		return mapGetNested(topValueAsMap, keys[1:])
	}
	return nil, false
}
