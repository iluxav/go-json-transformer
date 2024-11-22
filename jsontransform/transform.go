package jsontransform

import (
	"errors"
	"fmt"
	"strings"
)

func resolvePath(data interface{}, path string) (interface{}, error) {
	if strings.HasPrefix(path, "$root.") {
		path = strings.TrimPrefix(path, "$root.")
	}

	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		if part == "$element" {
			return current, nil
		}

		if currentMap, ok := current.(map[string]interface{}); ok {
			if val, exists := currentMap[part]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("key not found: %s", part)
			}
		} else {
			return nil, fmt.Errorf("invalid structure for key: %s", part)
		}
	}

	return current, nil
}

func projectObject(input interface{}, projection map[string]interface{}) (interface{}, error) {
	result := make(map[string]interface{})

	for key, value := range projection {
		switch v := value.(type) {
		case string:
			resolvedValue, err := resolvePath(input, v)
			if err != nil {
				return nil, fmt.Errorf("error resolving key %s: %w", key, err)
			}
			result[key] = resolvedValue

		case map[string]interface{}:
			if arrayPath, ok := v["$array"]; ok {
				resolvedArray, err := resolvePath(input, arrayPath.(string))
				if err != nil {
					return nil, fmt.Errorf("error resolving array path: %w", err)
				}

				arrayResults := []interface{}{}
				if array, ok := resolvedArray.([]interface{}); ok {
					elementProjection := make(map[string]interface{})
					for k, v := range v {
						if k != "$array" {
							elementProjection[k] = v
						}
					}

					for _, element := range array {
						projectedElement := make(map[string]interface{})
						for projKey, projValue := range elementProjection {
							if projValueStr, ok := projValue.(string); ok && strings.HasPrefix(projValueStr, "$element.") {
								fieldName := strings.TrimPrefix(projValueStr, "$element.")
								if elementMap, ok := element.(map[string]interface{}); ok {
									if val, exists := elementMap[fieldName]; exists {
										projectedElement[projKey] = val
									}
								}
							}
						}
						arrayResults = append(arrayResults, projectedElement)
					}
					result[key] = arrayResults
				} else {
					return nil, fmt.Errorf("resolved array is not valid: %v", resolvedArray)
				}
			} else {
				projectedValue, err := projectObject(input, v)
				if err != nil {
					return nil, fmt.Errorf("error projecting nested object for key %s: %w", key, err)
				}
				result[key] = projectedValue
			}

		default:
			return nil, fmt.Errorf("invalid projection type for key %s", key)
		}
	}

	return result, nil
}

func Transform(input []interface{}, projection []interface{}) ([]interface{}, error) {
	var output []interface{}

	for _, item := range input {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, errors.New("input items must be map[string]interface{}")
		}

		for _, proj := range projection {
			projMap, ok := proj.(map[string]interface{})
			if !ok {
				return nil, errors.New("projection must be map[string]interface{}")
			}

			projectedItem, err := projectObject(itemMap, projMap)
			if err != nil {
				return nil, err
			}

			output = append(output, projectedItem)
		}
	}

	return output, nil
}
