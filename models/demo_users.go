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
	"gopkg.in/nullbio/null.v6"
)

// DemoUser is an object representing the database table.
type DemoUser struct {
	RefID  int      `boil:"ref_id" json:"ref_id" toml:"ref_id" yaml:"ref_id"`
	DemoID null.Int `boil:"demo_id" json:"demo_id,omitempty" toml:"demo_id" yaml:"demo_id,omitempty"`
	UserID null.Int `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`

	R *demoUserR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L demoUserL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// demoUserR is where relationships are stored.
type demoUserR struct {
	Demo *Demo
	User *User
}

// demoUserL is where Load methods for each relationship are stored.
type demoUserL struct{}

var (
	demoUserColumns               = []string{"ref_id", "demo_id", "user_id"}
	demoUserColumnsWithoutDefault = []string{"demo_id", "user_id"}
	demoUserColumnsWithDefault    = []string{"ref_id"}
	demoUserPrimaryKeyColumns     = []string{"ref_id"}
)

type (
	// DemoUserSlice is an alias for a slice of pointers to DemoUser.
	// This should generally be used opposed to []DemoUser.
	DemoUserSlice []*DemoUser
	// DemoUserHook is the signature for custom DemoUser hook methods
	DemoUserHook func(boil.Executor, *DemoUser) error

	demoUserQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	demoUserType                 = reflect.TypeOf(&DemoUser{})
	demoUserMapping              = queries.MakeStructMapping(demoUserType)
	demoUserPrimaryKeyMapping, _ = queries.BindMapping(demoUserType, demoUserMapping, demoUserPrimaryKeyColumns)
	demoUserInsertCacheMut       sync.RWMutex
	demoUserInsertCache          = make(map[string]insertCache)
	demoUserUpdateCacheMut       sync.RWMutex
	demoUserUpdateCache          = make(map[string]updateCache)
	demoUserUpsertCacheMut       sync.RWMutex
	demoUserUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var demoUserBeforeInsertHooks []DemoUserHook
var demoUserBeforeUpdateHooks []DemoUserHook
var demoUserBeforeDeleteHooks []DemoUserHook
var demoUserBeforeUpsertHooks []DemoUserHook

var demoUserAfterInsertHooks []DemoUserHook
var demoUserAfterSelectHooks []DemoUserHook
var demoUserAfterUpdateHooks []DemoUserHook
var demoUserAfterDeleteHooks []DemoUserHook
var demoUserAfterUpsertHooks []DemoUserHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *DemoUser) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *DemoUser) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *DemoUser) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *DemoUser) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *DemoUser) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *DemoUser) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *DemoUser) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *DemoUser) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *DemoUser) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoUserAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddDemoUserHook registers your hook function for all future operations.
func AddDemoUserHook(hookPoint boil.HookPoint, demoUserHook DemoUserHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		demoUserBeforeInsertHooks = append(demoUserBeforeInsertHooks, demoUserHook)
	case boil.BeforeUpdateHook:
		demoUserBeforeUpdateHooks = append(demoUserBeforeUpdateHooks, demoUserHook)
	case boil.BeforeDeleteHook:
		demoUserBeforeDeleteHooks = append(demoUserBeforeDeleteHooks, demoUserHook)
	case boil.BeforeUpsertHook:
		demoUserBeforeUpsertHooks = append(demoUserBeforeUpsertHooks, demoUserHook)
	case boil.AfterInsertHook:
		demoUserAfterInsertHooks = append(demoUserAfterInsertHooks, demoUserHook)
	case boil.AfterSelectHook:
		demoUserAfterSelectHooks = append(demoUserAfterSelectHooks, demoUserHook)
	case boil.AfterUpdateHook:
		demoUserAfterUpdateHooks = append(demoUserAfterUpdateHooks, demoUserHook)
	case boil.AfterDeleteHook:
		demoUserAfterDeleteHooks = append(demoUserAfterDeleteHooks, demoUserHook)
	case boil.AfterUpsertHook:
		demoUserAfterUpsertHooks = append(demoUserAfterUpsertHooks, demoUserHook)
	}
}

