package models

import "database/sql"

type Schema struct {
	Name   string           `json:"name"`
	Tables map[string]Table `json:"tables"`
}

// TODO: Convert table and column slices to maps so I can search quickly by name

type Table struct {
	Name        string                  `json:"name"`
	Columns     map[string]Column       `json:"columns"`
	ForeignKeys map[string]ForeignKey   `json:"foreignKeys"`
	ChildTables map[string][]ForeignKey `json:"childTables"`
}

type Column struct {
	Field   string         `db:"Field" json:"field"`
	Type    string         `db:"Type" json:"type"`
	Null    string         `db:"Null" json:"null"`
	Key     string         `db:"Key" json:"key"`
	Default sql.NullString `db:"Default" json:"default"`
	Extra   string         `db:"Extra" json:"extra"`
}

type ForeignKey struct {
	TableName            string `db:"TABLE_NAME"`
	ColumnName           string `db:"COLUMN_NAME" json:"columnName"`
	ConstraintName       string `db:"CONSTRAINT_NAME" json:"constraintName"`
	ReferencedTableName  string `db:"REFERENCED_TABLE_NAME" json:"referencedTableName"`
	ReferencedColumnName string `db:"REFERENCED_COLUMN_NAME" json:"referencedColumnName"`
}
