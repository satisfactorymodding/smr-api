// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/satisfactorymodding/smr-api/generated/ent/guide"
	"github.com/satisfactorymodding/smr-api/generated/ent/guidetag"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/tag"
)

// GuideTagQuery is the builder for querying GuideTag entities.
type GuideTagQuery struct {
	config
	ctx        *QueryContext
	order      []guidetag.OrderOption
	inters     []Interceptor
	predicates []predicate.GuideTag
	withGuide  *GuideQuery
	withTag    *TagQuery
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the GuideTagQuery builder.
func (gtq *GuideTagQuery) Where(ps ...predicate.GuideTag) *GuideTagQuery {
	gtq.predicates = append(gtq.predicates, ps...)
	return gtq
}

// Limit the number of records to be returned by this query.
func (gtq *GuideTagQuery) Limit(limit int) *GuideTagQuery {
	gtq.ctx.Limit = &limit
	return gtq
}

// Offset to start from.
func (gtq *GuideTagQuery) Offset(offset int) *GuideTagQuery {
	gtq.ctx.Offset = &offset
	return gtq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (gtq *GuideTagQuery) Unique(unique bool) *GuideTagQuery {
	gtq.ctx.Unique = &unique
	return gtq
}

// Order specifies how the records should be ordered.
func (gtq *GuideTagQuery) Order(o ...guidetag.OrderOption) *GuideTagQuery {
	gtq.order = append(gtq.order, o...)
	return gtq
}

// QueryGuide chains the current query on the "guide" edge.
func (gtq *GuideTagQuery) QueryGuide() *GuideQuery {
	query := (&GuideClient{config: gtq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gtq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gtq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guidetag.Table, guidetag.GuideColumn, selector),
			sqlgraph.To(guide.Table, guide.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, guidetag.GuideTable, guidetag.GuideColumn),
		)
		fromU = sqlgraph.SetNeighbors(gtq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTag chains the current query on the "tag" edge.
func (gtq *GuideTagQuery) QueryTag() *TagQuery {
	query := (&TagClient{config: gtq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := gtq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := gtq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(guidetag.Table, guidetag.TagColumn, selector),
			sqlgraph.To(tag.Table, tag.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, guidetag.TagTable, guidetag.TagColumn),
		)
		fromU = sqlgraph.SetNeighbors(gtq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first GuideTag entity from the query.
// Returns a *NotFoundError when no GuideTag was found.
func (gtq *GuideTagQuery) First(ctx context.Context) (*GuideTag, error) {
	nodes, err := gtq.Limit(1).All(setContextOp(ctx, gtq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{guidetag.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (gtq *GuideTagQuery) FirstX(ctx context.Context) *GuideTag {
	node, err := gtq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// Only returns a single GuideTag entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one GuideTag entity is found.
// Returns a *NotFoundError when no GuideTag entities are found.
func (gtq *GuideTagQuery) Only(ctx context.Context) (*GuideTag, error) {
	nodes, err := gtq.Limit(2).All(setContextOp(ctx, gtq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{guidetag.Label}
	default:
		return nil, &NotSingularError{guidetag.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (gtq *GuideTagQuery) OnlyX(ctx context.Context) *GuideTag {
	node, err := gtq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// All executes the query and returns a list of GuideTags.
func (gtq *GuideTagQuery) All(ctx context.Context) ([]*GuideTag, error) {
	ctx = setContextOp(ctx, gtq.ctx, ent.OpQueryAll)
	if err := gtq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*GuideTag, *GuideTagQuery]()
	return withInterceptors[[]*GuideTag](ctx, gtq, qr, gtq.inters)
}

// AllX is like All, but panics if an error occurs.
func (gtq *GuideTagQuery) AllX(ctx context.Context) []*GuideTag {
	nodes, err := gtq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// Count returns the count of the given query.
func (gtq *GuideTagQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, gtq.ctx, ent.OpQueryCount)
	if err := gtq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, gtq, querierCount[*GuideTagQuery](), gtq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (gtq *GuideTagQuery) CountX(ctx context.Context) int {
	count, err := gtq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (gtq *GuideTagQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, gtq.ctx, ent.OpQueryExist)
	switch _, err := gtq.First(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (gtq *GuideTagQuery) ExistX(ctx context.Context) bool {
	exist, err := gtq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the GuideTagQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (gtq *GuideTagQuery) Clone() *GuideTagQuery {
	if gtq == nil {
		return nil
	}
	return &GuideTagQuery{
		config:     gtq.config,
		ctx:        gtq.ctx.Clone(),
		order:      append([]guidetag.OrderOption{}, gtq.order...),
		inters:     append([]Interceptor{}, gtq.inters...),
		predicates: append([]predicate.GuideTag{}, gtq.predicates...),
		withGuide:  gtq.withGuide.Clone(),
		withTag:    gtq.withTag.Clone(),
		// clone intermediate query.
		sql:  gtq.sql.Clone(),
		path: gtq.path,
	}
}

// WithGuide tells the query-builder to eager-load the nodes that are connected to
// the "guide" edge. The optional arguments are used to configure the query builder of the edge.
func (gtq *GuideTagQuery) WithGuide(opts ...func(*GuideQuery)) *GuideTagQuery {
	query := (&GuideClient{config: gtq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	gtq.withGuide = query
	return gtq
}

// WithTag tells the query-builder to eager-load the nodes that are connected to
// the "tag" edge. The optional arguments are used to configure the query builder of the edge.
func (gtq *GuideTagQuery) WithTag(opts ...func(*TagQuery)) *GuideTagQuery {
	query := (&TagClient{config: gtq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	gtq.withTag = query
	return gtq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		GuideID string `json:"guide_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.GuideTag.Query().
//		GroupBy(guidetag.FieldGuideID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (gtq *GuideTagQuery) GroupBy(field string, fields ...string) *GuideTagGroupBy {
	gtq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &GuideTagGroupBy{build: gtq}
	grbuild.flds = &gtq.ctx.Fields
	grbuild.label = guidetag.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		GuideID string `json:"guide_id,omitempty"`
//	}
//
//	client.GuideTag.Query().
//		Select(guidetag.FieldGuideID).
//		Scan(ctx, &v)
func (gtq *GuideTagQuery) Select(fields ...string) *GuideTagSelect {
	gtq.ctx.Fields = append(gtq.ctx.Fields, fields...)
	sbuild := &GuideTagSelect{GuideTagQuery: gtq}
	sbuild.label = guidetag.Label
	sbuild.flds, sbuild.scan = &gtq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a GuideTagSelect configured with the given aggregations.
func (gtq *GuideTagQuery) Aggregate(fns ...AggregateFunc) *GuideTagSelect {
	return gtq.Select().Aggregate(fns...)
}

func (gtq *GuideTagQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range gtq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, gtq); err != nil {
				return err
			}
		}
	}
	for _, f := range gtq.ctx.Fields {
		if !guidetag.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if gtq.path != nil {
		prev, err := gtq.path(ctx)
		if err != nil {
			return err
		}
		gtq.sql = prev
	}
	return nil
}

func (gtq *GuideTagQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*GuideTag, error) {
	var (
		nodes       = []*GuideTag{}
		_spec       = gtq.querySpec()
		loadedTypes = [2]bool{
			gtq.withGuide != nil,
			gtq.withTag != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*GuideTag).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &GuideTag{config: gtq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(gtq.modifiers) > 0 {
		_spec.Modifiers = gtq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, gtq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := gtq.withGuide; query != nil {
		if err := gtq.loadGuide(ctx, query, nodes, nil,
			func(n *GuideTag, e *Guide) { n.Edges.Guide = e }); err != nil {
			return nil, err
		}
	}
	if query := gtq.withTag; query != nil {
		if err := gtq.loadTag(ctx, query, nodes, nil,
			func(n *GuideTag, e *Tag) { n.Edges.Tag = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (gtq *GuideTagQuery) loadGuide(ctx context.Context, query *GuideQuery, nodes []*GuideTag, init func(*GuideTag), assign func(*GuideTag, *Guide)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*GuideTag)
	for i := range nodes {
		fk := nodes[i].GuideID
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(guide.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "guide_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (gtq *GuideTagQuery) loadTag(ctx context.Context, query *TagQuery, nodes []*GuideTag, init func(*GuideTag), assign func(*GuideTag, *Tag)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*GuideTag)
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

func (gtq *GuideTagQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := gtq.querySpec()
	if len(gtq.modifiers) > 0 {
		_spec.Modifiers = gtq.modifiers
	}
	_spec.Unique = false
	_spec.Node.Columns = nil
	return sqlgraph.CountNodes(ctx, gtq.driver, _spec)
}

func (gtq *GuideTagQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(guidetag.Table, guidetag.Columns, nil)
	_spec.From = gtq.sql
	if unique := gtq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if gtq.path != nil {
		_spec.Unique = true
	}
	if fields := gtq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		for i := range fields {
			_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
		}
		if gtq.withGuide != nil {
			_spec.Node.AddColumnOnce(guidetag.FieldGuideID)
		}
		if gtq.withTag != nil {
			_spec.Node.AddColumnOnce(guidetag.FieldTagID)
		}
	}
	if ps := gtq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := gtq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := gtq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := gtq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (gtq *GuideTagQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(gtq.driver.Dialect())
	t1 := builder.Table(guidetag.Table)
	columns := gtq.ctx.Fields
	if len(columns) == 0 {
		columns = guidetag.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if gtq.sql != nil {
		selector = gtq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if gtq.ctx.Unique != nil && *gtq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range gtq.modifiers {
		m(selector)
	}
	for _, p := range gtq.predicates {
		p(selector)
	}
	for _, p := range gtq.order {
		p(selector)
	}
	if offset := gtq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := gtq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (gtq *GuideTagQuery) Modify(modifiers ...func(s *sql.Selector)) *GuideTagSelect {
	gtq.modifiers = append(gtq.modifiers, modifiers...)
	return gtq.Select()
}

// GuideTagGroupBy is the group-by builder for GuideTag entities.
type GuideTagGroupBy struct {
	selector
	build *GuideTagQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (gtgb *GuideTagGroupBy) Aggregate(fns ...AggregateFunc) *GuideTagGroupBy {
	gtgb.fns = append(gtgb.fns, fns...)
	return gtgb
}

// Scan applies the selector query and scans the result into the given value.
func (gtgb *GuideTagGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gtgb.build.ctx, ent.OpQueryGroupBy)
	if err := gtgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuideTagQuery, *GuideTagGroupBy](ctx, gtgb.build, gtgb, gtgb.build.inters, v)
}

func (gtgb *GuideTagGroupBy) sqlScan(ctx context.Context, root *GuideTagQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(gtgb.fns))
	for _, fn := range gtgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*gtgb.flds)+len(gtgb.fns))
		for _, f := range *gtgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*gtgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gtgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// GuideTagSelect is the builder for selecting fields of GuideTag entities.
type GuideTagSelect struct {
	*GuideTagQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (gts *GuideTagSelect) Aggregate(fns ...AggregateFunc) *GuideTagSelect {
	gts.fns = append(gts.fns, fns...)
	return gts
}

// Scan applies the selector query and scans the result into the given value.
func (gts *GuideTagSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, gts.ctx, ent.OpQuerySelect)
	if err := gts.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*GuideTagQuery, *GuideTagSelect](ctx, gts.GuideTagQuery, gts, gts.inters, v)
}

func (gts *GuideTagSelect) sqlScan(ctx context.Context, root *GuideTagQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(gts.fns))
	for _, fn := range gts.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*gts.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := gts.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (gts *GuideTagSelect) Modify(modifiers ...func(s *sql.Selector)) *GuideTagSelect {
	gts.modifiers = append(gts.modifiers, modifiers...)
	return gts
}
