package strukt

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco"
)

type mapKeyI map[string]interface{}

func Test(t *testing.T) {
	t.Run(
		"detect unsupported types", func(t *testing.T) {
			type LeafType struct {
				v1      *float64
				v2      *int
				v3      *string
				invalid int
			}

			_, err := Provide(
				&LeafType{
					v1:      dsco.V(123.423),
					v2:      dsco.V(123),
					v3:      dsco.V("Haha"),
					invalid: 1,
				},
			)

			require.ErrorIs(t, err, ErrUnsupportedType)
			require.ErrorContains(t, err, "invalid")
			require.ErrorContains(t, err, "int")

		},
	)

	t.Run(
		"detect embedded struct", func(t *testing.T) {

			type Embedded struct {
				KEY1 *float64
			}

			type LeafType struct {
				*Embedded
				KEY2 *float64
				KEY3 *int
				KEY4 *string
			}

			val1 := dsco.V(1.124)
			val2 := dsco.V(123.423)
			val3 := dsco.V(123)
			val4 := dsco.V("Haha")

			v := &LeafType{
				Embedded: &Embedded{
					KEY1: val1,
				},
				KEY2: val2,
				KEY3: val3,
				KEY4: val4,
			}

			b, err := Provide(
				v,
			)

			require.NoError(t, err)
			b.checkValues(
				t,
				mapKeyI{
					"embedded-key1": val1,
					"key2":          val2,
					"key3":          val3,
					"key4":          val4,
				},
			)
		},
	)
}

func (ks *Binder) checkValues(
	t *testing.T,
	expectedKI mapKeyI,
) {
	t.Helper()

	ki := make(mapKeyI, len(ks.values))

	for k, e := range ks.values {
		require.False(t, e.Value.IsNil())
		ki[k] = e.Value.Interface()
	}

	require.Equal(t, expectedKI, ki)
}

/*

type GitTagOptions struct {
	Pattern *string `yaml:"pattern,omitempty"`
	Fmt     *string `yaml:"fmt,omitempty"`
}

type GitOptions struct {
	ScanCron *string        `yaml:"scan_cron,omitempty"`
	URL      *string        // `yaml:"url,omitempty"`
	Tag      *GitTagOptions `yaml:"tag,omitempty"`
	File     *string        `yaml:"file,omitempty"`
}

type S3Options struct {
	Bucket   *string  `yaml:"bucket,omitempty"`
	Prefix   *string  `yaml:"prefix,omitempty"`
	KeyFmt   *string  `yaml:"key_fmt,omitempty"`
	MyFloat  *float64 `yaml:"my_float,omitempty"`
	MyUint64 *uint64  `yaml:"my_uint64,omitempty"`
}

type SampleConf struct {
	WorkDir      *string        `yaml:"work_dir,omitempty"`
	Git          *GitOptions    `yaml:"git,omitempty"`
	AWSRegion    *string        `yaml:"aws_region,omitempty"`
	S3           *S3Options     `yaml:"s3,omitempty"`
	SomeDuration *time.Duration `yaml:"some_duration,omitempty"`
	SomeHash     *hash.Hash     `yaml:"some_hash,omitempty"`
}

func TestStruct(t *testing.T) {
	c := SampleConf{
		WorkDir: dsco.V("toto1"),
		Git: &GitOptions{
			ScanCron: dsco.V("toto2"),
			Tag: &GitTagOptions{
				Fmt: dsco.V("toto3"),
			},
			URL: dsco.V("beautiful.com"),
		},
		S3: &S3Options{
			MyFloat:  dsco.V(123.123123),
			MyUint64: dsco.V(uint64(12312315)),
		},
		SomeDuration: dsco.V(time.Second * 123),
	}

	s, err := Provide(&c)
	require.Nil(t, err)
	require.NotNil(t, s)

	for s2, val := range s.values {
		fmt.Println(s2, val)
	}

	{
		var k *uint64
		dstType := reflect.TypeOf(k)
		dstValue := reflect.ValueOf(k)

		origin, _, _, err := s.Bind("s3-my_uint64", true, dstType, &dstValue)
		require.Equal(t, ID, origin)
		require.NoError(t, err)

		fmt.Println(*dstValue.Interface().(*uint64))
	}

	{
		var k uint64
		dstType := reflect.TypeOf(k)
		dstValue := reflect.ValueOf(k)

		origin, _, _, err := s.Bind("s3-my_uint64", true, dstType, &dstValue)
		require.Equal(t, ID, origin)
		require.ErrorIs(t, err, ErrTypeMismatch)
	}
}
*/
