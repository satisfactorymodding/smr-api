// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/satisfactorymodding/smr-api/generated/ent/predicate"
	"github.com/satisfactorymodding/smr-api/generated/ent/user"
	"github.com/satisfactorymodding/smr-api/generated/ent/usersession"
)

// UserSessionQuery is the builder for querying UserSession entities.
type UserSessionQuery struct {
	config
	ctx        *QueryContext
	order      []usersession.OrderOption
	inters     []Interceptor
	predicates []predicate.UserSession
	withUser   *UserQuery
	withFKs    bool
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the UserSessionQuery builder.
func (usq *UserSessionQuery) Where(ps ...predicate.UserSession) *UserSessionQuery {
	usq.predicates = append(usq.predicates, ps...)
	return usq
}

// Limit the number of records to be returned by this query.
func (usq *UserSessionQuery) Limit(limit int) *UserSessionQuery {
	usq.ctx.Limit = &limit
	return usq
}

// Offset to start from.
func (usq *UserSessionQuery) Offset(offset int) *UserSessionQuery {
	usq.ctx.Offset = &offset
	return usq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (usq *UserSessionQuery) Unique(unique bool) *UserSessionQuery {
	usq.ctx.Unique = &unique
	return usq
}

// Order specifies how the records should be ordered.
func (usq *UserSessionQuery) Order(o ...usersession.OrderOption) *UserSessionQuery {
	usq.order = append(usq.order, o...)
	return usq
}

// QueryUser chains the current query on the "user" edge.
func (usq *UserSessionQuery) QueryUser() *UserQuery {
	query := (&UserClient{config: usq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := usq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := usq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(usersession.Table, usersession.FieldID, selector),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, usersession.UserTable, usersession.UserColumn),
		)
		fromU = sqlgraph.SetNeighbors(usq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first UserSession entity from the query.
// Returns a *NotFoundError when no UserSession was found.
func (usq *UserSessionQuery) First(ctx context.Context) (*UserSession, error) {
	nodes, err := usq.Limit(1).All(setContextOp(ctx, usq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{usersession.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (usq *UserSessionQuery) FirstX(ctx context.Context) *UserSession {
	node, err := usq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first UserSession ID from the query.
// Returns a *NotFoundError when no UserSession ID was found.
func (usq *UserSessionQuery) FirstID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = usq.Limit(1).IDs(setContextOp(ctx, usq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{usersession.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (usq *UserSessionQuery) FirstIDX(ctx context.Context) string {
	id, err := usq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single UserSession entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one UserSession entity is found.
// Returns a *NotFoundError when no UserSession entities are found.
func (usq *UserSessionQuery) Only(ctx context.Context) (*UserSession, error) {
	nodes, err := usq.Limit(2).All(setContextOp(ctx, usq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{usersession.Label}
	default:
		return nil, &NotSingularError{usersession.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (usq *UserSessionQuery) OnlyX(ctx context.Context) *UserSession {
	node, err := usq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only UserSession ID in the query.
// Returns a *NotSingularError when more than one UserSession ID is found.
// Returns a *NotFoundError when no entities are found.
func (usq *UserSessionQuery) OnlyID(ctx context.Context) (id string, err error) {
	var ids []string
	if ids, err = usq.Limit(2).IDs(setContextOp(ctx, usq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{usersession.Label}
	default:
		err = &NotSingularError{usersession.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (usq *UserSessionQuery) OnlyIDX(ctx context.Context) string {
	id, err := usq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of UserSessions.
func (usq *UserSessionQuery) All(ctx context.Context) ([]*UserSession, error) {
	ctx = setContextOp(ctx, usq.ctx, "All")
	if err := usq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*UserSession, *UserSessionQuery]()
	return withInterceptors[[]*UserSession](ctx, usq, qr, usq.inters)
}

// AllX is like All, but panics if an error occurs.
func (usq *UserSessionQuery) AllX(ctx context.Context) []*UserSession {
	nodes, err := usq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of UserSession IDs.
func (usq *UserSessionQuery) IDs(ctx context.Context) (ids []string, err error) {
	if usq.ctx.Unique == nil && usq.path != nil {
		usq.Unique(true)
	}
	ctx = setContextOp(ctx, usq.ctx, "IDs")
	if err = usq.Select(usersession.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (usq *UserSessionQuery) IDsX(ctx context.Context) []string {
	ids, err := usq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (usq *UserSessionQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, usq.ctx, "Count")
	if err := usq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, usq, querierCount[*UserSessionQuery](), usq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (usq *UserSessionQuery) CountX(ctx context.Context) int {
	count, err := usq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (usq *UserSessionQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, usq.ctx, "Exist")
	switch _, err := usq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (usq *UserSessionQuery) ExistX(ctx context.Context) bool {
	exist, err := usq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the UserSessionQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (usq *UserSessionQuery) Clone() *UserSessionQuery {
	if usq == nil {
		return nil
	}
	return &UserSessionQuery{
		config:     usq.config,
		ctx:        usq.ctx.Clone(),
		order:      append([]usersession.OrderOption{}, usq.order...),
		inters:     append([]Interceptor{}, usq.inters...),
		predicates: append([]predicate.UserSession{}, usq.predicates...),
		withUser:   usq.withUser.Clone(),
		// clone intermediate query.
		sql:  usq.sql.Clone(),
		path: usq.path,
	}
}

// WithUser tells the query-builder to eager-load the nodes that are connected to
// the "user" edge. The optional arguments are used to configure the query builder of the edge.
func (usq *UserSessionQuery) WithUser(opts ...func(*UserQuery)) *UserSessionQuery {
	query := (&UserClient{config: usq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	usq.withUser = query
	return usq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.UserSession.Query().
//		GroupBy(usersession.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (usq *UserSessionQuery) GroupBy(field string, fields ...string) *UserSessionGroupBy {
	usq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &UserSessionGroupBy{build: usq}
	grbuild.flds = &usq.ctx.Fields
	grbuild.label = usersession.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.UserSession.Query().
//		Select(usersession.FieldCreatedAt).
//		Scan(ctx, &v)
func (usq *UserSessionQuery) Select(fields ...string) *UserSessionSelect {
	usq.ctx.Fields = append(usq.ctx.Fields, fields...)
	sbuild := &UserSessionSelect{UserSessionQuery: usq}
	sbuild.label = usersession.Label
	sbuild.flds, sbuild.scan = &usq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a UserSessionSelect configured with the given aggregations.
func (usq *UserSessionQuery) Aggregate(fns ...AggregateFunc) *UserSessionSelect {
	return usq.Select().Aggregate(fns...)
}

func (usq *UserSessionQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range usq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, usq); err != nil {
				return err
			}
		}
	}
	for _, f := range usq.ctx.Fields {
		if !usersession.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if usq.path != nil {
		prev, err := usq.path(ctx)
		if err != nil {
			return err
		}
		usq.sql = prev
	}
	return nil
}

func (usq *UserSessionQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*UserSession, error) {
	var (
		nodes       = []*UserSession{}
		withFKs     = usq.withFKs
		_spec       = usq.querySpec()
		loadedTypes = [1]bool{
			usq.withUser != nil,
		}
	)
	if usq.withUser != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, usersession.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*UserSession).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &UserSession{config: usq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	if len(usq.modifiers) > 0 {
		_spec.Modifiers = usq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, usq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := usq.withUser; query != nil {
		if err := usq.loadUser(ctx, query, nodes, nil,
			func(n *UserSession, e *User) { n.Edges.User = e }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (usq *UserSessionQuery) loadUser(ctx context.Context, query *UserQuery, nodes []*UserSession, init func(*UserSession), assign func(*UserSession, *User)) error {
	ids := make([]string, 0, len(nodes))
	nodeids := make(map[string][]*UserSession)
	for i := range nodes {
		if nodes[i].user_id == nil {
			continue
		}
		fk := *nodes[i].user_id
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(user.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "user_id" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}

func (usq *UserSessionQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := usq.querySpec()
	if len(usq.modifiers) > 0 {
		_spec.Modifiers = usq.modifiers
	}
	_spec.Node.Columns = usq.ctx.Fields
	if len(usq.ctx.Fields) > 0 {
		_spec.Unique = usq.ctx.Unique != nil && *usq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, usq.driver, _spec)
}

func (usq *UserSessionQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(usersession.Table, usersession.Columns, sqlgraph.NewFieldSpec(usersession.FieldID, field.TypeString))
	_spec.From = usq.sql
	if unique := usq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if usq.path != nil {
		_spec.Unique = true
	}
	if fields := usq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, usersession.FieldID)
		for i := range fields {
			if fields[i] != usersession.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := usq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := usq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := usq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := usq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (usq *UserSessionQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(usq.driver.Dialect())
	t1 := builder.Table(usersession.Table)
	columns := usq.ctx.Fields
	if len(columns) == 0 {
		columns = usersession.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if usq.sql != nil {
		selector = usq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if usq.ctx.Unique != nil && *usq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range usq.modifiers {
		m(selector)
	}
	for _, p := range usq.predicates {
		p(selector)
	}
	for _, p := range usq.order {
		p(selector)
	}
	if offset := usq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := usq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (usq *UserSessionQuery) Modify(modifiers ...func(s *sql.Selector)) *UserSessionSelect {
	usq.modifiers = append(usq.modifiers, modifiers...)
	return usq.Select()
}

// UserSessionGroupBy is the group-by builder for UserSession entities.
type UserSessionGroupBy struct {
	selector
	build *UserSessionQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (usgb *UserSessionGroupBy) Aggregate(fns ...AggregateFunc) *UserSessionGroupBy {
	usgb.fns = append(usgb.fns, fns...)
	return usgb
}

// Scan applies the selector query and scans the result into the given value.
func (usgb *UserSessionGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, usgb.build.ctx, "GroupBy")
	if err := usgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*UserSessionQuery, *UserSessionGroupBy](ctx, usgb.build, usgb, usgb.build.inters, v)
}

func (usgb *UserSessionGroupBy) sqlScan(ctx context.Context, root *UserSessionQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(usgb.fns))
	for _, fn := range usgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*usgb.flds)+len(usgb.fns))
		for _, f := range *usgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*usgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := usgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// UserSessionSelect is the builder for selecting fields of UserSession entities.
type UserSessionSelect struct {
	*UserSessionQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (uss *UserSessionSelect) Aggregate(fns ...AggregateFunc) *UserSessionSelect {
	uss.fns = append(uss.fns, fns...)
	return uss
}

// Scan applies the selector query and scans the result into the given value.
func (uss *UserSessionSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, uss.ctx, "Select")
	if err := uss.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*UserSessionQuery, *UserSessionSelect](ctx, uss.UserSessionQuery, uss, uss.inters, v)
}

func (uss *UserSessionSelect) sqlScan(ctx context.Context, root *UserSessionQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(uss.fns))
	for _, fn := range uss.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*uss.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := uss.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (uss *UserSessionSelect) Modify(modifiers ...func(s *sql.Selector)) *UserSessionSelect {
	uss.modifiers = append(uss.modifiers, modifiers...)
	return uss
}
