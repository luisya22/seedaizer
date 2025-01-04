package dbscanner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/luisya22/seedaizer/internal/models"
)

type Config struct {
	DbUrl string
}

func connectDb(config Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", config.DbUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}

	return db, nil
}

// scan scans the db and returns a *models.schema. It fetches tables, columns and foreign key constraints.
func scan(config Config) (*models.Schema, error) {
	fmt.Println("Scanning")
	db, err := connectDb(config)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	splitUrl := strings.Split(config.DbUrl, "/")
	dbName := ""

	if len(splitUrl) > 1 {
		dbName = splitUrl[1]
	}

	var tablesName []string
	schema := models.Schema{
		Name:   "db",
		Tables: map[string]models.Table{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = db.Select(&tablesName, "SHOW TABLES")
	rows, err := db.QueryxContext(ctx, "SHOW TABLES")
	if err != nil {
		return nil, fmt.Errorf("error fetching tables: %w", err)
	}

	for rows.Next() {
		var table models.Table

		err := rows.Scan(&table.Name)
		if err != nil {
			return nil, fmt.Errorf("error scanning tables result: %w", err)
		}
		// FETCH COLUMNS

		query := fmt.Sprintf("DESCRIBE `%s`", table.Name)

		colRows, err := db.QueryxContext(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("error fetching columns: %w", err)
		}

		table.Columns = map[string]models.Column{}

		for colRows.Next() {
			var col models.Column

			err := colRows.StructScan(&col)
			if err != nil {
				return nil, fmt.Errorf("error scanning columns: %w", err)
			}

			table.Columns[col.Field] = col
		}

		table.ForeignKeys = make(map[string]models.ForeignKey)
		table.ChildTables = make(map[string][]models.ForeignKey)

		schema.Tables[table.Name] = table
	}

	// FETCH FOREIGN KEYS
	query := `SELECT table_name, column_name, referenced_table_name, referenced_column_name FROM information_schema.key_column_usage WHERE table_schema = ? and referenced_table_name is not null`

	fkRows, err := db.QueryxContext(ctx, query, dbName)
	if err != nil {
		return nil, fmt.Errorf("error fetching foreignkeys: %w", err)
	}

	for fkRows.Next() {
		var fk models.ForeignKey

		err := fkRows.StructScan(&fk)
		if err != nil {
			return nil, fmt.Errorf("error scanning foreign keys: %w", err)
		}

		if table, ok := schema.Tables[fk.TableName]; ok {
			table.ForeignKeys[fk.ColumnName] = fk
			schema.Tables[fk.TableName] = table
		}

		if table, ok := schema.Tables[fk.ReferencedTableName]; ok {
			table.ChildTables[fk.TableName] = append(table.ChildTables[fk.TableName], fk)
			schema.Tables[fk.ReferencedTableName] = table
		}
	}

	return &schema, nil
}

// ScanToJson uses scan to get the db schema and saves it directly to a json file
func ScanToJson(config Config) error {
	s, err := scan(config)
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(*s, "", "  ")
	if err != nil {
		return fmt.Errorf("error mapping json: %v", err)
	}

	f, err := os.Create("db.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	_, err = f.WriteString(string(jsonData))
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
