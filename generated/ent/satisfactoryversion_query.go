// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/satisfactoryversion"
)

// SatisfactoryVersionQuery is the builder for querying SatisfactoryVersion entities.
type SatisfactoryVersionQuery struct {
	config
	ctx        *QueryContext
	order      []satisfactoryversion.OrderOption
	inters     []Interceptor
	predicates []predicate.SatisfactoryVersion
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SatisfactoryVersionQuery builder.
func (svq *SatisfactoryVersionQuery) Where(ps ...predicate.SatisfactoryVersion) *SatisfactoryVersionQuery {
	svq.predicates = append(svq.predicates, ps...)
	return svq
}

// Limit the number of records to be returned by this query.
func (svq *SatisfactoryVersionQuery) Limit(limit int) *SatisfactoryVersionQuery {
	svq.ctx.Limit = &limit
	return svq
}

// Offset to start from.
func (svq *SatisfactoryVersionQuery) Offset(offset int) *SatisfactoryVersionQuery {
	svq.ctx.Offset = &offset
	return svq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (svq *SatisfactoryVersionQuery) Unique(unique bool) *SatisfactoryVersionQuery {
	svq.ctx.Unique = &unique
	return svq
}

// Order specifies how the records should be ordered.
func (svq *SatisfactoryVersionQuery) Order(o ...satisfactoryversion.OrderOption) *SatisfactoryVersionQuery {
	svq.order = append(svq.order, o...)
	return svq
}

// First returns the first SatisfactoryVersion entity from the query.
// Returns a *NotFoundError when no SatisfactoryVersion was found.
func (svq *SatisfactoryVersionQuery) First(ctx context.Context) (*SatisfactoryVersion, error) {
	nodes, err := svq.Limit(1).All(setContextOp(ctx, svq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{satisfactoryversion.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) FirstX(ctx context.Context) *SatisfactoryVersion {
	node, err := svq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first SatisfactoryVersion ID from the query.
// Returns a *NotFoundError when no SatisfactoryVersion ID was found.
func (svq *SatisfactoryVersionQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = svq.Limit(1).IDs(setContextOp(ctx, svq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{satisfactoryversion.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) FirstIDX(ctx context.Context) string {
	id, err := svq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single SatisfactoryVersion entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one SatisfactoryVersion entity is found.
// Returns a *NotFoundError when no SatisfactoryVersion entities are found.
func (svq *SatisfactoryVersionQuery) Only(ctx context.Context) (*SatisfactoryVersion, error) {
	nodes, err := svq.Limit(2).All(setContextOp(ctx, svq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{satisfactoryversion.Label}
	default:
		return nil, &NotSingularError{satisfactoryversion.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) OnlyX(ctx context.Context) *SatisfactoryVersion {
	node, err := svq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only SatisfactoryVersion ID in the query.
// Returns a *NotSingularError when more than one SatisfactoryVersion ID is found.
// Returns a *NotFoundError when no entities are found.
func (svq *SatisfactoryVersionQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = svq.Limit(2).IDs(setContextOp(ctx, svq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{satisfactoryversion.Label}
	default:
		err = &NotSingularError{satisfactoryversion.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) OnlyIDX(ctx context.Context) string {
	id, err := svq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SatisfactoryVersions.
func (svq *SatisfactoryVersionQuery) All(ctx context.Context) ([]*SatisfactoryVersion, error) {
	ctx = setContextOp(ctx, svq.ctx, ent.OpQueryAll)
	if err := svq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*SatisfactoryVersion, *SatisfactoryVersionQuery]()
	return withInterceptors[[]*SatisfactoryVersion](ctx, svq, qr, svq.inters)
}

// AllX is like All, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) AllX(ctx context.Context) []*SatisfactoryVersion {
	nodes, err := svq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of SatisfactoryVersion IDs.
func (svq *SatisfactoryVersionQuery) IDs(ctx context.Context) (ids []string, err error) {
	if svq.ctx.Unique == nil && svq.path != nil {
		svq.Unique(true)
	}
	ctx = setContextOp(ctx, svq.ctx, ent.OpQueryIDs)
	if err = svq.Select(satisfactoryversion.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) IDsX(ctx context.Context) []string {
	ids, err := svq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (svq *SatisfactoryVersionQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, svq.ctx, ent.OpQueryCount)
	if err := svq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, svq, querierCount[*SatisfactoryVersionQuery](), svq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) CountX(ctx context.Context) int {
	count, err := svq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (svq *SatisfactoryVersionQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, svq.ctx, ent.OpQueryExist)
	switch _, err := svq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (svq *SatisfactoryVersionQuery) ExistX(ctx context.Context) bool {
	exist, err := svq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SatisfactoryVersionQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (svq *SatisfactoryVersionQuery) Clone() *SatisfactoryVersionQuery {
	if svq == nil {
		return nil
	}
	return &SatisfactoryVersionQuery{
		config:     svq.config,
		ctx:        svq.ctx.Clone(),
		order:      append([]satisfactoryversion.OrderOption{}, svq.order...),
		inters:     append([]Interceptor{}, svq.inters...),
		predicates: append([]predicate.SatisfactoryVersion{}, svq.predicates...),
		// clone intermediate query.
		sql:  svq.sql.Clone(),
		path: svq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Version int `json:"version,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.SatisfactoryVersion.Query().
//		GroupBy(satisfactoryversion.FieldVersion).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (svq *SatisfactoryVersionQuery) GroupBy(field string, fields ...string) *SatisfactoryVersionGroupBy {
	svq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &SatisfactoryVersionGroupBy{build: svq}
	grbuild.flds = &svq.ctx.Fields
	grbuild.label = satisfactoryversion.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Version int `json:"version,omitempty"`
//	}
//
//	client.SatisfactoryVersion.Query().
//		Select(satisfactoryversion.FieldVersion).
//		Scan(ctx, &v)
func (svq *SatisfactoryVersionQuery) Select(fields ...string) *SatisfactoryVersionSelect {
	svq.ctx.Fields = append(svq.ctx.Fields, fields...)
	sbuild := &SatisfactoryVersionSelect{SatisfactoryVersionQuery: svq}
	sbuild.label = satisfactoryversion.Label
	sbuild.flds, sbuild.scan = &svq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a SatisfactoryVersionSelect configured with the given aggregations.
func (svq *SatisfactoryVersionQuery) Aggregate(fns ...AggregateFunc) *SatisfactoryVersionSelect {
	return svq.Select().Aggregate(fns...)
}

func (svq *SatisfactoryVersionQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range svq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, svq); err != nil {
				return err
			}
		}
	}
	for _, f := range svq.ctx.Fields {
		if !satisfactoryversion.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if svq.path != nil {
		prev, err := svq.path(ctx)
		if err != nil {
			return err
		}
		svq.sql = prev
	}
	return nil
}

func (svq *SatisfactoryVersionQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*SatisfactoryVersion, error) {
	var (
		nodes = []*SatisfactoryVersion{}
		_spec = svq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*SatisfactoryVersion).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &SatisfactoryVersion{config: svq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(svq.modifiers) > 0 {
		_spec.Modifiers = svq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, svq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (svq *SatisfactoryVersionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := svq.querySpec()
	if len(svq.modifiers) > 0 {
		_spec.Modifiers = svq.modifiers
	}
	_spec.Node.Columns = svq.ctx.Fields
	if len(svq.ctx.Fields) > 0 {
		_spec.Unique = svq.ctx.Unique != nil && *svq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, svq.driver, _spec)
}

func (svq *SatisfactoryVersionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(satisfactoryversion.Table, satisfactoryversion.Columns, sqlgraph.NewFieldSpec(satisfactoryversion.FieldID, field.TypeString))
	_spec.From = svq.sql
	if unique := svq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if svq.path != nil {
		_spec.Unique = true
	}
	if fields := svq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, satisfactoryversion.FieldID)
		for i := range fields {
			if fields[i] != satisfactoryversion.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := svq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := svq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := svq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := svq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (svq *SatisfactoryVersionQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(svq.driver.Dialect())
	t1 := builder.Table(satisfactoryversion.Table)
	columns := svq.ctx.Fields
	if len(columns) == 0 {
		columns = satisfactoryversion.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if svq.sql != nil {
		selector = svq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if svq.ctx.Unique != nil && *svq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range svq.modifiers {
		m(selector)
	}
	for _, p := range svq.predicates {
		p(selector)
	}
	for _, p := range svq.order {
		p(selector)
	}
	if offset := svq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := svq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (svq *SatisfactoryVersionQuery) Modify(modifiers ...func(s *sql.Selector)) *SatisfactoryVersionSelect {
	svq.modifiers = append(svq.modifiers, modifiers...)
	return svq.Select()
}

// SatisfactoryVersionGroupBy is the group-by builder for SatisfactoryVersion entities.
type SatisfactoryVersionGroupBy struct {
	selector
	build *SatisfactoryVersionQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (svgb *SatisfactoryVersionGroupBy) Aggregate(fns ...AggregateFunc) *SatisfactoryVersionGroupBy {
	svgb.fns = append(svgb.fns, fns...)
	return svgb
}

// Scan applies the selector query and scans the result into the given value.
func (svgb *SatisfactoryVersionGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, svgb.build.ctx, ent.OpQueryGroupBy)
	if err := svgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SatisfactoryVersionQuery, *SatisfactoryVersionGroupBy](ctx, svgb.build, svgb, svgb.build.inters, v)
}

func (svgb *SatisfactoryVersionGroupBy) sqlScan(ctx context.Context, root *SatisfactoryVersionQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(svgb.fns))
	for _, fn := range svgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*svgb.flds)+len(svgb.fns))
		for _, f := range *svgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*svgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := svgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// SatisfactoryVersionSelect is the builder for selecting fields of SatisfactoryVersion entities.
type SatisfactoryVersionSelect struct {
	*SatisfactoryVersionQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (svs *SatisfactoryVersionSelect) Aggregate(fns ...AggregateFunc) *SatisfactoryVersionSelect {
	svs.fns = append(svs.fns, fns...)
	return svs
}

// Scan applies the selector query and scans the result into the given value.
func (svs *SatisfactoryVersionSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, svs.ctx, ent.OpQuerySelect)
	if err := svs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SatisfactoryVersionQuery, *SatisfactoryVersionSelect](ctx, svs.SatisfactoryVersionQuery, svs, svs.inters, v)
}

func (svs *SatisfactoryVersionSelect) sqlScan(ctx context.Context, root *SatisfactoryVersionQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(svs.fns))
	for _, fn := range svs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*svs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := svs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (svs *SatisfactoryVersionSelect) Modify(modifiers ...func(s *sql.Selector)) *SatisfactoryVersionSelect {
	svs.modifiers = append(svs.modifiers, modifiers...)
	return svs
}
