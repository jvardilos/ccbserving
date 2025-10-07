package serving

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jvardilos/ccbapi"
)

type httpAPI struct{ base string }

func (h httpAPI) Authorize(c *ccbapi.Credentials) (*ccbapi.Token, error) { return &ccbapi.Token{}, nil }
func (h httpAPI) Call(method, path string, _ *ccbapi.Token, _ *ccbapi.Credentials) ([]byte, error) {
	req, _ := http.NewRequest(method, h.base+"/"+path, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

func TestIntegration_Roundtrip(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		_ = json.NewEncoder(w).Encode([]ccbGroup{{ID: 101, Name: name}})
	})
	mux.HandleFunc("/groups/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/members") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode([]ccbMember{{ID: 101, Individual: ccbIndividual{ID: 1001, Name: "Alice"}}})
	})
	mux.HandleFunc("/individuals/", func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/metrics/serving") {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode([]ccbServing{{ID: 1, Count: 2, Start: "2025-09-14"}})
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	api := httpAPI{base: srv.URL}

	g, err := getGroupID("Connections | Home Team", &ccbapi.Token{}, ccbapi.Credentials{}, api)
	if err != nil || g.ID != 101 {
		t.Fatalf("getGroupID: %v %#v", err, g)
	}

	ms, err := getGroupMembers(g.ID, &ccbapi.Token{}, ccbapi.Credentials{}, api)
	if err != nil || len(ms) != 1 {
		t.Fatalf("getGroupMembers: %v %#v", err, ms)
	}

	cs, err := getMemberServing(ms[0].Individual.ID, &ccbapi.Token{}, ccbapi.Credentials{}, "2025-09-01", "2025-10-01", api)
	if err != nil || len(cs) != 1 || cs[0].Count != 2 {
		t.Fatalf("getMemberServing: %v %#v", err, cs)
	}
}
