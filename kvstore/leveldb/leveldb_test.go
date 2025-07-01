package leveldb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLevelDBPutGet(t *testing.T) {
	// 准备测试数据库目录，避免影响真实数据
	testDBPath := filepath.Join(os.TempDir(), "leveldb_test")
	defer os.RemoveAll(testDBPath) // 测试结束后清理

	db, err := NewLevelDBStore(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open LevelDB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("Close failed: %v", err)
		}
	}()

	key := []byte("foo")
	value := []byte("bar")

	// 测试写入
	if err := db.Put(key, value); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// 测试读取
	val, err := db.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(val) != "bar" {
		t.Errorf("Expected 'bar', got '%s'", val)
	}

	// 测试 Has 方法
	has, err := db.Has(key)
	if err != nil {
		t.Fatalf("Has failed: %v", err)
	}
	if !has {
		t.Errorf("Expected key to exist")
	}

	// 测试 Delete 方法
	if err := db.Delete(key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	has, err = db.Has(key)
	if err != nil {
		t.Fatalf("Has after delete failed: %v", err)
	}
	if has {
		t.Errorf("Expected key to be deleted")
	}
}
