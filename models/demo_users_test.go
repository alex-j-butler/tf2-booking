package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testDemoUsers(t *testing.T) {
	t.Parallel()

	query := DemoUsers(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testDemoUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = demoUser.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemoUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = DemoUsers(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDemoUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := DemoUserSlice{demoUser}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testDemoUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := DemoUserExists(tx, demoUser.RefID)
	if err != nil {
		t.Errorf("Unable to check if DemoUser exists: %s", err)
	}
	if !e {
		t.Errorf("Expected DemoUserExistsG to return true, but got false.")
	}
}
func testDemoUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	demoUserFound, err := FindDemoUser(tx, demoUser.RefID)
	if err != nil {
		t.Error(err)
	}

	if demoUserFound == nil {
		t.Error("want a record, got nil")
	}
}
func testDemoUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = DemoUsers(tx).Bind(demoUser); err != nil {
		t.Error(err)
	}
}

func testDemoUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := DemoUsers(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testDemoUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUserOne := &DemoUser{}
	demoUserTwo := &DemoUser{}
	if err = randomize.Struct(seed, demoUserOne, demoUserDBTypes, false, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}
	if err = randomize.Struct(seed, demoUserTwo, demoUserDBTypes, false, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUserOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = demoUserTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := DemoUsers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testDemoUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	demoUserOne := &DemoUser{}
	demoUserTwo := &DemoUser{}
	if err = randomize.Struct(seed, demoUserOne, demoUserDBTypes, false, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}
	if err = randomize.Struct(seed, demoUserTwo, demoUserDBTypes, false, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUserOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = demoUserTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func demoUserBeforeInsertHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserAfterInsertHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserAfterSelectHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserBeforeUpdateHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserAfterUpdateHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserBeforeDeleteHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserAfterDeleteHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserBeforeUpsertHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func demoUserAfterUpsertHook(e boil.Executor, o *DemoUser) error {
	*o = DemoUser{}
	return nil
}

func testDemoUsersHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &DemoUser{}
	o := &DemoUser{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, demoUserDBTypes, false); err != nil {
		t.Errorf("Unable to randomize DemoUser object: %s", err)
	}

	AddDemoUserHook(boil.BeforeInsertHook, demoUserBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	demoUserBeforeInsertHooks = []DemoUserHook{}

	AddDemoUserHook(boil.AfterInsertHook, demoUserAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	demoUserAfterInsertHooks = []DemoUserHook{}

	AddDemoUserHook(boil.AfterSelectHook, demoUserAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	demoUserAfterSelectHooks = []DemoUserHook{}

	AddDemoUserHook(boil.BeforeUpdateHook, demoUserBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	demoUserBeforeUpdateHooks = []DemoUserHook{}

	AddDemoUserHook(boil.AfterUpdateHook, demoUserAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	demoUserAfterUpdateHooks = []DemoUserHook{}

	AddDemoUserHook(boil.BeforeDeleteHook, demoUserBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	demoUserBeforeDeleteHooks = []DemoUserHook{}

	AddDemoUserHook(boil.AfterDeleteHook, demoUserAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	demoUserAfterDeleteHooks = []DemoUserHook{}

	AddDemoUserHook(boil.BeforeUpsertHook, demoUserBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	demoUserBeforeUpsertHooks = []DemoUserHook{}

	AddDemoUserHook(boil.AfterUpsertHook, demoUserAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	demoUserAfterUpsertHooks = []DemoUserHook{}
}
func testDemoUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemoUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx, demoUserColumns...); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDemoUserToOneDemoUsingDemo(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local DemoUser
	var foreign Demo

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, demoDBTypes, true, demoColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Demo struct: %s", err)
	}

	local.DemoID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.DemoID.Int = foreign.DemoID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Demo(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.DemoID != foreign.DemoID {
		t.Errorf("want: %v, got %v", foreign.DemoID, check.DemoID)
	}

	slice := DemoUserSlice{&local}
	if err = local.L.LoadDemo(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Demo == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Demo = nil
	if err = local.L.LoadDemo(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Demo == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testDemoUserToOneUserUsingUser(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local DemoUser
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, true, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	local.UserID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.UserID.Int = foreign.UserID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.User(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.UserID != foreign.UserID {
		t.Errorf("want: %v, got %v", foreign.UserID, check.UserID)
	}

	slice := DemoUserSlice{&local}
	if err = local.L.LoadUser(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.User = nil
	if err = local.L.LoadUser(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.User == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testDemoUserToOneSetOpDemoUsingDemo(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a DemoUser
	var b, c Demo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Demo{&b, &c} {
		err = a.SetDemo(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Demo != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DemoUsers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DemoID.Int != x.DemoID {
			t.Error("foreign key was wrong value", a.DemoID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DemoID.Int))
		reflect.Indirect(reflect.ValueOf(&a.DemoID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DemoID.Int != x.DemoID {
			t.Error("foreign key was wrong value", a.DemoID.Int, x.DemoID)
		}
	}
}

func testDemoUserToOneRemoveOpDemoUsingDemo(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a DemoUser
	var b Demo

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, demoDBTypes, false, strmangle.SetComplement(demoPrimaryKeyColumns, demoColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetDemo(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveDemo(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Demo(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Demo != nil {
		t.Error("R struct entry should be nil")
	}

	if a.DemoID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.DemoUsers) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testDemoUserToOneSetOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a DemoUser
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
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
		err = a.SetUser(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.User != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DemoUsers[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.UserID.Int != x.UserID {
			t.Error("foreign key was wrong value", a.UserID.Int)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UserID.Int))
		reflect.Indirect(reflect.ValueOf(&a.UserID.Int)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.UserID.Int != x.UserID {
			t.Error("foreign key was wrong value", a.UserID.Int, x.UserID)
		}
	}
}

func testDemoUserToOneRemoveOpUserUsingUser(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a DemoUser
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, demoUserDBTypes, false, strmangle.SetComplement(demoUserPrimaryKeyColumns, demoUserColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetUser(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveUser(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.User(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.User != nil {
		t.Error("R struct entry should be nil")
	}

	if a.UserID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.DemoUsers) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testDemoUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = demoUser.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testDemoUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := DemoUserSlice{demoUser}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testDemoUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := DemoUsers(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	demoUserDBTypes = map[string]string{`DemoID`: `integer`, `RefID`: `integer`, `UserID`: `integer`}
	_               = bytes.MinRead
)

func testDemoUsersUpdate(t *testing.T) {
	t.Parallel()

	if len(demoUserColumns) == len(demoUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	if err = demoUser.Update(tx); err != nil {
		t.Error(err)
	}
}

func testDemoUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(demoUserColumns) == len(demoUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	demoUser := &DemoUser{}
	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, demoUser, demoUserDBTypes, true, demoUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(demoUserColumns, demoUserPrimaryKeyColumns) {
		fields = demoUserColumns
	} else {
		fields = strmangle.SetComplement(
			demoUserColumns,
			demoUserPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(demoUser))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := DemoUserSlice{demoUser}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testDemoUsersUpsert(t *testing.T) {
	t.Parallel()

	if len(demoUserColumns) == len(demoUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	demoUser := DemoUser{}
	if err = randomize.Struct(seed, &demoUser, demoUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = demoUser.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert DemoUser: %s", err)
	}

	count, err := DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &demoUser, demoUserDBTypes, false, demoUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DemoUser struct: %s", err)
	}

	if err = demoUser.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert DemoUser: %s", err)
	}

	count, err = DemoUsers(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
