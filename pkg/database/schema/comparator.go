package schema

import (
	"database/sql"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/graph"
	"github.com/production-grid/pgrid-core/pkg/loaders"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

/*
CompareOnly performs a non-destructive comparison between the ds and schema file.
*/
func CompareOnly(ds *sql.DB, loader loaders.ResourceLoader, schemaFiles []string, changers []SchemaChanger) ([]Change, error) {

	target, err := ReadModelFromSchemaFiles(loader, schemaFiles)
	if err != nil {
		return nil, err
	}

	dialecter := &PostgresDialect{}

	source, err := dialecter.ReadCurrentModel(ds)
	if err != nil {
		return nil, err
	}

	return Compare(*source, *target, changers, &PostgresDialect{})

}

func executeSchemaChangers(source Model, target Model, changers []SchemaChanger) ([]Change, error) {

	results := make([]Change, 0)

	for _, changer := range changers {
		changes, err := changer(source, target)
		if err != nil {
			return nil, err
		}
		if changes != nil {
			results = append(results, changes...)
		}
	}

	return results, nil

}

/*
Compare returns a slice of schema changes required to convert the targetModel
to the sourceModel.
*/
func Compare(sourceModel Model, targetModel Model, changers []SchemaChanger, dialecter Dialecter) ([]Change, error) {

	var results []Change

	var fkChanges []Change

	sourceMap := buildTableMap(sourceModel)
	targetMap := buildTableMap(targetModel)

	logging.Infoln("Source Tables: ", len(sourceMap))
	logging.Infoln("Target Tables: ", len(targetMap))

	//sortedTables := sortTables(targetMap)

	for _, targetTable := range targetMap {
		changes, fks := compareTable(targetTable, sourceMap)
		results = append(results, changes...)
		fkChanges = append(fkChanges, fks...)
	}

	results = append(results, fkChanges...)
	if changers != nil {
		changes, err := executeSchemaChangers(sourceModel, targetModel, changers)
		if err != nil {
			return nil, err
		}
		if changes != nil {
			results = append(results, changes...)
		}
	}

	createTriggers, dropTriggers := compareTriggers(sourceModel, targetModel)

	// Drops must be performed first, since they may conflict with new
	// triggers.
	for _, trigger := range dropTriggers {
		results = append(results, Change{
			ChangeType: Query,
			Query:      dialecter.DropTriggerQuery(trigger),
		})
	}
	for _, trigger := range createTriggers {
		results = append(results, Change{
			ChangeType: Query,
			Query:      dialecter.CreateTriggerQuery(trigger),
		})
	}

	return results, nil

}

func compareTriggers(sourceModel, targetModel Model) (createTriggers, dropTriggers []Trigger) {
	createTriggers = make([]Trigger, 0)
	dropTriggers = make([]Trigger, 0)

OUTER_CREATE:
	for _, targetTrigger := range targetModel.Triggers {
		for _, srcTrigger := range sourceModel.Triggers {
			if sameTrigger(srcTrigger, targetTrigger) {
				continue OUTER_CREATE
			}
		}

		createTriggers = append(createTriggers, targetTrigger)
	}

OUTER_DROP:
	for _, srcTrigger := range sourceModel.Triggers {
		for _, targetTrigger := range targetModel.Triggers {
			if sameTrigger(srcTrigger, targetTrigger) {
				continue OUTER_DROP
			}
		}

		dropTriggers = append(dropTriggers, srcTrigger)
	}

	return
}

func sameTrigger(srcTrigger, targetTrigger Trigger) bool {
	return srcTrigger.Name == targetTrigger.Name &&
		srcTrigger.Table == targetTrigger.Table &&
		srcTrigger.Timing == targetTrigger.Timing
}

func compareTable(targetTable Table, sourceMap map[string]Table) ([]Change, []Change) {

	var changes []Change

	var fkChanges []Change

	existingTable, exists := sourceMap[targetTable.Name]

	if !exists {
		changes = append(changes, generateCreateTable(targetTable))
		changes = append(changes, generateCreateIndices(targetTable)...)
		fkChanges = append(fkChanges, generateCreateForeignKeys(targetTable)...)
	} else {
		return compareTables(existingTable, targetTable)
	}

	return changes, fkChanges

}

func hasSize(dataType string) bool {
	switch dataType {
	case "TEXT", "BIT", "INTEGER":
		return false
	default:
		return true
	}
}

func isEquivalentDataType(targetType string, existingType string) bool {

	loudTarget := strings.ToUpper(targetType)
	loudExisting := strings.ToUpper(existingType)

	if loudTarget == loudExisting {
		return true
	}

	return strings.HasPrefix(loudTarget, loudExisting)

}

func compareTables(existingTable Table, targetTable Table) ([]Change, []Change) {

	results := make([]Change, 0)
	fkResults := make([]Change, 0)

	existingColumnMap := existingTable.ColumnMap()

	existingColumnsSeen := make(map[string]bool) //used to handle drop column events

	for _, targetCol := range targetTable.Columns {
		existingCol, exists := existingColumnMap[targetCol.Name]
		if !exists {
			results = append(results, Change{ChangeType: AddColumn, Table: targetTable, Column: targetCol, Reason: "column not present"})
		} else {
			existingColumnsSeen[targetCol.Name] = true
			if !isEquivalentDataType(targetCol.DataType, existingCol.DataType) {
				results = append(results, Change{ChangeType: ModifyColumn, Table: targetTable, Column: targetCol, OldColumn: existingCol, Reason: "data type mismatch"})
			} else if targetCol.Nullable != existingCol.Nullable {
				results = append(results, Change{ChangeType: ModifyColumn, Table: targetTable, Column: targetCol, OldColumn: existingCol, Reason: "nullable status"})
			} else if hasSize(targetCol.DataType) {
				if targetCol.Size != existingCol.Size {
					if targetCol.Size < existingCol.Size {
						logging.Warnf("Ignoring size change for %s since it might truncate data", targetCol.Name)
						continue
					}
					results = append(results, Change{ChangeType: ModifyColumn, Table: targetTable, Column: targetCol, OldColumn: existingCol, Reason: "size mismatch"})
				} else if targetCol.Decimal != existingCol.Decimal {
					logging.Error("Decimal Change")
					if targetCol.Decimal < existingCol.Decimal {
						logging.Warnf("Ignoring size change for %s since it might truncate data", targetCol.Name)
						continue
					}
					results = append(results, Change{ChangeType: ModifyColumn, Table: targetTable, Column: targetCol, OldColumn: existingCol, Reason: "decimal scale mismatch"})
				} // sizes not equal
			} // type check

			if targetCol.ForeignKey.Name != "" {
				if targetCol.ForeignKey.Name != existingCol.ForeignKey.Name {
					fkResults = append(fkResults, Change{ChangeType: AddFK, Table: targetTable, Column: targetCol, Reason: "foreign key missing"})
				}
			}

		} // exists or not
	}

	//compare indexes
	existingIndexMap := existingTable.IndexMap()

	for _, idx := range targetTable.Indices {
		_, exists := existingIndexMap[idx.Name]
		if !exists {
			results = append(results, Change{ChangeType: CreateIndex, Table: targetTable, Index: idx})
		}
	}

	//handle drop events
	for _, existingCol := range existingTable.Columns {
		_, seen := existingColumnsSeen[existingCol.Name]
		if !seen {
			results = append(results, Change{ChangeType: DropColumn, Table: targetTable, Column: existingCol})
		}
	}

	return results, fkResults

}

func generateCreateIndices(targetTable Table) []Change {

	results := make([]Change, 0)

	for _, idx := range targetTable.Indices {
		result := Change{ChangeType: CreateIndex, Table: targetTable, Index: idx}
		results = append(results, result)
	}

	return results
}

func generateCreateForeignKeys(targetTable Table) []Change {

	results := make([]Change, 0)

	for _, col := range targetTable.ForeignKeyColumns() {
		result := Change{ChangeType: AddFK, Table: targetTable, Column: col, PostMigrateOnly: true}
		results = append(results, result)
	}

	return results
}

func generateCreateTable(targetTable Table) Change {
	result := Change{ChangeType: CreateTable, Table: targetTable}
	return result
}

func sortTables(tables map[string]Table) []*Table {

	results := make([]*Table, 0)

	vertices := make(map[string]*graph.Vertex, len(tables))
	for _, t := range tables {
		vertices[t.Name] = &graph.Vertex{ID: t.Name}
	}

	arrVertices := make([]graph.Vertex, 0)

	for _, v := range vertices {
		table := tables[v.ID]
		fks := table.ForeignKeys()
		if len(fks) > 0 {
			parents := make([]*graph.Vertex, 0)
			for _, fk := range fks {
				if fk.TableName != v.ID {
					parents = append(parents, vertices[fk.TableName])
				}
			}
			v.ParentVertices = parents
		}
		arrVertices = append(arrVertices, *v)
	}

	tableNames := graph.TopographicSort(arrVertices)

	for _, tbl := range tableNames {
		table := tables[tbl]
		results = append(results, &table)
	}

	return results

}

func buildTableMap(model Model) map[string]Table {

	result := make(map[string]Table)

	if model.Tables != nil {
		for _, table := range model.Tables {
			result[table.Name] = table
		}
	}

	return result
}
