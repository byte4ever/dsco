package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/byte4ever/dsco/utils"
)

func Test_scan(t *testing.T) {
	t.Parallel()

	t.Run(
		"success", func(t *testing.T) {
			t.Parallel()

			type E1 struct {
				E1X *int
			}

			e1xType := reflect.TypeOf(utils.R(0))

			type E2 struct {
				E2X *string
			}

			e2xType := reflect.TypeOf(utils.R(""))

			type Sub1 struct {
				E1
				Sub1X *float32
				Sub1Y *float64
			}

			sub1xType := reflect.TypeOf(utils.R(float32(0.0)))
			sub1yType := reflect.TypeOf(utils.R(0.0))
			sub1Type := reflect.TypeOf(&Sub1{})

			type Sub2 struct {
				E2
				Sub2X *time.Duration
				Sub2Y *time.Time
			}

			sub2xType := reflect.TypeOf(utils.R(time.Duration(0)))
			sub2yType := reflect.TypeOf(&time.Time{})
			sub2Type := reflect.TypeOf(&Sub2{})

			type Root struct {
				E1
				E2
				S1 *Sub1
				S2 *Sub2
			}
			var maxUID uint
			vType := reflect.TypeOf(&Root{})
			node, mError := scan(&maxUID, "", vType)
			require.True(t, mError.None())
			require.Equal(t, uint(8), maxUID)

			expectedNode := &StructNode{
				Type: vType,
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node: &ValueNode{
							Type:        e1xType,
							VisiblePath: "E1X",
							UID:         0x0,
						},
						Index: []int{0, 0},
					},
					&IndexedSubNode{
						Node: &ValueNode{
							Type:        e2xType,
							VisiblePath: "E2X",
							UID:         0x1,
						},
						Index: []int{1, 0},
					},
					&IndexedSubNode{
						Node: &StructNode{
							Type: sub1Type,
							Index: IndexedSubNodes{
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        e1xType,
										VisiblePath: "S1.E1X",
										UID:         0x2,
									},
									Index: []int{0, 0},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub1xType,
										VisiblePath: "S1.Sub1X",
										UID:         0x3,
									},
									Index: []int{1},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub1yType,
										VisiblePath: "S1.Sub1Y",
										UID:         0x4,
									},
									Index: []int{2},
								},
							},
						},
						Index: []int{2},
					},
					&IndexedSubNode{
						Node: &StructNode{
							Type: sub2Type,
							Index: IndexedSubNodes{
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        e2xType,
										VisiblePath: "S2.E2X",
										UID:         0x5,
									},
									Index: []int{0, 0},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub2xType,
										VisiblePath: "S2.Sub2X",
										UID:         0x6,
									},
									Index: []int{1},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub2yType,
										VisiblePath: "S2.Sub2Y",
										UID:         0x7,
									},
									Index: []int{2},
								},
							},
						},
						Index: []int{3},
					},
				},
			}

			require.Equal(t, expectedNode, node)
		},
	)

	t.Run(
		"invalid type", func(t *testing.T) {
			t.Parallel()

			type E1 struct {
				E1X       *int
				E1Invalid string
			}

			e1xType := reflect.TypeOf(utils.R(0))
			e1invalidType := reflect.TypeOf("")

			type InvalidStruct struct{}

			type E2 struct { //nolint:govet // this field order is required
				E2X       *string
				E2Invalid InvalidStruct
			}

			e2xType := reflect.TypeOf(utils.R(""))
			e2invalidType := reflect.TypeOf(InvalidStruct{})

			type Sub1 struct { //nolint:govet // this field order is required
				E1
				Sub1X   *float32
				Sub1Y   *float64
				Invalid int
			}

			sub1xType := reflect.TypeOf(utils.R(float32(0.0)))
			sub1yType := reflect.TypeOf(utils.R(0.0))
			sub1invalidType := reflect.TypeOf(0)
			sub1Type := reflect.TypeOf(&Sub1{})

			type Sub2 struct { //nolint:govet // this field order is required
				E2
				Sub2X   *time.Duration
				Sub2Y   *time.Time
				Invalid float64
			}

			sub2xType := reflect.TypeOf(utils.R(time.Duration(0)))
			sub2yType := reflect.TypeOf(&time.Time{})
			sub2invalidType := reflect.TypeOf(23.32)
			sub2Type := reflect.TypeOf(&Sub2{})

			type Root struct { //nolint:govet // this field order is required
				E1
				E2
				S1 *Sub1
				S2 *Sub2
			}
			var maxUID uint
			vType := reflect.TypeOf(&Root{})
			node, mError := scan(&maxUID, "", vType)

			require.Len(t, mError, 6)

			for idx, expectedError := range []UnsupportedTypeError{
				{
					Path: "E1Invalid",
					Type: e1invalidType,
				},
				{
					Path: "E2Invalid",
					Type: e2invalidType,
				},
				{
					Path: "S1.E1Invalid",
					Type: e1invalidType,
				},
				{
					Path: "S1.Invalid",
					Type: sub1invalidType,
				},
				{
					Path: "S2.E2Invalid",
					Type: e2invalidType,
				},
				{
					Path: "S2.Invalid",
					Type: sub2invalidType,
				},
			} {
				var e UnsupportedTypeError
				require.ErrorAs(t, mError[idx], &e)
				require.Equal(t, expectedError, e)
			}

			expectedNode := &StructNode{
				Type: vType,
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node: &ValueNode{
							Type:        e1xType,
							VisiblePath: "E1X",
							UID:         0x0,
						},
						Index: []int{0, 0},
					},
					&IndexedSubNode{
						Node: &ValueNode{
							Type:        e2xType,
							VisiblePath: "E2X",
							UID:         0x1,
						},
						Index: []int{1, 0},
					},
					&IndexedSubNode{
						Node: &StructNode{
							Type: sub1Type,
							Index: IndexedSubNodes{
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        e1xType,
										VisiblePath: "S1.E1X",
										UID:         0x2,
									},
									Index: []int{0, 0},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub1xType,
										VisiblePath: "S1.Sub1X",
										UID:         0x3,
									},
									Index: []int{1},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub1yType,
										VisiblePath: "S1.Sub1Y",
										UID:         0x4,
									},
									Index: []int{2},
								},
							},
						},
						Index: []int{2},
					},
					&IndexedSubNode{
						Node: &StructNode{
							Type: sub2Type,
							Index: IndexedSubNodes{
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        e2xType,
										VisiblePath: "S2.E2X",
										UID:         0x5,
									},
									Index: []int{0, 0},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub2xType,
										VisiblePath: "S2.Sub2X",
										UID:         0x6,
									},
									Index: []int{1},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub2yType,
										VisiblePath: "S2.Sub2Y",
										UID:         0x7,
									},
									Index: []int{2},
								},
							},
						},
						Index: []int{3},
					},
				},
			}
			require.Equal(t, expectedNode, node)
		},
	)

	t.Run(
		"invalid type", func(t *testing.T) {
			t.Parallel()

			type E1 struct {
				E1X *int
			}

			e1xType := reflect.TypeOf(utils.R(0))

			type E2 struct {
				E2X *string
			}

			e2xType := reflect.TypeOf(utils.R(""))

			type Sub1 struct {
				E1
				Sub1X *float32
				Sub1Y *float64
			}

			sub1xType := reflect.TypeOf(utils.R(float32(0.0)))
			sub1yType := reflect.TypeOf(utils.R(0.0))
			sub1Type := reflect.TypeOf(&Sub1{})

			type Sub2 struct {
				E2
				Sub2X *time.Duration
				Sub2Y *time.Time
			}

			sub2xType := reflect.TypeOf(utils.R(time.Duration(0)))
			sub2yType := reflect.TypeOf(&time.Time{})
			sub2Type := reflect.TypeOf(&Sub2{})

			type Root struct {
				E1
				E2
				S1  *Sub1
				S2  *Sub2
				E1X *int
			}

			var maxUID uint
			vType := reflect.TypeOf(&Root{})
			node, mError := scan(&maxUID, "", vType)

			require.Len(t, mError, 1)

			for idx, expectedError := range []FieldNameCollisionError{
				{
					Path1: "E1X",
					Path2: "E1.E1X",
				},
			} {
				var e FieldNameCollisionError
				require.ErrorAs(t, mError[idx], &e)
				require.Equal(t, expectedError, e)
			}

			expectedNode := &StructNode{
				Type: vType,
				Index: IndexedSubNodes{
					&IndexedSubNode{
						Node: &ValueNode{
							Type:        e1xType,
							VisiblePath: "E1X",
							UID:         0x0,
						},
						Index: []int{0, 0},
					},
					&IndexedSubNode{
						Node: &ValueNode{
							Type:        e2xType,
							VisiblePath: "E2X",
							UID:         0x1,
						},
						Index: []int{1, 0},
					},
					&IndexedSubNode{
						Node: &StructNode{
							Type: sub1Type,
							Index: IndexedSubNodes{
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        e1xType,
										VisiblePath: "S1.E1X",
										UID:         0x2,
									},
									Index: []int{0, 0},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub1xType,
										VisiblePath: "S1.Sub1X",
										UID:         0x3,
									},
									Index: []int{1},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub1yType,
										VisiblePath: "S1.Sub1Y",
										UID:         0x4,
									},
									Index: []int{2},
								},
							},
						},
						Index: []int{2},
					},
					&IndexedSubNode{
						Node: &StructNode{
							Type: sub2Type,
							Index: IndexedSubNodes{
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        e2xType,
										VisiblePath: "S2.E2X",
										UID:         0x5,
									},
									Index: []int{0, 0},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub2xType,
										VisiblePath: "S2.Sub2X",
										UID:         0x6,
									},
									Index: []int{1},
								},
								&IndexedSubNode{
									Node: &ValueNode{
										Type:        sub2yType,
										VisiblePath: "S2.Sub2Y",
										UID:         0x7,
									},
									Index: []int{2},
								},
							},
						},
						Index: []int{3},
					},
				},
			}
			require.Equal(t, expectedNode, node)
		},
	)
}
