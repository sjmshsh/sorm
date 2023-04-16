package sorm

import (
	"database/sql"
	"sorm/dialect"
	"sorm/log"
	"sorm/session"
)

// Session负责与数据库的交互，那交互前的准备工作
// 交互后的收尾工作(关闭连接)等就交给Engine来负责了
// Engine是Sorm与用户交互的入口

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
	}
	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession 通过 Engine 实例创建会话，进而与数据库进行交互了
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
