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

// Booking is an object representing the database table.
type Booking struct {
	BookingID    int       `boil:"booking_id" json:"booking_id" toml:"booking_id" yaml:"booking_id"`
	BookerID     null.Int  `boil:"booker_id" json:"booker_id,omitempty" toml:"booker_id" yaml:"booker_id,omitempty"`
	ServerName   string    `boil:"server_name" json:"server_name" toml:"server_name" yaml:"server_name"`
	BookedTime   null.Time `boil:"booked_time" json:"booked_time,omitempty" toml:"booked_time" yaml:"booked_time,omitempty"`
	UnbookedTime null.Time `boil:"unbooked_time" json:"unbooked_time,omitempty" toml:"unbooked_time" yaml:"unbooked_time,omitempty"`

	R *bookingR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L bookingL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// bookingR is where relationships are stored.
type bookingR struct {
	Booker *User
	Demos  DemoSlice
}

// bookingL is where Load methods for each relationship are stored.
type bookingL struct{}

var (
	bookingColumns               = []string{"booking_id", "booker_id", "server_name", "booked_time", "unbooked_time"}
	bookingColumnsWithoutDefault = []string{"booker_id", "server_name", "booked_time", "unbooked_time"}
	bookingColumnsWithDefault    = []string{"booking_id"}
	bookingPrimaryKeyColumns     = []string{"booking_id"}
)

type (
	// BookingSlice is an alias for a slice of pointers to Booking.
	// This should generally be used opposed to []Booking.
	BookingSlice []*Booking
	// BookingHook is the signature for custom Booking hook methods
	BookingHook func(boil.Executor, *Booking) error

	bookingQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	bookingType                 = reflect.TypeOf(&Booking{})
	bookingMapping              = queries.MakeStructMapping(bookingType)
	bookingPrimaryKeyMapping, _ = queries.BindMapping(bookingType, bookingMapping, bookingPrimaryKeyColumns)
	bookingInsertCacheMut       sync.RWMutex
	bookingInsertCache          = make(map[string]insertCache)
	bookingUpdateCacheMut       sync.RWMutex
	bookingUpdateCache          = make(map[string]updateCache)
	bookingUpsertCacheMut       sync.RWMutex
	bookingUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var bookingBeforeInsertHooks []BookingHook
var bookingBeforeUpdateHooks []BookingHook
var bookingBeforeDeleteHooks []BookingHook
var bookingBeforeUpsertHooks []BookingHook

var bookingAfterInsertHooks []BookingHook
var bookingAfterSelectHooks []BookingHook
var bookingAfterUpdateHooks []BookingHook
var bookingAfterDeleteHooks []BookingHook
var bookingAfterUpsertHooks []BookingHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Booking) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Booking) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Booking) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Booking) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Booking) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Booking) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Booking) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Booking) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Booking) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range bookingAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddBookingHook registers your hook function for all future operations.
func AddBookingHook(hookPoint boil.HookPoint, bookingHook BookingHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		bookingBeforeInsertHooks = append(bookingBeforeInsertHooks, bookingHook)
	case boil.BeforeUpdateHook:
		bookingBeforeUpdateHooks = append(bookingBeforeUpdateHooks, bookingHook)
	case boil.BeforeDeleteHook:
		bookingBeforeDeleteHooks = append(bookingBeforeDeleteHooks, bookingHook)
	case boil.BeforeUpsertHook:
		bookingBeforeUpsertHooks = append(bookingBeforeUpsertHooks, bookingHook)
	case boil.AfterInsertHook:
		bookingAfterInsertHooks = append(bookingAfterInsertHooks, bookingHook)
	case boil.AfterSelectHook:
		bookingAfterSelectHooks = append(bookingAfterSelectHooks, bookingHook)
	case boil.AfterUpdateHook:
		bookingAfterUpdateHooks = append(bookingAfterUpdateHooks, bookingHook)
	case boil.AfterDeleteHook:
		bookingAfterDeleteHooks = append(bookingAfterDeleteHooks, bookingHook)
	case boil.AfterUpsertHook:
		bookingAfterUpsertHooks = append(bookingAfterUpsertHooks, bookingHook)
	}
}

