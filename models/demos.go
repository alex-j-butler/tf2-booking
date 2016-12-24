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

// Demo is an object representing the database table.
type Demo struct {
	DemoID        int       `boil:"demo_id" json:"demo_id" toml:"demo_id" yaml:"demo_id"`
	BookingID     null.Int  `boil:"booking_id" json:"booking_id,omitempty" toml:"booking_id" yaml:"booking_id,omitempty"`
	Name          string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	MapName       string    `boil:"map_name" json:"map_name" toml:"map_name" yaml:"map_name"`
	Configuration string    `boil:"configuration" json:"configuration" toml:"configuration" yaml:"configuration"`
	URL           string    `boil:"url" json:"url" toml:"url" yaml:"url"`
	UploadedTime  null.Time `boil:"uploaded_time" json:"uploaded_time,omitempty" toml:"uploaded_time" yaml:"uploaded_time,omitempty"`

	R *demoR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L demoL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// demoR is where relationships are stored.
type demoR struct {
	Booking   *Booking
	DemoUsers DemoUserSlice
}

// demoL is where Load methods for each relationship are stored.
type demoL struct{}

var (
	demoColumns               = []string{"demo_id", "booking_id", "name", "map_name", "configuration", "url", "uploaded_time"}
	demoColumnsWithoutDefault = []string{"booking_id", "name", "map_name", "configuration", "url", "uploaded_time"}
	demoColumnsWithDefault    = []string{"demo_id"}
	demoPrimaryKeyColumns     = []string{"demo_id"}
)

type (
	// DemoSlice is an alias for a slice of pointers to Demo.
	// This should generally be used opposed to []Demo.
	DemoSlice []*Demo
	// DemoHook is the signature for custom Demo hook methods
	DemoHook func(boil.Executor, *Demo) error

	demoQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	demoType                 = reflect.TypeOf(&Demo{})
	demoMapping              = queries.MakeStructMapping(demoType)
	demoPrimaryKeyMapping, _ = queries.BindMapping(demoType, demoMapping, demoPrimaryKeyColumns)
	demoInsertCacheMut       sync.RWMutex
	demoInsertCache          = make(map[string]insertCache)
	demoUpdateCacheMut       sync.RWMutex
	demoUpdateCache          = make(map[string]updateCache)
	demoUpsertCacheMut       sync.RWMutex
	demoUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var demoBeforeInsertHooks []DemoHook
var demoBeforeUpdateHooks []DemoHook
var demoBeforeDeleteHooks []DemoHook
var demoBeforeUpsertHooks []DemoHook

var demoAfterInsertHooks []DemoHook
var demoAfterSelectHooks []DemoHook
var demoAfterUpdateHooks []DemoHook
var demoAfterDeleteHooks []DemoHook
var demoAfterUpsertHooks []DemoHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Demo) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Demo) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range demoBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Demo) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range demoBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Demo) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Demo) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Demo) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range demoAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Demo) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range demoAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Demo) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range demoAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Demo) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range demoAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddDemoHook registers your hook function for all future operations.
func AddDemoHook(hookPoint boil.HookPoint, demoHook DemoHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		demoBeforeInsertHooks = append(demoBeforeInsertHooks, demoHook)
	case boil.BeforeUpdateHook:
		demoBeforeUpdateHooks = append(demoBeforeUpdateHooks, demoHook)
	case boil.BeforeDeleteHook:
		demoBeforeDeleteHooks = append(demoBeforeDeleteHooks, demoHook)
	case boil.BeforeUpsertHook:
		demoBeforeUpsertHooks = append(demoBeforeUpsertHooks, demoHook)
	case boil.AfterInsertHook:
		demoAfterInsertHooks = append(demoAfterInsertHooks, demoHook)
	case boil.AfterSelectHook:
		demoAfterSelectHooks = append(demoAfterSelectHooks, demoHook)
	case boil.AfterUpdateHook:
		demoAfterUpdateHooks = append(demoAfterUpdateHooks, demoHook)
	case boil.AfterDeleteHook:
		demoAfterDeleteHooks = append(demoAfterDeleteHooks, demoHook)
	case boil.AfterUpsertHook:
		demoAfterUpsertHooks = append(demoAfterUpsertHooks, demoHook)
	}
}

