package reporter

import (
	"errors"
	"testing"

	"github.com/cucumber/gherkin-go"
	"github.com/stretchr/testify/assert"
)

func TestReporter_Report(t *testing.T) {
	step := &gherkin.Step{}
	r := NewReporter()
	err := r.Report(step)
	assert.NoError(t, err)
	assert.Equal(t, 1, r.countSteps())
}

func TestReporter_Report_Invalid_Type(t *testing.T) {
	r := NewReporter()
	err := r.Report("this shouldn't go here")
	assert.Error(t, err)
	assert.Equal(t, 0, r.countSteps())
}

func TestReporter_Undefined(t *testing.T) {
	step := &gherkin.Step{}
	r := NewReporter()
	r.Undefined(step)
	assert.Equal(t, 1, r.countUndefined())
}

func TestReporter_Failed(t *testing.T) {
	step := &gherkin.Step{}
	r := NewReporter()
	r.Failed(step, errors.New("some error"))
	assert.Equal(t, 1, r.countFailed())
}

func TestReporter_Skip_Step(t *testing.T) {
	step := &gherkin.Step{}
	r := NewReporter()
	r.Skip(step)
	assert.Equal(t, 1, r.countSkippedSteps())
}