// OneP returns a single booking record from the query, and panics on error.
func (q bookingQuery) OneP() *Booking {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single booking record from the query.
func (q bookingQuery) One() (*Booking, error) {
	o := &Booking{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for bookings")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Booking records from the query, and panics on error.
func (q bookingQuery) AllP() BookingSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Booking records from the query.
func (q bookingQuery) All() (BookingSlice, error) {
	var o BookingSlice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Booking slice")
	}

	if len(bookingAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Booking records in the query, and panics on error.
func (q bookingQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Booking records in the query.
func (q bookingQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count bookings rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q bookingQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q bookingQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if bookings exists")
	}

	return count > 0, nil
}

// BookerG pointed to by the foreign key.
func (o *Booking) BookerG(mods ...qm.QueryMod) userQuery {
	return o.Booker(boil.GetDB(), mods...)
}

// Booker pointed to by the foreign key.
func (o *Booking) Booker(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("user_id=?", o.BookerID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// DemosG retrieves all the demo's demos.
func (o *Booking) DemosG(mods ...qm.QueryMod) demoQuery {
	return o.Demos(boil.GetDB(), mods...)
}

// Demos retrieves all the demo's demos with an executor.
func (o *Booking) Demos(exec boil.Executor, mods ...qm.QueryMod) demoQuery {
	queryMods := []qm.QueryMod{
		qm.Select("\"a\".*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"a\".\"booking_id\"=?", o.BookingID),
	)

	query := Demos(exec, queryMods...)
	queries.SetFrom(query.Query, "\"demos\" as \"a\"")
	return query
}

// LoadBooker allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (bookingL) LoadBooker(e boil.Executor, singular bool, maybeBooking interface{}) error {
	var slice []*Booking
	var object *Booking

	count := 1
	if singular {
		object = maybeBooking.(*Booking)
	} else {
		slice = *maybeBooking.(*BookingSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &bookingR{}
		}
		args[0] = object.BookerID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &bookingR{}
			}
			args[i] = obj.BookerID
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

	if len(bookingAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	if singular && len(resultSlice) != 0 {
		object.R.Booker = resultSlice[0]
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.BookerID.Int == foreign.UserID {
				local.R.Booker = foreign
				break
			}
		}
	}

	return nil
}

// LoadDemos allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (bookingL) LoadDemos(e boil.Executor, singular bool, maybeBooking interface{}) error {
	var slice []*Booking
	var object *Booking

	count := 1
	if singular {
		object = maybeBooking.(*Booking)
	} else {
		slice = *maybeBooking.(*BookingSlice)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &bookingR{}
		}
		args[0] = object.BookingID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &bookingR{}
			}
			args[i] = obj.BookingID
		}
	}

	query := fmt.Sprintf(
		"select * from \"demos\" where \"booking_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load demos")
	}
	defer results.Close()

	var resultSlice []*Demo
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice demos")
	}

	if len(demoAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Demos = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.BookingID == foreign.BookingID.Int {
				local.R.Demos = append(local.R.Demos, foreign)
				break
			}
		}
	}

	return nil
}

// SetBooker of the booking to the related item.
// Sets o.R.Booker to related.
// Adds o to related.R.BookerBookings.
func (o *Booking) SetBooker(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"bookings\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"booker_id"}),
		strmangle.WhereClause("\"", "\"", 2, bookingPrimaryKeyColumns),
	)
	values := []interface{}{related.UserID, o.BookingID}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.BookerID.Int = related.UserID
	o.BookerID.Valid = true

	if o.R == nil {
		o.R = &bookingR{
			Booker: related,
		}
	} else {
		o.R.Booker = related
	}

	if related.R == nil {
		related.R = &userR{
			BookerBookings: BookingSlice{o},
		}
	} else {
		related.R.BookerBookings = append(related.R.BookerBookings, o)
	}

	return nil
}

// RemoveBooker relationship.
// Sets o.R.Booker to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *Booking) RemoveBooker(exec boil.Executor, related *User) error {
	var err error

	o.BookerID.Valid = false
	if err = o.Update(exec, "booker_id"); err != nil {
		o.BookerID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.Booker = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.BookerBookings {
		if o.BookerID.Int != ri.BookerID.Int {
			continue
		}

		ln := len(related.R.BookerBookings)
		if ln > 1 && i < ln-1 {
			related.R.BookerBookings[i] = related.R.BookerBookings[ln-1]
		}
		related.R.BookerBookings = related.R.BookerBookings[:ln-1]
		break
	}
	return nil
}

