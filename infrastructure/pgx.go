package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/go-querystring/query"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/crypto-com/chainindex/adapter"
	"github.com/crypto-com/chainindex/internal/bignum"
	"github.com/crypto-com/chainindex/usecase"
)

var (
	BIGINT_Int0  = bignum.Int0()
	BIGINT_Int10 = bignum.Int10()
)

var PostgresStmtBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewPgxConn(config PgxConnConfig, logger usecase.Logger) (*pgx.Conn, error) {
	pgxConfig, err := pgx.ParseConfig(config.ToURL())
	if err != nil {
		return nil, err
	}
	pgxConfig.Logger = NewPgxLoggerAdapter(logger)

	return pgx.ConnectConfig(context.Background(), pgxConfig)
}

func NewPgxConnPool(config PgxConnPoolConfig, logger usecase.Logger) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(config.ToURL())
	if err != nil {
		return nil, err
	}
	pgxConfig.ConnConfig.Logger = NewPgxLoggerAdapter(logger)

	return pgxpool.ConnectConfig(context.Background(), pgxConfig)
}

type PgxConnLike interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type PgxConnConfig struct {
	Host     string
	Port     uint32
	Username string
	Password string
	Database string
	SSL      bool
}

type PgxConnPoolConfig struct {
	PgxConnConfig     `url:"-"`
	MaxConns          int32         `url:"pool_max_conns"`
	MinConns          int32         `url:"pool_min_conns"`
	MaxConnLifeTime   time.Duration `url:"pool_max_conn_lifetime"`
	MaxConnIdleTime   time.Duration `url:"pool_max_conn_idle_time"`
	HealthCheckPeriod time.Duration `url:"pool_health_check_period"`
}

func (config *PgxConnConfig) ToURL() string {
	var authStr string
	if config.Username != "" || config.Password != "" {
		authStr = config.Username + ":" + config.Password + "@"
	}

	connStr := fmt.Sprintf("postgres://%s%s:%d/%s", authStr, config.Host, config.Port, config.Database)
	if config.SSL {
		return connStr
	} else {
		return connStr + "?sslmode=disable"
	}
}

func (config *PgxConnPoolConfig) ToURL() string {
	var authStr string
	if config.Username != "" || config.Password != "" {
		authStr = config.Username + ":" + config.Password + "@"
	}
	connStr := fmt.Sprintf("postgres://%s%s:%d/%s", authStr, config.Host, config.Port, config.Database)

	queryValues, err := query.Values(config)
	if err != nil {
		panic(fmt.Sprintf("error parsing Pgx connection config: %v", err))
	}
	if !config.SSL {
		queryValues.Set("sslmode", "disable")
	}
	return connStr + "?" + queryValues.Encode()
}

type PgxRDbConn struct {
	pgxConn PgxConnLike
}

func NewPgxRDbConn(pgxConn PgxConnLike) *PgxRDbConn {
	return &PgxRDbConn{
		pgxConn,
	}
}

func (conn *PgxRDbConn) Begin() (adapter.RDbTx, error) {
	tx, err := conn.pgxConn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	return &PgxRDbTx{
		tx,
	}, nil
}
func (conn *PgxRDbConn) Exec(sql string, args ...interface{}) (adapter.RDbExecResult, error) {
	result, err := conn.pgxConn.Exec(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	return &PgxRDbExecResult{
		result,
	}, nil
}
func (conn *PgxRDbConn) Query(sql string, args ...interface{}) (adapter.RDbRowsResult, error) {
	rows, err := conn.pgxConn.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	return &PgxRDbRowsResult{
		rows,
	}, nil
}
func (conn *PgxRDbConn) QueryRow(sql string, args ...interface{}) adapter.RDbRowResult {
	return &PgxRDbRowResult{
		row: conn.pgxConn.QueryRow(context.Background(), sql, args...),
	}
}

type PgxRDbTx struct {
	tx pgx.Tx
}

func (tx *PgxRDbTx) Exec(sql string, args ...interface{}) (adapter.RDbExecResult, error) {
	commandTag, err := tx.tx.Exec(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	return &PgxRDbExecResult{
		commandTag,
	}, nil
}

func (tx *PgxRDbTx) Query(sql string, args ...interface{}) (adapter.RDbRowsResult, error) {
	rows, err := tx.tx.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	return &PgxRDbRowsResult{
		rows,
	}, nil
}
func (tx *PgxRDbTx) QueryRow(sql string, args ...interface{}) adapter.RDbRowResult {
	return &PgxRDbRowResult{
		row: tx.tx.QueryRow(context.Background(), sql, args...),
	}
}
func (tx *PgxRDbTx) Commit() error {
	return tx.tx.Commit(context.Background())
}
func (tx *PgxRDbTx) Rollback() error {
	return tx.tx.Rollback(context.Background())
}

type PgxRDbExecResult struct {
	commandTag pgconn.CommandTag
}

func (result *PgxRDbExecResult) RowsAffected() int64 {
	return result.commandTag.RowsAffected()
}
func (result *PgxRDbExecResult) IsInsert() bool {
	return result.commandTag.Insert()
}
func (result *PgxRDbExecResult) IsUpdate() bool {
	return result.commandTag.Update()
}
func (result *PgxRDbExecResult) IsDelete() bool {
	return result.commandTag.Delete()
}
func (result *PgxRDbExecResult) IsSelect() bool {
	return result.commandTag.Select()
}
func (result *PgxRDbExecResult) String() string {
	return result.commandTag.String()
}

type PgxRDbRowsResult struct {
	rows pgx.Rows
}

func (result *PgxRDbRowsResult) Close() {
	result.rows.Close()
}
func (result *PgxRDbRowsResult) Err() error {
	return result.rows.Err()
}
func (result *PgxRDbRowsResult) ExecResult() adapter.RDbExecResult {
	return &PgxRDbExecResult{
		commandTag: result.rows.CommandTag(),
	}
}
func (result *PgxRDbRowsResult) Next() bool {
	return result.rows.Next()
}
func (result *PgxRDbRowsResult) Scan(dest ...interface{}) error {
	err := result.rows.Scan(dest...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return adapter.ErrNoRows
		}
		return err
	}
	return nil
}

type PgxRDbRowResult struct {
	row pgx.Row
}

func (result *PgxRDbRowResult) Scan(dest ...interface{}) error {
	err := result.row.Scan(dest...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return adapter.ErrNoRows
		}
		return err
	}
	return nil
}

