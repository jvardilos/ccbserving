package serving

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func GetServing(ctx context.Context, cmd *cli.Command) error {
	if err := getServing(cmd, realAPI{}); err != nil {
		return fmt.Errorf("getting serving: %w", err)
	}
	return nil
}
