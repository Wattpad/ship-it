package ecrconsumer

import (
	"reflect"
)

func search(v reflect.Value, needle string) []string {
	if v.Kind() != reflect.Map {
		return []string{}
	}

	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key().String()
		if key == needle {
			return []string{key}
		}

		val := reflect.ValueOf(iter.Value().Interface())
		if val.Kind() == reflect.Map {
			path := search(val, needle)
			if len(path) == 0 {
				continue
			}
			return append([]string{key}, path...)
		}
	}

	return []string{}
}
