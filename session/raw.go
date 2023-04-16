package session

import (
	"database/sql"
	"sorm/clause"
	"sorm/dialect"
	"sorm/log"
	"sorm/schema"
	"strings"
)

// Session 封装有两个目的，一个是统一打印日志(包括SQL语句和错误日志)
// 第二是情况sql和sqlVars两个变量，这样Session可以复用，开启一次会话可以执行多个sql
// 每一次查询都对应着一个Session
type Session struct {
	db       *sql.DB         // 使用sql.Open()方法连接数据库成功之后返回的指针
	sql      strings.Builder // 拼接SQL语句
	sqlVars  []interface{}   // sql语句中的参数
	refTable *schema.Schema
	dialect  dialect.Dialect
	clause   clause.Clause
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		dialect: dialect,
		db:      db,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
