package dsco

import (
	"errors"
	"time"
)

type T1Sub3 struct {
	SubKey3 *T1Sub1
	KEY1    *float64
}

type T1Sub2 struct {
	SubKey2 *T1Sub3
	KEY1    *float64
}

type T1Sub1 struct {
	SubKey1 *T1Sub2
	KEY1    *float64
}

type T1Root struct {
	SubKeyRoot *T1Sub1
	KEY2       *float64
	KEY3       *int
	KEY4       *string
}

type T2Sub3 struct {
	KEY1 *float64
}

type T2Sub2 struct {
	SubKey2 *T2Sub3
	KEY1    *float64
}

type T2Sub1 struct {
	SubKey1 *T2Sub2
	SubKey2 *T2Sub2
	KEY1    *float64
}

type T2Root struct {
	SubKeyRoot *T2Sub1
	KEY2       *float64
	KEY3       *int
	KEY4       *string
}

type T3Sub3 struct {
	KEY1      *float64
	CycleRoot *T3Root
}

type T3Sub2 struct {
	SubKey2 *T3Sub3
	KEY1    *float64
}

type T3Sub1 struct {
	SubKey1 *T3Sub2
	KEY1    *float64
}

type T3Root struct {
	SubKeyRoot *T3Sub1
	KEY2       *float64
	KEY3       *int
	KEY4       *string
}

type T4Root struct {
	KEY2 *float64
	KEY3 *int
	KEY4 int
}

type T5Root struct {
	KEY2 *float64
	KEY3 *int
	KEY4 map[string]string
}

type T6Root struct {
	KEY2 *float64
	KEY3 *int
	KEY4 map[string]string `yaml:"renamed"`
}

type OkSub struct {
	Key1 *string
	Key2 *uint8
}

type OkRoot struct {
	Sub  *OkSub
	Key1 *int
	Key2 *float64
	Key3 *time.Time
}

var err1 = errors.New("mocked error 1")
var err2 = errors.New("mocked error 2")
var err3 = errors.New("mocked error 3")
var err4 = errors.New("mocked error 4")
var err5 = errors.New("mocked error 5")
var err6 = errors.New("mocked error 6")
