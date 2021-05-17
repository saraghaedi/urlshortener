package migrate

import (
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/migrations/bindata/postgres"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main(cfg config.Database) {
	source := bindata.Resource(postgres.AssetNames(), postgres.Asset)

	if err := database.Migrate(source, cfg.MasterConnStr); err != nil {
		logrus.Fatalf("failed to run database migrations: %s", err.Error())
	}
}

// Register registers migrate command for urlshortener binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg.Database)

			cmd.Println("migrations ran successfully")
		},
	}

	root.AddCommand(cmd)
}
