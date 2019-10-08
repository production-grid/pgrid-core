package relational

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/production-grid/pgrid-core/pkg/ids"
	"github.com/production-grid/pgrid-core/pkg/logging"
)

// enumerates the constants for databas types
const (
	PRIMARY = "primary"
	REPLICA = "replica"
)

// ErrNoResults is returned when a finder does not return results.
var ErrNoResults = errors.New("no results")

var modelCache map[string]mappingModel

type mappingModel struct {
	TableName       string
	IDField         reflect.StructField
	FieldsWithID    []reflect.StructField
	Fields          []reflect.StructField
	Type            reflect.Type
	InsertQuery     string
	UpdateQuery     string
	HardDeleteQuery string
	SoftDeleteQuery string
	FindByIDQuery   string
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

// FindByID locates an entity by ID
func FindByID(dbType string, tableName string, id string, target interface{}) error {

	model, err := resolveMappingModel(target, tableName)

	if err != nil {
		return err
	}

	logging.Infoln("SQL:", model.FindByIDQuery, id)

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

	el := reflect.ValueOf(target).Elem()

	for idx, fld := range model.FieldsWithID {
		val := el.FieldByIndex(fld.Index)
		logging.Errorln(val.Type())
		i := val.Interface()
		targets[idx] = &i
	}

	err := rows.Scan(targets...)

	if err != nil {
		return err
	}

	logging.LogJSON(targets)

	/*

		el := reflect.ValueOf(target).Elem()

		for idx, fld := range model.FieldsWithID {
			val := el.FieldByIndex(fld.Index)
			v := val.Interface()
			scanVal := reflect.ValueOf(targets[idx])
			switch v.(type) {
			case string:
				val.SetString(targets[idx].(string))
			case bool:
				val.SetBool(scanVal.Bool())
			case time.Time:
				val.Set(scanVal)
			}
		}
	*/

	return err
}

// FindByIDWithTx locates an entity by ID
func FindByIDWithTx(tx *sql.Tx, tableName string, id string, target interface{}) error {

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

	logging.Infoln("SQL:", model.InsertQuery)

	_, err = tx.Exec(model.InsertQuery, params...)

	return err
}

func update(tx *sql.Tx, el reflect.Value, model *mappingModel, id string) error {

	params, err := resolveMutationParams(el, model, id)

	if err != nil {
		return err
	}

	logging.Infoln("SQL:", model.UpdateQuery)

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
