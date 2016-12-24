package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testSchemaMigrations(t *testing.T) {
	t.Parallel()

	query := SchemaMigrations(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testSchemaMigrationsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = schemaMigration.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSchemaMigrationsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = SchemaMigrations(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSchemaMigrationsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SchemaMigrationSlice{schemaMigration}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testSchemaMigrationsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := SchemaMigrationExists(tx, schemaMigration.Version)
	if err != nil {
		t.Errorf("Unable to check if SchemaMigration exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SchemaMigrationExistsG to return true, but got false.")
	}
}
func testSchemaMigrationsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	schemaMigrationFound, err := FindSchemaMigration(tx, schemaMigration.Version)
	if err != nil {
		t.Error(err)
	}

	if schemaMigrationFound == nil {
		t.Error("want a record, got nil")
	}
}
func testSchemaMigrationsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = SchemaMigrations(tx).Bind(schemaMigration); err != nil {
		t.Error(err)
	}
}

func testSchemaMigrationsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := SchemaMigrations(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testSchemaMigrationsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigrationOne := &SchemaMigration{}
	schemaMigrationTwo := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigrationOne, schemaMigrationDBTypes, false, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}
	if err = randomize.Struct(seed, schemaMigrationTwo, schemaMigrationDBTypes, false, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigrationOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = schemaMigrationTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := SchemaMigrations(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testSchemaMigrationsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	schemaMigrationOne := &SchemaMigration{}
	schemaMigrationTwo := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigrationOne, schemaMigrationDBTypes, false, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}
	if err = randomize.Struct(seed, schemaMigrationTwo, schemaMigrationDBTypes, false, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigrationOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = schemaMigrationTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func schemaMigrationBeforeInsertHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationAfterInsertHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationAfterSelectHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationBeforeUpdateHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationAfterUpdateHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationBeforeDeleteHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationAfterDeleteHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationBeforeUpsertHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func schemaMigrationAfterUpsertHook(e boil.Executor, o *SchemaMigration) error {
	*o = SchemaMigration{}
	return nil
}

func testSchemaMigrationsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &SchemaMigration{}
	o := &SchemaMigration{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, schemaMigrationDBTypes, false); err != nil {
		t.Errorf("Unable to randomize SchemaMigration object: %s", err)
	}

	AddSchemaMigrationHook(boil.BeforeInsertHook, schemaMigrationBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	schemaMigrationBeforeInsertHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.AfterInsertHook, schemaMigrationAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	schemaMigrationAfterInsertHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.AfterSelectHook, schemaMigrationAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	schemaMigrationAfterSelectHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.BeforeUpdateHook, schemaMigrationBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	schemaMigrationBeforeUpdateHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.AfterUpdateHook, schemaMigrationAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	schemaMigrationAfterUpdateHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.BeforeDeleteHook, schemaMigrationBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	schemaMigrationBeforeDeleteHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.AfterDeleteHook, schemaMigrationAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	schemaMigrationAfterDeleteHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.BeforeUpsertHook, schemaMigrationBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	schemaMigrationBeforeUpsertHooks = []SchemaMigrationHook{}

	AddSchemaMigrationHook(boil.AfterUpsertHook, schemaMigrationAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	schemaMigrationAfterUpsertHooks = []SchemaMigrationHook{}
}
func testSchemaMigrationsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSchemaMigrationsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx, schemaMigrationColumns...); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSchemaMigrationsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = schemaMigration.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testSchemaMigrationsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := SchemaMigrationSlice{schemaMigration}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testSchemaMigrationsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := SchemaMigrations(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	schemaMigrationDBTypes = map[string]string{`Version`: `integer`}
	_                      = bytes.MinRead
)

func testSchemaMigrationsUpdate(t *testing.T) {
	t.Parallel()

	if len(schemaMigrationColumns) == len(schemaMigrationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	if err = schemaMigration.Update(tx); err != nil {
		t.Error(err)
	}
}

func testSchemaMigrationsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(schemaMigrationColumns) == len(schemaMigrationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	schemaMigration := &SchemaMigration{}
	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, schemaMigration, schemaMigrationDBTypes, true, schemaMigrationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(schemaMigrationColumns, schemaMigrationPrimaryKeyColumns) {
		fields = schemaMigrationColumns
	} else {
		fields = strmangle.SetComplement(
			schemaMigrationColumns,
			schemaMigrationPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(schemaMigration))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := SchemaMigrationSlice{schemaMigration}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testSchemaMigrationsUpsert(t *testing.T) {
	t.Parallel()

	if len(schemaMigrationColumns) == len(schemaMigrationPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	schemaMigration := SchemaMigration{}
	if err = randomize.Struct(seed, &schemaMigration, schemaMigrationDBTypes, true); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = schemaMigration.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert SchemaMigration: %s", err)
	}

	count, err := SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &schemaMigration, schemaMigrationDBTypes, false, schemaMigrationPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize SchemaMigration struct: %s", err)
	}

	if err = schemaMigration.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert SchemaMigration: %s", err)
	}

	count, err = SchemaMigrations(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