// OneP returns a single demo record from the query, and panics on error.
func (q demoQuery) OneP() *Demo {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single demo record from the query.
func (q demoQuery) One() (*Demo, error) {
	o := &Demo{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for demos")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Demo records from the query, and panics on error.
func (q demoQuery) AllP() DemoSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Demo records from the query.
func (q demoQuery) All() (DemoSlice, error) {
	var o DemoSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Demo slice")
	}

	if len(demoAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Demo records in the query, and panics on error.
func (q demoQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Demo records in the query.
func (q demoQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count demos rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q demoQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q demoQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if demos exists")
	}

	return count > 0, nil
}

// BookingG pointed to by the foreign key.
func (o *Demo) BookingG(mods ...qm.QueryMod) bookingQuery {
	return o.Booking(boil.GetDB(), mods...)
}

// Booking pointed to by the foreign key.
func (o *Demo) Booking(exec boil.Executor, mods ...qm.QueryMod) bookingQuery {
	queryMods := []qm.QueryMod{
		qm.Where("booking_id=?", o.BookingID),
	}

	queryMods = append(queryMods, mods...)

	query := Bookings(exec, queryMods...)
	queries.SetFrom(query.Query, "\"bookings\"")

	return query
}

// DemoUsersG retrieves all the demo_user's demo users.
func (o *Demo) DemoUsersG(mods ...qm.QueryMod) demoUserQuery {
	return o.DemoUsers(boil.GetDB(), mods...)
}

// DemoUsers retrieves all the demo_user's demo users with an executor.
func (o *Demo) DemoUsers(exec boil.Executor, mods ...qm.QueryMod) demoUserQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"demo_id\"=?", o.DemoID),
	)

	query := DemoUsers(exec, queryMods...)
	queries.SetFrom(query.Query, "\"demo_users\" as \"a\"")
	return query
}

// LoadBooking allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (demoL) LoadBooking(e boil.Executor, singular bool, maybeDemo interface{}) error {
	var slice []*Demo
	var object *Demo

	count := 1
	if singular {
		object = maybeDemo.(*Demo)
	} else {
		slice = *maybeDemo.(*DemoSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &demoR{}
		}
		args[0] = object.BookingID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &demoR{}
			}
			args[i] = obj.BookingID
		}
	}

	query := fmt.Sprintf(
		"select * from \"bookings\" where \"booking_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Booking")
	}
	defer results.Close()

	var resultSlice []*Booking
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Booking")
	}

	if len(demoAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Booking = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.BookingID.Int == foreign.BookingID {
				local.R.Booking = foreign
				break
			}
		}
	}

	return nil
}

// LoadDemoUsers allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (demoL) LoadDemoUsers(e boil.Executor, singular bool, maybeDemo interface{}) error {
	var slice []*Demo
	var object *Demo

	count := 1
	if singular {
		object = maybeDemo.(*Demo)
	} else {
		slice = *maybeDemo.(*DemoSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &demoR{}
		}
		args[0] = object.DemoID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &demoR{}
			}
			args[i] = obj.DemoID
		}
	}

	query := fmt.Sprintf(
		"select * from \"demo_users\" where \"demo_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load demo_users")
	}
	defer results.Close()

	var resultSlice []*DemoUser
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice demo_users")
	}

	if len(demoUserAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.DemoUsers = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.DemoID == foreign.DemoID.Int {
				local.R.DemoUsers = append(local.R.DemoUsers, foreign)
				break
			}
		}
	}

	return nil
}

// SetBooking of the demo to the related item.
// Sets o.R.Booking to related.
// Adds o to related.R.Demos.
func (o *Demo) SetBooking(exec boil.Executor, insert bool, related *Booking) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"demos\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"booking_id"}),
		strmangle.WhereClause("\"", "\"", 2, demoPrimaryKeyColumns),
	)
	values := []interface{}{related.BookingID, o.DemoID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.BookingID.Int = related.BookingID
	o.BookingID.Valid = true

	if o.R == nil {
		o.R = &demoR{
			Booking: related,
		}
	} else {
		o.R.Booking = related
	}

	if related.R == nil {
		related.R = &bookingR{
			Demos: DemoSlice{o},
		}
	} else {
		related.R.Demos = append(related.R.Demos, o)
	}

	return nil
}

