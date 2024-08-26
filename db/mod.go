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

	if len(filter.IDs) > 0 {
		query = query.Where(mod.IDIn(filter.IDs...))
	} else if len(filter.References) > 0 {
		query = query.Where(mod.ModReferenceIn(filter.References...))
	} else if filter != nil {
		if !count {
			query = query.
				Limit(*filter.Limit).
				Offset(*filter.Offset)
		}

		if filter.OrderBy != nil && *filter.OrderBy != generated.ModFieldsSearch {
			if string(*filter.OrderBy) == "last_version_date" {
				query = query.Order(func(s *sql.Selector) {
					s.OrderExpr(sql.ExprP("case when last_version_date is null then 1 else 0 end"))
				})
			}
			query = query.Order(sql.OrderByField(
				filter.OrderBy.String(),
				OrderToOrder(filter.Order.String()),
			).ToFunc())
		}

		if filter.Search != nil && *filter.Search != "" {
			cleanSearch := strings.ReplaceAll(strings.TrimSpace(*filter.Search), " ", " & ")

			query = query.Where(func(s *sql.Selector) {
				join := sql.Select("id")
				join = join.AppendSelectExprAs(
					sql.P(func(builder *sql.Builder) {
						builder.WriteString("similarity(name, ").Arg(cleanSearch).WriteString(") * 2").
							WriteString(" + ").
							WriteString("similarity(short_description, ").Arg(cleanSearch).WriteString(")").
							WriteString(" + ").
							WriteString("similarity(full_description, ").Arg(cleanSearch).WriteString(") * 0.5")
					}),
					"s",
				)
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

		if len(filter.TagIDs) > 0 {
			query = query.Where(mod.HasModTagsWith(modtag.TagIDIn(filter.TagIDs...)))
		}
	}

	query = query.Where(mod.Approved(!unapproved), mod.Denied(false))

	return query
}
