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

var modelCache map[string]mappingModel

type mappingModel struct {
	TableName       string
	IDField         reflect.StructField
	Fields          []reflect.StructField
	Type            reflect.Type
	InsertQuery     string
	UpdateQuery     string
	HardDeleteQuery string
	SoftDeleteQuery string
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
		return insert(domain, model, id)
	}
	return update(domain, model, id)

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

	logging.Errorln("UPDATE:", sb)

	return sb

}

func buildInsertQuery(model *mappingModel) string {

	sb := "insert into "
	sb += model.TableName
	sb += " (id,"
	sb += strings.Join(model.Columns(), ",")
	sb += (") VALUES ($1")
	for i := 2; i < len(model.Fields); i++ {
		sb += ",$"
		sb += strconv.Itoa(i)
	}
	sb += ")"

	logging.Errorln("INSERT:", sb)

	return sb

}

func buildSoftDeleteQuery(model *mappingModel) string {

	sb := "UPDATE "
	sb += model.TableName
	sb += " SET is_deleted = true "
	sb += " WHERE id = $1"

	logging.Errorln("SOFT DELETE:", sb)

	return sb

}

func buildHardDeleteQuery(model *mappingModel) string {

	sb := "delete from  "
	sb += model.TableName
	sb += " WHERE id = $1"

	logging.Errorln("HARD DELETE:", sb)

	return sb

}

func insert(domain interface{}, model *mappingModel, id string) (string, error) {

	return id, nil
}

func update(domain interface{}, model *mappingModel, id string) (string, error) {

	return id, nil
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
	}

	model.Fields = fields
	model.InsertQuery = buildInsertQuery(model)
	model.UpdateQuery = buildUpdateQuery(model)
	model.HardDeleteQuery = buildHardDeleteQuery(model)
	model.SoftDeleteQuery = buildSoftDeleteQuery(model)

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
