// Code generated .* DO NOT EDIT.
package context

import (
	"strconv"
	"testing"
	"fmt"
)

func TestContext_GetFloat32Param(t *testing.T) {
	ctx := New()
	given := float32(123.5)
	bs := float64ToByte(float64(given))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetFloat32Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_GetFloat64Param(t *testing.T) {
	ctx := New()
	given := float64(123.5)
	bs := float64ToByte(float64(given))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetFloat64Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func float64ToByte(f float64) []byte {
	return []byte(fmt.Sprintf("%f", f))
}

func TestContext_GetIntParam(t *testing.T) {
	ctx := New()
	given := int(123)
	bs := []byte(strconv.Itoa(int(given)))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetIntParam(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_GetIntParam_PanicsOnNoSuchParam(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetIntParam should panic")
        }
    }()
	_ = ctx.GetIntParam(0)
}

func TestContext_GetInt8Param(t *testing.T) {
	ctx := New()
	given := int8(123)
	bs := []byte(strconv.Itoa(int(given)))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetInt8Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_GetInt8Param_PanicsOnNoSuchParam(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt8Param should panic")
        }
    }()
	_ = ctx.GetInt8Param(0)
}

func TestContext_GetInt16Param(t *testing.T) {
	ctx := New()
	given := int16(123)
	bs := []byte(strconv.Itoa(int(given)))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetInt16Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_GetInt16Param_PanicsOnNoSuchParam(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt16Param should panic")
        }
    }()
	_ = ctx.GetInt16Param(0)
}

func TestContext_GetInt32Param(t *testing.T) {
	ctx := New()
	given := int32(123)
	bs := []byte(strconv.Itoa(int(given)))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetInt32Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_GetInt32Param_PanicsOnNoSuchParam(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt32Param should panic")
        }
    }()
	_ = ctx.GetInt32Param(0)
}

func TestContext_GetInt64Param(t *testing.T) {
	ctx := New()
	given := int64(123)
	bs := []byte(strconv.Itoa(int(given)))
	ctx.SetParams([][]byte{
		bs,
	})

	received := ctx.GetInt64Param(0)
	if received != given {
		t.Errorf("expected %+v bug %+v received", given, received)
	}
}

func TestContext_GetInt64Param_PanicsOnNoSuchParam(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt64Param should panic")
        }
    }()
	_ = ctx.GetInt64Param(0)
}

