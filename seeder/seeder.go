package seeder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/luisya22/seedaizer/internal/models"
)

type Config struct {
	DbUrl     string
	FilePath  string
	OpenAiKey string
}

type Options struct {
	TableName     string
	Quantity      int
	AddChildren   bool
	SkipChildrens string
	BatchSize     int
	OutputMode    string
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

func Seed(config Config, options Options) error {
	schema, err := getSchema(config.FilePath)
	if err != nil {
		return err
	}

	// queryTables := getTables(schema, options.TableName)

	// tables := findReferencedAndChildTables(schema, queryTables)

	// tablesJson, err := json.MarshalIndent(tables, "", "\t")
	// if err != nil {
	// 	return err
	// }

	// userPrompt, err := buildPrompt(options.TableName, string(tablesJson))
	// if err != nil {
	// 	return err
	// }
	//
	// openaiService := NewOpenAIService(config.OpenAiKey)

	fmt.Println("querying llm")

	_, err = generateQueries(config, schema, options)
	if err != nil {
		return err
	}

	// response, err := openaiService.queryllm(systemPrompt, userPrompt)
	// if err != nil {
	// 	return err
	// }
	//
	// err = output(config, outputMode, response)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func generateQueries(config Config, schema *models.Schema, options Options) ([]string, error) {
	queries := []string{}

	table, ok := schema.Tables[options.TableName]
	if !ok {
		return nil, fmt.Errorf("table not found: %s", options.TableName)
	}

	query, columns := constructTableQuery(table)

	columnsString, err := json.MarshalIndent(columns, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("error marshaling columns: %w", err)
	}

	promptData := fmt.Sprintf(seederPrompt2, options.TableName, string(columnsString), strconv.Itoa(options.BatchSize))

	openaiService := NewOpenAIService(config.OpenAiKey)

	// Generate n data for this columns

	// TODO: Generate data for a batch
	// TODO: add to queries using the query
	// TODO: Do same for childs
	// TODO: How to generate better data with llm

	fmt.Println("Generating queries", options.Quantity%options.BatchSize, promptData)

	for x := 0; x <= options.Quantity/options.BatchSize; x++ {
		response, err := openaiService.queryllm(systempPrompt2, promptData)
		if err != nil {
			return nil, fmt.Errorf("error generating data: %w", err)
		}

		data := make([]map[string]any, options.BatchSize)
		err = json.Unmarshal([]byte(response), &data)
		if err != nil {
			return nil, fmt.Errorf("error encoding llm data: %w", err)
		}

		for _, d := range data {

			values := []string{}

			for _, c := range columns {
				values = append(values, fmt.Sprintf("%v", d[c.Field]))
			}

			fmt.Printf(query, strings.Join(values, ", "))
			fmt.Println("")
		}
		fmt.Println("This is x", x, options.BatchSize, options.Quantity%options.BatchSize, options.BatchSize%options.Quantity)

	}

	fmt.Println("Outside loop")

	// for x := 0; x < amount; x++ {
	// 	queries = append(queries, fmt.Sprintf(query))
	// }

	return queries, nil
}

func constructTableQuery(table models.Table) (string, []models.Column) {

	insertColumns, columns := getInsertColumns(table)

	query := fmt.Sprintf("INSERT INTO %s (%s)", table.Name, insertColumns)

	query += " VALUES (%s)"

	return query, columns
}

func getInsertColumns(table models.Table) (string, []models.Column) {
	var columnsNames []string
	var columns []models.Column

	for _, c := range table.Columns {
		if !strings.Contains(c.Extra, "auto_increment") || !strings.Contains(c.Extra, "DEFAULT_GENERATED") {
			columnsNames = append(columnsNames, c.Field)
			columns = append(columns, c)
		}
	}

	return strings.Join(columnsNames, ", "), columns
}

func output(config Config, outputMode string, response string) error {
	var queries []string

	err := json.Unmarshal([]byte(response), &queries)
	if err != nil {
		return fmt.Errorf("error unmarshaling queries: %w", err)
	}

	switch outputMode {
	case "print":
		fmt.Println(queries)
	case "sql":
		err := saveToSql(queries)
		if err != nil {
			return err
		}
	case "execute":
		err := executeQuery(config, queries)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("error output method not valid")
	}

	return nil
}

func saveToSql(queries []string) error {
	f, err := os.Create("seedaizer.sql")
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}

	for _, q := range queries {
		_, err := f.WriteString(q + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	log.Println("queries written to seedaizer.sql file")

	return nil
}

func executeQuery(config Config, queries []string) error {
	fmt.Println("connecting to db")

	db, err := sqlx.Connect("mysql", config.DbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Duration(len(queries))*time.Second)
	defer cancel()

	for _, q := range queries {
		_, err := db.ExecContext(ctx, q)
		if err != nil {
			return fmt.Errorf("error failed to execute query: (%v); %w", q, err)
		}
	}

	fmt.Println("queries executed")

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

	Thinking if it could be a good idea for me to create the queries and the LLM will only generate the data and return as json
	Instead of querying it should be something like: seedup -t users -n 10
	Children will be added by default but you can avoid children inserts with: -c false
	Or skip just some childrens with: -skip-childrens=reviews,orders
	Let the user decide batch size

	Later with something like BubbleTea this can be transform to selection where you can checkmark the childrens you want to add or skip
**/

/**
Handle duplicates:
	duplicate_handling:
		mode: "regenerate", # options: "skip", "update", "regenerate"
**/