// AddDemos adds the given related objects to the existing relationships
// of the booking, optionally inserting them as new records.
// Appends related to o.R.Demos.
// Sets related.R.Booking appropriately.
func (o *Booking) AddDemos(exec boil.Executor, insert bool, related ...*Demo) error {
	var err error
	for _, rel := range related {
		rel.BookingID.Int = o.BookingID
		rel.BookingID.Valid = true
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			if err = rel.Update(exec, "booking_id"); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}
		}
	}

	if o.R == nil {
		o.R = &bookingR{
			Demos: related,
		}
	} else {
		o.R.Demos = append(o.R.Demos, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &demoR{
				Booking: o,
			}
		} else {
			rel.R.Booking = o
		}
	}
	return nil
}

// SetDemos removes all previously related items of the
// booking replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Booking's Demos accordingly.
// Replaces o.R.Demos with related.
// Sets related.R.Booking's Demos accordingly.
func (o *Booking) SetDemos(exec boil.Executor, insert bool, related ...*Demo) error {
	query := "update \"demos\" set \"booking_id\" = null where \"booking_id\" = $1"
	values := []interface{}{o.BookingID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Demos {
			rel.BookingID.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.Booking = nil
		}

		o.R.Demos = nil
	}
	return o.AddDemos(exec, insert, related...)
}

// RemoveDemos relationships from objects passed in.
// Removes related items from R.Demos (uses pointer comparison, removal does not keep order)
// Sets related.R.Booking.
func (o *Booking) RemoveDemos(exec boil.Executor, related ...*Demo) error {
	var err error
	for _, rel := range related {
		rel.BookingID.Valid = false
		if rel.R != nil {
			rel.R.Booking = nil
		}
		if err = rel.Update(exec, "booking_id"); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Demos {
			if rel != ri {
				continue
			}

			ln := len(o.R.Demos)
			if ln > 1 && i < ln-1 {
				o.R.Demos[i] = o.R.Demos[ln-1]
			}
			o.R.Demos = o.R.Demos[:ln-1]
			break
		}
	}

	return nil
}

// BookingsG retrieves all records.
func BookingsG(mods ...qm.QueryMod) bookingQuery {
	return Bookings(boil.GetDB(), mods...)
}

// Bookings retrieves all the records using an executor.
func Bookings(exec boil.Executor, mods ...qm.QueryMod) bookingQuery {
	mods = append(mods, qm.From("\"bookings\""))
	return bookingQuery{NewQuery(exec, mods...)}
}

// FindBookingG retrieves a single record by ID.
func FindBookingG(bookingID int, selectCols ...string) (*Booking, error) {
	return FindBooking(boil.GetDB(), bookingID, selectCols...)
}

