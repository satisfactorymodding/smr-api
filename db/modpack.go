package db

import (
	"strings"

	"entgo.io/ent/dialect/sql"

	"github.com/satisfactorymodding/smr-api/generated"
	"github.com/satisfactorymodding/smr-api/generated/ent"
	"github.com/satisfactorymodding/smr-api/generated/ent/modpack"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
	"github.com/satisfactorymodding/smr-api/models"
)

func ConvertModpackFilter(query *ent.ModpackQuery, filter *models.ModpackFilter, count bool) *ent.ModpackQuery {
	if len(filter.IDs) > 0 {
		query = query.Where(modpack.IDIn(filter.IDs...))
	} else if filter != nil {
		if !count {
			query = query.
				Limit(*filter.Limit).
				Offset(*filter.Offset)
		}

		if filter.OrderBy != nil && *filter.OrderBy != generated.ModpackFieldsSearch {
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
				join.From(sql.Table(modpack.Table)).As("t1")
				s.Join(join).On(s.C(modpack.FieldID), join.C("id"))
			})

			query = query.Where(func(s *sql.Selector) {
				s.Where(sql.ExprP(`"t1"."s" > 0.2`))
			})

			if !count && *filter.OrderBy == generated.ModpackFieldsSearch {
				query = query.Order(func(s *sql.Selector) {
					s.OrderExpr(sql.ExprP(`"t1"."s" DESC`))
				})
			}
		}

		if filter.Hidden == nil || !(*filter.Hidden) {
			query = query.Where(modpack.Hidden(false))
		}

		if len(filter.TagIDs) > 0 {
			query = query.Where(modpack.HasTagsWith(tag.IDIn(filter.TagIDs...)))
		}
	}
	return query
}
