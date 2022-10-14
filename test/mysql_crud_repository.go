package test

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/bloock/go-kit/client"
	"github.com/huandu/go-sqlbuilder"
)

type MysqlCrudRepository struct {
	client    *client.MysqlClient
	table     string
	sqlStruct *sqlbuilder.Struct
}

func NewMysqlCrudRepository(client *client.MysqlClient, table string, sqlStruct *sqlbuilder.Struct) MysqlCrudRepository {
	return MysqlCrudRepository{
		client:    client,
		table:     table,
		sqlStruct: sqlStruct,
	}
}

func (r MysqlCrudRepository) Create(value interface{}) error {
	query, args := r.sqlStruct.InsertInto(r.table, value).Build()

	_, err := r.client.DB().Exec(query, args...)
	return err
}

func (r MysqlCrudRepository) List(res interface{}) error {
	sb := r.sqlStruct.SelectFrom(r.table)
	query, args := sb.Build()

	rows, err := r.client.DB().Query(query, args...)
	if err != nil {
		return err
	}

	err = r.decodeSlice(rows, res)
	if err != nil {
		return err
	}

	return nil
}

func (r MysqlCrudRepository) Retrieve(id int, columnName string, res interface{}) error {
	sb := r.sqlStruct.SelectFrom(r.table)
	query, args := sb.Where(sb.Equal(columnName, id)).Build()

	row := r.client.DB().QueryRow(query, args...)
	if row.Err() != nil {
		return row.Err()
	}

	err := r.decodeObject(row, res)
	if err != nil {
		return err
	}

	return nil
}

func (r MysqlCrudRepository) Update(id int, columnName string, value interface{}) error {
	ub := r.sqlStruct.Update(r.table, value)
	ub.Where(ub.Equal(columnName, id))
	query, args := ub.Build()

	_, err := r.client.DB().Exec(query, args...)
	return err
}

func (r MysqlCrudRepository) Delete(id int, columnName string) error {
	db := r.sqlStruct.DeleteFrom(r.table)
	query, args := db.Where(db.Equal(columnName, id)).Build()

	_, err := r.client.DB().Exec(query, args...)
	return err
}

func (r MysqlCrudRepository) Truncate() error {
	query := fmt.Sprintf("SET FOREIGN_KEY_CHECKS=0; TRUNCATE %s; SET FOREIGN_KEY_CHECKS = 1", r.table)

	if _, err := r.client.DB().Exec(query); err != nil {
		return err
	}

	return nil
}

func (r MysqlCrudRepository) decodeSlice(rows *sql.Rows, res interface{}) error {
	resultsVal := reflect.ValueOf(res)
	if resultsVal.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultsVal.Kind())
	}

	sliceVal := resultsVal.Elem()
	if sliceVal.Kind() == reflect.Interface {
		sliceVal = sliceVal.Elem()
	}

	if sliceVal.Kind() != reflect.Slice {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a pointer to %s", sliceVal.Kind())
	}

	elementType := sliceVal.Type().Elem()

	index := 0

	for rows.Next() {
		if sliceVal.Len() == index {
			// slice is full
			newElem := reflect.New(elementType)
			sliceVal = reflect.Append(sliceVal, newElem.Elem())
			sliceVal = sliceVal.Slice(0, sliceVal.Cap())
		}

		currElem := sliceVal.Index(index).Addr().Interface()
		err := rows.Scan(r.sqlStruct.Addr(&currElem)...)
		if err != nil {
			return err
		}

		index++
	}

	resultsVal.Elem().Set(sliceVal.Slice(0, index))
	return nil
}

func (r MysqlCrudRepository) decodeObject(row *sql.Row, res interface{}) error {
	resultVal := reflect.ValueOf(res)
	if resultVal.Kind() != reflect.Ptr {
		return fmt.Errorf("results argument must be a pointer to a slice, but was a %s", resultVal.Kind())
	}

	objectVal := resultVal.Elem()
	if objectVal.Kind() == reflect.Interface {
		objectVal = objectVal.Elem()
	}

	if objectVal.Kind() != reflect.Struct {
		return fmt.Errorf("results argument must be a pointer to a struct, but was a pointer to %s", objectVal.Kind())
	}

	elem := objectVal.Addr().Interface()

	err := row.Scan(r.sqlStruct.Addr(&elem)...)
	if err != nil {
		return err
	}

	return nil
}
