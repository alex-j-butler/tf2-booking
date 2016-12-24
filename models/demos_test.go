package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testDemos(t *testing.T) {
	t.Parallel()

	query := Demos(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testDemosDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = demo.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemosQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Demos(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemosSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := DemoSlice{demo}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testDemosExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := DemoExists(tx, demo.DemoID)
	if err != nil {
		t.Errorf("Unable to check if Demo exists: %s", err)
	}
	if !e {
		t.Errorf("Expected DemoExistsG to return true, but got false.")
	}
}
func testDemosFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	demoFound, err := FindDemo(tx, demo.DemoID)
	if err != nil {
		t.Error(err)
	}

	if demoFound == nil {
		t.Error("want a record, got nil")
	}
}
func testDemosBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Demos(tx).Bind(demo); err != nil {
		t.Error(err)
	}
}

func testDemosOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Demos(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testDemosAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoOne := &Demo{}
	demoTwo := &Demo{}
	if err = randomize.Struct(seed, demoOne, demoDBTypes, false, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}
	if err = randomize.Struct(seed, demoTwo, demoDBTypes, false, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = demoTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Demos(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testDemosCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	demoOne := &Demo{}
	demoTwo := &Demo{}
	if err = randomize.Struct(seed, demoOne, demoDBTypes, false, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}
	if err = randomize.Struct(seed, demoTwo, demoDBTypes, false, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = demoTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func demoBeforeInsertHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoAfterInsertHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoAfterSelectHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoBeforeUpdateHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoAfterUpdateHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoBeforeDeleteHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoAfterDeleteHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoBeforeUpsertHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func demoAfterUpsertHook(e boil.Executor, o *Demo) error {
	*o = Demo{}
	return nil
}

func testDemosHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Demo{}
	o := &Demo{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, demoDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Demo object: %s", err)
	}

	AddDemoHook(boil.BeforeInsertHook, demoBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	demoBeforeInsertHooks = []DemoHook{}

	AddDemoHook(boil.AfterInsertHook, demoAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	demoAfterInsertHooks = []DemoHook{}

	AddDemoHook(boil.AfterSelectHook, demoAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	demoAfterSelectHooks = []DemoHook{}

	AddDemoHook(boil.BeforeUpdateHook, demoBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	demoBeforeUpdateHooks = []DemoHook{}

	AddDemoHook(boil.AfterUpdateHook, demoAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	demoAfterUpdateHooks = []DemoHook{}

	AddDemoHook(boil.BeforeDeleteHook, demoBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	demoBeforeDeleteHooks = []DemoHook{}

	AddDemoHook(boil.AfterDeleteHook, demoAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	demoAfterDeleteHooks = []DemoHook{}

	AddDemoHook(boil.BeforeUpsertHook, demoBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	demoBeforeUpsertHooks = []DemoHook{}

	AddDemoHook(boil.AfterUpsertHook, demoAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	demoAfterUpsertHooks = []DemoHook{}
}
func testDemosInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemosInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx, demoColumns...); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemoToManyDemoUsers(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Demo
	var b, c DemoUser

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, demoUserDBTypes, false, demoUserColumnsWithDefault...)
	randomize.Struct(seed, &c, demoUserDBTypes, false, demoUserColumnsWithDefault...)
	b.DemoID.Valid = true
	c.DemoID.Valid = true
	b.DemoID.Int = a.DemoID
	c.DemoID.Int = a.DemoID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	demoUser, err := a.DemoUsers(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range demoUser {
		if v.DemoID.Int == b.DemoID.Int {
			bFound = true
		}
		if v.DemoID.Int == c.DemoID.Int {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := DemoSlice{&a}
	if err = a.L.LoadDemoUsers(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DemoUsers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DemoUsers = nil
	if err = a.L.LoadDemoUsers(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DemoUsers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", demoUser)
	}
}

func testDemoToManyAddOpDemoUsers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Demo
	var b, c, d, e DemoUser

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*DemoUser{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*DemoUser{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDemoUsers(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.DemoID != first.DemoID.Int {
			t.Error("foreign key was wrong value", a.DemoID, first.DemoID.Int)
		}
		if a.DemoID != second.DemoID.Int {
			t.Error("foreign key was wrong value", a.DemoID, second.DemoID.Int)
		}

		if first.R.Demo != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Demo != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DemoUsers[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DemoUsers[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DemoUsers(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testDemoToManySetOpDemoUsers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Demo
	var b, c, d, e DemoUser

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*DemoUser{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
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

	err = a.SetDemoUsers(tx, false, &b, &c)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DemoUsers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	err = a.SetDemoUsers(tx, true, &d, &e)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DemoUsers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DemoID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DemoID.Valid {
		t.Error("want c's foreign key value to be nil")
	}
	if a.DemoID != d.DemoID.Int {
		t.Error("foreign key was wrong value", a.DemoID, d.DemoID.Int)
	}
	if a.DemoID != e.DemoID.Int {
		t.Error("foreign key was wrong value", a.DemoID, e.DemoID.Int)
	}

	if b.R.Demo != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Demo != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Demo != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}
	if e.R.Demo != &a {
		t.Error("relationship was not added properly to the foreign struct")
	}

	if a.R.DemoUsers[0] != &d {
		t.Error("relationship struct slice not set to correct value")
	}
	if a.R.DemoUsers[1] != &e {
		t.Error("relationship struct slice not set to correct value")
	}
}

func testDemoToManyRemoveOpDemoUsers(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Demo
	var b, c, d, e DemoUser

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*DemoUser{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	err = a.AddDemoUsers(tx, true, foreigners...)
	if err != nil {
		t.Fatal(err)
	}

	count, err := a.DemoUsers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Error("count was wrong:", count)
	}

	err = a.RemoveDemoUsers(tx, foreigners[:2]...)
	if err != nil {
		t.Fatal(err)
	}

	count, err = a.DemoUsers(tx).Count()
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Error("count was wrong:", count)
	}

	if b.DemoID.Valid {
		t.Error("want b's foreign key value to be nil")
	}
	if c.DemoID.Valid {
		t.Error("want c's foreign key value to be nil")
	}

	if b.R.Demo != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if c.R.Demo != nil {
		t.Error("relationship was not removed properly from the foreign struct")
	}
	if d.R.Demo != &a {
		t.Error("relationship to a should have been preserved")
	}
	if e.R.Demo != &a {
		t.Error("relationship to a should have been preserved")
	}

	if len(a.R.DemoUsers) != 2 {
		t.Error("should have preserved two relationships")
	}

	// Removal doesn't do a stable deletion for performance so we have to flip the order
	if a.R.DemoUsers[1] != &d {
		t.Error("relationship to d should have been preserved")
	}
	if a.R.DemoUsers[0] != &e {
		t.Error("relationship to e should have been preserved")
	}
}

func testDemoToOneBookingUsingBooking(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Demo
	var foreign Booking

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, bookingDBTypes, true, bookingColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Booking struct: %s", err)
	}

	local.BookingID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.BookingID.Int = foreign.BookingID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Booking(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.BookingID != foreign.BookingID {
		t.Errorf("want: %v, got %v", foreign.BookingID, check.BookingID)
	}

	slice := DemoSlice{&local}
	if err = local.L.LoadBooking(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Booking == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Booking = nil
	if err = local.L.LoadBooking(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Booking == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testDemoToOneSetOpBookingUsingBooking(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Demo
	var b, c Booking

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Booking{&b, &c} {
		err = a.SetBooking(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Booking != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Demos[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.BookingID.Int != x.BookingID {
			t.Error("foreign key was wrong value", a.BookingID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.BookingID.Int))
		reflect.Indirect(reflect.ValueOf(&a.BookingID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.BookingID.Int != x.BookingID {
			t.Error("foreign key was wrong value", a.BookingID.Int, x.BookingID)
		}
	}
}

func testDemoToOneRemoveOpBookingUsingBooking(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Demo
	var b Booking

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, bookingDBTypes, false, strmangle.SetComplement(bookingPrimaryKeyColumns, bookingColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetBooking(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveBooking(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Booking(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Booking != nil {
		t.Error("R struct entry should be nil")
	}

	if a.BookingID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.Demos) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testDemosReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = demo.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testDemosReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := DemoSlice{demo}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testDemosSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Demos(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	demoDBTypes = map[string]string{`BookingID`: `integer`, `Configuration`: `character varying`, `DemoID`: `integer`, `MapName`: `character varying`, `Name`: `character varying`, `URL`: `character varying`, `UploadedTime`: `timestamp without time zone`}
	_           = bytes.MinRead
)

func testDemosUpdate(t *testing.T) {
	t.Parallel()

	if len(demoColumns) == len(demoPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	if err = demo.Update(tx); err != nil {
		t.Error(err)
	}
}

func testDemosSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(demoColumns) == len(demoPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	demo := &Demo{}
	if err = randomize.Struct(seed, demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, demo, demoDBTypes, true, demoPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(demoColumns, demoPrimaryKeyColumns) {
		fields = demoColumns
	} else {
		fields = strmangle.SetComplement(
			demoColumns,
			demoPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(demo))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := DemoSlice{demo}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testDemosUpsert(t *testing.T) {
	t.Parallel()

	if len(demoColumns) == len(demoPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	demo := Demo{}
	if err = randomize.Struct(seed, &demo, demoDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demo.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Demo: %s", err)
	}

	count, err := Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &demo, demoDBTypes, false, demoPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	if err = demo.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Demo: %s", err)
	}

	count, err = Demos(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
