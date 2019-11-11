package relational

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/database/schema"
	"github.com/production-grid/pgrid-core/pkg/ids"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

// enumerates the constants for database types
const (
	PRIMARY = "primary"
	REPLICA = "replica"
)

// EntityFactory is used to produce new entities
type EntityFactory func() interface{}

// EntityFinder models the generic form of a finder
type EntityFinder interface {
	FindInterfaceByID(dbType string, id string) (interface{}, error)
}

// Entity models the generic form of an entity
type Entity interface {
	Save() (id string, err error)
	Delete() error
}

// ErrNoResults is returned when a finder does not return results.
var ErrNoResults = errors.New("no results")

var modelCache map[string]mappingModel

type mappingModel struct {
	TableName       string
	SoftDeleted     bool
	IDField         reflect.StructField
	FieldsWithID    []reflect.StructField
	Fields          []reflect.StructField
	Type            reflect.Type
	InsertQuery     string
	UpdateQuery     string
	HardDeleteQuery string
	SoftDeleteQuery string
	FindByIDQuery   string
	FindCountQuery  string
}

func (model *mappingModel) Columns() []string {

	columns := make([]string, len(model.Fields))

	for idx, fld := range model.Fields {
		columns[idx] = fld.Tag.Get("col")
	}

	return columns

}

// NewWritableTx returns a writable transation for the primary database.
func NewWritableTx() (*sql.Tx, error) {
	return Primary.Begin()
}

//HardDeleteWithTx deletes an entity from the database with a transaction.
func HardDeleteWithTx(tx *sql.Tx, domain interface{}, table string) error {

	model, err := resolveMappingModel(domain, table)

	if err != nil {
		return err
	}

	if model == nil {
		return errors.New("unable to resolve mapping model")
	}

	el := reflect.ValueOf(domain).Elem()

	id, err := resolveID(el, model)

	if err != nil {
		return err
	}

	_, err = tx.Exec(model.HardDeleteQuery, id)

	return err
}

//SoftDelete updates the is_deleted flag for an entity
func SoftDelete(domain interface{}, table string) error {

	model, err := resolveMappingModel(domain, table)

	if err != nil {
		return err
	}

	if model == nil {
		return errors.New("unable to resolve mapping model")
	}

	el := reflect.ValueOf(domain).Elem()

	id, err := resolveID(el, model)

	if err != nil {
		return err
	}

	logging.Traceln("SQL", model.SoftDeleteQuery)

	_, err = Primary.Exec(model.SoftDeleteQuery, id)

	return err
}

//SoftDeleteWithTx updates the is_deleted flag for an entity
func SoftDeleteWithTx(tx *sql.Tx, domain interface{}, table string) error {

	model, err := resolveMappingModel(domain, table)

	if err != nil {
		return err
	}

	if model == nil {
		return errors.New("unable to resolve mapping model")
	}

	el := reflect.ValueOf(domain).Elem()

	id, err := resolveID(el, model)

	if err != nil {
		return err
	}

	_, err = tx.Exec(model.SoftDeleteQuery, id)

	return err
}

func resolveDatabaseType(dbType string) *sql.DB {
	switch dbType {
	case REPLICA:
		return Replica
	default:
		return Primary
	}
}

