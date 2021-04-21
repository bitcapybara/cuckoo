package main

type Transport interface {
	Trigger() error
}
