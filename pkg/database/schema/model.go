package schema

import "database/sql"

var lastTargetModel *Model

// TableModel locates table model
func TableModel(tableName string) *Table {

	for _, table := range lastTargetModel.Tables {
		if table.Name == tableName {
			return &table
		}
	}

	return nil

}

// IsSoftDeleted returns the soft deleted status
func IsSoftDeleted(tableName string) bool {

	table := TableModel(tableName)

	for _, col := range table.Columns {
		if col.Name == "is_deleted" {
			return true
		}
	}

	return false

}

/*
Model models a database model.
*/
type Model struct {
	Tables   []Table            `json:"tables"`
	Queries  []ConditionalQuery `json:"queries"`
	Triggers []Trigger          `json:"triggers"`
}

// Trigger is a callback function performed when a condition is matched.
type Trigger struct {
	Name      string `json:"name"`
	Table     string `json:"table"`
	Statement string `json:"statement"`
	Timing    string `json:"timing"`
	Event     string `json:"event"`
}

/*
Table models a database table.
*/
type Table struct {
	Name       string   `json:"name"`
	CloneNames string   `json:"cloneNames"`
	Columns    []Column `json:"columns"`
	Indices    []Index  `json:"indices"`
}

/*
Column models a database column.
*/
type Column struct {
	Name         string     `json:"name"`
	DataType     string     `json:"type"`
	Size         int        `json:"size"`
	Decimal      int        `json:"decimal"`
	Nullable     bool       `json:"nullable"`
	PrimaryKey   bool       `json:"primaryKey"`
	DefaultValue string     `json:"defaultValue"`
	ForeignKey   ForeignKey `json:"foreignKey"`
}

/*
ForeignKey models a foreign key.
*/
type ForeignKey struct {
	Name      string `json:"name"`
	TableName string `json:"tableName"`
}

/*
Index models a database index.
*/
type Index struct {
	Name        string   `json:"name"`
	Unique      bool     `json:"unique"`
	ColumnNames []string `json:"columnNames"`
}

/*
ConditionalQuery models a free form query that will be run if the condition returns
any rows.
*/
type ConditionalQuery struct {
	Condition       string `json:"condition"`
	TableExists     string `json:"tableExists"`
	ColumnExists    string `json:"columnExists"`
	Query           string `json:"query"`
	PostMigrateOnly bool   `json:"postMigrateOnly"`
}

/*
ForeignKeyHaver allows model objects to declare that they have foreign keys.
*/
type ForeignKeyHaver interface {
	ForeignKeys() []ForeignKey
	ForeignKeyColumns() []Column
}

/*
ColumnHaver allows model objects to declare that they have columns.
*/
type ColumnHaver interface {
	ColumnMap() map[string]Column
}

/*
IndexHaver allows model objects to declare that they have indices.
*/
type IndexHaver interface {
	IndexMap() map[string]Index
}

/*
Dialecter defines the interface for RDBMS specific dialects.
*/
type Dialecter interface {
	ModifyColumn(Change) string
	ColumnDefinition(Column) string
	PrimaryKeyDefinition(Table) string
	ForeignKeyDefinition(Column) string
	IndexDefinition(Table, Index) string
	ReadCurrentModel(ds *sql.DB) (*Model, error)
	CreateTriggerQuery(Trigger) string
	DropTriggerQuery(Trigger) string
}

//HasTable returns true if the model contains a reference to the table.
func (model *Model) HasTable(tableName string) bool {

	for _, table := range model.Tables {
		if table.Name == tableName {
			return true
		}
	}

	return false

}

/*
IndexMap returns all indices on the table keyed by name.
*/
func (table *Table) IndexMap() map[string]Index {

	results := make(map[string]Index, 0)

	for _, idx := range table.Indices {
		results[idx.Name] = idx
	}

	return results

}

/*
ColumnMap returns all columns on the table keyed by name.
*/
func (table *Table) ColumnMap() map[string]Column {
	results := make(map[string]Column, 0)
	for _, col := range table.Columns {
		results[col.Name] = col
	}
	return results
}

/*
ForeignKeys returns all foreign keys on the table.
*/
func (table *Table) ForeignKeys() []ForeignKey {
	results := make([]ForeignKey, 0)

	for _, col := range table.Columns {
		if col.ForeignKey.Name != "" {
			results = append(results, col.ForeignKey)
		}
	}

	return results
}

/*
ForeignKeyColumns returns all columns on the table with foreign keys.
*/
func (table *Table) ForeignKeyColumns() []Column {
	results := make([]Column, 0)

	for _, col := range table.Columns {
		if col.ForeignKey.Name != "" {
			results = append(results, col)
		}
	}

	return results
}
