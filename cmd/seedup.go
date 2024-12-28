/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/luisya22/seedaizer/seeder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// seedupCmd represents the seedup command
var seedupCmd = &cobra.Command{
	Use:   "seedup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("seedup called")

		filePath := "db.json"
		query := viper.GetString("query")

		databaseUrl := viper.GetString("database.url")
		if databaseUrl == "" {
			log.Fatal("No database URL found in the configuration.")
		} else {
			fmt.Println("Connecting to database")
		}

		config := seeder.Config{
			DbUrl:    databaseUrl,
			FilePath: filePath,
		}

		err := seeder.Seed(config, query)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(seedupCmd)

	seedupCmd.PersistentFlags().String("query", "", "Query to generate necessary sql commands")

	err := viper.BindPFlag("query", seedupCmd.PersistentFlags().Lookup("query"))
	if err != nil {
		log.Fatal("error binding query flag: %w", err)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
