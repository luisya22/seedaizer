/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/luisya22/seedaizer/dbscanner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scandbCmd represents the scandb command
var scandbCmd = &cobra.Command{
	Use:   "scandb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scandb called")

		databaseUrl := viper.GetString("database.url")
		if databaseUrl == "" {
			log.Fatal("No database URL found in the configuration.")
		} else {
			fmt.Println("Connecting to database")
		}

		config := dbscanner.Config{DbUrl: databaseUrl}

		err := dbscanner.ScanToJson(config)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(scandbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scandbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scandbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
