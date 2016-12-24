package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
)

// SchemaMigration is an object representing the database table.
type SchemaMigration struct {
	Version int `boil:"version" json:"version" toml:"version" yaml:"version"`

	R *schemaMigrationR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L schemaMigrationL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// schemaMigrationR is where relationships are stored.
type schemaMigrationR struct {
}

// schemaMigrationL is where Load methods for each relationship are stored.
type schemaMigrationL struct{}

var (
	schemaMigrationColumns               = []string{"version"}
	schemaMigrationColumnsWithoutDefault = []string{"version"}
	schemaMigrationColumnsWithDefault    = []string{}
	schemaMigrationPrimaryKeyColumns     = []string{"version"}
)

type (
	// SchemaMigrationSlice is an alias for a slice of pointers to SchemaMigration.
	// This should generally be used opposed to []SchemaMigration.
	SchemaMigrationSlice []*SchemaMigration
	// SchemaMigrationHook is the signature for custom SchemaMigration hook methods
	SchemaMigrationHook func(boil.Executor, *SchemaMigration) error

	schemaMigrationQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	schemaMigrationType                 = reflect.TypeOf(&SchemaMigration{})
	schemaMigrationMapping              = queries.MakeStructMapping(schemaMigrationType)
	schemaMigrationPrimaryKeyMapping, _ = queries.BindMapping(schemaMigrationType, schemaMigrationMapping, schemaMigrationPrimaryKeyColumns)
	schemaMigrationInsertCacheMut       sync.RWMutex
	schemaMigrationInsertCache          = make(map[string]insertCache)
	schemaMigrationUpdateCacheMut       sync.RWMutex
	schemaMigrationUpdateCache          = make(map[string]updateCache)
	schemaMigrationUpsertCacheMut       sync.RWMutex
	schemaMigrationUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var schemaMigrationBeforeInsertHooks []SchemaMigrationHook
var schemaMigrationBeforeUpdateHooks []SchemaMigrationHook
var schemaMigrationBeforeDeleteHooks []SchemaMigrationHook
var schemaMigrationBeforeUpsertHooks []SchemaMigrationHook

var schemaMigrationAfterInsertHooks []SchemaMigrationHook
var schemaMigrationAfterSelectHooks []SchemaMigrationHook
var schemaMigrationAfterUpdateHooks []SchemaMigrationHook
var schemaMigrationAfterDeleteHooks []SchemaMigrationHook
var schemaMigrationAfterUpsertHooks []SchemaMigrationHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *SchemaMigration) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *SchemaMigration) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *SchemaMigration) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *SchemaMigration) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *SchemaMigration) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *SchemaMigration) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *SchemaMigration) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *SchemaMigration) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *SchemaMigration) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range schemaMigrationAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSchemaMigrationHook registers your hook function for all future operations.
func AddSchemaMigrationHook(hookPoint boil.HookPoint, schemaMigrationHook SchemaMigrationHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		schemaMigrationBeforeInsertHooks = append(schemaMigrationBeforeInsertHooks, schemaMigrationHook)
	case boil.BeforeUpdateHook:
		schemaMigrationBeforeUpdateHooks = append(schemaMigrationBeforeUpdateHooks, schemaMigrationHook)
	case boil.BeforeDeleteHook:
		schemaMigrationBeforeDeleteHooks = append(schemaMigrationBeforeDeleteHooks, schemaMigrationHook)
	case boil.BeforeUpsertHook:
		schemaMigrationBeforeUpsertHooks = append(schemaMigrationBeforeUpsertHooks, schemaMigrationHook)
	case boil.AfterInsertHook:
		schemaMigrationAfterInsertHooks = append(schemaMigrationAfterInsertHooks, schemaMigrationHook)
	case boil.AfterSelectHook:
		schemaMigrationAfterSelectHooks = append(schemaMigrationAfterSelectHooks, schemaMigrationHook)
	case boil.AfterUpdateHook:
		schemaMigrationAfterUpdateHooks = append(schemaMigrationAfterUpdateHooks, schemaMigrationHook)
	case boil.AfterDeleteHook:
		schemaMigrationAfterDeleteHooks = append(schemaMigrationAfterDeleteHooks, schemaMigrationHook)
	case boil.AfterUpsertHook:
		schemaMigrationAfterUpsertHooks = append(schemaMigrationAfterUpsertHooks, schemaMigrationHook)
	}
}

