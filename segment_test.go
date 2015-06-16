package AppNexus

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSegmentService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/segment/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"response":
            {"status":"OK",
            "count": 1,
            "start_element": 0,
            "num_elements": 100,
            "segment": {
                "id": 1,
                "code": "123",
                "state": "active",
                "short_name": "test_segment"
            }}}`)
	})

	actual, err := client.Segments.Get(1, 1)
	if err != nil {
		t.Errorf("Segment.Get returned error: %v", err)
	}

	expected := Segment{
		ID:        1,
		ShortName: "test_segment",
		State:     "active",
		Code:      "123",
	}

	if actual.ID != expected.ID || actual.ShortName != expected.ShortName {
		t.Errorf("Segment.Get returned %+v, expected %+v", actual, expected)
	}
}

func TestSegmentService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/segment/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"response":
            {"status":"OK",
            "count": 3,
            "start_element": 2,
            "num_elements": 2,
            "segments": [{
                "id": 3,
                "code": "1234",
                "state": "active",
                "short_name": "test_segment_3"
            }]}}`)
	})

	actual, _, err := client.Segments.List(1, &ListOptions{StartElement: 2, NumElements: 2})
	if err != nil {
		t.Errorf("Segments.List returned error: %v", err)
	}

	expected := make([]Segment, 1)
	expected[0] = Segment{
		ID:        3,
		ShortName: "test_segment_3",
		State:     "active",
		Code:      "1234",
	}

	if actual[0].ID != expected[0].ID || actual[0].ShortName != expected[0].ShortName {
		t.Errorf("Segment.Get returned %+v, \nexpected %+v", actual, expected)
	}
}

func TestSegmentService_Add(t *testing.T) {
	setup()
	defer teardown()

	data := Segment{
		ShortName: "Hello Seggy",
		Code:      "seg1",
	}

	mux.HandleFunc("/segment/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"response": {"status":"OK", "id": 4 }}`)
	})

	actual, err := client.Segments.Add(1, &data)
	if err != nil {
		t.Errorf("Segments.Add returned error: %v", err)
	}

	expected := Response{}
	expected.Obj.ID = 4
	expected.Obj.Status = "OK"

	if expected.Obj.ID != actual.Obj.ID || actual.Obj.Status != "OK" {
		t.Errorf("Segment.Get returned %+v, \nexpected %+v", actual, expected)
	}
}

func TestSegmentService_Update(t *testing.T) {
	setup()
	defer teardown()

	data := Segment{
		ID:        4,
		ShortName: "Hello Seggy 2",
		Code:      "seg1",
	}

	mux.HandleFunc("/segment/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"response": {"status":"OK" }}`)
	})

	actual, err := client.Segments.Update(1, data)
	if err != nil {
		t.Errorf("Segments.Update returned error: %v", err)
	}

	expected := Response{}
	expected.Obj.Status = "OK"

	if actual.Obj.Status != "OK" {
		t.Errorf("Segment.Update returned %+v, \nexpected %+v", actual, expected)
	}
}

func TestSegmentService_Delete(t *testing.T) {
	setup()
	defer teardown()

	data := Segment{
		ID:        4,
		ShortName: "Hello Seggy 2",
		Code:      "seg1",
	}

	mux.HandleFunc("/segment/1", func(w http.ResponseWriter, r *http.Request) {

	})

	err := client.Segments.Delete(1, data)
	if err != nil {
		t.Errorf("Segments.Delete returned error: %v", err)
	}
}
