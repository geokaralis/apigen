package create

import (
	"context"

	"github.com/geokaralis/apigen/internal/client"
	"github.com/geokaralis/apigen/internal/config"
	"github.com/urfave/cli/v3"
)

func Client(ctx context.Context, cmd *cli.Command) error {
	config := config.Config{
		Schema: cmd.StringArgs("schema")[0],
		Output: cmd.String("output"),
	}

	if err := config.Validate(); err != nil {
		return err
	}

	return client.Generate(ctx, config)
}
