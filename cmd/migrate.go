package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	internalDB "alielgamal.com/myservice/internal/db"
)

func migrateCmd(db *internalDB.SQLDB) *cobra.Command {
	result := &cobra.Command{
		Use:   "migrate [version]",
		Short: "Run migrations on the database",
		Long:  "Run migrations on the database by passing a target DB version OR latest version if no version is specified",
		Args:  cobra.MaximumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Starting Migrate Command")

			force, err := cmd.Flags().GetBool("force")
			if err != nil {
				log.Fatalf("Invalid force flat: %v\n", err)
			}
			if force && len(args) == 0 {
				log.Fatal("Target version must be specified with 'force' flag is specified")
			}

			var version uint
			switch len(args) {
			case 0:
				log.Println("Running all pending DB migrations...")
				version, err = internalDB.UpgradeDB(db.DB)
			default:
				var targetVersion int
				targetVersion, err = strconv.Atoi(args[0])
				if err == nil {
					if force {
						log.Printf("Migrating DB to version %v...\n", targetVersion)
						version = uint(targetVersion)
						err = internalDB.Force(db.DB, targetVersion)
					} else {
						log.Printf("Migrating DB to version %v...\n", targetVersion)
						version, err = internalDB.MigrateDBTo(db.DB, uint(targetVersion))
					}
				}
			}

			if err != nil {
				log.Fatalf("Error while running migrations: %v\n", err)
			}
			log.Printf("The DB is now at version %v\n", version)
		},
	}

	result.Flags().BoolP("force", "f", false, "Force DB version without running migrations")
	return result
}
