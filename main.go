package main

import (
	"encoding/json"
	"fmt"
	"go-json-tr/jsontransform"
)

func main() {

	// objArrayData, err := os.ReadFile("obj_array.json")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	// projData, err := os.ReadFile("proj.json")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }
	// var proj []interface{}
	// if err := json.Unmarshal(projData, &proj); err != nil {
	// 	fmt.Println("Error parsing JSON:", err)
	// 	return
	// }

	// var objArray []interface{}
	// if err := json.Unmarshal(objArrayData, &objArray); err != nil {
	// 	fmt.Println("Error parsing JSON:", err)
	// 	return
	// }

	input := []interface{}{
		map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": "baz",
				"tags": []interface{}{
					map[string]interface{}{"name": "tag1"},
					map[string]interface{}{"name": "tag2"},
				},
			},
		},
	}

	projection := []interface{}{
		map[string]interface{}{
			"data": map[string]interface{}{
				"bar": "$root.foo.bar",
				"tags": map[string]interface{}{
					"$array":   "$root.foo.tags",
					"tag_name": "$element.name",
				},
			},
		},
	}
	p, _ := json.MarshalIndent(projection, "", "  ")
	fmt.Println(string(p))
	output, err := jsontransform.Transform(input, projection)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	outputJSON, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println(string(outputJSON))
}
