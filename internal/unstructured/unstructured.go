// Package unstructured provides utility functions for working with unstructured
// objects. For example, YAML or JSON data that cannot be strongly typed.
package unstructured

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
