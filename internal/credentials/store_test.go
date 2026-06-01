package credentials_test

import (
	"errors"
	"testing"

	"github.com/Scale-Flow/trello-cli/internal/credentials"
)

func TestErrNotConfiguredSentinel(t *testing.T) {
	err := credentials.ErrNotConfigured
	if err == nil {
		t.Fatal("ErrNotConfigured should not be nil")
	}
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Error("errors.Is should match ErrNotConfigured")
	}
}

func TestCredentialsStruct(t *testing.T) {
	creds := credentials.Credentials{
		APIKey:   "test-key",
		Token:    "test-token",
		AuthMode: "manual",
	}
	if creds.APIKey != "test-key" {
		t.Errorf("APIKey = %q, want %q", creds.APIKey, "test-key")
	}
	if creds.Token != "test-token" {
		t.Errorf("Token = %q, want %q", creds.Token, "test-token")
	}
	if creds.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", creds.AuthMode, "manual")
	}
}

func TestMemoryStoreSetAndGet(t *testing.T) {
	store := credentials.NewMemoryStore()
	creds := credentials.Credentials{
		APIKey:   "key1",
		Token:    "tok1",
		AuthMode: "manual",
	}

	if err := store.Set("default", creds); err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.APIKey != "key1" {
		t.Errorf("APIKey = %q, want %q", got.APIKey, "key1")
	}
	if got.Token != "tok1" {
		t.Errorf("Token = %q, want %q", got.Token, "tok1")
	}
	if got.AuthMode != "manual" {
		t.Errorf("AuthMode = %q, want %q", got.AuthMode, "manual")
	}
}

func TestMemoryStoreGetMissing(t *testing.T) {
	store := credentials.NewMemoryStore()
	_, err := store.Get("nonexistent")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() error = %v, want ErrNotConfigured", err)
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	store := credentials.NewMemoryStore()
	creds := credentials.Credentials{APIKey: "key1", Token: "tok1", AuthMode: "manual"}
	store.Set("default", creds)

	if err := store.Delete("default"); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() after Delete() error = %v, want ErrNotConfigured", err)
	}
}

func TestMemoryStoreDeleteMissing(t *testing.T) {
	store := credentials.NewMemoryStore()
	// Deleting a nonexistent profile should not error
	if err := store.Delete("nonexistent"); err != nil {
		t.Errorf("Delete() on missing profile returned error: %v", err)
	}
}

func TestMemoryStoreOverwrite(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{APIKey: "old", Token: "old", AuthMode: "manual"})
	store.Set("default", credentials.Credentials{APIKey: "new", Token: "new", AuthMode: "interactive"})

	got, _ := store.Get("default")
	if got.APIKey != "new" {
		t.Errorf("APIKey = %q, want %q after overwrite", got.APIKey, "new")
	}
	if got.AuthMode != "interactive" {
		t.Errorf("AuthMode = %q, want %q after overwrite", got.AuthMode, "interactive")
	}
}

func TestMemoryStoreMultipleProfiles(t *testing.T) {
	store := credentials.NewMemoryStore()
	store.Set("default", credentials.Credentials{APIKey: "k1", Token: "t1", AuthMode: "manual"})
	store.Set("work", credentials.Credentials{APIKey: "k2", Token: "t2", AuthMode: "interactive"})

	got1, _ := store.Get("default")
	got2, _ := store.Get("work")

	if got1.APIKey != "k1" {
		t.Errorf("default APIKey = %q, want %q", got1.APIKey, "k1")
	}
	if got2.APIKey != "k2" {
		t.Errorf("work APIKey = %q, want %q", got2.APIKey, "k2")
	}
}

func TestKeyringStoreServiceName(t *testing.T) {
	svc := credentials.KeyringServiceName("default")
	if svc != "trello-cli/default" {
		t.Errorf("service name = %q, want %q", svc, "trello-cli/default")
	}

	svc2 := credentials.KeyringServiceName("work")
	if svc2 != "trello-cli/work" {
		t.Errorf("service name = %q, want %q", svc2, "trello-cli/work")
	}
}

func TestEnvStoreGetWithEnvVars(t *testing.T) {
	t.Setenv("TRELLO_API_KEY", "env-key")
	t.Setenv("TRELLO_TOKEN", "env-token")

	store := credentials.NewEnvStore()
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got.APIKey != "env-key" {
		t.Errorf("APIKey = %q, want %q", got.APIKey, "env-key")
	}
	if got.Token != "env-token" {
		t.Errorf("Token = %q, want %q", got.Token, "env-token")
	}
	if got.AuthMode != "env" {
		t.Errorf("AuthMode = %q, want %q", got.AuthMode, "env")
	}
}

func TestEnvStoreGetMissingKey(t *testing.T) {
	t.Setenv("TRELLO_API_KEY", "")
	t.Setenv("TRELLO_TOKEN", "")

	store := credentials.NewEnvStore()
	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() error = %v, want ErrNotConfigured", err)
	}
}

func TestEnvStoreGetPartialCreds(t *testing.T) {
	t.Setenv("TRELLO_API_KEY", "env-key")
	t.Setenv("TRELLO_TOKEN", "")

	store := credentials.NewEnvStore()
	_, err := store.Get("default")
	if !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() with only API key error = %v, want ErrNotConfigured", err)
	}
}