// FindOneByWhereClause returns the first object matching the where clause
func FindOneByWhereClause(dbType string, tableName string, whereClause string, target interface{}, params ...interface{}) error {

	model, err := resolveMappingModel(target, tableName)

	if err != nil {
		return err
	}

	sql := "select id, "
	sql += strings.Join(model.Columns(), ",")
	sql += " from "
	sql += tableName
	sql += " "
	sql += whereClause
	if model.SoftDeleted {
		sql += " and is_deleted = false"
	}

	logging.Tracef("SQL: %v", sql)

	rows, err := resolveDatabaseType(dbType).Query(sql, params...)
	defer rows.Close()

	if err != nil {
		logging.Errorf("Find failed: %v", err)
		return err
	}

	if rows.Next() {
		err := scan(model, rows, target)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrNoResults

}

// FindByWhereClause returns all objects in the given table matching the where clause.
func FindByWhereClause(dbType string, entityFactory EntityFactory, tableName string, whereClause string, params ...interface{}) ([]interface{}, error) {

	prototype := entityFactory()

	model, err := resolveMappingModel(prototype, tableName)

	if err != nil {
		return nil, err
	}

	sql := "select id, "
	sql += strings.Join(model.Columns(), ",")
	sql += " from "
	sql += tableName
	sql += " "
	sql += whereClause
	if model.SoftDeleted {
		sql += " and is_deleted = false"
	}

	logging.Tracef("SQL: %v", sql)

	results := make([]interface{}, 0)

	rows, err := resolveDatabaseType(dbType).Query(sql, params...)
	defer rows.Close()

	if err != nil {
		logging.Errorf("Find failed: %v", err)
		return results, err
	}

	for rows.Next() {
		target := entityFactory()
		err := scan(model, rows, target)
		if err != nil {
			return results, err
		}
		results = append(results, target)
	}

	return results, nil

}

// FindAll returns all objects in the given table matching the where clause.
func FindAll(dbType string, entityFactory EntityFactory, tableName string, orderClause string) ([]interface{}, error) {

	prototype := entityFactory()

	model, err := resolveMappingModel(prototype, tableName)

	if err != nil {
		return nil, err
	}

	sql := "select id, "
	sql += strings.Join(model.Columns(), ",")
	sql += " from "
	sql += tableName
	if model.SoftDeleted {
		sql += " where is_deleted = false"
	}
	if orderClause != "" {
		sql += " "
		sql += orderClause
	}

	logging.Tracef("SQL: %v", sql)

	results := make([]interface{}, 0)

	rows, err := resolveDatabaseType(dbType).Query(sql)
	defer rows.Close()

	if err != nil {
		logging.Errorf("Find failed: %v", err)
		return results, err
	}

	for rows.Next() {
		target := entityFactory()
		err := scan(model, rows, target)
		if err != nil {
			return results, err
		}
		results = append(results, target)
	}

	return results, nil

}

// FindCount returns the number of activate rows in the table
func FindCount(dbType string, tableName string, target interface{}) (int, error) {

	model, err := resolveMappingModel(target, tableName)

	if err != nil {
		return 0, err
	}

	logging.Traceln("SQL:", model.FindCountQuery)

	rows, err := resolveDatabaseType(dbType).Query(model.FindCountQuery)

	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() {
		var result *int
		err := rows.Scan(&result)
		if err != nil {
			return 0, err
		}
		return *result, nil
	}

	return 0, nil
}

// FindByID locates an entity by ID
func FindByID(dbType string, tableName string, id string, target interface{}) error {

	model, err := resolveMappingModel(target, tableName)

	if err != nil {
		return err
	}

	logging.Traceln("SQL:", model.FindByIDQuery, id)

	rows, err := resolveDatabaseType(dbType).Query(model.FindByIDQuery, id)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err := scan(model, rows, target)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrNoResults
}

func scan(model *mappingModel, rows *sql.Rows, target interface{}) error {

	targets := make([]interface{}, len(model.FieldsWithID))

	v := reflect.ValueOf(target)

	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}

	for idx, fld := range model.FieldsWithID {
		targets[idx] = v.Elem().FieldByIndex(fld.Index).Addr().Interface()
	}

	err := rows.Scan(targets...)

	if err != nil {
		return err
	}

	return err
}

// FindByIDWithTx locates an entity by ID
func FindByIDWithTx(tx *sql.Tx, tableName string, id string, target interface{}) error {

	model, err := resolveMappingModel(target, tableName)

	if err != nil {
		return err
	}

	logging.Traceln("SQL:", model.FindByIDQuery, id)

	rows, err := tx.Query(model.FindByIDQuery, id)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		err := scan(model, rows, target)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrNoResults
}

//LoadChildren loads child values from the database in accordance with the prophecy
func LoadChildren(dbType string, tableName string, entityCol string, childCol string, id string) ([]string, error) {

	selectQuery := "select " + childCol + " from " + tableName + " where " + entityCol + " = $1"

	rows, err := resolveDatabaseType(dbType).Query(selectQuery, id)

	if err != nil {
		return nil, err
	}

	results := make([]string, 0)
	for rows.Next() {
		var childValue string
		err = rows.Scan(&childValue)
		if err != nil {
			return nil, err
		}
		results = append(results, childValue)
	}

	return results, nil
}

//SaveChildren saves a collection of child values to the database
func SaveChildren(tx *sql.Tx, table string, domain interface{}, entityCol string, childCol string, children []string) error {

	model, err := resolveMappingModel(domain, table)

	if err != nil {
		return err
	}

	if model == nil {
		return errors.New("unable to resolve mapping model")
	}

	el := reflect.ValueOf(domain).Elem()

	id, err := resolveID(el, model)

	if err != nil {
		return err
	}

	deleteQuery := "delete from " + table + " where " + entityCol + " = $1"

	_, err = tx.Exec(deleteQuery, id)

	if err != nil {
		return err
	}

	if children != nil {

		insertQuery := "insert into " + table + "(id," + entityCol + ", " + childCol + ")"
		insertQuery += " VALUES ($1, $2, $3)"

		for _, child := range children {
			_, err = tx.Exec(insertQuery, ids.NewSecureID(), id, child)
			if err != nil {
				return err
			}
		}

	}

	return nil

}

//SaveWithTx saves a domain object to the database with a transaction.
func SaveWithTx(tx *sql.Tx, domain interface{}, table string) (string, error) {

	model, err := resolveMappingModel(domain, table)

	if err != nil {
		return "", err
	}

	if model == nil {
		return "", errors.New("unable to resolve mapping model")
	}

	el := reflect.ValueOf(domain).Elem()

	id, err := resolveID(el, model)

	if err != nil {
		return "", err
	}

	if id == "" {
		id = ids.NewSecureID()
		err = setID(el, model, id)
		if err != nil {
			return "", nil
		}
		err = insert(tx, el, model, id)
	} else {
		err = update(tx, el, model, id)
	}

	return id, err

}

