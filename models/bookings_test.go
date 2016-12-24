package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testBookings(t *testing.T) {
	t.Parallel()

	query := Bookings(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testBookingsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = booking.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBookingsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Bookings(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBookingsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := BookingSlice{booking}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testBookingsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := BookingExists(tx, booking.BookingID)
	if err != nil {
		t.Errorf("Unable to check if Booking exists: %s", err)
	}
	if !e {
		t.Errorf("Expected BookingExistsG to return true, but got false.")
	}
}
func testBookingsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	bookingFound, err := FindBooking(tx, booking.BookingID)
	if err != nil {
		t.Error(err)
	}

	if bookingFound == nil {
		t.Error("want a record, got nil")
	}
}
func testBookingsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Bookings(tx).Bind(booking); err != nil {
		t.Error(err)
	}
}

func testBookingsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Bookings(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testBookingsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	bookingOne := &Booking{}
	bookingTwo := &Booking{}
	if err = randomize.Struct(seed, bookingOne, bookingDBTypes, false, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}
	if err = randomize.Struct(seed, bookingTwo, bookingDBTypes, false, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = bookingOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = bookingTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Bookings(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testBookingsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	bookingOne := &Booking{}
	bookingTwo := &Booking{}
	if err = randomize.Struct(seed, bookingOne, bookingDBTypes, false, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}
	if err = randomize.Struct(seed, bookingTwo, bookingDBTypes, false, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = bookingOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = bookingTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func bookingBeforeInsertHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingAfterInsertHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingAfterSelectHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingBeforeUpdateHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingAfterUpdateHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingBeforeDeleteHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingAfterDeleteHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingBeforeUpsertHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func bookingAfterUpsertHook(e boil.Executor, o *Booking) error {
	*o = Booking{}
	return nil
}

func testBookingsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Booking{}
	o := &Booking{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, bookingDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Booking object: %s", err)
	}

	AddBookingHook(boil.BeforeInsertHook, bookingBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	bookingBeforeInsertHooks = []BookingHook{}

	AddBookingHook(boil.AfterInsertHook, bookingAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	bookingAfterInsertHooks = []BookingHook{}

	AddBookingHook(boil.AfterSelectHook, bookingAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	bookingAfterSelectHooks = []BookingHook{}

	AddBookingHook(boil.BeforeUpdateHook, bookingBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	bookingBeforeUpdateHooks = []BookingHook{}

	AddBookingHook(boil.AfterUpdateHook, bookingAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	bookingAfterUpdateHooks = []BookingHook{}

	AddBookingHook(boil.BeforeDeleteHook, bookingBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	bookingBeforeDeleteHooks = []BookingHook{}

	AddBookingHook(boil.AfterDeleteHook, bookingAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	bookingAfterDeleteHooks = []BookingHook{}

	AddBookingHook(boil.BeforeUpsertHook, bookingBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	bookingBeforeUpsertHooks = []BookingHook{}

	AddBookingHook(boil.AfterUpsertHook, bookingAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	bookingAfterUpsertHooks = []BookingHook{}
}
func testBookingsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBookingsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx, bookingColumns...); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBookingToManyDemos(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Booking
	var b, c Demo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, demoDBTypes, false, demoColumnsWithDefault...)
	randomize.Struct(seed, &c, demoDBTypes, false, demoColumnsWithDefault...)
	b.BookingID.Valid = true
	c.BookingID.Valid = true
	b.BookingID.Int = a.BookingID
	c.BookingID.Int = a.BookingID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	demo, err := a.Demos(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range demo {
		if v.BookingID.Int == b.BookingID.Int {
			bFound = true
		}
		if v.BookingID.Int == c.BookingID.Int {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BookingSlice{&a}
	if err = a.L.LoadDemos(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Demos); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Demos = nil
	if err = a.L.LoadDemos(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Demos); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", demo)
	}
}

func testBookingToManyAddOpDemos(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Booking
	var b, c, d, e Demo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Demo{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Demo{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDemos(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.BookingID != first.BookingID.Int {
			t.Error("foreign key was wrong value", a.BookingID, first.BookingID.Int)
		}
		if a.BookingID != second.BookingID.Int {
			t.Error("foreign key was wrong value", a.BookingID, second.BookingID.Int)
		}

		if first.R.Booking != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Booking != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Demos[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Demos[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Demos(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testBookingToManySetOpDemos(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Booking
	var b, c, d, e Demo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Demo{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.SetDemos(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Demos(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDemos(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Demos(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.BookingID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.BookingID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.BookingID != d.BookingID.Int {
		t.Error("foreign key was wrong value", a.BookingID, d.BookingID.Int)
	}
	if a.BookingID != e.BookingID.Int {
		t.Error("foreign key was wrong value", a.BookingID, e.BookingID.Int)
	}

	if b.R.Booking != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Booking != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Booking != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Booking != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.Demos[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.Demos[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testBookingToManyRemoveOpDemos(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Booking
	var b, c, d, e Demo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Demo{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddDemos(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.Demos(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDemos(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.Demos(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.BookingID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.BookingID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Booking != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Booking != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Booking != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Booking != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.Demos) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.Demos[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.Demos[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testBookingToOneUserUsingBooker(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Booking
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	local.BookerID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.BookerID.Int = foreign.UserID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Booker(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.UserID != foreign.UserID {
		t.Errorf("want: %v, got %v", foreign.UserID, check.UserID)
	}

	slice := BookingSlice{&local}
	if err = local.L.LoadBooker(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Booker == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Booker = nil
	if err = local.L.LoadBooker(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Booker == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testBookingToOneSetOpUserUsingBooker(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Booking
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetBooker(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Booker != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.BookerBookings[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.BookerID.Int != x.UserID {
			t.Error("foreign key was wrong value", a.BookerID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.BookerID.Int))
		reflect.Indirect(reflect.ValueOf(&a.BookerID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.BookerID.Int != x.UserID {
			t.Error("foreign key was wrong value", a.BookerID.Int, x.UserID)
		}
	}
}

func testBookingToOneRemoveOpUserUsingBooker(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Booking
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetBooker(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveBooker(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Booker(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Booker != nil {
		t.Error("R struct entry should be nil")
	}

	if a.BookerID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.BookerBookings) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testBookingsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = booking.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testBookingsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := BookingSlice{booking}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testBookingsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Bookings(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	bookingDBTypes = map[string]string{`BookedTime`: `timestamp without time zone`, `BookerID`: `integer`, `BookingID`: `integer`, `ServerName`: `character varying`, `UnbookedTime`: `timestamp without time zone`}
	_              = bytes.MinRead
)

func testBookingsUpdate(t *testing.T) {
	t.Parallel()

	if len(bookingColumns) == len(bookingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	if err = booking.Update(tx); err != nil {
		t.Error(err)
	}
}

func testBookingsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(bookingColumns) == len(bookingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	booking := &Booking{}
	if err = randomize.Struct(seed, booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, booking, bookingDBTypes, true, bookingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(bookingColumns, bookingPrimaryKeyColumns) {
		fields = bookingColumns
	} else {
		fields = strmangle.SetComplement(
			bookingColumns,
			bookingPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(booking))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := BookingSlice{booking}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testBookingsUpsert(t *testing.T) {
	t.Parallel()

	if len(bookingColumns) == len(bookingPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	booking := Booking{}
	if err = randomize.Struct(seed, &booking, bookingDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = booking.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Booking: %s", err)
	}

	count, err := Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &booking, bookingDBTypes, false, bookingPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	if err = booking.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Booking: %s", err)
	}

	count, err = Bookings(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
