package db

import "entgo.io/ent/dialect/sql"

func OrderToOrder(order string) sql.OrderTermOption {
	if order == "asc" || order == "ascending" {
		return sql.OrderAsc()
	}

	return sql.OrderDesc()
}
