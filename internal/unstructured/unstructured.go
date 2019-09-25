// Package unstructured provides utility functions for working with unstructured
// objects. For example, YAML or JSON data that cannot be strongly typed.
package unstructured

import "gopkg.in/yaml.v2"

// FindAll recursively finds every instance of the 'key' in the object, and
// calls the provided callback function for each value at the matching keys. It
// will not recurse through the values of a matching key.
func FindAll(obj map[string]interface{}, key string, cb func(interface{})) {
	for k, v := range obj {
		if key == k {
			cb(v)
			continue
		}

		if nested, ok := v.(map[string]interface{}); ok {
			FindAll(nested, key, cb)
		}
	}
}

// VisitOne recursively finds the first item in the MapSlice that satisfies the
// predicate, and calls the provided callback function on the selected item.
func VisitOne(obj yaml.MapSlice, pred func(yaml.MapItem) bool, visit func(*yaml.MapItem)) {
	for i := range obj {
		if pred(obj[i]) {
			visit(&obj[i])
			break
		}

		if nested, ok := obj[i].Value.(yaml.MapSlice); ok {
			VisitOne(nested, pred, visit)
		}
	}
}
