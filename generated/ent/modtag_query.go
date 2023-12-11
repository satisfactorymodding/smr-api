// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/mod"
	"github.com/satisfactorymodding/smr-api/generated/ent/modtag"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
)

// ModTagQuery is the builder for querying ModTag entities.
type ModTagQuery struct {
	config
	ctx        *QueryContext
	order      []modtag.OrderOption
	inters     []Interceptor
	predicates []predicate.ModTag
	withMod    *ModQuery
	withTag    *TagQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ModTagQuery builder.
func (mtq *ModTagQuery) Where(ps ...predicate.ModTag) *ModTagQuery {
	mtq.predicates = append(mtq.predicates, ps...)
	return mtq
}

// Limit the number of records to be returned by this query.
func (mtq *ModTagQuery) Limit(limit int) *ModTagQuery {
	mtq.ctx.Limit = &limit
	return mtq
}

// Offset to start from.
func (mtq *ModTagQuery) Offset(offset int) *ModTagQuery {
	mtq.ctx.Offset = &offset
	return mtq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (mtq *ModTagQuery) Unique(unique bool) *ModTagQuery {
	mtq.ctx.Unique = &unique
	return mtq
}

// Order specifies how the records should be ordered.
func (mtq *ModTagQuery) Order(o ...modtag.OrderOption) *ModTagQuery {
	mtq.order = append(mtq.order, o...)
	return mtq
}

// QueryMod chains the current query on the "mod" edge.
func (mtq *ModTagQuery) QueryMod() *ModQuery {
	query := (&ModClient{config: mtq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mtq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mtq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(modtag.Table, modtag.ModColumn, selector),
			sqlgraph.To(mod.Table, mod.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, modtag.ModTable, modtag.ModColumn),
		)
		fromU = sqlgraph.SetNeighbors(mtq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTag chains the current query on the "tag" edge.
func (mtq *ModTagQuery) QueryTag() *TagQuery {
	query := (&TagClient{config: mtq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := mtq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := mtq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(modtag.Table, modtag.TagColumn, selector),
			sqlgraph.To(tag.Table, tag.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, modtag.TagTable, modtag.TagColumn),
		)
		fromU = sqlgraph.SetNeighbors(mtq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first ModTag entity from the query.
// Returns a *NotFoundError when no ModTag was found.
func (mtq *ModTagQuery) First(ctx context.Context) (*ModTag, error) {
	nodes, err := mtq.Limit(1).All(setContextOp(ctx, mtq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{modtag.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (mtq *ModTagQuery) FirstX(ctx context.Context) *ModTag {
	node, err := mtq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// Only returns a single ModTag entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one ModTag entity is found.
// Returns a *NotFoundError when no ModTag entities are found.
func (mtq *ModTagQuery) Only(ctx context.Context) (*ModTag, error) {
	nodes, err := mtq.Limit(2).All(setContextOp(ctx, mtq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{modtag.Label}
	default:
		return nil, &NotSingularError{modtag.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (mtq *ModTagQuery) OnlyX(ctx context.Context) *ModTag {
	node, err := mtq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// All executes the query and returns a list of ModTags.
func (mtq *ModTagQuery) All(ctx context.Context) ([]*ModTag, error) {
	ctx = setContextOp(ctx, mtq.ctx, "All")
	if err := mtq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*ModTag, *ModTagQuery]()
	return withInterceptors[[]*ModTag](ctx, mtq, qr, mtq.inters)
}

// AllX is like All, but panics if an error occurs.
func (mtq *ModTagQuery) AllX(ctx context.Context) []*ModTag {
	nodes, err := mtq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// Count returns the count of the given query.
func (mtq *ModTagQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, mtq.ctx, "Count")
	if err := mtq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, mtq, querierCount[*ModTagQuery](), mtq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (mtq *ModTagQuery) CountX(ctx context.Context) int {
	count, err := mtq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (mtq *ModTagQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, mtq.ctx, "Exist")
	switch _, err := mtq.First(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (mtq *ModTagQuery) ExistX(ctx context.Context) bool {
	exist, err := mtq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ModTagQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (mtq *ModTagQuery) Clone() *ModTagQuery {
	if mtq == nil {
		return nil
	}
	return &ModTagQuery{
		config:     mtq.config,
		ctx:        mtq.ctx.Clone(),
		order:      append([]modtag.OrderOption{}, mtq.order...),
		inters:     append([]Interceptor{}, mtq.inters...),
		predicates: append([]predicate.ModTag{}, mtq.predicates...),
		withMod:    mtq.withMod.Clone(),
		withTag:    mtq.withTag.Clone(),
		// clone intermediate query.
		sql:  mtq.sql.Clone(),
		path: mtq.path,
	}
}

// WithMod tells the query-builder to eager-load the nodes that are connected to
// the "mod" edge. The optional arguments are used to configure the query builder of the edge.
func (mtq *ModTagQuery) WithMod(opts ...func(*ModQuery)) *ModTagQuery {
	query := (&ModClient{config: mtq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mtq.withMod = query
	return mtq
}

// WithTag tells the query-builder to eager-load the nodes that are connected to
// the "tag" edge. The optional arguments are used to configure the query builder of the edge.
func (mtq *ModTagQuery) WithTag(opts ...func(*TagQuery)) *ModTagQuery {
	query := (&TagClient{config: mtq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	mtq.withTag = query
	return mtq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		ModID string `json:"mod_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.ModTag.Query().
//		GroupBy(modtag.FieldModID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (mtq *ModTagQuery) GroupBy(field string, fields ...string) *ModTagGroupBy {
	mtq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &ModTagGroupBy{build: mtq}
	grbuild.flds = &mtq.ctx.Fields
	grbuild.label = modtag.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		ModID string `json:"mod_id,omitempty"`
//	}
//
//	client.ModTag.Query().
//		Select(modtag.FieldModID).
//		Scan(ctx, &v)
func (mtq *ModTagQuery) Select(fields ...string) *ModTagSelect {
	mtq.ctx.Fields = append(mtq.ctx.Fields, fields...)
	sbuild := &ModTagSelect{ModTagQuery: mtq}
	sbuild.label = modtag.Label
	sbuild.flds, sbuild.scan = &mtq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ModTagSelect configured with the given aggregations.
func (mtq *ModTagQuery) Aggregate(fns ...AggregateFunc) *ModTagSelect {
	return mtq.Select().Aggregate(fns...)
}

func (mtq *ModTagQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range mtq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, mtq); err != nil {
				return err
			}
		}
	}
	for _, f := range mtq.ctx.Fields {
		if !modtag.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if mtq.path != nil {
		prev, err := mtq.path(ctx)
		if err != nil {
			return err
		}
		mtq.sql = prev
	}
	return nil
}

func (mtq *ModTagQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*ModTag, error) {
	var (
		nodes       = []*ModTag{}
		_spec       = mtq.querySpec()
		loadedTypes = [2]bool{
			mtq.withMod != nil,
			mtq.withTag != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*ModTag).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &ModTag{config: mtq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(mtq.modifiers) > 0 {
		_spec.Modifiers = mtq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, mtq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := mtq.withMod; query != nil {
		if err := mtq.loadMod(ctx, query, nodes, nil,
			func(n *ModTag, e *Mod) { n.Edges.Mod = e }); err != nil {
			return nil, err
		}
	}
	if query := mtq.withTag; query != nil {
		if err := mtq.loadTag(ctx, query, nodes, nil,
			func(n *ModTag, e *Tag) { n.Edges.Tag = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (mtq *ModTagQuery) loadMod(ctx context.Context, query *ModQuery, nodes []*ModTag, init func(*ModTag), assign func(*ModTag, *Mod)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*ModTag)
	for i := range nodes {
		fk := nodes[i].ModID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(mod.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "mod_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (mtq *ModTagQuery) loadTag(ctx context.Context, query *TagQuery, nodes []*ModTag, init func(*ModTag), assign func(*ModTag, *Tag)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*ModTag)
	for i := range nodes {
		fk := nodes[i].TagID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(tag.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "tag_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (mtq *ModTagQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := mtq.querySpec()
	if len(mtq.modifiers) > 0 {
		_spec.Modifiers = mtq.modifiers
	}
	_spec.Unique = false
	_spec.Node.Columns = nil
	return sqlgraph.CountNodes(ctx, mtq.driver, _spec)
}

func (mtq *ModTagQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(modtag.Table, modtag.Columns, nil)
	_spec.From = mtq.sql
	if unique := mtq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if mtq.path != nil {
		_spec.Unique = true
	}
	if fields := mtq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		for i := range fields {
			_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
		}
		if mtq.withMod != nil {
			_spec.Node.AddColumnOnce(modtag.FieldModID)
		}
		if mtq.withTag != nil {
			_spec.Node.AddColumnOnce(modtag.FieldTagID)
		}
	}
	if ps := mtq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := mtq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := mtq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := mtq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (mtq *ModTagQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(mtq.driver.Dialect())
	t1 := builder.Table(modtag.Table)
	columns := mtq.ctx.Fields
	if len(columns) == 0 {
		columns = modtag.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if mtq.sql != nil {
		selector = mtq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if mtq.ctx.Unique != nil && *mtq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range mtq.modifiers {
		m(selector)
	}
	for _, p := range mtq.predicates {
		p(selector)
	}
	for _, p := range mtq.order {
		p(selector)
	}
	if offset := mtq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := mtq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (mtq *ModTagQuery) Modify(modifiers ...func(s *sql.Selector)) *ModTagSelect {
	mtq.modifiers = append(mtq.modifiers, modifiers...)
	return mtq.Select()
}

// ModTagGroupBy is the group-by builder for ModTag entities.
type ModTagGroupBy struct {
	selector
	build *ModTagQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (mtgb *ModTagGroupBy) Aggregate(fns ...AggregateFunc) *ModTagGroupBy {
	mtgb.fns = append(mtgb.fns, fns...)
	return mtgb
}

// Scan applies the selector query and scans the result into the given value.
func (mtgb *ModTagGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mtgb.build.ctx, "GroupBy")
	if err := mtgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ModTagQuery, *ModTagGroupBy](ctx, mtgb.build, mtgb, mtgb.build.inters, v)
}

func (mtgb *ModTagGroupBy) sqlScan(ctx context.Context, root *ModTagQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(mtgb.fns))
	for _, fn := range mtgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*mtgb.flds)+len(mtgb.fns))
		for _, f := range *mtgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*mtgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mtgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ModTagSelect is the builder for selecting fields of ModTag entities.
type ModTagSelect struct {
	*ModTagQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (mts *ModTagSelect) Aggregate(fns ...AggregateFunc) *ModTagSelect {
	mts.fns = append(mts.fns, fns...)
	return mts
}

// Scan applies the selector query and scans the result into the given value.
func (mts *ModTagSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, mts.ctx, "Select")
	if err := mts.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ModTagQuery, *ModTagSelect](ctx, mts.ModTagQuery, mts, mts.inters, v)
}

func (mts *ModTagSelect) sqlScan(ctx context.Context, root *ModTagQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(mts.fns))
	for _, fn := range mts.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*mts.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := mts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (mts *ModTagSelect) Modify(modifiers ...func(s *sql.Selector)) *ModTagSelect {
	mts.modifiers = append(mts.modifiers, modifiers...)
	return mts
}