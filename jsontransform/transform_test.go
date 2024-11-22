package jsontransform

import (
	"reflect"
	"testing"
)

func TestSimplifiedProjection(t *testing.T) {
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

	expected := []interface{}{
		map[string]interface{}{
			"data": map[string]interface{}{
				"bar": "baz",
				"tags": []interface{}{
					map[string]interface{}{"tag_name": "tag1"},
					map[string]interface{}{"tag_name": "tag2"},
				},
			},
		},
	}

	output, err := Transform(input, projection)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Errorf("expected %v, got %v", expected, output)
	}
}
