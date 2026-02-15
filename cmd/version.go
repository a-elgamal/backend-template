package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"alielgamal.com/myservice/internal"
	internalDB "alielgamal.com/myservice/internal/db"
)

func versionCmd(db *internalDB.SQLDB) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of the service",
		Run: func(cmd *cobra.Command, _ []string) {
			out := cmd.OutOrStdout()
			dbVersion, dirty, err := internalDB.Version(db.DB)
			fmt.Fprintf(out, "DB Version: %v\n", dbVersion)
			fmt.Fprintf(out, "DB Dirty: %v\n", dirty)
			fmt.Fprintf(out, "DB Error: %v\n", err)
			fmt.Fprintf(out, "App Version: %v\n", internal.Version)
			fmt.Fprintf(out, "Git Tag: %v\n", internal.GitTag)
			fmt.Fprintf(out, "Git Commit: %v\n", internal.GitCommit)
			fmt.Fprintf(out, "Build Date: %v\n", internal.BuildDate)
		},
	}
}
