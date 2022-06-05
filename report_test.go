package dsco

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReport_addEntry(t *testing.T) {
	t.Run(
		"nil case", func(t *testing.T) {
			var report Report

			report.addEntry(
				ReportEntry{
					Key: "k1",
				},
			)
			require.Len(t, report, 1)
			require.Equal(
				t, Report{
					ReportEntry{
						Key: "k1",
					},
				}, report,
			)
		},
	)
	t.Run(
		"empty case", func(t *testing.T) {
			report := Report{}

			report.addEntry(
				ReportEntry{
					Key: "k1",
				},
			)
			require.Len(t, report, 1)
			require.Equal(
				t, Report{
					ReportEntry{
						Key: "k1",
					},
				}, report,
			)
		},
	)
	t.Run(
		"success", func(t *testing.T) {
			report := Report{
				ReportEntry{
					Key: "k1",
				},
			}

			report.addEntry(
				ReportEntry{
					Key: "k2",
				},
			)
			require.Len(t, report, 2)
			require.Equal(
				t, Report{
					ReportEntry{
						Key: "k1",
					},
					ReportEntry{
						Key: "k2",
					},
				}, report,
			)
		},
	)
}

func TestReport_perEntryReport(t *testing.T) {
	t.Run(
		"", func(t *testing.T) {
			report := Report{
				ReportEntry{
					Idx: 5,
				},
				ReportEntry{
					Idx: 3,
				},
				ReportEntry{
					Key: "k3",
					Idx: -1,
				}, ReportEntry{
					ExternalKey: "k1",
					Errors:      []error{err1, err2},
				},
				ReportEntry{
					Idx: 0,
				},
				ReportEntry{
					ExternalKey: "k2",
					Errors:      []error{err3},
				},
				ReportEntry{
					Idx: 1,
				},
			}

			errs := report.perEntryReport()
			require.Len(t, errs, 4)
			require.Equal(t, []error{err1, err2, err3}, errs[1:])
			require.ErrorIs(t, errs[0], ErrUninitialized)
			require.ErrorContains(t, errs[0], "k3")
		},
	)
}

func TestReportEntry_isFound(t *testing.T) {
	type fields struct {
		Value       reflect.Value
		Key         string
		ExternalKey string
		Idx         int
		Errors      []error
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "not found",
			fields: fields{
				Idx: -1,
			},
			want: false,
		},
		{
			name: "found edge case",
			fields: fields{
				Idx: 0,
			},
			want: true,
		},
		{
			name: "found ",
			fields: fields{
				Idx: 123,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				re := &ReportEntry{
					Value:       tt.fields.Value,
					Key:         tt.fields.Key,
					ExternalKey: tt.fields.ExternalKey,
					Idx:         tt.fields.Idx,
					Errors:      tt.fields.Errors,
				}
				if got := re.isFound(); got != tt.want {
					t.Errorf("isFound() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
