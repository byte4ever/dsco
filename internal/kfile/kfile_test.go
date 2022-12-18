package kfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/svalue"
)

type MockFS struct {
	afero.Fs
	t               *testing.T
	errPermissionOn map[string]struct{} //nolint:revive //dgas
}

//nolint:ireturn //dgas
func (m *MockFS) Create(name string) (
	afero.File,
	error,
) { //nolint:ireturn //dgas
	m.t.Helper()
	m.t.Logf("Call Create(%q):", name)

	return m.Fs.Create(name) //nolint:wrapcheck // dgas
}

func (m *MockFS) Mkdir(name string, perm os.FileMode) error {
	m.t.Helper()
	m.t.Logf("Call Mkdir(%q, %q):", name, perm)

	return m.Fs.Mkdir(name, perm) //nolint:wrapcheck // dgas
}

func (m *MockFS) MkdirAll(path string, perm os.FileMode) error {
	m.t.Helper()
	m.t.Logf("Call MkdirAll(%q, %q):", path, perm)

	return m.Fs.MkdirAll(path, perm) //nolint:wrapcheck // dgas
}

func (m *MockFS) Open(name string) (afero.File, error) { //nolint:ireturn //dgas
	m.t.Helper()
	m.t.Logf("Call Open(%q):", name)

	if _, ok := m.errPermissionOn[name]; ok {
		return nil, os.ErrPermission
	}

	return m.Fs.Open(name) //nolint:wrapcheck // dgas
}

//nolint:ireturn //dgas
func (m *MockFS) OpenFile(
	name string,
	flag int,
	perm os.FileMode,
) (
	afero.File, error,
) {
	m.t.Helper()
	m.t.Logf("Call OpenFile(%q, %d, %q):", name, flag, perm)

	return m.Fs.OpenFile(name, flag, perm) //nolint:wrapcheck // dgas
}

func (m *MockFS) Remove(name string) error {
	m.t.Helper()
	m.t.Logf("Call Remove(%q):", name)

	return m.Fs.Remove(name) //nolint:wrapcheck // dgas
}

func (m *MockFS) RemoveAll(path string) error {
	m.t.Helper()
	m.t.Logf("Call RemoveAll(%q):", path)

	return m.Fs.RemoveAll(path) //nolint:wrapcheck // dgas
}

func (m *MockFS) Rename(oldname, newname string) error {
	m.t.Helper()
	m.t.Logf("Call Rename(%q, %q):", oldname, newname)

	return m.Fs.Rename(oldname, newname) //nolint:wrapcheck // dgas
}

func (m *MockFS) Stat(name string) (os.FileInfo, error) {
	m.t.Helper()
	m.t.Logf("Call Stat(%q):", name)

	return m.Fs.Stat(name) //nolint:wrapcheck // dgas
}

func (m *MockFS) Name() string {
	m.t.Helper()
	m.t.Log("Call Name()")

	return m.Fs.Name() //nolint:wrapcheck // dgas
}

func (m *MockFS) Chmod(name string, mode os.FileMode) error {
	m.t.Helper()
	m.t.Logf("Call Chmod(%q, %o):", name, mode)

	return m.Fs.Chmod(name, mode) //nolint:wrapcheck // dgas
}

func (m *MockFS) Chown(name string, uid, gid int) error {
	m.t.Helper()
	m.t.Logf("Call Chown(%q, %d, %d):", name, uid, gid)

	return m.Fs.Chown(name, uid, gid) //nolint:wrapcheck // dgas
}

func (m *MockFS) Chtimes(name string, atime time.Time, mtime time.Time) error {
	m.t.Helper()
	m.t.Logf("Call Chtimes(%q, %s, %s):", name, atime, mtime)

	return m.Fs.Chtimes(name, atime, mtime) //nolint:wrapcheck // dgas
}

