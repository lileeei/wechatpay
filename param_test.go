package wechatpay

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	ps := make(Params)

	mInt8 := make(map[string]int8)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("int8_%v", i)
		mInt8[key] = int8(i)
	}

	mInt16 := make(map[string]int16)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("int16_%v", i)
		mInt16[key] = int16(i)
	}

	mInt32 := make(map[string]int32)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("int32_%v", i)
		mInt32[key] = int32(i)
	}

	mInt64 := make(map[string]int64)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("int64_%v", i)
		mInt64[key] = int64(i)
	}

	mUint8 := make(map[string]uint8)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("int8_%v", i)
		mUint8[key] = uint8(i)
	}

	mUint16 := make(map[string]uint16)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("uint16_%v", i)
		mUint16[key] = uint16(i)
	}

	mUint32 := make(map[string]uint32)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("uint32_%v", i)
		mUint32[key] = uint32(i)
	}

	mUint64 := make(map[string]uint64)
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("uint64_%v", i)
		mUint64[key] = uint64(i)
	}

	for k, v := range mInt8 {
		ps.Set(k, v)
	}
	for k, v := range mInt16 {
		ps.Set(k, v)
	}
	for k, v := range mInt32 {
		ps.Set(k, v)
	}
	for k, v := range mInt64 {
		ps.Set(k, v)
	}

	for k, v := range mUint8 {
		ps.Set(k, v)
	}
	for k, v := range mUint16 {
		ps.Set(k, v)
	}
	for k, v := range mUint32 {
		ps.Set(k, v)
	}
	for k, v := range mUint64 {
		ps.Set(k, v)
	}

	fmt.Printf("ps: %#v. \n", ps)

}
