package adapter

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/crypto-com/chainindex/internal/bignum"
)

// Relational database interface

var (
	ErrBuildSQLStmt = fmt.Errorf("error building SQL statement: %w", ErrRepoPrepare)
	ErrTypeConv     = fmt.Errorf("error converting between types: %w", ErrRepoPrepare)

	// when trying to scan a null row
	ErrNoRows = errors.New("no rows in result set")
)

type RDbConn interface {
	Begin() (RDbTx, error)
	Exec(sql string, args ...interface{}) (RDbExecResult, error)
	Query(sql string, args ...interface{}) (RDbRowsResult, error)
	QueryRow(sql string, args ...interface{}) RDbRowResult
}

type RDbTx interface {
	Exec(sql string, args ...interface{}) (RDbExecResult, error)
	Query(sql string, args ...interface{}) (RDbRowsResult, error)
	QueryRow(sql string, args ...interface{}) RDbRowResult
	Commit() error
	Rollback() error
}

// Implementing RDbConn and RDbTx interface automatically fullfills RDbRunner
type RDbRunner interface {
	Exec(sql string, args ...interface{}) (RDbExecResult, error)
	Query(sql string, args ...interface{}) (RDbRowsResult, error)
	QueryRow(sql string, args ...interface{}) RDbRowResult
}

type RDbExecResult interface {
	RowsAffected() int64
	IsInsert() bool
	IsUpdate() bool
	IsDelete() bool
	IsSelect() bool
	String() string
}

type RDbRowsResult interface {
	Close()
	Err() error
	ExecResult() RDbExecResult
	Next() bool
	Scan(dest ...interface{}) error
}

type RDbRowResult interface {
	Scan(dest ...interface{}) error
}

type RDbTypeConv interface {
	// convert big.Int to native number pointer. Return nil if big.Int is nil
	Bton(*big.Int) interface{}
	Iton(int) interface{}
	// native number reader and parser to big.Int
	NtobReader() RDbNtobReader

	// convert time.Time to native time format pointer. Return nil if time
	// is nil
	Tton(*time.Time) interface{}
	NtotReader() RDbNtotReader
}

type RDbNtobReader interface {
	// returns reference to a scannable type
	ScannableArg() interface{}
	// parse the scannable argument reference to big.Int
	Parse() (*big.Int, error)
	ParseW() (*bignum.WBigInt, error)
}

type RDbNtotReader interface {
	// returns reference to a scannable type
	ScannableArg() interface{}
	// parse the scannable argument reference to time.Time
	Parse() (*time.Time, error)
}
