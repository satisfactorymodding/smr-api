package db

import (
	"strings"

	"entgo.io/ent/dialect/sql"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/modtag"
	"github.com/satisfactorymodding/smr-api/models"
)

func ConvertModFilter(query *ent.ModQuery, filter *models.ModFilter, count bool, unapproved bool) *ent.ModQuery {
	query = query.WithTags()

	if len(filter.Ids) > 0 {
		query = query.Where(mod.IDIn(filter.Ids...))
	} else if len(filter.References) > 0 {
		query = query.Where(mod.ModReferenceIn(filter.References...))
	} else if filter != nil {
		query = query.
			Limit(*filter.Limit).
			Offset(*filter.Offset)

		if *filter.OrderBy != generated.ModFieldsSearch {
			if string(*filter.OrderBy) == "last_version_date" {
				query = query.Modify(func(s *sql.Selector) {
					s.OrderExpr(sql.ExprP("case when last_version_date is null then 1 else 0 end, last_version_date"))
				}).Clone()
			} else {
				query = query.Order(sql.OrderByField(
					filter.OrderBy.String(),
					OrderToOrder(filter.Order.String()),
				).ToFunc())
			}
		}

		if filter.Search != nil && *filter.Search != "" {
			cleanSearch := strings.ReplaceAll(strings.TrimSpace(*filter.Search), " ", " & ")

			query = query.Where(func(s *sql.Selector) {
				join := sql.SelectExpr(sql.ExprP("id, (similarity(name, ?) * 2 + similarity(short_description, ?) + similarity(full_description, ?) * 0.5) as s", cleanSearch, cleanSearch, cleanSearch))
				join.From(sql.Table(mod.Table)).As("t1")
				s.Join(join).On(s.C(mod.FieldID), join.C("id"))
			})

			query = query.Where(func(s *sql.Selector) {
				s.Where(sql.ExprP(`"t1"."s" > 0.2`))
			})

			if !count && *filter.OrderBy == generated.ModFieldsSearch {
				query = query.Order(func(s *sql.Selector) {
					s.OrderExpr(sql.ExprP(`"t1"."s" DESC`))
				})
			}
		}

		if filter.Hidden == nil || !(*filter.Hidden) {
			query = query.Where(mod.Hidden(false))
		}

		if filter.TagIDs != nil && len(filter.TagIDs) > 0 {
			query = query.Where(func(s *sql.Selector) {
				t := sql.Table(modtag.Table)
				s.Join(t).OnP(sql.ExprP("mod_tags.tag_id in ? AND mod_tags.mod_id = mods.id", filter.TagIDs))
			})
		}
	}

	query = query.Where(mod.Approved(!unapproved), mod.Denied(false))

	return query
}
