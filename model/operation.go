package model

import "reflect"

func Insert(tb interface{}) bool {
	n, err := engine.Insert(tb)
	if err != nil || n == 0 {
		return false
	}
	return true
}

func Update(tb interface{}) bool {
	n, err := engine.Update(tb)
	if err != nil || n == 0 {
		return false
	}
	return true
}

func Find(tb interface{}) bool {
	err := engine.Find(tb)
	if err != nil {
		return false
	}
	return true
}

func Get(tb interface{}) bool {
	ok, err := engine.Get(tb)
	if err != nil || !ok {
		return false
	}
	return true
}

func Delete(tb interface{}) bool {
	n, err := engine.Delete(tb)
	if err != nil || n == 0 {
		return false
	}
	return true
}

func In(tb interface{}, field string, conditions []interface{}) []interface{} {
	rows, err := engine.In(field, conditions...).Desc(field).Rows(tb)
	if err != nil {
		return nil
	}
	items := make([]interface{}, 0, len(conditions))
	defer rows.Close()
	for rows.Next() {
		item := newInterface(tb)
		err = rows.Scan(item)
		if err != nil {
			break
		}
		items = append(items, item)
	}
	if err != nil {
		return nil
	}
	return items
}

func newInterface(tb interface{}) interface{} {
	t := reflect.ValueOf(tb).Type()
	return reflect.New(t).Interface()
}
