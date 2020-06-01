package adaptermock

import (
	"math/big"

	"github.com/stretchr/testify/mock"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/internal/bignum"
)

type MockRDbConn struct {
	mock.Mock
}

func (conn *MockRDbConn) Begin() (adapter.RDbTx, error) {
	args := conn.Called()
	return args.Get(0).(*MockRDbTx), args.Error(1)
}
func (conn *MockRDbConn) Exec(sql string, args ...interface{}) (adapter.RDbExecResult, error) {
	mockArgs := conn.Called(append([]interface{}{sql}, args...)...)
	result, _ := mockArgs.Get(0).(*MockRDbExecResult)
	return result, mockArgs.Error(1)
}
func (conn *MockRDbConn) Query(sql string, args ...interface{}) (adapter.RDbRowsResult, error) {
	mockArgs := conn.Called(append([]interface{}{sql}, args...)...)
	result, _ := mockArgs.Get(0).(*MockRDbRowsResult)
	return result, mockArgs.Error(1)
}
func (conn *MockRDbConn) QueryRow(sql string, args ...interface{}) adapter.RDbRowResult {
	mockArgs := conn.Called(append([]interface{}{sql}, args...)...)
	result, _ := mockArgs.Get(0).(*MockRDbRowResult)
	return result
}

type MockRDbTx struct {
	mock.Mock
}

func appendStrToISlice(str string, slice ...interface{}) []interface{} {
	result := make([]interface{}, len(slice)+1)
	result[0] = str
	for i, value := range slice {
		result[i+1] = value
	}

	return result
}

func (tx *MockRDbTx) Exec(sql string, args ...interface{}) (adapter.RDbExecResult, error) {
	mockArgs := tx.Called(appendStrToISlice(sql, args...)...)
	result, _ := mockArgs.Get(0).(*MockRDbExecResult)
	return result, mockArgs.Error(1)
}
func (tx *MockRDbTx) Query(sql string, args ...interface{}) (adapter.RDbRowsResult, error) {
	mockArgs := tx.Called(appendStrToISlice(sql, args...)...)
	result, _ := mockArgs.Get(0).(*MockRDbRowsResult)
	return result, mockArgs.Error(1)
}
func (tx *MockRDbTx) QueryRow(sql string, args ...interface{}) adapter.RDbRowResult {
	mockArgs := tx.Called(appendStrToISlice(sql, args...)...)
	result, _ := mockArgs.Get(0).(*MockRDbRowResult)
	return result
}
func (tx *MockRDbTx) Commit() error {
	args := tx.Called()
	return args.Error(0)
}
func (tx *MockRDbTx) Rollback() error {
	args := tx.Called()
	return args.Error(0)
}

type MockRDbRowsResult struct {
	mock.Mock
}

func (result *MockRDbRowsResult) Close() {
	_ = result.Called()
}
func (result *MockRDbRowsResult) Err() error {
	args := result.Called()
	return args.Error(0)
}
func (result *MockRDbRowsResult) ExecResult() adapter.RDbExecResult {
	args := result.Called()
	execResult, _ := args.Get(0).(*MockRDbExecResult)
	return execResult
}
func (result *MockRDbRowsResult) Next() bool {
	args := result.Called()
	return args.Bool(0)
}
func (result *MockRDbRowsResult) Scan(dest ...interface{}) error {
	args := result.Called(dest...)
	return args.Error(0)
}

type MockRDbExecResult struct {
	mock.Mock
}

func (result *MockRDbExecResult) RowsAffected() int64 {
	args := result.Called()
	return args.Get(0).(int64)
}
func (result *MockRDbExecResult) IsInsert() bool {
	args := result.Called()
	return args.Bool(0)
}
func (result *MockRDbExecResult) IsUpdate() bool {
	args := result.Called()
	return args.Bool(0)
}
func (result *MockRDbExecResult) IsDelete() bool {
	args := result.Called()
	return args.Bool(0)
}
func (result *MockRDbExecResult) IsSelect() bool {
	args := result.Called()
	return args.Bool(0)
}
func (result *MockRDbExecResult) String() string {
	args := result.Called()
	return args.String(0)
}

type MockRDbRowResult struct {
	mock.Mock
}

func (result *MockRDbRowResult) Scan(dest ...interface{}) error {
	args := result.Called(dest...)

	return args.Error(0)
}

func MockSQLWithAnyArgs(sql string, argsCount int) []interface{} {
	args := make([]interface{}, argsCount+1)
	args[0] = sql
	for i := 1; i <= argsCount; i += 1 {
		args[i] = mock.Anything
	}

	return args
}

type MockRDbTypeConv struct {
	mock.Mock
}

func (conv *MockRDbTypeConv) Bton(bigInt *big.Int) interface{} {
	args := conv.Called(bigInt)
	return args.Get(0)
}
func (conv *MockRDbTypeConv) Iton(i int) interface{} {
	args := conv.Called(i)
	return args.Get(0)
}
func (conv *MockRDbTypeConv) NtobReader() adapter.RDbNtobReader {
	args := conv.Called()
	result, _ := args.Get(0).(*MockRDbNtobReader)
	return result
}

type MockRDbNtobReader struct {
	mock.Mock
}

func (reader *MockRDbNtobReader) ScannableArg() interface{} {
	args := reader.Called()
	return args.Get(0)
}
func (reader *MockRDbNtobReader) ParseW() (*bignum.WBigInt, error) {
	args := reader.Called()
	result, _ := args.Get(0).(*bignum.WBigInt)
	return result, args.Error(1)
}
func (reader *MockRDbNtobReader) Parse() (*big.Int, error) {
	args := reader.Called()
	result, _ := args.Get(0).(*big.Int)
	return result, args.Error(1)
}