// RemoveBooking relationship.
// Sets o.R.Booking to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Demo) RemoveBooking(exec boil.Executor, related *Booking) error {
	var err error

	o.BookingID.Valid = false
	if err = o.Update(exec, "booking_id"); err != nil {
		o.BookingID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Booking = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.Demos {
		if o.BookingID.Int != ri.BookingID.Int {
			continue
		}

		ln := len(related.R.Demos)
		if ln > 1 && i < ln-1 {
			related.R.Demos[i] = related.R.Demos[ln-1]
		}
		related.R.Demos = related.R.Demos[:ln-1]
		break
	}
	return nil
}

// AddDemoUsers adds the given related objects to the existing relationships
// of the demo, optionally inserting them as new records.
// Appends related to o.R.DemoUsers.
// Sets related.R.Demo appropriately.
func (o *Demo) AddDemoUsers(exec boil.Executor, insert bool, related ...*DemoUser) error {
	var err error
	for _, rel := range related {
		rel.DemoID.Int = o.DemoID
		rel.DemoID.Valid = true
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			if err = rel.Update(exec, "demo_id"); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}
		}
	}

	if o.R == nil {
		o.R = &demoR{
			DemoUsers: related,
		}
	} else {
		o.R.DemoUsers = append(o.R.DemoUsers, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &demoUserR{
				Demo: o,
			}
		} else {
			rel.R.Demo = o
		}
	}
	return nil
}

// SetDemoUsers removes all previously related items of the
// demo replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Demo's DemoUsers accordingly.
// Replaces o.R.DemoUsers with related.
// Sets related.R.Demo's DemoUsers accordingly.
func (o *Demo) SetDemoUsers(exec boil.Executor, insert bool, related ...*DemoUser) error {
	query := "update \"demo_users\" set \"demo_id\" = null where \"demo_id\" = $1"
	values := []interface{}{o.DemoID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.DemoUsers {
			rel.DemoID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Demo = nil
		}

		o.R.DemoUsers = nil
	}
	return o.AddDemoUsers(exec, insert, related...)
}

// RemoveDemoUsers relationships from objects passed in.
// Removes related items from R.DemoUsers (uses pointer comparison, removal does not keep order)
// Sets related.R.Demo.
func (o *Demo) RemoveDemoUsers(exec boil.Executor, related ...*DemoUser) error {
	var err error
	for _, rel := range related {
		rel.DemoID.Valid = false
		if rel.R != nil {
			rel.R.Demo = nil
		}
		if err = rel.Update(exec, "demo_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.DemoUsers {
			if rel != ri {
				continue
			}

			ln := len(o.R.DemoUsers)
			if ln > 1 && i < ln-1 {
				o.R.DemoUsers[i] = o.R.DemoUsers[ln-1]
			}
			o.R.DemoUsers = o.R.DemoUsers[:ln-1]
			break
		}
	}

	return nil
}

// DemosG retrieves all records.
func DemosG(mods ...qm.QueryMod) demoQuery {
	return Demos(boil.GetDB(), mods...)
}

// Demos retrieves all the records using an executor.
func Demos(exec boil.Executor, mods ...qm.QueryMod) demoQuery {
	mods = append(mods, qm.From("\"demos\""))
	return demoQuery{NewQuery(exec, mods...)}
}

// FindDemoG retrieves a single record by ID.
func FindDemoG(demoID int, selectCols ...string) (*Demo, error) {
	return FindDemo(boil.GetDB(), demoID, selectCols...)
}

