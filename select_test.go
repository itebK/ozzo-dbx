// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectQuery(t *testing.T) {
	db := getDB()

	// minimal select query
	q := db.Select().From("users").Build()
	expected := "SELECT * FROM `users`"
	assert.Equal(t, q.SQL(), expected, "t1")
	assert.Equal(t, len(q.Params()), 0, "t2")

	// a full select query
	q = db.Select("id", "name").
		AndSelect("age").
		Distinct(true).
		SelectOption("CALC").
		From("users").
		Where(NewExp("age>30")).
		AndWhere(NewExp("status=1")).
		OrWhere(NewExp("type=2")).
		InnerJoin("profile", NewExp("user.id=profile.id")).
		OrderBy("age DESC", "type").
		AndOrderBy("id").
		GroupBy("id").
		AndGroupBy("age").
		Having(NewExp("id>10")).
		AndHaving(NewExp("id<20")).
		OrHaving(NewExp("type=3")).
		Limit(10).
		Offset(20).
		Bind(Params{"id": 1}).
		AndBind(Params{"age": 30}).
		Build()

	expected = "SELECT DISTINCT CALC `id`, `name`, `age` FROM `users` INNER JOIN `profile` ON user.id=profile.id WHERE ((age>30) AND (status=1)) OR (type=2) GROUP BY `id`, `age` HAVING ((id>10) AND (id<20)) OR (type=3) ORDER BY `age` DESC, `type`, `id` LIMIT 10 OFFSET 20"
	assert.Equal(t, q.SQL(), expected, "t3")
	assert.Equal(t, len(q.Params()), 2, "t4")

	// union
	q1 := db.Select().From("users").Build()
	q2 := db.Select().From("posts").Build()
	q = db.Select().From("profiles").Union(q1).UnionAll(q2).Build()
	expected = "(SELECT * FROM `profiles`) UNION (SELECT * FROM `users`) UNION ALL (SELECT * FROM `posts`)"
	assert.Equal(t, q.SQL(), expected, "t5")
}

func TestSelectQuery_Data(t *testing.T) {
	db := getPreparedDB()
	defer db.Close()

	q := db.Select("id", "email").From("customer").OrderBy("id")

	var customer Customer
	q.One(&customer)
	assert.Equal(t, customer.Email, "user1@example.com", "customer.Email")

	var customers []Customer
	q.All(&customers)
	assert.Equal(t, len(customers), 3, "len(customers)")

	rows, _ := q.Rows()
	customer.Email = ""
	rows.one(&customer)
	assert.Equal(t, customer.Email, "user1@example.com", "customer.Email")

	var id, email string
	q.Row(&id, &email)
	assert.Equal(t, id, "1", "id")
	assert.Equal(t, email, "user1@example.com", "email")
}
