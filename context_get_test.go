// Code generated .* DO NOT EDIT.	
package gobdd

import "testing"
import "errors"

func TestContext_GetError(t *testing.T) {
	ctx := NewContext()
	expected := errors.New("new err")
	ctx.Set("test", expected)
	received, err := ctx.GetError("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}


func TestContext_GetString(t *testing.T) {
	ctx := NewContext()
	expected := string("example text")
	ctx.Set("test", expected)
	received, err := ctx.GetString("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetString_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := string("example text")
	received, err := ctx.GetString("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetString_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetString("test", "example text", "example text")
	if err == nil  {
		t.Error("the GetString should return an error")
	}
}

func TestContext_GetString_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetString("test")
	if err == nil  {
		t.Error("the GetString should return an error")
	}
}

func TestContext_GetInt(t *testing.T) {
	ctx := NewContext()
	expected := int(123)
	ctx.Set("test", expected)
	received, err := ctx.GetInt("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := int(123)
	received, err := ctx.GetInt("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt("test", 123, 123)
	if err == nil  {
		t.Error("the GetInt should return an error")
	}
}

func TestContext_GetInt_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt("test")
	if err == nil  {
		t.Error("the GetInt should return an error")
	}
}

func TestContext_GetInt8(t *testing.T) {
	ctx := NewContext()
	expected := int8(123)
	ctx.Set("test", expected)
	received, err := ctx.GetInt8("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt8_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := int8(123)
	received, err := ctx.GetInt8("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt8_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt8("test", 123, 123)
	if err == nil  {
		t.Error("the GetInt8 should return an error")
	}
}

func TestContext_GetInt8_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt8("test")
	if err == nil  {
		t.Error("the GetInt8 should return an error")
	}
}

func TestContext_GetInt16(t *testing.T) {
	ctx := NewContext()
	expected := int16(123)
	ctx.Set("test", expected)
	received, err := ctx.GetInt16("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt16_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := int16(123)
	received, err := ctx.GetInt16("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt16_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt16("test", 123, 123)
	if err == nil  {
		t.Error("the GetInt16 should return an error")
	}
}

func TestContext_GetInt16_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt16("test")
	if err == nil  {
		t.Error("the GetInt16 should return an error")
	}
}

func TestContext_GetInt32(t *testing.T) {
	ctx := NewContext()
	expected := int32(123)
	ctx.Set("test", expected)
	received, err := ctx.GetInt32("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt32_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := int32(123)
	received, err := ctx.GetInt32("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt32_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt32("test", 123, 123)
	if err == nil  {
		t.Error("the GetInt32 should return an error")
	}
}

func TestContext_GetInt32_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt32("test")
	if err == nil  {
		t.Error("the GetInt32 should return an error")
	}
}

func TestContext_GetInt64(t *testing.T) {
	ctx := NewContext()
	expected := int64(123)
	ctx.Set("test", expected)
	received, err := ctx.GetInt64("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetInt64_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := int64(123)
	received, err := ctx.GetInt64("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetInt64_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt64("test", 123, 123)
	if err == nil  {
		t.Error("the GetInt64 should return an error")
	}
}

func TestContext_GetInt64_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetInt64("test")
	if err == nil  {
		t.Error("the GetInt64 should return an error")
	}
}

func TestContext_GetFloat32(t *testing.T) {
	ctx := NewContext()
	expected := float32(123.5)
	ctx.Set("test", expected)
	received, err := ctx.GetFloat32("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetFloat32_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := float32(123.5)
	received, err := ctx.GetFloat32("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetFloat32_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetFloat32("test", 123.5, 123.5)
	if err == nil  {
		t.Error("the GetFloat32 should return an error")
	}
}

func TestContext_GetFloat32_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetFloat32("test")
	if err == nil  {
		t.Error("the GetFloat32 should return an error")
	}
}

func TestContext_GetFloat64(t *testing.T) {
	ctx := NewContext()
	expected := float64(123.5)
	ctx.Set("test", expected)
	received, err := ctx.GetFloat64("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetFloat64_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := float64(123.5)
	received, err := ctx.GetFloat64("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetFloat64_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetFloat64("test", 123.5, 123.5)
	if err == nil  {
		t.Error("the GetFloat64 should return an error")
	}
}

func TestContext_GetFloat64_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetFloat64("test")
	if err == nil  {
		t.Error("the GetFloat64 should return an error")
	}
}

func TestContext_GetBool(t *testing.T) {
	ctx := NewContext()
	expected := bool(false)
	ctx.Set("test", expected)
	received, err := ctx.GetBool("test")
	if err != nil {
		t.Error(err)
	}
	if received != expected {
		t.Errorf("expected %+v but received %+v", expected, received)
	}
}

func TestContext_GetBool_WithDefaultValue(t *testing.T) {
	ctx := NewContext()
	defaultValue := bool(false)
	received, err := ctx.GetBool("test", defaultValue)
	if err != nil {
		t.Error(err)
	}
	if received != defaultValue {
		t.Errorf("expected %+v but received %+v", defaultValue, received)
	}
}

func TestContext_GetBool_ShouldReturnErrorWhenMoreThanOneDefaultValue(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetBool("test", false, false)
	if err == nil  {
		t.Error("the GetBool should return an error")
	}
}

func TestContext_GetBool_ErrorOnNotFound(t *testing.T) {
	ctx := NewContext()
	_, err := ctx.GetBool("test")
	if err == nil  {
		t.Error("the GetBool should return an error")
	}
}
	