type PgxLoggerAdapter struct {
	logger usecase.Logger
}

func NewPgxLoggerAdapter(logger usecase.Logger) *PgxLoggerAdapter {
	return &PgxLoggerAdapter{
		logger: logger.WithFields(usecase.LogFields{
			"module": "pgx",
		}),
	}
}

func (logger *PgxLoggerAdapter) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	contextedLogger := logger.logger.WithFields(data)

	switch level {
	case pgx.LogLevelError:
		contextedLogger.Error(msg)
	case pgx.LogLevelWarn:
		contextedLogger.Info(msg)
	case pgx.LogLevelInfo:
		fallthrough
	case pgx.LogLevelDebug:
		fallthrough
	case pgx.LogLevelNone:
		fallthrough
	default:
		contextedLogger.Debug(msg)
	}
}

type PgxRDbTypeConv struct{}

func (conv *PgxRDbTypeConv) Bton(b *big.Int) interface{} {
	var err error

	var nilValue pgtype.Numeric
	_ = nilValue.Set(nil)
	if b == nil {
		return nilValue
	}

	var num pgtype.Numeric
	if err = num.Set(b.String()); err != nil {
		panic(fmt.Sprintf("cannot convert %v to numeric: %v", b, err))
	}

	return num
}
func (conv *PgxRDbTypeConv) Iton(i int) interface{} {
	var num pgtype.Numeric
	if err := num.Set(i); err != nil {
		panic(fmt.Sprintf("cannot convert %v to numeric: %v", i, err))
	}

	return num
}
func (conv *PgxRDbTypeConv) NtobReader() adapter.RDbNtobReader {
	return new(PgxRDbNtobReader)
}

type PgxRDbNtobReader struct {
	num pgtype.Numeric
}

func (reader *PgxRDbNtobReader) ScannableArg() interface{} {
	return &reader.num
}
func (reader *PgxRDbNtobReader) ParseW() (*bignum.WBigInt, error) {
	var b bignum.WBigInt
	i, err := reader.Parse()
	if err != nil {
		return nil, err
	}
	return b.FromBigInt(i), nil
}
func (reader *PgxRDbNtobReader) Parse() (*big.Int, error) {
	var err error

	value, err := reader.num.Value()
	if err != nil {
		return nil, err
	}
	switch str := value.(type) {
	case string:
		// pgtype.Numeric.Value() returns scientific notation e.g. "1000e0".
		i, err := sciToBigIntPtr(str)
		if err != nil {
			return nil, err
		}
		return i, nil
	case nil:
		return nil, nil
	default:
		return nil, errors.New("unknown pgtype value")
	}
}
func sciToBigIntPtr(sci string) (*big.Int, error) {
	var err error
	var ok bool

	parts := strings.Split(sci, "e")
	if len(parts) != 2 {
		return nil, errors.New("non-scientific notation value error")
	}

	intval, ok := new(big.Int).SetString(parts[0], 10)
	if !ok {
		return nil, errors.New("integer part to bigInt error")
	}
	if parts[1] == "0" {
		return intval, nil
	}

	exp, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("exponent part to bigInt error: %v", err)
	}

	num := &big.Int{}
	num.Set(intval)
	if exp > 0 {
		mul := &big.Int{}
		mul.Exp(BIGINT_Int10, big.NewInt(exp), nil)
		num.Mul(num, mul)
		return num, nil
	}

	div := &big.Int{}
	div.Exp(BIGINT_Int10, big.NewInt(int64(-exp)), nil)
	remainder := &big.Int{}
	num.DivMod(num, div, remainder)
	if remainder.Cmp(BIGINT_Int0) != 0 {
		return nil, fmt.Errorf("cannot convert %v to bigInt", sci)
	}
	return num, nil
}

func (conv *PgxRDbTypeConv) Tton(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.UTC().UnixNano()
}

func (conv *PgxRDbTypeConv) NtotReader() adapter.RDbNtotReader {
	return new(PgxRDbNtotReader)
}

type PgxRDbNtotReader struct {
	unixNano *int64
}

func NewPgxRDbNtotReader() *PgxRDbNtotReader {
	var i int64
	return &PgxRDbNtotReader{
		unixNano: &i,
	}
}

func (reader *PgxRDbNtotReader) ScannableArg() interface{} {
	return &reader.unixNano
}
func (reader *PgxRDbNtotReader) Parse() (*time.Time, error) {
	if reader.unixNano == nil {
		return nil, nil
	}
	t := time.Unix(0, *reader.unixNano)
	return &t, nil
}
