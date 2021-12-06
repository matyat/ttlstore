package ttlstore_test

import (
	"testing"
	"time"

	"github.com/matyat/ttlstore"
)

func TestGet(t *testing.T) {
	s := ttlstore.New(5 * time.Second)
	if s.Get("key") != nil {
		t.Error("non-nil key on empty store")
	}
}

func TestSetGet(t *testing.T) {
	s := ttlstore.New(5 * time.Second)
	s.Set("key-1", 5)
	s.Set("key-2", "hello")

	v1, ok := s.Get("key-1").(int)
	if !ok {
		t.Error("key-1 missing")
	} else if v1 != 5 {
		t.Errorf("unexpected value for key-1 '%d", v1)
	}

	v2, ok := s.Get("key-2").(string)
	if !ok {
		t.Error("key-2 missing")
	} else if v2 != "hello" {
		t.Errorf("unexpected value for key-2 '%s", v2)
	}
}

func TestTTL(t *testing.T) {
	s := ttlstore.New(time.Second)
	s.Set("key-1", 5)

	v, ok := s.Get("key-1").(int)
	if !ok {
		t.Error("key-1 missing")
	} else if v != 5 {
		t.Errorf("unexpected value for key-1 '%d", v)
	}

	time.Sleep(2 * time.Second)

	if s.Get("key-1") != nil {
		t.Errorf("key not expired")
	}
}

func TestCustomTTL(t *testing.T) {
	s := ttlstore.New(time.Second)
	s.Set("key-1", 5)
	s.SetWithTTL("key-2", 6, 3*time.Second)

	time.Sleep(2 * time.Second)

	if s.Get("key-1") != nil {
		t.Errorf("key not expired")
	}

	if s.Get("key-2") == nil {
		t.Errorf("key with custom ttl expired")
	}
}
