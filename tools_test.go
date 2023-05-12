package servicenameallocator

import (
	"context"
	"testing"
	"time"
)

func getName(t *testing.T, ctx context.Context) {
	ns, err := NewNameService(ctx,
		"127.0.0.123",
		[]string{
			"127.0.0.1",
		},
		"/name_allocator/test_tasks")
	if err != nil {
		t.Fail()
		return
	}
	name, err := ns.GetName()
	if err != nil {
		logger.Errorf("get name err:%v", err)
		return
	}
	logger.Infof("get name:%v", name)
	go ns.KeepAlive(name)
}

// TestGetName ...
func TestGetName(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getName(t, ctx)
	getName(t, ctx)
	getName(t, ctx)
	getName(t, ctx)

	time.Sleep(time.Second * 15)
	cancel()
	time.Sleep(time.Second * 5)
}

// TestNameService_InitNameNodes ...
func TestNameService_InitNameNodes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ns, err := NewNameService(ctx,
		"127.0.0.123",
		[]string{
			"127.0.0.1",
		},
		"/name_allocator/test_tasks")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	err = ns.InitNameNodes([]string{
		"svc_1",
		"svc_2",
		"svc_3",
		"svc_4",
	})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
