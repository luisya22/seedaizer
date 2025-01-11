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
		tableName := viper.GetString("table-name")
		quantity := viper.GetInt("quantity")
		addChildren := viper.GetBool("add-chilren")
		skipChildrens := viper.GetString("skip-childrens")
		batchSize := viper.GetInt("batch-size")
		outputMode := viper.GetString("output-mode")

		databaseUrl := viper.GetString("database.url")
		if databaseUrl == "" {
			log.Fatal("No database URL found in the configuration.")
		} else {
			fmt.Println("Connecting to database")
		}

		openAiKey := viper.GetString("openaikey")

		config := seeder.Config{
			DbUrl:     databaseUrl,
			FilePath:  filePath,
			OpenAiKey: openAiKey,
		}

		options := seeder.Options{
			TableName:     tableName,
			Quantity:      quantity,
			AddChildren:   addChildren,
			SkipChildrens: skipChildrens,
			BatchSize:     batchSize,
			OutputMode:    outputMode,
		}

		err := seeder.Seed(config, options)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(seedupCmd)

	seedupCmd.PersistentFlags().StringP("table-name", "t", "", "Name of the main table where data will be inserted.")
	seedupCmd.PersistentFlags().IntP("quantity", "q", 1, "Number of records to generate.")
	seedupCmd.PersistentFlags().BoolP("add-children", "c", false, "Include child table data during generation.")
	seedupCmd.PersistentFlags().StringP("skip-childrens", "s", "", "Comma-separated list of child tables to skip.")
	seedupCmd.PersistentFlags().IntP("batch-size", "b", 10, "Number of records to generate per batch.")
	seedupCmd.PersistentFlags().StringP("output-mode", "o", "print", "How to handle output: 'print' to display in terminal, 'json' to save as a file, or 'execute' to apply directly to the database.")

	flagsToBind := []string{"table-name", "quantity", "add-children", "skip-childrens", "batch-size", "output-mode"}
	for _, flag := range flagsToBind {
		err := viper.BindPFlag(flag, seedupCmd.PersistentFlags().Lookup(flag))
		if err != nil {
			log.Fatalf("Error binding %s flag: %v", flag, err)
		}
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