// FindDemoGP retrieves a single record by ID, and panics on error.
func FindDemoGP(demoID int, selectCols ...string) *Demo {
	retobj, err := FindDemo(boil.GetDB(), demoID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindDemo retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindDemo(exec boil.Executor, demoID int, selectCols ...string) (*Demo, error) {
	demoObj := &Demo{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"demos\" where \"demo_id\"=$1", sel,
	)

	q := queries.Raw(exec, query, demoID)

	err := q.Bind(demoObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from demos")
	}

	return demoObj, nil
}

// FindDemoP retrieves a single record by ID with an executor, and panics on error.
func FindDemoP(exec boil.Executor, demoID int, selectCols ...string) *Demo {
	retobj, err := FindDemo(exec, demoID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Demo) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Demo) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Demo) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Demo) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no demos provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(demoColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	demoInsertCacheMut.RLock()
	cache, cached := demoInsertCache[key]
	demoInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			demoColumns,
			demoColumnsWithDefault,
			demoColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(demoType, demoMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(demoType, demoMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"demos\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into demos")
	}

	if !cached {
		demoInsertCacheMut.Lock()
		demoInsertCache[key] = cache
		demoInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Demo record. See Update for
// whitelist behavior description.
func (o *Demo) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Demo record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Demo) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Demo, and panics on error.
// See Update for whitelist behavior description.
func (o *Demo) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Demo.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Demo) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	demoUpdateCacheMut.RLock()
	cache, cached := demoUpdateCache[key]
	demoUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(demoColumns, demoPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update demos, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"demos\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, demoPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(demoType, demoMapping, append(wl, demoPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update demos row")
	}

	if !cached {
		demoUpdateCacheMut.Lock()
		demoUpdateCache[key] = cache
		demoUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q demoQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q demoQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for demos")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o DemoSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o DemoSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o DemoSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o DemoSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), demoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"demos\" SET %s WHERE (\"demo_id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(demoPrimaryKeyColumns), len(colNames)+1, len(demoPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in demo slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Demo) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Demo) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Demo) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Demo) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no demos provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(demoColumnsWithDefault, o)

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

	demoUpsertCacheMut.RLock()
	cache, cached := demoUpsertCache[key]
	demoUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			demoColumns,
			demoColumnsWithDefault,
			demoColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			demoColumns,
			demoPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert demos, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(demoPrimaryKeyColumns))
			copy(conflict, demoPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"demos\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(demoType, demoMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(demoType, demoMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for demos")
	}

	if !cached {
		demoUpsertCacheMut.Lock()
		demoUpsertCache[key] = cache
		demoUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Demo record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Demo) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Demo record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Demo) DeleteG() error {
	if o == nil {
		return errors.New("models: no Demo provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Demo record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Demo) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Demo record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Demo) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Demo provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), demoPrimaryKeyMapping)
	sql := "DELETE FROM \"demos\" WHERE \"demo_id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from demos")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q demoQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q demoQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no demoQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from demos")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o DemoSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o DemoSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Demo slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o DemoSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o DemoSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Demo slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(demoBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), demoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"demos\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, demoPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(demoPrimaryKeyColumns), 1, len(demoPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from demo slice")
	}

	if len(demoAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Demo) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Demo) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Demo) ReloadG() error {
	if o == nil {
		return errors.New("models: no Demo provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Demo) Reload(exec boil.Executor) error {
	ret, err := FindDemo(exec, o.DemoID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DemoSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *DemoSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DemoSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty DemoSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *DemoSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	demos := DemoSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), demoPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"demos\".* FROM \"demos\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, demoPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(demoPrimaryKeyColumns), 1, len(demoPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&demos)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in DemoSlice")
	}

	*o = demos

	return nil
}

// DemoExists checks if the Demo row exists.
func DemoExists(exec boil.Executor, demoID int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"demos\" where \"demo_id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, demoID)
	}

	row := exec.QueryRow(sql, demoID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if demos exists")
	}

	return exists, nil
}

// DemoExistsG checks if the Demo row exists.
func DemoExistsG(demoID int) (bool, error) {
	return DemoExists(boil.GetDB(), demoID)
}

// DemoExistsGP checks if the Demo row exists. Panics on error.
func DemoExistsGP(demoID int) bool {
	e, err := DemoExists(boil.GetDB(), demoID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// DemoExistsP checks if the Demo row exists. Panics on error.
func DemoExistsP(exec boil.Executor, demoID int) bool {
	e, err := DemoExists(exec, demoID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
