package schema

/*
SchemaChanger can execute arbitrary code and generate schema changes.
*/
type SchemaChanger func(sourceModel Model, targetModel Model) ([]SchemaChange, error)

/*
PostMigrator is designed to clean up post migration.
*/
type PostMigrator func(sourceModel Model, targetModel Model) error

/*
SchemaChange models a single schema change operation.
*/
type SchemaChange struct {
	PostMigrateOnly bool
	ChangeType      ChangeType
	Table           Table
	Index           Index
	Column          Column
	Query           string
}

/*
ChangeType is a wrapper datatype for database change types.
*/
type ChangeType string

/*
Schema change types.
*/
const (
	CreateTable  ChangeType = "CREATE_TABLE"
	CreateIndex  ChangeType = "CREATE_INDEX"
	AddColumn    ChangeType = "ADD_COLUMN"
	ModifyColumn ChangeType = "MODIFY_COLUMN"
	DropColumn   ChangeType = "DROP_COLUMN"
	AddFK        ChangeType = "ADD_FK"
	Query        ChangeType = "QUERY"
)