// OneP returns a single demoUser record from the query, and panics on error.
func (q demoUserQuery) OneP() *DemoUser {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single demoUser record from the query.
func (q demoUserQuery) One() (*DemoUser, error) {
	o := &DemoUser{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for demo_users")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all DemoUser records from the query, and panics on error.
func (q demoUserQuery) AllP() DemoUserSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all DemoUser records from the query.
func (q demoUserQuery) All() (DemoUserSlice, error) {
	var o DemoUserSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to DemoUser slice")
	}

	if len(demoUserAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all DemoUser records in the query, and panics on error.
func (q demoUserQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all DemoUser records in the query.
func (q demoUserQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count demo_users rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q demoUserQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q demoUserQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if demo_users exists")
	}

	return count > 0, nil
}

// DemoG pointed to by the foreign key.
func (o *DemoUser) DemoG(mods ...qm.QueryMod) demoQuery {
	return o.Demo(boil.GetDB(), mods...)
}

// Demo pointed to by the foreign key.
func (o *DemoUser) Demo(exec boil.Executor, mods ...qm.QueryMod) demoQuery {
	queryMods := []qm.QueryMod{
		qm.Where("demo_id=?", o.DemoID),
	}

	queryMods = append(queryMods, mods...)

	query := Demos(exec, queryMods...)
	queries.SetFrom(query.Query, "\"demos\"")

	return query
}

// UserG pointed to by the foreign key.
func (o *DemoUser) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *DemoUser) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("user_id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadDemo allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (demoUserL) LoadDemo(e boil.Executor, singular bool, maybeDemoUser interface{}) error {
	var slice []*DemoUser
	var object *DemoUser

	count := 1
	if singular {
		object = maybeDemoUser.(*DemoUser)
	} else {
		slice = *maybeDemoUser.(*DemoUserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &demoUserR{}
		}
		args[0] = object.DemoID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &demoUserR{}
			}
			args[i] = obj.DemoID
		}
	}

	query := fmt.Sprintf(
		"select * from \"demos\" where \"demo_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Demo")
	}
	defer results.Close()

	var resultSlice []*Demo
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Demo")
	}

	if len(demoUserAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Demo = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.DemoID.Int == foreign.DemoID {
				local.R.Demo = foreign
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (demoUserL) LoadUser(e boil.Executor, singular bool, maybeDemoUser interface{}) error {
	var slice []*DemoUser
	var object *DemoUser

	count := 1
	if singular {
		object = maybeDemoUser.(*DemoUser)
	} else {
		slice = *maybeDemoUser.(*DemoUserSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &demoUserR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &demoUserR{}
			}
			args[i] = obj.UserID
		}
	}

	query := fmt.Sprintf(
		"select * from \"users\" where \"user_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}
	defer results.Close()

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if len(demoUserAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.User = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.UserID.Int == foreign.UserID {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// SetDemo of the demo_user to the related item.
// Sets o.R.Demo to related.
// Adds o to related.R.DemoUsers.
func (o *DemoUser) SetDemo(exec boil.Executor, insert bool, related *Demo) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"demo_users\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"demo_id"}),
		strmangle.WhereClause("\"", "\"", 2, demoUserPrimaryKeyColumns),
	)
	values := []interface{}{related.DemoID, o.RefID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.DemoID.Int = related.DemoID
	o.DemoID.Valid = true

	if o.R == nil {
		o.R = &demoUserR{
			Demo: related,
		}
	} else {
		o.R.Demo = related
	}

	if related.R == nil {
		related.R = &demoR{
			DemoUsers: DemoUserSlice{o},
		}
	} else {
		related.R.DemoUsers = append(related.R.DemoUsers, o)
	}

	return nil
}

// RemoveDemo relationship.
// Sets o.R.Demo to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *DemoUser) RemoveDemo(exec boil.Executor, related *Demo) error {
	var err error

	o.DemoID.Valid = false
	if err = o.Update(exec, "demo_id"); err != nil {
		o.DemoID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Demo = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.DemoUsers {
		if o.DemoID.Int != ri.DemoID.Int {
			continue
		}

		ln := len(related.R.DemoUsers)
		if ln > 1 && i < ln-1 {
			related.R.DemoUsers[i] = related.R.DemoUsers[ln-1]
		}
		related.R.DemoUsers = related.R.DemoUsers[:ln-1]
		break
	}
	return nil
}

// SetUser of the demo_user to the related item.
// Sets o.R.User to related.
// Adds o to related.R.DemoUsers.
func (o *DemoUser) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"demo_users\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, demoUserPrimaryKeyColumns),
	)
	values := []interface{}{related.UserID, o.RefID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID.Int = related.UserID
	o.UserID.Valid = true

	if o.R == nil {
		o.R = &demoUserR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			DemoUsers: DemoUserSlice{o},
		}
	} else {
		related.R.DemoUsers = append(related.R.DemoUsers, o)
	}

	return nil
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *DemoUser) RemoveUser(exec boil.Executor, related *User) error {
	var err error

	o.UserID.Valid = false
	if err = o.Update(exec, "user_id"); err != nil {
		o.UserID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.User = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.DemoUsers {
		if o.UserID.Int != ri.UserID.Int {
			continue
		}

		ln := len(related.R.DemoUsers)
		if ln > 1 && i < ln-1 {
			related.R.DemoUsers[i] = related.R.DemoUsers[ln-1]
		}
		related.R.DemoUsers = related.R.DemoUsers[:ln-1]
		break
	}
	return nil
}

// DemoUsersG retrieves all records.
func DemoUsersG(mods ...qm.QueryMod) demoUserQuery {
	return DemoUsers(boil.GetDB(), mods...)
}

// DemoUsers retrieves all the records using an executor.
func DemoUsers(exec boil.Executor, mods ...qm.QueryMod) demoUserQuery {
	mods = append(mods, qm.From("\"demo_users\""))
	return demoUserQuery{NewQuery(exec, mods...)}
}

// FindDemoUserG retrieves a single record by ID.
func FindDemoUserG(refID int, selectCols ...string) (*DemoUser, error) {
	return FindDemoUser(boil.GetDB(), refID, selectCols...)
}

// FindDemoUserGP retrieves a single record by ID, and panics on error.
func FindDemoUserGP(refID int, selectCols ...string) *DemoUser {
	retobj, err := FindDemoUser(boil.GetDB(), refID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindDemoUser retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindDemoUser(exec boil.Executor, refID int, selectCols ...string) (*DemoUser, error) {
	demoUserObj := &DemoUser{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"demo_users\" where \"ref_id\"=$1", sel,
	)

	q := queries.Raw(exec, query, refID)

	err := q.Bind(demoUserObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from demo_users")
	}

	return demoUserObj, nil
}

// FindDemoUserP retrieves a single record by ID with an executor, and panics on error.
func FindDemoUserP(exec boil.Executor, refID int, selectCols ...string) *DemoUser {
	retobj, err := FindDemoUser(exec, refID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *DemoUser) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *DemoUser) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *DemoUser) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *DemoUser) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no demo_users provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(demoUserColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	demoUserInsertCacheMut.RLock()
	cache, cached := demoUserInsertCache[key]
	demoUserInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			demoUserColumns,
			demoUserColumnsWithDefault,
			demoUserColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(demoUserType, demoUserMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(demoUserType, demoUserMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"demo_users\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into demo_users")
	}

	if !cached {
		demoUserInsertCacheMut.Lock()
		demoUserInsertCache[key] = cache
		demoUserInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single DemoUser record. See Update for
// whitelist behavior description.
func (o *DemoUser) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single DemoUser record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *DemoUser) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the DemoUser, and panics on error.
// See Update for whitelist behavior description.
func (o *DemoUser) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the DemoUser.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *DemoUser) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	demoUserUpdateCacheMut.RLock()
	cache, cached := demoUserUpdateCache[key]
	demoUserUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(demoUserColumns, demoUserPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update demo_users, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"demo_users\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, demoUserPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(demoUserType, demoUserMapping, append(wl, demoUserPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update demo_users row")
	}

	if !cached {
		demoUserUpdateCacheMut.Lock()
		demoUserUpdateCache[key] = cache
		demoUserUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q demoUserQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q demoUserQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for demo_users")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o DemoUserSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o DemoUserSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o DemoUserSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o DemoUserSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), demoUserPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"demo_users\" SET %s WHERE (\"ref_id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(demoUserPrimaryKeyColumns), len(colNames)+1, len(demoUserPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in demoUser slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *DemoUser) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *DemoUser) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *DemoUser) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *DemoUser) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no demo_users provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(demoUserColumnsWithDefault, o)

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

	demoUserUpsertCacheMut.RLock()
	cache, cached := demoUserUpsertCache[key]
	demoUserUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			demoUserColumns,
			demoUserColumnsWithDefault,
			demoUserColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			demoUserColumns,
			demoUserPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert demo_users, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(demoUserPrimaryKeyColumns))
			copy(conflict, demoUserPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"demo_users\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(demoUserType, demoUserMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(demoUserType, demoUserMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for demo_users")
	}

	if !cached {
		demoUserUpsertCacheMut.Lock()
		demoUserUpsertCache[key] = cache
		demoUserUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single DemoUser record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *DemoUser) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single DemoUser record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *DemoUser) DeleteG() error {
	if o == nil {
		return errors.New("models: no DemoUser provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single DemoUser record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *DemoUser) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single DemoUser record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *DemoUser) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no DemoUser provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), demoUserPrimaryKeyMapping)
	sql := "DELETE FROM \"demo_users\" WHERE \"ref_id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from demo_users")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q demoUserQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q demoUserQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no demoUserQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from demo_users")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o DemoUserSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o DemoUserSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no DemoUser slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o DemoUserSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o DemoUserSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no DemoUser slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(demoUserBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), demoUserPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"demo_users\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, demoUserPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(demoUserPrimaryKeyColumns), 1, len(demoUserPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from demoUser slice")
	}

	if len(demoUserAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *DemoUser) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *DemoUser) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *DemoUser) ReloadG() error {
	if o == nil {
		return errors.New("models: no DemoUser provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *DemoUser) Reload(exec boil.Executor) error {
	ret, err := FindDemoUser(exec, o.RefID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DemoUserSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DemoUserSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DemoUserSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty DemoUserSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DemoUserSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	demoUsers := DemoUserSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), demoUserPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"demo_users\".* FROM \"demo_users\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, demoUserPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(demoUserPrimaryKeyColumns), 1, len(demoUserPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&demoUsers)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in DemoUserSlice")
	}

	*o = demoUsers

	return nil
}

// DemoUserExists checks if the DemoUser row exists.
func DemoUserExists(exec boil.Executor, refID int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"demo_users\" where \"ref_id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, refID)
	}

	row := exec.QueryRow(sql, refID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if demo_users exists")
	}

	return exists, nil
}

// DemoUserExistsG checks if the DemoUser row exists.
func DemoUserExistsG(refID int) (bool, error) {
	return DemoUserExists(boil.GetDB(), refID)
}

// DemoUserExistsGP checks if the DemoUser row exists. Panics on error.
func DemoUserExistsGP(refID int) bool {
	e, err := DemoUserExists(boil.GetDB(), refID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// DemoUserExistsP checks if the DemoUser row exists. Panics on error.
func DemoUserExistsP(exec boil.Executor, refID int) bool {
	e, err := DemoUserExists(exec, refID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
