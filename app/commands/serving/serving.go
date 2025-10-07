package serving

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/jvardilos/ccbapi"
	"github.com/jvardilos/ccbserving/pkg/ccbtime"
	"github.com/urfave/cli/v3"
)

func getServing(cmd *cli.Command, api ccb) error {
	if api == nil {
		api = realAPI{}
	}

	n := time.Now()

	c := ccbapi.Credentials{
		Subdomain: cmd.String("subdomain"),
		Client:    cmd.String("client_id"),
		Secret:    cmd.String("secret"),
	}

	t, err := api.Authorize(&c)
	if err != nil {
		return err
	}

	foundGroup, err := getGroupID(cmd.String("group"), t, c, api)
	if err != nil {
		return err
	}

	members, err := getGroupMembers(foundGroup.ID, t, c, api)
	if err != nil {
		return err
	}

	start := ccbtime.FormatLastMonthDate(n)
	end := ccbtime.FormatDate(n)

	for _, m := range members {
		fmt.Println("---{start}")
		fmt.Println(m.Individual.Name, m.Individual.ID)
		serving, err := getMemberServing(m.Individual.ID, t, c, start, end, api)
		if err != nil {
			return err
		}
		fmt.Println(serving)
		fmt.Println("---{end}")
	}
	return nil
}

func getGroupID(group string, t *ccbapi.Token, c ccbapi.Credentials, api ccb) (ccbGroup, error) {
	body, err := api.Call("GET", "groups?name="+url.QueryEscape(group), t, &c)
	if err != nil {
		return ccbGroup{}, err
	}

	var cg []ccbGroup
	if err := json.Unmarshal(body, &cg); err != nil {
		return ccbGroup{}, err
	}
	if len(cg) == 0 {
		return ccbGroup{}, fmt.Errorf("group not found")
	}
	return cg[0], nil
}

// TODO: paging
func getGroupMembers(id int, t *ccbapi.Token, c ccbapi.Credentials, api ccb) ([]ccbMember, error) {
	body, err := api.Call("GET", fmt.Sprintf("groups/%d/members", id), t, &c)
	if err != nil {
		return nil, err
	}
	var cm []ccbMember
	if err := json.Unmarshal(body, &cm); err != nil {
		return nil, err
	}
	return cm, nil
}

func getMemberServing(id int, t *ccbapi.Token, c ccbapi.Credentials, start, end string, api ccb) ([]ccbServing, error) {
	path := fmt.Sprintf("individuals/%d/metrics/serving?start=%s&end=%s", id, start, end)
	body, err := api.Call("GET", path, t, &c)
	if err != nil {
		return nil, err
	}
	var cs []ccbServing
	if err := json.Unmarshal(body, &cs); err != nil {
		return nil, err
	}
	return cs, nil
}
