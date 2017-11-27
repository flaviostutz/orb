package geojson

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/paulmach/orb"
)

func TestNewFeature(t *testing.T) {
	f := NewFeature(orb.Point{1, 2})

	if f.Type != "Feature" {
		t.Errorf("incorrect feature: %v != Feature", f.Type)
	}
}

func TestFeatureMarshalJSON(t *testing.T) {
	f := NewFeature(orb.Point{1, 2})
	blob, err := f.MarshalJSON()

	if err != nil {
		t.Fatalf("error marshalling to json: %v", err)
	}

	if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestFeatureMarshalJSON_Bound(t *testing.T) {
	f := NewFeature(orb.Bound{Min: orb.Point{1, 1}, Max: orb.Point{2, 2}})
	blob, err := f.MarshalJSON()

	if err != nil {
		t.Fatalf("error marshalling to json: %v", err)
	}

	if !bytes.Contains(blob, []byte(`"type":"Polygon"`)) {
		t.Errorf("should set type to polygon")
	}

	if !bytes.Contains(blob, []byte(`"coordinates":[[[1,1],[2,1],[2,2],[1,2],[1,1]]]`)) {
		t.Errorf("should set type to polygon coords: %v", string(blob))
	}
}

func TestFeatureMarshal(t *testing.T) {
	f := NewFeature(orb.Point{1, 2})
	blob, err := json.Marshal(f)

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestFeatureMarshalValue(t *testing.T) {
	f := NewFeature(orb.Point{1, 2})
	blob, err := json.Marshal(*f)

	if err != nil {
		t.Fatalf("should marshal to json just fine but got %v", err)
	}

	if !bytes.Contains(blob, []byte(`"properties":null`)) {
		t.Errorf("json should set properties to null if there are none")
	}
}

func TestUnmarshalFeature(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]},
	    "properties": {"prop0": "value0"}
	  }`

	f, err := UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if f.Type != "Feature" {
		t.Errorf("should have type of Feature got: %v", f.Type)
	}

	if len(f.Properties) != 1 {
		t.Errorf("should have 1 property but got: %v", f.Properties)
	}
}

func TestMarshalFeatureID(t *testing.T) {
	f := &Feature{
		ID: "asdf",
	}

	data, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)
	}

	if !bytes.Equal(data, []byte(`{"id":"asdf","type":"Feature","geometry":null,"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}

	f.ID = 123
	data, err = f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)

	}

	if !bytes.Equal(data, []byte(`{"id":123,"type":"Feature","geometry":null,"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
}

func TestUnmarshalFeatureID(t *testing.T) {
	rawJSON := `
	  { "type": "Feature",
	    "id": 123,
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]}
	  }`

	f, err := UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature without issue, err %v", err)
	}

	if v, ok := f.ID.(float64); !ok || v != 123 {
		t.Errorf("should parse id as number, got %T %f", f.ID, v)
	}

	rawJSON = `
	  { "type": "Feature",
	    "id": "abcd",
	    "geometry": {"type": "Point", "coordinates": [102.0, 0.5]}
	  }`

	f, err = UnmarshalFeature([]byte(rawJSON))
	if err != nil {
		t.Fatalf("should unmarshal feature without issue, err %v", err)
	}

	if v, ok := f.ID.(string); !ok || v != "abcd" {
		t.Errorf("should parse id as string, got %T %s", f.ID, v)
	}
}

func TestMarshalRing(t *testing.T) {
	ring := orb.Ring{{0, 0}, {1, 1}, {2, 1}, {0, 0}}

	f := NewFeature(ring)
	data, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("should marshal, %v", err)
	}

	if !bytes.Equal(data, []byte(`{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[0,0],[1,1],[2,1],[0,0]]]},"properties":null}`)) {
		t.Errorf("data not correct")
		t.Logf("%v", string(data))
	}
}
