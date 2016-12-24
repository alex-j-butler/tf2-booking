package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Demos", testDemos)
	t.Run("Bookings", testBookings)
	t.Run("Users", testUsers)
	t.Run("DemoUsers", testDemoUsers)
}

func TestDelete(t *testing.T) {
	t.Run("Demos", testDemosDelete)
	t.Run("Bookings", testBookingsDelete)
	t.Run("Users", testUsersDelete)
	t.Run("DemoUsers", testDemoUsersDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Demos", testDemosQueryDeleteAll)
	t.Run("Bookings", testBookingsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
	t.Run("DemoUsers", testDemoUsersQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Demos", testDemosSliceDeleteAll)
	t.Run("Bookings", testBookingsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
	t.Run("DemoUsers", testDemoUsersSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Demos", testDemosExists)
	t.Run("Bookings", testBookingsExists)
	t.Run("Users", testUsersExists)
	t.Run("DemoUsers", testDemoUsersExists)
}

func TestFind(t *testing.T) {
	t.Run("Demos", testDemosFind)
	t.Run("Bookings", testBookingsFind)
	t.Run("Users", testUsersFind)
	t.Run("DemoUsers", testDemoUsersFind)
}

func TestBind(t *testing.T) {
	t.Run("Demos", testDemosBind)
	t.Run("Bookings", testBookingsBind)
	t.Run("Users", testUsersBind)
	t.Run("DemoUsers", testDemoUsersBind)
}

func TestOne(t *testing.T) {
	t.Run("Demos", testDemosOne)
	t.Run("Bookings", testBookingsOne)
	t.Run("Users", testUsersOne)
	t.Run("DemoUsers", testDemoUsersOne)
}

func TestAll(t *testing.T) {
	t.Run("Demos", testDemosAll)
	t.Run("Bookings", testBookingsAll)
	t.Run("Users", testUsersAll)
	t.Run("DemoUsers", testDemoUsersAll)
}

func TestCount(t *testing.T) {
	t.Run("Demos", testDemosCount)
	t.Run("Bookings", testBookingsCount)
	t.Run("Users", testUsersCount)
	t.Run("DemoUsers", testDemoUsersCount)
}

func TestHooks(t *testing.T) {
	t.Run("Demos", testDemosHooks)
	t.Run("Bookings", testBookingsHooks)
	t.Run("Users", testUsersHooks)
	t.Run("DemoUsers", testDemoUsersHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Demos", testDemosInsert)
	t.Run("Demos", testDemosInsertWhitelist)
	t.Run("Bookings", testBookingsInsert)
	t.Run("Bookings", testBookingsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
	t.Run("DemoUsers", testDemoUsersInsert)
	t.Run("DemoUsers", testDemoUsersInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("DemoToBookingUsingBooking", testDemoToOneBookingUsingBooking)
	t.Run("BookingToUserUsingBooker", testBookingToOneUserUsingBooker)
	t.Run("DemoUserToDemoUsingDemo", testDemoUserToOneDemoUsingDemo)
	t.Run("DemoUserToUserUsingUser", testDemoUserToOneUserUsingUser)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("DemoToDemoUsers", testDemoToManyDemoUsers)
	t.Run("BookingToDemos", testBookingToManyDemos)
	t.Run("UserToBookerBookings", testUserToManyBookerBookings)
	t.Run("UserToDemoUsers", testUserToManyDemoUsers)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("DemoToBookingUsingBooking", testDemoToOneSetOpBookingUsingBooking)
	t.Run("BookingToUserUsingBooker", testBookingToOneSetOpUserUsingBooker)
	t.Run("DemoUserToDemoUsingDemo", testDemoUserToOneSetOpDemoUsingDemo)
	t.Run("DemoUserToUserUsingUser", testDemoUserToOneSetOpUserUsingUser)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("DemoToBookingUsingBooking", testDemoToOneRemoveOpBookingUsingBooking)
	t.Run("BookingToUserUsingBooker", testBookingToOneRemoveOpUserUsingBooker)
	t.Run("DemoUserToDemoUsingDemo", testDemoUserToOneRemoveOpDemoUsingDemo)
	t.Run("DemoUserToUserUsingUser", testDemoUserToOneRemoveOpUserUsingUser)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("DemoToDemoUsers", testDemoToManyAddOpDemoUsers)
	t.Run("BookingToDemos", testBookingToManyAddOpDemos)
	t.Run("UserToBookerBookings", testUserToManyAddOpBookerBookings)
	t.Run("UserToDemoUsers", testUserToManyAddOpDemoUsers)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("DemoToDemoUsers", testDemoToManySetOpDemoUsers)
	t.Run("BookingToDemos", testBookingToManySetOpDemos)
	t.Run("UserToBookerBookings", testUserToManySetOpBookerBookings)
	t.Run("UserToDemoUsers", testUserToManySetOpDemoUsers)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("DemoToDemoUsers", testDemoToManyRemoveOpDemoUsers)
	t.Run("BookingToDemos", testBookingToManyRemoveOpDemos)
	t.Run("UserToBookerBookings", testUserToManyRemoveOpBookerBookings)
	t.Run("UserToDemoUsers", testUserToManyRemoveOpDemoUsers)
}

func TestReload(t *testing.T) {
	t.Run("Demos", testDemosReload)
	t.Run("Bookings", testBookingsReload)
	t.Run("Users", testUsersReload)
	t.Run("DemoUsers", testDemoUsersReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Demos", testDemosReloadAll)
	t.Run("Bookings", testBookingsReloadAll)
	t.Run("Users", testUsersReloadAll)
	t.Run("DemoUsers", testDemoUsersReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Demos", testDemosSelect)
	t.Run("Bookings", testBookingsSelect)
	t.Run("Users", testUsersSelect)
	t.Run("DemoUsers", testDemoUsersSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Demos", testDemosUpdate)
	t.Run("Bookings", testBookingsUpdate)
	t.Run("Users", testUsersUpdate)
	t.Run("DemoUsers", testDemoUsersUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Demos", testDemosSliceUpdateAll)
	t.Run("Bookings", testBookingsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
	t.Run("DemoUsers", testDemoUsersSliceUpdateAll)
}

func TestUpsert(t *testing.T) {
	t.Run("Demos", testDemosUpsert)
	t.Run("Bookings", testBookingsUpsert)
	t.Run("Users", testUsersUpsert)
	t.Run("DemoUsers", testDemoUsersUpsert)
}
