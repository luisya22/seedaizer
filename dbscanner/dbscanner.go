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

		// FETCH FOREIGN KEYS
		query = `SELECT table_name, column_name, referenced_table_name, referenced_column_name FROM information_schema.key_column_usage WHERE referenced_table_name = ? and table_schema = ?`

		fkRows, err := db.QueryxContext(ctx, query, table.Name, dbName)
		if err != nil {
			return nil, fmt.Errorf("error fetching foreignkeys: %w", err)
		}

		table.ForeignKeys = map[string]models.ForeignKey{}

		for fkRows.Next() {
			var fk models.ForeignKey

			err := fkRows.StructScan(&fk)
			if err != nil {
				return nil, fmt.Errorf("error scanning columns: %w", err)
			}

			table.ForeignKeys[fk.ColumnName] = fk
		}

		schema.Tables[table.Name] = table
	}

	return &schema, nil
}

func ScanToJson(config Config) error {
	s, err := scan(config)
	if err != nil {
		return err
	}

	fmt.Println(s)

	jsonData, err := json.MarshalIndent(*s, "", "  ")
	if err != nil {
		return fmt.Errorf("error mapping json: %v", err)
	}

	fmt.Println(string(jsonData))

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
