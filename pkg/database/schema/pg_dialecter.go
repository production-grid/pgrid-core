package schema

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	//used to bring in the postgres driver
	_ "github.com/lib/pq"
)

// PostgresDialect handles postgress specific DDL quirks.
type PostgresDialect struct {
}

func (psql *PostgresDialect) isTypeChange(change Change) bool {

	if change.Column.DataType != psql.translateDDLDataType(change.OldColumn) {
		return true
	}
	if change.Column.Size != change.OldColumn.Size {
		return true
	}
	if change.Column.Decimal != change.OldColumn.Decimal {
		return true
	}
	return false

}

//ModifyColumn assembles a modify column command
func (psql *PostgresDialect) ModifyColumn(change Change) string {

	sql := "ALTER TABLE "
	sql += change.Table.Name
	sql += " ALTER COLUMN "
	sql += change.Column.Name
	if psql.isTypeChange(change) {
		sql += " TYPE "
		sql += psql.translateDDLDataType(change.Column)
		if change.Column.Size > 0 {
			sql += "("
			sql += strconv.Itoa(change.Column.Size)
			if change.Column.Decimal > 0 {
				sql += ","
				sql += strconv.Itoa(change.Column.Decimal)
			}
			sql += ")"
		}
	}
	if change.Column.Nullable != change.OldColumn.Nullable {
		if change.Column.Nullable {
			sql += " DROP NOT NULL"
		} else {
			sql += " SET NOT NULL"
		}
	}

	return sql

}

func (psql *PostgresDialect) translateDDLDataType(col Column) string {

	if col.DataType == "TIMESTAMP" {
		return "TIMESTAMP WITH TIME ZONE"
	}
	return col.DataType

}

func (psql *PostgresDialect) columnTypeDefinition(col Column) string {

	ddl := psql.translateDDLDataType(col)

	/*
		We permit data types to be declared with sizes.
		For example DECIMAL (10, 7).  This check bypasses
		the default size stuff if needed.
	*/
	if !strings.Contains(col.DataType, "(") {
		if col.Size > 0 {
			ddl += "("
			ddl += strconv.Itoa(col.Size)
			if col.Decimal > 0 {
				ddl += ","
				ddl += strconv.Itoa(col.Decimal)
			}
			ddl += ") "
		}
	}
	if !col.Nullable {
		ddl += " NOT NULL"
	}
	if len(col.DefaultValue) > 0 {
		ddl += " DEFAULT "
		if col.DataType != "VARCHAR" {
			ddl += col.DefaultValue
		} else {
			ddl += "'" + col.DefaultValue + "'"
		}
	}

	return ddl
}

/*
ColumnDefinition assembles a column definition DDL fragment.
*/
func (psql *PostgresDialect) ColumnDefinition(col Column) string {

	//TODO convert to string templates if practical

	ddl := col.Name
	ddl += " "
	ddl += psql.columnTypeDefinition(col)

	return ddl
}

/*
PrimaryKeyDefinition assembles the DDL query to create a new foreign key.
*/
func (psql *PostgresDialect) PrimaryKeyDefinition(table Table) string {

	pkCols := make([]string, 0)

	for _, col := range table.Columns {
		if col.PrimaryKey {
			pkCols = append(pkCols, col.Name)
		}
	}

	ddl := ""
	if len(pkCols) > 0 {
		ddl += "PRIMARY KEY ("
		ddl += strings.Join(pkCols, ",")
		ddl += ")"
	}

	return ddl

}

/*
ForeignKeyDefinition assembles the DDL query to create a new foreign key.
*/
func (psql *PostgresDialect) ForeignKeyDefinition(col Column) string {

	ddl := "CONSTRAINT "
	ddl += col.ForeignKey.Name
	ddl += " FOREIGN KEY ("
	ddl += col.Name
	ddl += ") REFERENCES "
	ddl += col.ForeignKey.TableName
	ddl += "(id)"
	return ddl

}

/*
IndexDefinition assembles the DDL query to create a new index.
*/
func (psql *PostgresDialect) IndexDefinition(table Table, index Index) string {

	ddl := ""
	if index.Unique {
		ddl += "UNIQUE "
	}
	ddl += "INDEX "
	ddl += index.Name
	ddl += " ON "
	ddl += table.Name
	ddl += "("
	ddl += strings.Join(index.ColumnNames, ",")
	ddl += ")"
	return ddl

}

