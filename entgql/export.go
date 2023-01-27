package entgql

import (
	"entgo.io/contrib/entgql/runtime"
)

type (
	// TxOpener represents types than can open transactions.
	TxOpener = runtime.TxOpener
	// The TxOpenerFunc type is an adapter to allow the use of
	// ordinary functions as tx openers.
	TxOpenerFunc = runtime.TxOpenerFunc
	// Transactioner for graphql mutations.
	Transactioner = runtime.Transactioner
)

var (
	// ErrNodeNotFound creates a node not found graphql error.
	ErrNodeNotFound = runtime.ErrNodeNotFound
)
