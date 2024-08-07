package cmd

import (
	"github.com/nutcas3/my-ticko/db/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/nutcas3/my-ticko/db"
)

var createAdminUserCmd = &cobra.Command{
	Use:   "create-admin-user",
	Short: "Create admin user",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")

		if username == "" {
			return errors.New("username cannot be empty")
		}

		logger, err := getLogger()
		if err != nil {
			return err
		}

		dbConfig, err := db.InitConfig()
		if err != nil {
			return err
		}

		db, err := db.New(dbConfig, logger)
		if err != nil {
			return err
		}
		defer db.Close()

		// Create the user
		_, err = db.CreateUser(username, model.Admin)
		if err != nil {
			return err
		}
		logger.Infof("Admin user created")

		return nil
	},
}

func init() {
	createAdminUserCmd.Flags().StringP("username", "u", "", "Username")
	_ = createAdminUserCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(createAdminUserCmd)
}