func Test_scanDirectory(t *testing.T) {
	t.Run(
		"invalid directory", func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()

			provider, err := newProvider(fs, "/tmp", &options{})

			var ep PathErrors
			require.ErrorAs(t, err, &ep)
			require.Len(t, ep, 1)

			e := ep[0]
			require.Equal(t, e.path, "/tmp")
			require.ErrorIs(t, e.err, os.ErrNotExist)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"empty directory", func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()

			require.NoError(t, fs.MkdirAll("/test", 0755))

			provider, err := newProvider(fs, "/test", &options{})

			require.NoError(t, err)
			require.NotNil(t, provider)
			require.Nil(t, provider.values)
		},
	)

	t.Run(
		"some keys", func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()

			require.NoError(t, fs.MkdirAll("/test", 0755))

			require.NoError(
				t,
				afero.WriteFile(
					fs, "/test/K1", []byte("content1"), 0755,
				),
			)

			require.NoError(
				t,
				afero.WriteFile(
					fs, "/test/K2", []byte("content2"), 0755,
				),
			)

			require.NoError(
				t,
				afero.WriteFile(
					fs, "/test/b/K3", []byte("content3"), 0755,
				),
			)

			provider, err := newProvider(fs, "/test", &options{})

			require.NoError(t, err)

			values := provider.values
			require.Len(t, values, 3)
			require.Contains(t, values, "k1")
			require.Contains(t, values, "k2")
			require.Contains(t, values, "k3")

			{
				v := values["k1"]
				require.Equal(t, "content1", v.Value)
				require.Equal(t, "kfile[/test/K1]", v.Location)
			}

			{
				v := values["k2"]
				require.Equal(t, "content2", v.Value)
				require.Equal(t, "kfile[/test/K2]", v.Location)
			}

			{
				v := values["k3"]
				require.Equal(t, "content3", v.Value)
				require.Equal(t, "kfile[/test/b/K3]", v.Location)
			}
		},
	)

	t.Run(
		"permission on file", func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()

			ff2 := &MockFS{
				Fs: fs,
				t:  t,
				errPermissionOn: map[string]struct{}{
					"/a/b/c/K1": {},
				},
			}

			require.NoError(t, fs.MkdirAll("/a/b/c", 0000))
			require.NoError(
				t,
				afero.WriteFile(
					fs, "/a/b/c/K1", []byte{}, 0755,
				),
			)

			provider, err := newProvider(ff2, "/", &options{})

			var ep PathErrors
			require.ErrorAs(t, err, &ep)
			require.Len(t, ep, 1)

			e := ep[0]
			require.Equal(t, e.path, "/a/b/c/K1")
			require.ErrorIs(t, e.err, os.ErrPermission)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"permission on directory", func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()

			ff2 := &MockFS{
				Fs: fs,
				t:  t,
				errPermissionOn: map[string]struct{}{
					"/a/b/c": {},
				},
			}

			require.NoError(t, fs.MkdirAll("/a/b/c", 0000))
			provider, err := newProvider(ff2, "/", &options{})

			var ep PathErrors
			require.ErrorAs(t, err, &ep)
			require.Len(t, ep, 1)

			e := ep[0]
			require.Equal(t, e.path, "/a/b/c")
			require.ErrorIs(t, e.err, os.ErrPermission)
			require.Nil(t, provider)
		},
	)

	t.Run(
		"permission on file", func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()

			require.NoError(
				t,
				afero.WriteFile(
					fs, "/a/b/c/zobby", []byte{}, 0755,
				),
			)

			provider, err := newProvider(fs, "/", &options{})

			var ep PathErrors
			require.ErrorAs(t, err, &ep)
			require.Len(t, ep, 1)

			e := ep[0]
			require.Equal(t, e.path, "/a/b/c/zobby")
			require.ErrorIs(t, e.err, ErrInvalidFileName)
			require.Nil(t, provider)
		},
	)
}

func Test(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileName string
		match    bool
	}{
		{
			name:     "match1",
			fileName: "A4S3D",
			match:    true,
		},
		{
			name:     "match2",
			fileName: "A1SD_A1SD",
			match:    true,
		},
		{
			name:     "match3",
			fileName: "A1SD-Z1AY_ASD1-Q1WE",
			match:    true,
		},
		{
			name:     "unmatch1",
			fileName: "12ASD",
		},
		{
			name:     "unmatch2",
			fileName: "PQR+",
		},
		{
			name:     "unmatch2",
			fileName: "dPQR",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(
			test.name, func(t *testing.T) {
				t.Parallel()
				require.Equal(
					t,
					test.match,
					reFileName.Match([]byte(test.fileName)),
				)
			},
		)
	}
}

func TestEntriesProvider_GetStringValues(t *testing.T) {
	t.Parallel()

	values := make(svalue.Values)

	values["k1"] = &svalue.Value{
		Location: "l1",
		Value:    "v1",
	}

	ep := &EntriesProvider{
		values: values,
	}

	fv := ep.GetStringValues()
	require.Equal(t, values, fv)
}

func TestNewEntriesProvider(t *testing.T) {
	t.Parallel()

	t.Run(
		"ok", func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			fs := afero.NewBasePathFs(
				afero.NewOsFs(),
				tempDir,
			)

			require.NoError(
				t, fs.MkdirAll(
					"/a/b/c",
					0777,
				),
			)

			require.NoError(
				t, afero.WriteFile(
					fs,
					"/a/b/c/K1",
					[]byte("content1"),
					0777,
				),
			)

			provider, err := NewEntriesProvider(tempDir)

			require.NoError(t, err)
			require.NotNil(t, provider)
			require.Len(t, provider.values, 1)
		},
	)

	t.Run(
		"ko", func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			provider, err := NewEntriesProvider(
				filepath.Join(tempDir, "a/b/c/fail"),
			)

			require.Error(t, err)
			require.Nil(t, provider)
		},
	)
}