// OneP returns a single schemaMigration record from the query, and panics on error.
func (q schemaMigrationQuery) OneP() *SchemaMigration {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single schemaMigration record from the query.
func (q schemaMigrationQuery) One() (*SchemaMigration, error) {
	o := &SchemaMigration{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for schema_migrations")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all SchemaMigration records from the query, and panics on error.
func (q schemaMigrationQuery) AllP() SchemaMigrationSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all SchemaMigration records from the query.
func (q schemaMigrationQuery) All() (SchemaMigrationSlice, error) {
	var o SchemaMigrationSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SchemaMigration slice")
	}

	if len(schemaMigrationAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all SchemaMigration records in the query, and panics on error.
func (q schemaMigrationQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all SchemaMigration records in the query.
func (q schemaMigrationQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count schema_migrations rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q schemaMigrationQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q schemaMigrationQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if schema_migrations exists")
	}

	return count > 0, nil
}

// SchemaMigrationsG retrieves all records.
func SchemaMigrationsG(mods ...qm.QueryMod) schemaMigrationQuery {
	return SchemaMigrations(boil.GetDB(), mods...)
}

// SchemaMigrations retrieves all the records using an executor.
func SchemaMigrations(exec boil.Executor, mods ...qm.QueryMod) schemaMigrationQuery {
	mods = append(mods, qm.From("\"schema_migrations\""))
	return schemaMigrationQuery{NewQuery(exec, mods...)}
}

// FindSchemaMigrationG retrieves a single record by ID.
func FindSchemaMigrationG(version int, selectCols ...string) (*SchemaMigration, error) {
	return FindSchemaMigration(boil.GetDB(), version, selectCols...)
}

// FindSchemaMigrationGP retrieves a single record by ID, and panics on error.
func FindSchemaMigrationGP(version int, selectCols ...string) *SchemaMigration {
	retobj, err := FindSchemaMigration(boil.GetDB(), version, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindSchemaMigration retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSchemaMigration(exec boil.Executor, version int, selectCols ...string) (*SchemaMigration, error) {
	schemaMigrationObj := &SchemaMigration{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"schema_migrations\" where \"version\"=$1", sel,
	)

	q := queries.Raw(exec, query, version)

	err := q.Bind(schemaMigrationObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from schema_migrations")
	}

	return schemaMigrationObj, nil
}

// FindSchemaMigrationP retrieves a single record by ID with an executor, and panics on error.
func FindSchemaMigrationP(exec boil.Executor, version int, selectCols ...string) *SchemaMigration {
	retobj, err := FindSchemaMigration(exec, version, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *SchemaMigration) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *SchemaMigration) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *SchemaMigration) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *SchemaMigration) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no schema_migrations provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(schemaMigrationColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	schemaMigrationInsertCacheMut.RLock()
	cache, cached := schemaMigrationInsertCache[key]
	schemaMigrationInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			schemaMigrationColumns,
			schemaMigrationColumnsWithDefault,
			schemaMigrationColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(schemaMigrationType, schemaMigrationMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(schemaMigrationType, schemaMigrationMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"schema_migrations\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

		if len(cache.retMapping) != 0 {
			cache.query += fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into schema_migrations")
	}

	if !cached {
		schemaMigrationInsertCacheMut.Lock()
		schemaMigrationInsertCache[key] = cache
		schemaMigrationInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single SchemaMigration record. See Update for
// whitelist behavior description.
func (o *SchemaMigration) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single SchemaMigration record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *SchemaMigration) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the SchemaMigration, and panics on error.
// See Update for whitelist behavior description.
func (o *SchemaMigration) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the SchemaMigration.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *SchemaMigration) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	schemaMigrationUpdateCacheMut.RLock()
	cache, cached := schemaMigrationUpdateCache[key]
	schemaMigrationUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(schemaMigrationColumns, schemaMigrationPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update schema_migrations, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"schema_migrations\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, schemaMigrationPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(schemaMigrationType, schemaMigrationMapping, append(wl, schemaMigrationPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update schema_migrations row")
	}

	if !cached {
		schemaMigrationUpdateCacheMut.Lock()
		schemaMigrationUpdateCache[key] = cache
		schemaMigrationUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q schemaMigrationQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q schemaMigrationQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for schema_migrations")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SchemaMigrationSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o SchemaMigrationSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o SchemaMigrationSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SchemaMigrationSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), schemaMigrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"schema_migrations\" SET %s WHERE (\"version\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(schemaMigrationPrimaryKeyColumns), len(colNames)+1, len(schemaMigrationPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in schemaMigration slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *SchemaMigration) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *SchemaMigration) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *SchemaMigration) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *SchemaMigration) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no schema_migrations provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(schemaMigrationColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	schemaMigrationUpsertCacheMut.RLock()
	cache, cached := schemaMigrationUpsertCache[key]
	schemaMigrationUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			schemaMigrationColumns,
			schemaMigrationColumnsWithDefault,
			schemaMigrationColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			schemaMigrationColumns,
			schemaMigrationPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert schema_migrations, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(schemaMigrationPrimaryKeyColumns))
			copy(conflict, schemaMigrationPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"schema_migrations\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(schemaMigrationType, schemaMigrationMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(schemaMigrationType, schemaMigrationMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert for schema_migrations")
	}

	if !cached {
		schemaMigrationUpsertCacheMut.Lock()
		schemaMigrationUpsertCache[key] = cache
		schemaMigrationUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single SchemaMigration record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *SchemaMigration) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single SchemaMigration record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *SchemaMigration) DeleteG() error {
	if o == nil {
		return errors.New("models: no SchemaMigration provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single SchemaMigration record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *SchemaMigration) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single SchemaMigration record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SchemaMigration) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no SchemaMigration provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), schemaMigrationPrimaryKeyMapping)
	sql := "DELETE FROM \"schema_migrations\" WHERE \"version\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from schema_migrations")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q schemaMigrationQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q schemaMigrationQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no schemaMigrationQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from schema_migrations")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o SchemaMigrationSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o SchemaMigrationSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no SchemaMigration slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o SchemaMigrationSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SchemaMigrationSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no SchemaMigration slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(schemaMigrationBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), schemaMigrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"schema_migrations\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, schemaMigrationPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(schemaMigrationPrimaryKeyColumns), 1, len(schemaMigrationPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from schemaMigration slice")
	}

	if len(schemaMigrationAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *SchemaMigration) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *SchemaMigration) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *SchemaMigration) ReloadG() error {
	if o == nil {
		return errors.New("models: no SchemaMigration provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SchemaMigration) Reload(exec boil.Executor) error {
	ret, err := FindSchemaMigration(exec, o.Version)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SchemaMigrationSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SchemaMigrationSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SchemaMigrationSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty SchemaMigrationSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SchemaMigrationSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	schemaMigrations := SchemaMigrationSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), schemaMigrationPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"schema_migrations\".* FROM \"schema_migrations\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, schemaMigrationPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(schemaMigrationPrimaryKeyColumns), 1, len(schemaMigrationPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&schemaMigrations)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SchemaMigrationSlice")
	}

	*o = schemaMigrations

	return nil
}

// SchemaMigrationExists checks if the SchemaMigration row exists.
func SchemaMigrationExists(exec boil.Executor, version int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"schema_migrations\" where \"version\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, version)
	}

	row := exec.QueryRow(sql, version)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if schema_migrations exists")
	}

	return exists, nil
}

// SchemaMigrationExistsG checks if the SchemaMigration row exists.
func SchemaMigrationExistsG(version int) (bool, error) {
	return SchemaMigrationExists(boil.GetDB(), version)
}

// SchemaMigrationExistsGP checks if the SchemaMigration row exists. Panics on error.
func SchemaMigrationExistsGP(version int) bool {
	e, err := SchemaMigrationExists(boil.GetDB(), version)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// SchemaMigrationExistsP checks if the SchemaMigration row exists. Panics on error.
func SchemaMigrationExistsP(exec boil.Executor, version int) bool {
	e, err := SchemaMigrationExists(exec, version)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
