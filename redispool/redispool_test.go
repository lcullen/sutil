// Copyright 2014 The sutil Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package redispool

import (
	"context"
	"testing"
	"github.com/fzzy/radix/redis"
	"github.com/shawnfeng/sutil/slog/slog"
)

func TestLuaLoad(t *testing.T) {
	pool := NewRedisPool(10)

	err := pool.LoadLuaFile("Test", "./test.luad")
	slog.Infoln(context.TODO(), err)
	if err == nil || err.Error() != "open ./test.luad: no such file or directory" {
		t.Errorf("error here")
	}

	err = pool.LoadLuaFile("Test", "./test.lua")
	if err != nil {
		t.Errorf("error here")
	}

	addr := "localhost:9600"
    args := []interface{}{
		2,
		"key1",
		"key2",
		"argv1",
		"argv2",
	}

	rp := pool.EvalSingle(addr, "Nothave", args)

	slog.Infoln(context.TODO(), rp)

	if "get lua sha1 add:localhost:9600 key:Nothave err:lua not find" != rp.String() {
		t.Errorf("error here")
	}

	rp = pool.EvalSingle(addr, "Test", args)

	slog.Infoln(context.TODO(), rp)
	if rp.Type == redis.ErrorReply {
		t.Errorf("error here")
	}

	if rp.String() != "key1key2argv1argv222" {
		t.Errorf("error here")
	}


}
