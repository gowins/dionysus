package distributedlock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
)

func TestNew(t *testing.T) {
	lockKey := "testLockKet"
	testExpiration := 11 * time.Second
	testRetryTTL := 2 * time.Second
	db, _ := redismock.NewClientMock()
	rlock := New(db, lockKey, WithExpiration(testExpiration), WithRetryTTL(testRetryTTL), WithWatchDog(false))
	if rlock.watchDogEnable != false {
		t.Errorf("want watchDogEnable false get true")
		return
	}
	if rlock.expiration != testExpiration {
		t.Errorf("want expiration %v, get expiration %v", testExpiration, rlock.expiration)
		return
	}
	if rlock.retryTTL != testRetryTTL {
		t.Errorf("want retryTTL %v, get retryTTL %v", testRetryTTL, rlock.retryTTL)
		return
	}
}

func TestRedisLock_Lock(t *testing.T) {
	testLock1 := "testKey1111"
	testLock2 := "testKey2222"
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	rlock := New(db, testLock1, WithExpiration(0), WithRetryTTL(0), WithWatchDog(true))
	va, err := getLockValue()
	if err != nil {
		t.Errorf("want get lock value want error nil, get %v", err)
		return
	}
	mock.ExpectSetNX(testLock1, va, 0).SetVal(true)
	_, err = rlock.Lock(context.Background())
	if err != nil {
		t.Errorf("1st get lock want error nil, get %v", err)
		return
	}

	mock.ClearExpect()
	mock.ExpectSetNX(testLock1, va, 0).SetVal(false)
	_, err = rlock.Lock(context.Background())
	if err == nil {
		t.Errorf("2nd get lock want error not nil, get nil")
		return
	}

	mock.ClearExpect()
	mock.ExpectSetNX(testLock1, va, 0).SetErr(fmt.Errorf("redis error"))
	rlock2 := New(db, testLock2, WithExpiration(0), WithRetryTTL(0), WithWatchDog(true))
	_, err = rlock2.Lock(context.Background())
	if err == nil {
		t.Errorf("2nd get lock want error not nil, get nil")
		return
	}
}

func TestWithWatchDog(t *testing.T) {
	expiration := 10 * time.Second
	va, err := getLockValue()
	if err != nil {
		t.Errorf("want get lock value want error nil, get %v", err)
		return
	}
	testLockKey := "testKeyWD"
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectSetNX(testLockKey, va, expiration).SetVal(true)
	mock.ExpectEvalSha("a2c2e4c111924caec00216d4881ed37a644435ce", []string{testLockKey}, va, expiration.Seconds()).SetVal(1)
	mock.ExpectEvalSha("a2c2e4c111924caec00216d4881ed37a644435ce", []string{testLockKey}, va, expiration.Seconds()).SetErr(nil)
	rlock := New(db, testLockKey, WithExpiration(expiration), WithWatchDog(true))
	timeStart := time.Now()
	ctx, err := rlock.Lock(context.Background())
	if err != nil {
		t.Errorf("want lock error nil, get %v", err)
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	select {
	case <-ctxTimeout.Done():
		t.Errorf("want not timeout")
		return
	case <-ctx.Done():
		fmt.Printf("time now is %v, time start is %v, sub %v", time.Now().String(), timeStart.String(), time.Now().Sub(timeStart).Seconds())
		if time.Now().Sub(timeStart).Seconds() < 13 {
			t.Errorf("want failed at 2nd")
		}
	}
}

func TestRedisLock_Unlock(t *testing.T) {
	testLock1 := "testKeyUnlock"
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	rlock := New(db, testLock1, WithExpiration(0), WithRetryTTL(0), WithWatchDog(true), WithWatchDog(true))
	va, err := getLockValue()
	if err != nil {
		t.Errorf("want get lock value want error nil, get %v", err)
		return
	}
	mock.ExpectSetNX(testLock1, va, 0).SetVal(true)
	_, err = rlock.Lock(context.Background())
	if err != nil {
		t.Errorf("1st get lock want error nil, get %v", err)
		return
	}

	mock.ClearExpect()
	mock.ExpectEvalSha("cf0e94b2e9ffc7e04395cf88f7583fc309985910", []string{testLock1}, va).SetVal(int64(1))
	mock.ExpectEvalSha("cf0e94b2e9ffc7e04395cf88f7583fc309985910", []string{testLock1}, va).SetErr(nil)
	err = rlock.Unlock(context.Background())
	if err != nil {
		t.Errorf("want error nil, get error %v", err)
		return
	}

	err = rlock.Unlock(context.Background())
	if err == nil {
		t.Errorf("want error, but get error nil")
	}
}

func TestRedisLock_TTL(t *testing.T) {
	testLock1 := "testKeyTTL"
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	rlock := New(db, testLock1, WithExpiration(0), WithRetryTTL(0), WithWatchDog(true))
	va, err := getLockValue()
	if err != nil {
		t.Errorf("want get lock value want error nil, get %v", err)
		return
	}
	_, err = rlock.TTL(context.Background())
	if err == nil {
		t.Errorf("want get ttl error not nil")
		return
	}
	mock.ExpectEvalSha("6484da7d58920897fdc22b7f9afd1a3d47524ea8", []string{testLock1}, va).SetVal(1000)
	mock.ExpectEvalSha("6484da7d58920897fdc22b7f9afd1a3d47524ea8", []string{testLock1}, va).SetErr(nil)
	_, err = rlock.TTL(context.Background())
	if err != nil {
		t.Errorf("want get ttl error nil, get error %v", err)
		return
	}
}

func TestRedisLock_ClearForce(t *testing.T) {
	testLock1 := "testKeyClearForce"
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	rlock := New(db, testLock1, WithExpiration(0), WithRetryTTL(0), WithWatchDog(true))
	mock.ExpectDel(testLock1).SetVal(int64(0))
	mock.ExpectDel(testLock1).SetErr(nil)
	re, err := rlock.ClearForce(context.Background())
	if re != 0 || err != nil {
		t.Errorf("want result 0, get %v, want err nil, get %v", re, err)
	}
}

func TestRedisLock_GetLockIdAndTTL(t *testing.T) {
	testLock1 := "testKeyClearForce"
	testLockId := "testLockId"
	testTTL := 10 * time.Second
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	rlock := New(db, testLock1, WithExpiration(0), WithRetryTTL(0), WithWatchDog(true))
	mock.ExpectGet(testLock1).SetVal(testLockId)
	mock.ExpectDel(testLock1).SetErr(nil)
	mock.ExpectTTL(testLock1).SetVal(testTTL)
	mock.ExpectDel(testLock1).SetErr(nil)
	testid, ttl, err := rlock.GetLockIdAndTTL(context.Background())
	if err != nil {
		t.Errorf("want err nil, get %v", err)
		return
	}
	if testid != testLockId {
		t.Errorf("want id testLockId get %v", testid)
		return
	}
	if ttl != testTTL {
		t.Errorf("want ttl %v get %v", testTTL, ttl)
	}
}