func buildUpdateQuery(model *mappingModel) string {

	sb := "update "
	sb += model.TableName
	sb += " SET "
	for idx, fld := range model.Fields {
		if idx > 0 {
			sb += ", "
		}
		sb += fld.Tag.Get("col")
		sb += " = $"
		sb += strconv.Itoa(idx + 2)
	}
	sb += " WHERE id = $1"

	return sb

}

func buildFindByIDQuery(model *mappingModel) string {

	sb := "SELECT id, "
	sb += strings.Join(model.Columns(), ",")
	sb += " FROM "
	sb += model.TableName
	sb += " WHERE id = $1"
	if model.SoftDeleted {
		sb += " AND is_deleted = false"
	}

	return sb

}

func buildFindCountQuery(model *mappingModel) string {

	sb := "SELECT count(*) FROM "
	sb += model.TableName
	if model.SoftDeleted {
		sb += " WHERE is_deleted = false"
	}

	return sb

}

func buildInsertQuery(model *mappingModel) string {

	sb := "insert into "
	sb += model.TableName
	sb += " (id,"
	sb += strings.Join(model.Columns(), ",")
	sb += (") VALUES ($1")
	for i := 0; i < len(model.Fields); i++ {
		sb += ",$"
		sb += strconv.Itoa(i + 2)
	}
	sb += ")"

	return sb

}

func buildSoftDeleteQuery(model *mappingModel) string {

	sb := "UPDATE "
	sb += model.TableName
	sb += " SET is_deleted = true "
	sb += " WHERE id = $1"

	return sb

}

func buildHardDeleteQuery(model *mappingModel) string {

	sb := "delete from  "
	sb += model.TableName
	sb += " WHERE id = $1"

	return sb

}

func resolveMutationParams(el reflect.Value, model *mappingModel, id string) ([]interface{}, error) {

	results := make([]interface{}, len(model.Fields)+1)
	results[0] = id

	for idx, fld := range model.Fields {
		results[idx+1] = el.FieldByIndex(fld.Index).Interface()
	}

	return results, nil

}

func insert(tx *sql.Tx, el reflect.Value, model *mappingModel, id string) error {

	params, err := resolveMutationParams(el, model, id)

	if err != nil {
		return err
	}

	logging.Traceln("SQL:", model.InsertQuery)

	_, err = tx.Exec(model.InsertQuery, params...)

	return err
}

func update(tx *sql.Tx, el reflect.Value, model *mappingModel, id string) error {

	params, err := resolveMutationParams(el, model, id)

	if err != nil {
		return err
	}

	logging.Traceln("SQL:", model.UpdateQuery)

	_, err = tx.Exec(model.UpdateQuery, params...)

	return err

}

func resolveID(el reflect.Value, model *mappingModel) (string, error) {

	val := el.FieldByIndex(model.IDField.Index)
	return val.String(), nil

}

func setID(el reflect.Value, model *mappingModel, id string) error {

	val := el.FieldByIndex(model.IDField.Index)
	val.SetString(id)

	return nil

}

func resolveModelCache(t reflect.Type) (*mappingModel, error) {

	if modelCache != nil {
		mappingModel, ok := modelCache[t.String()]
		if ok {
			return &mappingModel, nil
		}
	}

	return nil, nil

}

func resolveMappingModel(domain interface{}, table string) (*mappingModel, error) {

	el := reflect.ValueOf(domain).Elem()
	t := el.Type()

	model, err := resolveModelCache(t)

	if err != nil {
		return nil, err
	}

	if model != nil {
		return model, nil
	}

	fld, ok := t.FieldByName("ID")

	if !ok {
		return nil, errors.New("entity must have an ID field")
	}

	model = &mappingModel{
		TableName: table,
		Type:      t,
		IDField:   fld,
	}

	model.SoftDeleted = schema.IsSoftDeleted(table)

	fieldsWithID := make([]reflect.StructField, 0)
	fieldsWithID = append(fieldsWithID, fld)

	fields := make([]reflect.StructField, 0)

	for i := 0; i < t.NumField(); i++ {
		fldCand := t.Field(i)
		column := fldCand.Tag.Get("col")
		if column == "" {
			continue
		}
		if column == "id" {
			continue
		}
		fields = append(fields, fldCand)
		fieldsWithID = append(fieldsWithID, fldCand)
	}

	model.Fields = fields
	model.FieldsWithID = fieldsWithID
	model.InsertQuery = buildInsertQuery(model)
	model.UpdateQuery = buildUpdateQuery(model)
	model.HardDeleteQuery = buildHardDeleteQuery(model)
	model.SoftDeleteQuery = buildSoftDeleteQuery(model)
	model.FindByIDQuery = buildFindByIDQuery(model)
	model.FindCountQuery = buildFindCountQuery(model)

	if err != nil {
		return nil, err
	}

	//cache model
	if modelCache == nil {
		modelCache = make(map[string]mappingModel)
	}
	modelCache[t.String()] = *model

	return model, nil
}

//Save saves an entity to the db without a transaction.
func Save(domain interface{}, table string) (string, error) {
	tx, err := NewWritableTx()
	if err != nil {
		return "", err
	}
	id, err := SaveWithTx(tx, domain, table)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return id, nil
}
