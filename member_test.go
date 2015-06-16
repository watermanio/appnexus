package AppNexus

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMemberService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/member/1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"response":
            {"status":"OK",
            "count": 1,
            "start_element": 0,
            "num_elements": 100,
            "member": {
                "ID": 1,
                "Name": "Test Member",
                "State": "active",
                "default_currency": "GBP"
            }}}`)
	})

	actual, err := client.Members.Get(1)
	if err != nil {
		t.Errorf("Members.Get returned error: %v", err)
	}

	expected := Member{
		ID:             1,
		Name:           "Test Member",
		State:          "active",
		DefaltCurrency: "GBP",
	}

	if actual.ID != expected.ID || actual.Name != expected.Name {
		t.Errorf("Members.Get returned %+v, expected %+v", actual, expected)
	}
}
