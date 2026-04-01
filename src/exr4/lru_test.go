package lru

import "testing"

func TestSetGet_Basic(t *testing.T) {
	c := NewCache[int, string](2)

	c.Set(1, "one")
	c.Set(2, "two")

	if got, ok := c.Get(1); !ok || got != "one" {
		t.Fatalf("expected key=1 value=one, got (%q, %v)", got, ok)
	}
	if got, ok := c.Get(2); !ok || got != "two" {
		t.Fatalf("expected key=2 value=two, got (%q, %v)", got, ok)
	}
}

func TestGet_NotFoundReturnsZeroValue(t *testing.T) {
	c := NewCache[int, int](2)

	got, ok := c.Get(42)
	if ok {
		t.Fatalf("expected key=42 to be absent")
	}
	if got != 0 {
		t.Fatalf("expected zero-value 0, got %d", got)
	}
}

func TestSet_ExistingKeyUpdatesValueAndKeepsSingleEntry(t *testing.T) {
	c := NewCache[string, int](2)

	c.Set("a", 1)
	c.Set("a", 100)

	if got, ok := c.Get("a"); !ok || got != 100 {
		t.Fatalf("expected updated value 100, got (%d, %v)", got, ok)
	}

	if c.evict.Len() != 1 {
		t.Fatalf("expected list len=1 for same key updates, got %d", c.evict.Len())
	}
}

func TestLRU_EvictsLeastRecentlyUsed_OnCapacityOverflow(t *testing.T) {
	c := NewCache[int, string](2)

	c.Set(1, "one")
	c.Set(2, "two")

	// Делаем 1 недавно использованным => LRU должен стать 2
	if _, ok := c.Get(1); !ok {
		t.Fatalf("expected key 1 to exist")
	}

	// Переполнение, должен вытесниться ключ 2
	c.Set(3, "three")

	if _, ok := c.Get(2); ok {
		t.Fatalf("expected key 2 to be evicted as LRU")
	}
	if got, ok := c.Get(1); !ok || got != "one" {
		t.Fatalf("expected key 1 to remain, got (%q, %v)", got, ok)
	}
	if got, ok := c.Get(3); !ok || got != "three" {
		t.Fatalf("expected key 3 to exist, got (%q, %v)", got, ok)
	}
}

func TestClear_RemovesAllElements(t *testing.T) {
	c := NewCache[int, string](3)

	c.Set(1, "one")
	c.Set(2, "two")
	c.Set(3, "three")

	c.Clear()

	if c.evict.Len() != 0 {
		t.Fatalf("expected empty list after Clear, got len=%d", c.evict.Len())
	}
	if len(c.items) != 0 {
		t.Fatalf("expected empty map after Clear, got len=%d", len(c.items))
	}

	if _, ok := c.Get(1); ok {
		t.Fatalf("expected key 1 to be absent after Clear")
	}
	if _, ok := c.Get(2); ok {
		t.Fatalf("expected key 2 to be absent after Clear")
	}
	if _, ok := c.Get(3); ok {
		t.Fatalf("expected key 3 to be absent after Clear")
	}
}