func TestEnvStoreSetReturnsError(t *testing.T) {
	store := credentials.NewEnvStore()
	err := store.Set("default", credentials.Credentials{})
	if err == nil {
		t.Error("EnvStore.Set() should return an error (read-only)")
	}
}

func TestEnvStoreDeleteReturnsError(t *testing.T) {
	store := credentials.NewEnvStore()
	err := store.Delete("default")
	if err == nil {
		t.Error("EnvStore.Delete() should return an error (read-only)")
	}
}

func TestFallbackStoreUsesFirst(t *testing.T) {
	primary := credentials.NewMemoryStore()
	secondary := credentials.NewMemoryStore()
	want := credentials.Credentials{APIKey: "primary-key", Token: "primary-token", AuthMode: "manual"}
	if err := primary.Set("default", want); err != nil {
		t.Fatalf("primary.Set() returned error: %v", err)
	}

	store := credentials.NewFallbackStore(primary, secondary)
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got != want {
		t.Errorf("Get() = %+v, want %+v", got, want)
	}
}

func TestFallbackStoreUsesSecondWhenFirstMissing(t *testing.T) {
	primary := credentials.NewMemoryStore()
	secondary := credentials.NewMemoryStore()
	want := credentials.Credentials{APIKey: "secondary-key", Token: "secondary-token", AuthMode: "manual"}
	if err := secondary.Set("default", want); err != nil {
		t.Fatalf("secondary.Set() returned error: %v", err)
	}

	store := credentials.NewFallbackStore(primary, secondary)
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got != want {
		t.Errorf("Get() = %+v, want %+v", got, want)
	}
}

// errStore is a Store whose Get always fails with a non-ErrNotConfigured
// error, simulating an unavailable keyring backend (e.g. no Secret Service
// running on a headless Linux box).
type errStore struct {
	err error
}

func (e errStore) Get(string) (credentials.Credentials, error) {
	return credentials.Credentials{}, e.err
}
func (e errStore) Set(string, credentials.Credentials) error { return e.err }
func (e errStore) Delete(string) error                       { return e.err }

func TestFallbackStoreFallsBackWhenPrimaryUnavailable(t *testing.T) {
	// Primary fails with a backend error that is NOT ErrNotConfigured, the way
	// go-keyring reports a missing Secret Service.
	primary := errStore{err: errors.New("org.freedesktop.secrets was not provided by any .service files")}
	secondary := credentials.NewMemoryStore()
	want := credentials.Credentials{APIKey: "env-key", Token: "env-token", AuthMode: "env"}
	if err := secondary.Set("default", want); err != nil {
		t.Fatalf("secondary.Set() returned error: %v", err)
	}

	store := credentials.NewFallbackStore(primary, secondary)
	got, err := store.Get("default")
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if got != want {
		t.Errorf("Get() = %+v, want %+v", got, want)
	}
}

func TestFallbackStoreReturnsPrimaryErrWhenSecondaryEmpty(t *testing.T) {
	// Primary fails with a backend error and the secondary has nothing — the
	// informative primary error should surface, not ErrNotConfigured.
	primaryErr := errors.New("org.freedesktop.secrets was not provided by any .service files")
	primary := errStore{err: primaryErr}
	secondary := credentials.NewMemoryStore()

	store := credentials.NewFallbackStore(primary, secondary)
	if _, err := store.Get("default"); !errors.Is(err, primaryErr) {
		t.Errorf("Get() error = %v, want %v", err, primaryErr)
	}
}

func TestFallbackStoreReturnsErrWhenBothMissing(t *testing.T) {
	primary := credentials.NewMemoryStore()
	secondary := credentials.NewMemoryStore()

	store := credentials.NewFallbackStore(primary, secondary)
	if _, err := store.Get("default"); !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("Get() error = %v, want ErrNotConfigured", err)
	}
}

func TestFallbackStoreSetWritesToPrimary(t *testing.T) {
	primary := credentials.NewMemoryStore()
	secondary := credentials.NewMemoryStore()
	creds := credentials.Credentials{APIKey: "set-key", Token: "set-token", AuthMode: "manual"}

	store := credentials.NewFallbackStore(primary, secondary)
	if err := store.Set("default", creds); err != nil {
		t.Fatalf("Set() returned error: %v", err)
	}

	got, err := primary.Get("default")
	if err != nil {
		t.Fatalf("primary.Get() returned error: %v", err)
	}
	if got != creds {
		t.Errorf("primary.Get() = %+v, want %+v", got, creds)
	}
}

func TestFallbackStoreDeleteFromPrimary(t *testing.T) {
	primary := credentials.NewMemoryStore()
	secondary := credentials.NewMemoryStore()
	creds := credentials.Credentials{APIKey: "del-key", Token: "del-token", AuthMode: "manual"}
	if err := primary.Set("default", creds); err != nil {
		t.Fatalf("primary.Set() returned error: %v", err)
	}

	store := credentials.NewFallbackStore(primary, secondary)
	if err := store.Delete("default"); err != nil {
		t.Fatalf("Delete() returned error: %v", err)
	}

	if _, err := primary.Get("default"); !errors.Is(err, credentials.ErrNotConfigured) {
		t.Errorf("primary.Get() error = %v, want ErrNotConfigured", err)
	}
}
