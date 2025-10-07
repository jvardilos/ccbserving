package serving

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/jvardilos/ccbapi"
	"github.com/jvardilos/ccbserving/pkg/ccbtime"
	"github.com/urfave/cli/v3"
)

func GetServing(ctx context.Context, cli *cli.Command) error {
	n := time.Now()

	c := ccbapi.Credentials{
		Subdomain: cli.String("subdomain"),
		Client:    cli.String("client_id"),
		Secret:    cli.String("secret"),
	}

	t, err := ccbapi.Authorize(&c)
	if err != nil {
		return err
	}

	foundGroup, err := getGroupID(cli.String("group"), t, c)
	if err != nil {
		return err
	}

	members, err := getGroupMembers(foundGroup.ID, t, c)
	if err != nil {
		return err
	}

	start := ccbtime.FormatLastMonthDate(n)
	end := ccbtime.FormatDate(n)

	for _, m := range members {
		fmt.Println("---{start}")
		fmt.Println(m.Individual.Name, m.Individual.ID)
		serving, err := getMemberServing(m.Individual.ID, t, c, start, end)
		if err != nil {
			return err
		}
		fmt.Println(serving)
		fmt.Println("---{end}")
	}

	return nil
}

func getGroupID(group string, t *ccbapi.Token, c ccbapi.Credentials) (ccbGroup, error) {
	body, err := ccbapi.Call("GET", "groups?name="+url.QueryEscape(group), t, &c)
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

	// return top result cuz I'm lazy
	return cg[0], nil
}

// TODO: paging
func getGroupMembers(id int, t *ccbapi.Token, c ccbapi.Credentials) ([]ccbMember, error) {
	path := fmt.Sprintf("groups/%d/members", id)
	body, err := ccbapi.Call("GET", path, t, &c)
	if err != nil {
		return nil, err
	}

	var cm []ccbMember
	if err := json.Unmarshal(body, &cm); err != nil {
		return nil, err
	}

	return cm, nil
}

func getMemberServing(id int, t *ccbapi.Token, c ccbapi.Credentials, start, end string) ([]ccbServing, error) {
	path := fmt.Sprintf("individuals/%d/metrics/serving?start=%s&end=%s", id, start, end)
	fmt.Println("PATH:", path)
	body, err := ccbapi.Call("GET", path, t, &c)
	if err != nil {
		return nil, err
	}

	var cs []ccbServing
	if err := json.Unmarshal(body, &cs); err != nil {
		return nil, err
	}

	return cs, nil
}
