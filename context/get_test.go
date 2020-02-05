// Code generated .* DO NOT EDIT.	
package context

import "testing"

func TestContext_GetString(t *testing.T) {
	ctx := New()
	expected := string("example text")
	ctx.Set("test", expected)
	received := ctx.GetString("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetString_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := string("example text")
	received := ctx.GetString("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetString_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetString should panic")
        }
    }()
	_ = ctx.GetString("test", "example text", "example text")
}

func TestContext_GetString_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetString should panic")
        }
    }()
	_ = ctx.GetString("test")
}

func TestContext_GetInt(t *testing.T) {
	ctx := New()
	expected := int(123)
	ctx.Set("test", expected)
	received := ctx.GetInt("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := int(123)
	received := ctx.GetInt("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt should panic")
        }
    }()
	_ = ctx.GetInt("test", 123, 123)
}

func TestContext_GetInt_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt should panic")
        }
    }()
	_ = ctx.GetInt("test")
}

func TestContext_GetInt8(t *testing.T) {
	ctx := New()
	expected := int8(123)
	ctx.Set("test", expected)
	received := ctx.GetInt8("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt8_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := int8(123)
	received := ctx.GetInt8("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt8_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt8 should panic")
        }
    }()
	_ = ctx.GetInt8("test", 123, 123)
}

func TestContext_GetInt8_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt8 should panic")
        }
    }()
	_ = ctx.GetInt8("test")
}

func TestContext_GetInt16(t *testing.T) {
	ctx := New()
	expected := int16(123)
	ctx.Set("test", expected)
	received := ctx.GetInt16("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt16_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := int16(123)
	received := ctx.GetInt16("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt16_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt16 should panic")
        }
    }()
	_ = ctx.GetInt16("test", 123, 123)
}

func TestContext_GetInt16_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt16 should panic")
        }
    }()
	_ = ctx.GetInt16("test")
}

func TestContext_GetInt32(t *testing.T) {
	ctx := New()
	expected := int32(123)
	ctx.Set("test", expected)
	received := ctx.GetInt32("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt32_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := int32(123)
	received := ctx.GetInt32("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt32_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt32 should panic")
        }
    }()
	_ = ctx.GetInt32("test", 123, 123)
}

func TestContext_GetInt32_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt32 should panic")
        }
    }()
	_ = ctx.GetInt32("test")
}

func TestContext_GetInt64(t *testing.T) {
	ctx := New()
	expected := int64(123)
	ctx.Set("test", expected)
	received := ctx.GetInt64("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt64_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := int64(123)
	received := ctx.GetInt64("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt64_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt64 should panic")
        }
    }()
	_ = ctx.GetInt64("test", 123, 123)
}

func TestContext_GetInt64_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetInt64 should panic")
        }
    }()
	_ = ctx.GetInt64("test")
}

func TestContext_GetFloat32(t *testing.T) {
	ctx := New()
	expected := float32(123.5)
	ctx.Set("test", expected)
	received := ctx.GetFloat32("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetFloat32_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := float32(123.5)
	received := ctx.GetFloat32("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetFloat32_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetFloat32 should panic")
        }
    }()
	_ = ctx.GetFloat32("test", 123.5, 123.5)
}

func TestContext_GetFloat32_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetFloat32 should panic")
        }
    }()
	_ = ctx.GetFloat32("test")
}

func TestContext_GetFloat64(t *testing.T) {
	ctx := New()
	expected := float64(123.5)
	ctx.Set("test", expected)
	received := ctx.GetFloat64("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetFloat64_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := float64(123.5)
	received := ctx.GetFloat64("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetFloat64_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetFloat64 should panic")
        }
    }()
	_ = ctx.GetFloat64("test", 123.5, 123.5)
}

func TestContext_GetFloat64_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetFloat64 should panic")
        }
    }()
	_ = ctx.GetFloat64("test")
}

func TestContext_GetBool(t *testing.T) {
	ctx := New()
	expected := bool(false)
	ctx.Set("test", expected)
	received := ctx.GetBool("test")
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetBool_WithDefaultValue(t *testing.T) {
	ctx := New()
	defaultValue := bool(false)
	received := ctx.GetBool("test", defaultValue)
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetBool_PanicOnMoreThanOneDefaultValue(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetBool should panic")
        }
    }()
	_ = ctx.GetBool("test", false, false)
}

func TestContext_GetBool_PanicOnNotFound(t *testing.T) {
	ctx := New()
	defer func() {
        if r := recover(); r == nil {
            t.Error("the GetBool should panic")
        }
    }()
	_ = ctx.GetBool("test")
}
	
