package transfig_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/transfig"
)

func Test_MapDeepCopy(t *testing.T) {
	original := map[KeyString]interface{}{
		"key1": "value1",
		"key2": map[KeyString]interface{}{
			"key3": "value3",
		},
	}
	copy := MapDeepCopy(original)
	assert.Equal(t, original, copy)
	original["key1"] = "value4"
	copy["key1"] = "value5"
	assert.Equal(t, "value4", original["key1"])
	assert.Equal(t, "value5", copy["key1"])
	original["key2"].(map[KeyString]interface{})["key3"] = "value6"
	copy["key2"].(map[KeyString]interface{})["key3"] = "value7"
	assert.Equal(t, "value6", original["key2"].(map[KeyString]interface{})["key3"])
	assert.Equal(t, "value7", copy["key2"].(map[KeyString]interface{})["key3"])
}

func Test_MapSetNested(t *testing.T) {
	m := map[KeyString]interface{}{
		"key1": map[KeyString]interface{}{
			"key2": "value2",
			"key3": map[KeyString]interface{}{
				"key4": "value4",
			},
		},
	}
	MapSetNested(m, []KeyString{"key1", "key3", "key5"}, "value5")
	assert.Equal(t, "value2", m["key1"].(map[KeyString]interface{})["key2"])
	assert.Equal(t, "value5", m["key1"].(map[KeyString]interface{})["key3"].(map[KeyString]interface{})["key5"])

}

func Test_MapSetNested_OverrideExistingValue(t *testing.T) {
	m := map[KeyString]interface{}{
		"key1": map[KeyString]interface{}{
			"key2": "value2",
		},
	}
	MapSetNested(m, []KeyString{"key1", "key2"}, "value3")
	assert.Equal(t, "value3", m["key1"].(map[KeyString]interface{})["key2"])
}

func Test_MapSetNested_OverrideExistingNonMapValue(t *testing.T) {
	m := map[KeyString]interface{}{
		"key1": "value1",
		"key2": map[KeyString]interface{}{"key3": "value3"},
	}
	MapSetNested(m, []KeyString{"key2", "key3", "key4"}, "value4")
	assert.Equal(t, map[KeyString]interface{}{
		"key1": "value1",
		"key2": map[KeyString]interface{}{
			"key3": map[KeyString]interface{}{
				"key4": "value4",
			},
		},
	}, m)
}

func Test_MapGetNested_Existing(t *testing.T) {
	m := map[KeyString]interface{}{
		"key1": map[KeyString]interface{}{
			"key3": map[KeyString]interface{}{
				"key4": "value4",
			},
		},
	}
	value, found := MapGetNested(m, []KeyString{"key1", "key3", "key4"})
	assert.True(t, found)
	assert.Equal(t, "value4", value)
}

func Test_MapGetNested_Missing(t *testing.T) {
	m := map[KeyString]interface{}{"key1": map[KeyString]interface{}{}}
	value, found := MapGetNested(m, []KeyString{"key1", "key2"})
	assert.False(t, found)
	assert.Nil(t, value)
}

func Test_MapGesNested_MissingNonMap(t *testing.T) {
	m := map[KeyString]interface{}{"key1": "value1"}
	value, found := MapGetNested(m, []KeyString{"key1", "key2"})
	assert.False(t, found)
	assert.Nil(t, value)
}
