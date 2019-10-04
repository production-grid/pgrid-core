package schema

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/logging"

	//used to bring in the postgres driver
	_ "github.com/lib/pq"
)

// PostgresDialect handles postgress specific DDL quirks.
type PostgresDialect struct {
}

/*
ColumnDefinition assembles a column definition DDL fragment.
*/
func (mysql *PostgresDialect) ColumnDefinition(col Column) string {

	//TODO convert to string templates if practical

	ddl := col.Name
	ddl += " "
	if col.DataType == "AUTOINCREMENT" {
		ddl += "INT"
	} else {
		ddl += col.DataType
	}
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
	if col.DataType == "AUTOINCREMENT" {
		ddl += " AUTO_INCREMENT UNIQUE"
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
PrimaryKeyDefinition assembles the DDL query to create a new foreign key.
*/
func (mysql *PostgresDialect) PrimaryKeyDefinition(table Table) string {

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
func (mysql *PostgresDialect) ForeignKeyDefinition(col Column) string {

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
func (mysql *PostgresDialect) IndexDefinition(table Table, index Index) string {

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
func (mysql *PostgresDialect) ReadCurrentModel(schemaName string, db *sql.DB) (*Model, error) {

	var model Model

	tables, err := mysql.ReadTableList(schemaName, db)
	if err != nil {
		return nil, err
	}

	tableModels := make([]Table, 0)

	for _, table := range tables {
		tableModel, err := mysql.ReadTableModel(db, table)
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

/*
ReadTableModel reads the full table model for a given table name, which means columns, etc.
*/
func (mysql *PostgresDialect) ReadTableModel(db *sql.DB, tableName string) (Table, error) {

	result := Table{Name: tableName}

	rows, err := db.Query("describe " + tableName)
	defer rows.Close()
	if err != nil {
		return result, err
	}
	cols := make([]Column, 0)
	for rows.Next() {
		var colName, datatype, nullable, pk, defaultValue, extra *string
		err = rows.Scan(&colName, &datatype, &nullable, &pk, &defaultValue, &extra)
		if err != nil {
			return result, err
		}
		col := Column{Name: *colName}
		col.PrimaryKey = (*pk) == "PRI"
		col.Nullable = (*nullable) == "YES"

		//parse the datatype stuff
		tokens := strings.FieldsFunc(*datatype, parenSplit)
		col.DataType = tokens[0]
		if len(tokens) > 1 {
			col.Size, err = strconv.Atoi(tokens[1])
		}
		if len(tokens) > 2 {
			col.Decimal, err = strconv.Atoi(tokens[2])
		}

		if (extra != nil) && (*extra == "auto_increment") {
			col.DataType = "AUTOINCREMENT"
		}

		if err != nil {
			return result, err
		}
		cols = append(cols, col)
	}

	result.Columns = cols

	result, err = readForeignKeys(db, result)

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

	rows, err := db.Query("show index from " + table.Name)
	defer rows.Close()
	if err != nil {
		return table, err
	}

	fkMap := make(map[string]bool)
	//pre cache foreign keys so we know to skip them
	for _, col := range table.ForeignKeyColumns() {
		if col.ForeignKey.Name != "" {
			fkMap[col.ForeignKey.Name] = true
		}
	}

	idxMap := make(map[string]Index)

	for rows.Next() {
		var tableName, keyName, columnName, collation, subPart, nullable, indexType, comment, indexComment *string
		var nonUnique, packed *bool
		var seqInIndex, cardinality *int
		err = rows.Scan(&tableName, &nonUnique, &keyName, &seqInIndex, &columnName, &collation, &cardinality, &subPart, &packed, &nullable, &indexType, &comment, &indexComment)
		if err != nil {
			return table, err
		}
		//skip the primary
		if *keyName == "PRIMARY" {
			continue
		}
		_, exists := fkMap[*keyName]
		if exists {
			continue
		}

		idx, exists := idxMap[*keyName]
		if !exists {
			idx = Index{}
			idx.Name = *keyName
			idx.Unique = !*nonUnique
			idx.ColumnNames = make([]string, 0)
		}
		idx.ColumnNames = append(idx.ColumnNames, *columnName)
		idxMap[idx.Name] = idx

	}

	idxSlice := make([]Index, 0)
	for _, v := range idxMap {
		idxSlice = append(idxSlice, v)
	}
	table.Indices = idxSlice
	return table, nil
}

func readForeignKeys(db *sql.DB, table Table) (Table, error) {
	return processCreateTable(db, table)
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
func (mysql *PostgresDialect) ReadTableList(schemaName string, db *sql.DB) ([]string, error) {

	results := make([]string, 0)

	rows, err := db.Query("select tablename from pg_catalog.pg_tables where schemaname = $1", schemaName)
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
	logging.LogJSON(results)
	return results, nil
}

// ReadTriggers reads the triggers for a table.
func (mysql *PostgresDialect) ReadTriggers(db *sql.DB) ([]Trigger, error) {
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
func (mysql *PostgresDialect) CreateTriggerQuery(trigger Trigger) string {
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
func (mysql *PostgresDialect) DropTriggerQuery(trigger Trigger) string {
	return fmt.Sprintf("DROP TRIGGER %s", trigger.Name)
}
