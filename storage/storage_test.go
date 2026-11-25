package storage

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func createTestStore(t *testing.T) *Store {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		db.Close()
	})

	store, err := New(db)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}

	return store
}

// Test 1: проверка создания хранилища.
func TestNewStore(t *testing.T) {
	store := createTestStore(t)

	if store == nil {
		t.Fatal("expected non-nil store")
	}
}

// Test 2: проверка сохранения и чтения конфигурации.
func TestSaveAndGetConfig(t *testing.T) {
	store := createTestStore(t)

	input := Config{
		Name:     "Test STM32",
		MCU:      "STM32F103C8",
		Settings: `{"gpio":"PA0","mode":"input"}`,
	}

	id, err := store.Save(input)
	if err != nil {
		t.Fatalf("save configuration: %v", err)
	}

	if id <= 0 {
		t.Fatalf("expected positive ID, got %d", id)
	}

	got, err := store.Get(id)
	if err != nil {
		t.Fatalf("get configuration: %v", err)
	}

	if got.Name != input.Name {
		t.Errorf("name: expected %q, got %q", input.Name, got.Name)
	}

	if got.MCU != input.MCU {
		t.Errorf("MCU: expected %q, got %q", input.MCU, got.MCU)
	}

	if got.Settings != input.Settings {
		t.Errorf("settings: expected %q, got %q",
			input.Settings, got.Settings)
	}
}

// Test 3: проверка отклонения некорректного JSON.
func TestSaveInvalidJSON(t *testing.T) {
	store := createTestStore(t)

	input := Config{
		Name:     "Invalid config",
		MCU:      "STM32F103C8",
		Settings: `{"gpio":}`,
	}

	_, err := store.Save(input)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

// Test 4: проверка обязательности имени конфигурации.
func TestSaveWithoutName(t *testing.T) {
	store := createTestStore(t)

	input := Config{
		MCU:      "STM32F103C8",
		Settings: `{}`,
	}

	_, err := store.Save(input)
	if err == nil {
		t.Fatal("expected error when configuration name is empty")
	}
}
