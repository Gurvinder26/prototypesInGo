package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {

	numberOfReqs := 1000
	maxDbConnectionLimit := 1000
	poolConnLimit := 10
	benchMarkPool(numberOfReqs, maxDbConnectionLimit, poolConnLimit)
	benchMarkNonPool(numberOfReqs, maxDbConnectionLimit)
}

type cpool struct {
	mu      *sync.Mutex
	channel chan interface{}
	conns   []string
	maxConn int
}

func NewCPool(connLimit int, db *dbMock) *cpool {
	var mu = sync.Mutex{}

	pool := &cpool{
		mu:      &mu,
		conns:   make([]string, db.maxConnectionLimit),
		maxConn: connLimit,
		channel: make(chan interface{}, db.maxConnectionLimit),
	}

	for i := 0; i < db.maxConnectionLimit; i++ {
		dbConn, err := db.newConnection()
		if err != nil {
			panic(err)
		}
		pool.conns = append(pool.conns, dbConn)
		pool.channel <- nil
	}
	return pool
}

func (pool *cpool) Close() {
	close(pool.channel)
}

func (pool *cpool) Get() string {
	<-pool.channel

	pool.mu.Lock()
	c := pool.conns[0]
	pool.conns = pool.conns[1:]
	pool.mu.Unlock()

	return c
}

func (pool *cpool) Put(c string) {
	pool.mu.Lock()
	pool.conns = append(pool.conns, c)
	pool.mu.Unlock()
	pool.channel <- nil
}

func benchMarkPool(numOfReq int, maxDbConnLimit int, poolConnLimit int) {
	var mu = sync.Mutex{}
	db := dbMock{mu: &mu, maxConnectionLimit: maxDbConnLimit, currCount: 0}
	startTime := time.Now()
	pool := NewCPool(poolConnLimit, &db)

	var wg sync.WaitGroup

	for i := 0; i < numOfReq; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			conn := pool.Get()
			// replicate executing a query time
			time.Sleep(10 * time.Millisecond)
			pool.Put(conn)
		}()
	}
	wg.Wait()
	fmt.Println("BenchMark Connection Pool", time.Since(startTime))
	pool.Close()

	// close all the db connections
	for j := 0; j < db.maxConnectionLimit; j++ {
		db.closeConnection()
	}
}

func benchMarkNonPool(numOfReq int, maxDbConnLimit int) {
	startTime := time.Now()
	var wg sync.WaitGroup
	var mu = sync.Mutex{}
	db := &dbMock{mu: &mu, maxConnectionLimit: maxDbConnLimit, currCount: 0}

	for i := 0; i < numOfReq; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer db.closeConnection()
			db.newConnection()
			// replicate executing a query time
			time.Sleep(1 * time.Second)
		}()
	}
	wg.Wait()
	fmt.Println("BenchMark Non Connection Pool", time.Since(startTime))
}

// Db mock
type dbMock struct {
	maxConnectionLimit int
	currCount          int
	mu                 *sync.Mutex
}

func (db *dbMock) newConnection() (string, error) {
	db.mu.Lock()
	if db.currCount < db.maxConnectionLimit {
		time.Sleep(10 * time.Millisecond)
		db.currCount = db.currCount + 1
		db.mu.Unlock()
		return "connection" + strconv.Itoa(rand.Int()), nil
	} else {
		db.mu.Unlock()
		panic("Db connection limit reached")
	}
}

func (db *dbMock) closeConnection() {
	db.mu.Lock()

	if db.currCount > 0 {
		time.Sleep(10 * time.Millisecond)
		db.currCount = db.currCount - 1
	}
	db.mu.Unlock()
}