// FindBookingGP retrieves a single record by ID, and panics on error.
func FindBookingGP(bookingID int, selectCols ...string) *Booking {
	retobj, err := FindBooking(boil.GetDB(), bookingID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindBooking retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBooking(exec boil.Executor, bookingID int, selectCols ...string) (*Booking, error) {
	bookingObj := &Booking{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"bookings\" where \"booking_id\"=$1", sel,
	)

	q := queries.Raw(exec, query, bookingID)

	err := q.Bind(bookingObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from bookings")
	}

	return bookingObj, nil
}

// FindBookingP retrieves a single record by ID with an executor, and panics on error.
func FindBookingP(exec boil.Executor, bookingID int, selectCols ...string) *Booking {
	retobj, err := FindBooking(exec, bookingID, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Booking) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Booking) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Booking) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Booking) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no bookings provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(bookingColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	bookingInsertCacheMut.RLock()
	cache, cached := bookingInsertCache[key]
	bookingInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			bookingColumns,
			bookingColumnsWithDefault,
			bookingColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(bookingType, bookingMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(bookingType, bookingMapping, returnColumns)
		if err != nil {
			return err
		}
		cache.query = fmt.Sprintf("INSERT INTO \"bookings\" (\"%s\") VALUES (%s)", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))

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
		return errors.Wrap(err, "models: unable to insert into bookings")
	}

	if !cached {
		bookingInsertCacheMut.Lock()
		bookingInsertCache[key] = cache
		bookingInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Booking record. See Update for
// whitelist behavior description.
func (o *Booking) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Booking record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Booking) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Booking, and panics on error.
// See Update for whitelist behavior description.
func (o *Booking) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Booking.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Booking) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	bookingUpdateCacheMut.RLock()
	cache, cached := bookingUpdateCache[key]
	bookingUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(bookingColumns, bookingPrimaryKeyColumns, whitelist)
		if len(wl) == 0 {
			return errors.New("models: unable to update bookings, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"bookings\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, bookingPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(bookingType, bookingMapping, append(wl, bookingPrimaryKeyColumns...))
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
		return errors.Wrap(err, "models: unable to update bookings row")
	}

	if !cached {
		bookingUpdateCacheMut.Lock()
		bookingUpdateCache[key] = cache
		bookingUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q bookingQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q bookingQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for bookings")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o BookingSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o BookingSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o BookingSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BookingSlice) UpdateAll(exec boil.Executor, cols M) error {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bookingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"UPDATE \"bookings\" SET %s WHERE (\"booking_id\") IN (%s)",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(bookingPrimaryKeyColumns), len(colNames)+1, len(bookingPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in booking slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Booking) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Booking) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Booking) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Booking) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no bookings provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(bookingColumnsWithDefault, o)

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

	bookingUpsertCacheMut.RLock()
	cache, cached := bookingUpsertCache[key]
	bookingUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		var ret []string
		whitelist, ret = strmangle.InsertColumnSet(
			bookingColumns,
			bookingColumnsWithDefault,
			bookingColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)
		update := strmangle.UpdateColumnSet(
			bookingColumns,
			bookingPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert bookings, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(bookingPrimaryKeyColumns))
			copy(conflict, bookingPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"bookings\"", updateOnConflict, ret, update, conflict, whitelist)

		cache.valueMapping, err = queries.BindMapping(bookingType, bookingMapping, whitelist)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(bookingType, bookingMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert for bookings")
	}

	if !cached {
		bookingUpsertCacheMut.Lock()
		bookingUpsertCache[key] = cache
		bookingUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Booking record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Booking) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Booking record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Booking) DeleteG() error {
	if o == nil {
		return errors.New("models: no Booking provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Booking record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Booking) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Booking record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Booking) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Booking provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), bookingPrimaryKeyMapping)
	sql := "DELETE FROM \"bookings\" WHERE \"booking_id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from bookings")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q bookingQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q bookingQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no bookingQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from bookings")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o BookingSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o BookingSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Booking slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o BookingSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BookingSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Booking slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(bookingBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bookingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM \"bookings\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, bookingPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o)*len(bookingPrimaryKeyColumns), 1, len(bookingPrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from booking slice")
	}

	if len(bookingAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Booking) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Booking) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Booking) ReloadG() error {
	if o == nil {
		return errors.New("models: no Booking provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Booking) Reload(exec boil.Executor) error {
	ret, err := FindBooking(exec, o.BookingID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *BookingSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *BookingSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BookingSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty BookingSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BookingSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	bookings := BookingSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bookingPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT \"bookings\".* FROM \"bookings\" WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, bookingPrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o)*len(bookingPrimaryKeyColumns), 1, len(bookingPrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&bookings)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in BookingSlice")
	}

	*o = bookings

	return nil
}

// BookingExists checks if the Booking row exists.
func BookingExists(exec boil.Executor, bookingID int) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from \"bookings\" where \"booking_id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, bookingID)
	}

	row := exec.QueryRow(sql, bookingID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if bookings exists")
	}

	return exists, nil
}

// BookingExistsG checks if the Booking row exists.
func BookingExistsG(bookingID int) (bool, error) {
	return BookingExists(boil.GetDB(), bookingID)
}

// BookingExistsGP checks if the Booking row exists. Panics on error.
func BookingExistsGP(bookingID int) bool {
	e, err := BookingExists(boil.GetDB(), bookingID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// BookingExistsP checks if the Booking row exists. Panics on error.
func BookingExistsP(exec boil.Executor, bookingID int) bool {
	e, err := BookingExists(exec, bookingID)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
