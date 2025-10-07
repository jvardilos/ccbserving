package serving

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/jvardilos/ccbapi"
)

type fakeAPI struct {
	call map[string][]byte
	err  map[string]error
}

func (f *fakeAPI) Authorize(c *ccbapi.Credentials) (*ccbapi.Token, error) {
	return &ccbapi.Token{}, nil
}
func (f *fakeAPI) Call(method, path string, t *ccbapi.Token, c *ccbapi.Credentials) ([]byte, error) {
	key := method + " " + path
	if err := f.err[key]; err != nil {
		return nil, err
	}
	return f.call[key], nil
}

func mustJSON(v any) []byte { b, _ := json.Marshal(v); return b }

func TestGetGroupID_Success(t *testing.T) {
	api := &fakeAPI{
		call: map[string][]byte{
			"GET groups?name=" + url.QueryEscape("Connections | Home Team"): mustJSON([]ccbGroup{{ID: 42, Name: "Connections | Home Team"}}),
		},
	}
	got, err := getGroupID("Connections | Home Team", &ccbapi.Token{}, ccbapi.Credentials{}, api)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if got.ID != 42 {
		t.Fatalf("want 42 got %d", got.ID)
	}
}

func TestGetGroupMembers_Success(t *testing.T) {
	api := &fakeAPI{
		call: map[string][]byte{
			"GET groups/99/members": mustJSON([]ccbMember{
				{ID: 99, Individual: ccbIndividual{ID: 1, Name: "Alice"}},
				{ID: 99, Individual: ccbIndividual{ID: 2, Name: "Bob"}},
			}),
		},
	}
	ms, err := getGroupMembers(99, &ccbapi.Token{}, ccbapi.Credentials{}, api)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(ms) != 2 || ms[0].Individual.Name != "Alice" {
		t.Fatalf("bad: %#v", ms)
	}
}

func TestGetMemberServing_Success(t *testing.T) {
	api := &fakeAPI{
		call: map[string][]byte{
			"GET individuals/7/metrics/serving?start=2025-09-01&end=2025-10-01": mustJSON([]ccbServing{
				{ID: 1, Count: 3, Start: "2025-09-07"},
			}),
		},
	}
	cs, err := getMemberServing(7, &ccbapi.Token{}, ccbapi.Credentials{}, "2025-09-01", "2025-10-01", api)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if len(cs) != 1 || cs[0].Count != 3 {
		t.Fatalf("bad: %#v", cs)
	}
}

// func TestServing_Orchestrates(t *testing.T) {
// 	api := &fakeAPI{
// 		call: map[string][]byte{
// 			"GET groups?name=My+Group": mustJSON([]ccbGroup{{ID: 7, Name: "My Group"}}),
// 			"GET groups/7/members":     mustJSON([]ccbMember{{ID: 7, Individual: ccbIndividual{ID: 10, Name: "Alice"}}}),
// 			"GET individuals/10/metrics/serving?start=2025-09-01&end=2025-10-01": mustJSON([]ccbServing{{ID: 1, Count: 2, Start: "2025-09-15"}}),
// 		},
// 	}

// 	cmd := &cli.Command{}
// 	ctx := cli.NewContext(cmd, &cli.FlagValues{
// 		Values: map[string]string{
// 			"subdomain": "demo",
// 			"client_id": "abc",
// 			"secret":    "xyz",
// 			"group":     "My Group",
// 		},
// 	}, nil)

// 	if err := getServing(context.Background(), ctx, Deps{API: api}); err != nil {
// 		t.Fatal(err)
// 	}
// }
