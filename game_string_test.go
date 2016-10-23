package main

import "testing"

func TestGameStringZero(t *testing.T) {
	gameString := GetGameString(0)
	expected := "No servers available"

	if gameString != expected {
		t.Errorf("TestGameStringZero: Expected \"%s\", got \"%s\"", expected, gameString)
	}
}

func TestGameStringOne(t *testing.T) {
	gameString := GetGameString(1)
	expected := "1 server available"

	if gameString != expected {
		t.Errorf("TestGameStringOne: Expected \"%s\", got \"%s\"", expected, gameString)
	}
}

func TestGameStringFive(t *testing.T) {
	gameString := GetGameString(5)
	expected := "5 servers available"

	if gameString != expected {
		t.Errorf("TestGameStringFive: Expected \"%s\", got \"%s\"", expected, gameString)
	}
}

func TestGameStringFifteen(t *testing.T) {
	gameString := GetGameString(15)
	expected := "15 servers available"

	if gameString != expected {
		t.Errorf("TestGameStringFifteen: Expected \"%s\", got \"%s\"", expected, gameString)
	}
}