/*
ReadCurrentModel reads the full database model for the given db.
*/
func (psql *PostgresDialect) ReadCurrentModel(db *sql.DB) (*Model, error) {

	var model Model

	tables, err := psql.ReadTableList(db)
	if err != nil {
		return nil, err
	}

	tableModels := make([]Table, 0)

	for _, table := range tables {
		tableModel, err := psql.ReadTableModel(db, table)
		if err != nil {
			return nil, err
		}
		tableModels = append(tableModels, tableModel)
	}

	model.Tables = tableModels

	/*
		model.Triggers, err = mysql.ReadTriggers(db)
		if err != nil {
			return nil, err
		}
	*/

	return &model, nil
}

func (psql *PostgresDialect) translateDataType(psqlType string) string {

	switch psqlType {
	case "character":
		return "CHAR"
	case "character varying":
		return "VARCHAR"
	case "timestamp", "timestamp with time zone":
		return "TIMESTAMP"
	case "numeric":
		return "DECIMAL"
	}

	return psqlType

}

/*
ReadTableModel reads the full table model for a given table name, which means columns, keys, etc.
*/
func (psql *PostgresDialect) ReadTableModel(db *sql.DB, tableName string) (Table, error) {

	result := Table{Name: tableName}

	pks, err := getPrimaryKeys(db, &result)

	if err != nil {
		return result, err
	}

	fks, err := readForeignKeys(db, result)

	if err != nil {
		return result, err
	}

	rows, err := db.Query("SELECT column_name, data_type, is_nullable, column_default, character_maximum_length, numeric_precision, numeric_scale FROM information_schema.columns where TABLE_NAME = $1", tableName)
	defer rows.Close()
	if err != nil {
		return result, err
	}
	cols := make([]Column, 0)
	for rows.Next() {
		var colName, datatype, nullable, defaultValue *string
		var charLen, numPrec, numScale *int
		err = rows.Scan(&colName, &datatype, &nullable, &defaultValue, &charLen, &numPrec, &numScale)
		if err != nil {
			return result, err
		}
		col := Column{Name: *colName}
		col.Nullable = (*nullable) == "YES"
		col.DataType = psql.translateDataType(*datatype)
		if charLen != nil && *charLen > 0 {
			col.Size = *charLen
		}
		if numPrec != nil && *numPrec > 0 {
			col.Size = *numPrec
		}
		if numScale != nil && *numScale > 0 {
			col.Decimal = *numScale
		}

		for _, pkCol := range pks {
			if pkCol == col.Name {
				col.PrimaryKey = true
			}
		}
		fk, ok := fks[col.Name]
		if ok {
			col.ForeignKey = fk
		}

		if err != nil {
			return result, err
		}
		cols = append(cols, col)
	}

	result.Columns = cols

	if err != nil {
		return result, err
	}

	result, err = readIndexes(db, result)

	if err != nil {
		return result, err
	}

	return result, nil

}

func readIndexes(db *sql.DB, table Table) (Table, error) {

	query := `select indexname, indexdef
						from pg_indexes where tablename = $1`

	rows, err := db.Query(query, table.Name)
	defer rows.Close()
	if err != nil {
		return table, err
	}

	idxSlice := make([]Index, 0)
	var indexName, indexDef *string
	for rows.Next() {
		err = rows.Scan(&indexName, &indexDef)
		if err != nil {
			return table, err
		}
		if strings.HasSuffix(*indexName, "_pkey") {
			//skip primary keys
			continue
		}
		index := Index{
			Name:   *indexName,
			Unique: strings.Contains(*indexDef, "UNIQUE"),
		}
		//parse column names
		derefDef := *indexDef
		parenIndex := strings.Index(derefDef, "(")
		columnList := derefDef[parenIndex+1 : len(derefDef)-1]
		columns := strings.Split(columnList, ",")

		for idx, col := range columns {
			columns[idx] = strings.TrimSpace(col)
		}
		index.ColumnNames = columns

		idxSlice = append(idxSlice, index)
	}

	table.Indices = idxSlice
	return table, nil
}

func getPrimaryKeys(db *sql.DB, table *Table) ([]string, error) {

	query := `SELECT c.column_name
						FROM information_schema.table_constraints tc
						JOIN information_schema.constraint_column_usage AS ccu USING (constraint_schema, constraint_name)
						JOIN information_schema.columns AS c ON c.table_schema = tc.constraint_schema
						AND tc.table_name = c.table_name AND ccu.column_name = c.column_name
						WHERE constraint_type = 'PRIMARY KEY' and tc.table_name = $1`

	rows, err := db.Query(query, table.Name)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	results := make([]string, 0)
	var name *string
	for rows.Next() {
		rows.Scan(&name)
		results = append(results, *name)
	}

	return results, nil

}

