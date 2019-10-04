package schema

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"

	// brining in the mysql driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	migratorLockKey           = "lock:schematron"
	migratorLockTimeout int32 = 300 // Seconds
)

var (
	insideParens = regexp.MustCompile(`\(([^()]*)\)`)
)

/*
DefaultMigrator is the default implementation of the Migrator interface.
*/
type DefaultMigrator struct {
	Datasource *sql.DB
	Dialecter  Dialecter
	Profile    string
}

/*
Migrator defines the interface for datase migrators.
*/
type Migrator interface {
	Migrate(*os.File) error
}

// Migrate migrates the local database to the format specified in
// the schema file.
func (migrator *DefaultMigrator) Migrate(loader loaders.ResourceLoader, schemaFiles []string, schemaName string, changers []SchemaChanger, postMigrators []PostMigrator) error {
	// Take the migrator lock or return
	/*
		if err := addMigratorLock(filepath.Base(schemaPath)); err == nil {
			logging.Warn("Took schematron migrator lock")
			defer removeMigratorLock(filepath.Base(schemaPath))
		} else {
			logging.Warn("Schematron lock is taken - doing nothing")
			return nil
		}
	*/

	target, err := ReadModelFromSchemaFiles(loader, schemaFiles, migrator.Profile)
	if err != nil {
		return err
	}

	source, err := migrator.Dialecter.ReadCurrentModel(schemaName, migrator.Datasource)
	if err != nil {
		return err
	}

	changes, err := Compare(*source, *target, changers, migrator.Dialecter)

	if err != nil {
		return err
	}

	err = migrator.executeChangeSets(changes)
	if err != nil {
		return err
	}
	if postMigrators != nil {
		for _, migrator := range postMigrators {
			err = migrator(*source, *target)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (migrator *DefaultMigrator) runSchemaChangers(changers []SchemaChanger) []SchemaChange {

	changes := make([]SchemaChange, 0)

	return changes

}

//PreMigrate updates the database changes that can be safely executed before a software update
func (migrator *DefaultMigrator) PreMigrate(schemaFiles string, changers []SchemaChanger) error {

	return nil
}

//PostMigrate updates the database changes that can be safely executed before a software update
func (migrator *DefaultMigrator) PostMigrate(schemaFiles string, changers []SchemaChanger) error {

	return nil
}

//GeneratePreMigrationScript generates a script with the database changes that can be run with the update
func (migrator *DefaultMigrator) GeneratePreMigrationScript(changes []SchemaChange, changers []SchemaChanger) (*string, error) {

	return nil, nil
}

//GeneratePostMigrationScript generates a script with the changes that can be run once all software nodes are updated.
func (migrator *DefaultMigrator) GeneratePostMigrationScript(changes []SchemaChange, changers []SchemaChanger) (*string, error) {

	return nil, nil
}

func (migrator *DefaultMigrator) executeChangeSets(changes []SchemaChange) error {

	for _, change := range changes {
		var err error
		switch change.ChangeType {
		case CreateTable:
			err = migrator.executeCreateTable(change)
		case CreateIndex:
			err = migrator.executeCreateIndex(change)
		case AddColumn:
			err = migrator.executeAddColumn(change)
		case ModifyColumn:
			err = migrator.executeModifyColumn(change)
		case AddFK:
			err = migrator.executeAddForeignKey(change)
		case Query:
			err = migrator.executeQuery(change)
		case DropColumn:
			return fmt.Errorf("Destructive changes are not allowed. Drop column: %s.%s manually", change.Table.Name, change.Column.Name)
		}

		if err != nil {
			return err
		}

	}

	return nil

}

func (migrator *DefaultMigrator) executeCreateTable(change SchemaChange) error {

	//TODO method too long - should be caught once we add uncle bob style checking

	logging.Infof("Creating Table: %s\n", change.Table.Name)

	pkCols := make([]string, 0)

	sql := "create table " + change.Table.Name
	sql += "("
	for idx, col := range change.Table.Columns {
		if col.PrimaryKey {
			pkCols = append(pkCols, col.Name)
		}
		if idx > 0 {
			sql += ", "
		}
		sql += migrator.Dialecter.ColumnDefinition(col)
	}
	pkDef := migrator.Dialecter.PrimaryKeyDefinition(change.Table)
	if pkDef != "" {
		sql += ", "
		sql += pkDef
	}
	sql += ")"
	logging.Infoln("Executing:", sql)

	_, err := migrator.Datasource.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (migrator *DefaultMigrator) executeQuery(change SchemaChange) error {

	logging.Infof("Executing Query Change: %s", change.Query)

	_, err := migrator.Datasource.Exec(change.Query)
	if err != nil {
		logging.Errorf("Failed to execute query: %v", err)
		return err
	}

	return nil

}

func (migrator *DefaultMigrator) executeAddForeignKey(change SchemaChange) error {

	logging.Infof("Adding Foreign Key: %s:%s\n", change.Table.Name, change.Column.ForeignKey.Name)

	sql := "alter table "
	sql += change.Table.Name
	sql += " add "
	sql += migrator.Dialecter.ForeignKeyDefinition(change.Column)

	logging.Infoln("Executing:", sql)

	_, err := migrator.Datasource.Exec(sql)
	if err != nil {
		logging.Errorf("Failed to add foreign key: %v", err)
		return err
	}

	return nil
}

func (migrator *DefaultMigrator) executeModifyColumn(change SchemaChange) error {

	logging.Infof("Modifying Column: %s:%s\n", change.Table.Name, change.Column.Name)

	sql := "alter table "
	sql += change.Table.Name
	sql += " modify column "
	sql += migrator.Dialecter.ColumnDefinition(change.Column)

	logging.Infoln("Executing:", sql)

	_, err := migrator.Datasource.Exec(sql)
	if err != nil {
		logging.Errorf("Failed to modify column: %v", err)
		return err
	}

	return nil
}

func (migrator *DefaultMigrator) executeAddColumn(change SchemaChange) error {

	logging.Infof("Adding Column: %s:%s\n", change.Table.Name, change.Column.Name)

	sql := "alter table "
	sql += change.Table.Name
	sql += " add column "
	sql += migrator.Dialecter.ColumnDefinition(change.Column)

	logging.Infoln("Executing:", sql)

	_, err := migrator.Datasource.Exec(sql)
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	return nil
}

func (migrator *DefaultMigrator) executeCreateIndex(change SchemaChange) error {

	logging.Infof("Creating Table: %s\n", change.Table.Name)

	sql := "create " + migrator.Dialecter.IndexDefinition(change.Table, change.Index)

	logging.Infoln("Executing:", sql)

	_, err := migrator.Datasource.Exec(sql)
	if err != nil {
		logging.Errorf("Failed to create index: %v", err)
		return err
	}

	return nil
}

/*
ReadModelFromSchemaFiles reads a schema file and loads it into a model struct.
*/
func ReadModelFromSchemaFiles(loader loaders.ResourceLoader, schemaFiles []string, profile string) (*Model, error) {

	tables := make([]Table, 0)

	for _, schemaFile := range schemaFiles {
		newTables, err := ReadModelFromSchemaFile(loader, schemaFile, profile)
		if err != nil {
			return nil, err
		}
		if newTables != nil {
			tables = append(tables, newTables...)
		}
	}

	model := Model{}
	model.Tables = tables

	model = loadProfile(model, profile)

	return &model, nil

}

/*
ReadModelFromSchemaFile reads a schema file and loads it into a model struct.
*/
func ReadModelFromSchemaFile(loader loaders.ResourceLoader, schemaFile string, profile string) ([]Table, error) {
	b, err := loader.Bytes(schemaFile)
	if err != nil {
		return nil, err
	}
	model := Model{}
	if err := json.Unmarshal(b, &model); err != nil {
		return nil, err
	}

	logging.Debugf("Tables in schema file: %+v", len(model.Tables))

	//validate model
	for _, tbl := range model.Tables {
		if tbl.Name == "" {
			return nil, errors.New("empty table named")
		}
	}

	//handle duping table clones
	tables := make([]Table, 0)

	for _, table := range model.Tables {
		table = sanitizeTable(table)
		tables = append(tables, table)
		if table.CloneNames != "" {
			tableNames := strings.Split(table.CloneNames, ",")
			for _, tableName := range tableNames {
				cloneTable := table
				cloneTable.CloneNames = ""
				cloneTable.Name = tableName

				newCols := make([]Column, 0)
				cols := cloneTable.Columns
				for _, col := range cols {
					if col.ForeignKey.Name != "" {
						col.ForeignKey.Name = col.ForeignKey.Name + "_" + cloneTable.Name
						if col.ForeignKey.TableName == table.Name {
							col.ForeignKey.TableName = tableName
						}
					}
					newCols = append(newCols, col)
				}
				cloneTable.Columns = newCols
				tables = append(tables, cloneTable)
			}
		}
	}

	return tables, nil

}

// loadProfile loads profile specific changes to the target schema.
func loadProfile(model Model, profile string) Model {
	if profile == "" {
		return model
	}

	profileOpts, ok := model.Profiles[profile]
	if !ok {
		logging.Warnf("Unrecognized schematron profile: %s", profile)
		return model
	}

	logging.Infof("Using schematron profile: %+v", profile)

	// Add a trigger to each table with an is_deleted column to prevent hard
	// deletes.
	if profileOpts.EnforceSoftDeletes {
		for _, table := range model.Tables {
			if isSoftDeleteTable(table) {
				model.Triggers = append(model.Triggers, Trigger{
					Name:      table.Name + "_enforce_soft_delete",
					Event:     "DELETE",
					Table:     table.Name,
					Statement: "SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'DELETE not allowed'",
					Timing:    "BEFORE",
				})
			}
		}
	}

	return model
}

func isSoftDeleteTable(table Table) bool {
	if table.AllowHardDelete {
		return false
	}

	const softDeleteColumn = "is_deleted"

	for _, column := range table.Columns {
		if column.Name == softDeleteColumn {
			return true
		}
	}

	return false
}

//sanitizeTable handles weird post processing edge cases
func sanitizeTable(table Table) Table {

	columns := make([]Column, 0)

	for _, col := range table.Columns {
		/*
			This is designed to normalize columns
			that use inline size declaration like
			DECIMAL (10, 7).  We check for an open
			paren just to save a regex.
		*/
		if strings.Contains(col.DataType, "(") {
			rs := insideParens.FindStringSubmatch(col.DataType)
			//for some reason the contents of parentheses
			//are returned as the second element of the array
			if len(rs) == 2 {
				tokens := strings.Split(rs[1], ",")
				if len(tokens) > 0 {
					size, err := strconv.Atoi(tokens[0])
					if err != nil {
						logging.Warnln("error parsing size token", err)
					} else {
						col.Size = size
					}
				}
				if len(tokens) > 1 {
					decimal, err := strconv.Atoi(tokens[1])
					if err != nil {
						logging.Warnln("error parsing decimal precision token", err)
					} else {
						col.Decimal = decimal
					}
				}
				if len(tokens) > 2 {
					logging.Warnln("data type has more size tokens than expected")
				}
			}
		}
		columns = append(columns, col)
	}

	table.Columns = columns

	return table

}
