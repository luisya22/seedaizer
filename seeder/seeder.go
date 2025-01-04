package seeder

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/luisya22/seedaizer/internal/models"
)

type Config struct {
	DbUrl     string
	FilePath  string
	OpenAiKey string
}

func getSchema(filePath string) (*models.Schema, error) {
	// TODO: Load schema from file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var schema models.Schema
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&schema); err != nil {
		return nil, fmt.Errorf("error decoding json file: %w", err)
	}

	return &schema, nil
}

func getTables(schema *models.Schema, query string) []models.Table {
	tables := []models.Table{}

	// TODO: Handle when tables not found
	re := regexp.MustCompile(`\{([^{}]*)\}`)

	matches := re.FindAllStringSubmatch(query, -1)

	for _, match := range matches {
		if len(match) > 1 {
			if table, ok := schema.Tables[match[1]]; ok {
				tables = append(tables, table)
			}
		}
	}

	return tables
}

func findReferencedAndChildTables(schema *models.Schema, tables []models.Table) []models.Table {
	foundTables := make(map[string]struct{})
	var result []models.Table

	for _, rootTable := range tables {
		if _, alreadyFound := foundTables[rootTable.Name]; alreadyFound {
			continue
		}

		foundTables[rootTable.Name] = struct{}{}

		result = append(result, rootTable)

		referencedTables := getReferencedTables(schema, rootTable, foundTables)
		result = append(tables, referencedTables...)

		childTables := getChildTables(schema, rootTable, foundTables)
		result = append(result, childTables...)
	}

	return result
}

func getReferencedTables(schema *models.Schema, rootTable models.Table, foundTables map[string]struct{}) []models.Table {
	var result []models.Table

	for _, foreignKeys := range rootTable.ForeignKeys {
		referencedTable, ok := schema.Tables[foreignKeys.ReferencedTableName]
		if !ok {
			continue
		}

		if _, alreadyFound := foundTables[referencedTable.Name]; alreadyFound {
			continue
		}

		foundTables[referencedTable.Name] = struct{}{}
		result = append(result, referencedTable)

		referencedTables := getReferencedTables(schema, referencedTable, foundTables)
		result = append(result, referencedTables...)
	}

	return result
}

func getChildTables(schema *models.Schema, rootTable models.Table, foundTables map[string]struct{}) []models.Table {
	var result []models.Table

	for childTableName := range rootTable.ChildTables {
		childTable, ok := schema.Tables[childTableName]
		if !ok {
			continue
		}

		if _, alreadyFound := foundTables[childTableName]; alreadyFound {
			continue
		}

		foundTables[childTableName] = struct{}{}
		result = append(result, childTable)

		referencedTables := getReferencedTables(schema, childTable, foundTables)
		result = append(result, referencedTables...)

		childTables := getChildTables(schema, childTable, foundTables)
		result = append(result, childTables...)
	}

	return result
}

func Seed(config Config, query string) error {
	schema, err := getSchema(config.FilePath)
	if err != nil {
		return err
	}

	queryTables := getTables(schema, query)

	tables := findReferencedAndChildTables(schema, queryTables)

	tablesJson, err := json.MarshalIndent(tables, "", "\t")
	if err != nil {
		return err
	}

	userPrompt, err := buildPrompt(query, string(tablesJson))
	if err != nil {
		return err
	}

	openaiService := NewOpenAIService(config.OpenAiKey)

	response, err := openaiService.queryllm(systemPrompt, userPrompt)
	if err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}

func buildPrompt(query string, tablesJson string) (string, error) {

	p := fmt.Sprintf(seederPrompt, query, tablesJson, "")

	return p, nil
}

/**
TODO:
	(DONE) queryllm - will query llm receiving a string
	(DONE) get tables - return a list of table objects to add to the prompt
	(DONE) identify tables - get table names from the prompt and the childrens
	(DONE) buildprompt - build prompt with tables and user instructions
	seed - main interface it will call everything and either print results, save to .sql, or execute directly
		It'll need user_query, result_type, filepath
**/