func readForeignKeys(db *sql.DB, table Table) (map[string]ForeignKey, error) {

	query := `SELECT
						    tc.constraint_name,
						    kcu.column_name,
						    ccu.table_name AS foreign_table_name,
						    ccu.column_name AS foreign_column_name
						FROM
						    information_schema.table_constraints AS tc
						    JOIN information_schema.key_column_usage AS kcu
						      ON tc.constraint_name = kcu.constraint_name
						      AND tc.table_schema = kcu.table_schema
						    JOIN information_schema.constraint_column_usage AS ccu
						      ON ccu.constraint_name = tc.constraint_name
						      AND ccu.table_schema = tc.table_schema
						WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name=$1`

	rows, err := db.Query(query, table.Name)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	fks := make(map[string]ForeignKey)
	var name, columnName, foreignTable, foreignColumn *string
	for rows.Next() {
		rows.Scan(&name, &columnName, &foreignTable, &foreignColumn)
		if *foreignColumn != "id" {
			return nil, errors.New("invalid foreign key column (must be 'id')")
		}
		fk := ForeignKey{
			Name:      *name,
			TableName: *foreignTable,
		}
		fks[*columnName] = fk
	}

	return fks, nil

}

func processCreateTable(db *sql.DB, table Table) (Table, error) {

	rows, err := db.Query("show create table " + table.Name)
	defer rows.Close()
	if err != nil {
		return table, err
	}
	var name, ddl string
	if rows.Next() {
		rows.Scan(&name, &ddl)
	}

	//TODO this is a terrible hack - replace this with a real lexer/parser with an AST later

	lines := strings.Split(ddl, "\n")

	for _, line := range lines {
		if strings.Contains(line, "FOREIGN KEY") {
			table = addForeignKeyToModel(table, line)
		}
	}

	return table, nil
}

func addForeignKeyToModel(table Table, rawLine string) Table {

	tokens := strings.Split(strings.TrimSpace(rawLine), " ")

	fkName := stripDelims(tokens[1])
	fkColumn := stripDelims(tokens[4])
	fkTable := stripDelims(tokens[6])

	for idx, col := range table.Columns {
		if col.Name == fkColumn {
			var fk ForeignKey
			fk.Name = fkName
			fk.TableName = fkTable
			col.ForeignKey = fk
			table.Columns[idx] = col
		}
	}

	return table

}

func stripDelims(rawValue string) string {
	rawValue = strings.Replace(rawValue, "(", "", -1)
	rawValue = strings.Replace(rawValue, ")", "", -1)
	rawValue = strings.Replace(rawValue, "`", "", -1)
	return rawValue
}

func parenSplit(r rune) bool {
	return r == '(' || r == ')' || r == ','
}

/*
ReadTableList returns a slice of all tables in the given DB.
*/
func (psql *PostgresDialect) ReadTableList(db *sql.DB) ([]string, error) {

	results := make([]string, 0)

	rows, err := db.Query("select tablename from pg_catalog.pg_tables where schemaname = 'public'")
	defer rows.Close()
	if err != nil {
		fmt.Println(err.Error())
		return results, err
	}
	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return results, err
		}
		results = append(results, tableName)
	}
	return results, nil
}

// ReadTriggers reads the triggers for a table.
func (psql *PostgresDialect) ReadTriggers(db *sql.DB) ([]Trigger, error) {
	results := make([]Trigger, 0)

	rows, err := db.Query(
		`SELECT TRIGGER_NAME, EVENT_MANIPULATION, EVENT_OBJECT_TABLE, ACTION_STATEMENT, ACTION_TIMING
FROM information_schema.TRIGGERS
WHERE TRIGGER_SCHEMA = DATABASE()`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name, event, table, statement, timing string
		if err := rows.Scan(&name, &event, &table, &statement, &timing); err != nil {
			return nil, err
		}

		results = append(results, Trigger{
			Name:      name,
			Event:     event,
			Table:     table,
			Statement: statement,
			Timing:    timing,
		})
	}

	return results, nil
}

// CreateTriggerQuery generates a query for creating a trigger.
func (psql *PostgresDialect) CreateTriggerQuery(trigger Trigger) string {
	return fmt.Sprintf(
		`CREATE TRIGGER %s
%s %s
ON %s
FOR EACH ROW
BEGIN
	%s;
END`,
		trigger.Name, trigger.Timing, trigger.Event, trigger.Table, trigger.Statement)
}

// DropTriggerQuery generates a query for deleting a trigger.
func (psql *PostgresDialect) DropTriggerQuery(trigger Trigger) string {
	return fmt.Sprintf("DROP TRIGGER %s", trigger.Name)
}
