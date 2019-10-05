package relational

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/production-grid/pgrid-core/pkg/logging"
)

// enumerates the constants for databas types
const (
	PRIMARY = "primary"
	REPLICA = "replica"
)

var modelCache map[string]mappingModel

type mappingModel struct {
	IDField reflect.StructField
	Fields  []reflect.StructField
}

// NewWritableTx returns a writable transation for the primary database.
func NewWritableTx() (*sql.Tx, error) {
	return Primary.Begin()
}

//SaveWithTx saves a domain object to the database with a transaction.
func SaveWithTx(tx *sql.Tx, domain interface{}, table string) (string, error) {

	_, err := resolveMappingModel(domain)

	if err != nil {
		return "", err
	}

	return "", nil
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

func resolveMappingModel(domain interface{}) (*mappingModel, error) {

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
		IDField: fld,
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

	logging.LogJSON(model)

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
