package techan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("SimpleDT:SimpleDT", func(t *testing.T) {
		parseable := "01/20/2009T12:00:00:01/20/2017T12:00:00"
		timePeriod, err := Parse(parseable)
		assert.NoError(t, err)

		assert.EqualValues(t, 2009, timePeriod.Start.Year())
		assert.EqualValues(t, 01, timePeriod.Start.Month())
		assert.EqualValues(t, 20, timePeriod.Start.Day())
		assert.EqualValues(t, 12, timePeriod.Start.Hour())
		assert.EqualValues(t, 0, timePeriod.Start.Minute())
		assert.EqualValues(t, 0, timePeriod.Start.Second())

		assert.EqualValues(t, 2017, timePeriod.End.Year())
		assert.EqualValues(t, 01, timePeriod.End.Month())
		assert.EqualValues(t, 20, timePeriod.End.Day())
		assert.EqualValues(t, 12, timePeriod.End.Hour())
		assert.EqualValues(t, 0, timePeriod.End.Minute())
		assert.EqualValues(t, 0, timePeriod.End.Second())
	})

	// this has the potential to be flaky, on account of the slight difference in time between
	// now and the now created in Parse
	t.Run("SimpleDT:", func(t *testing.T) {
		parseable := "08/15/1991T20:30:00:"
		timePeriod, err := Parse(parseable)
		now := time.Now()
		assert.NoError(t, err)

		assert.EqualValues(t, 1991, timePeriod.Start.Year())
		assert.EqualValues(t, 8, timePeriod.Start.Month())
		assert.EqualValues(t, 15, timePeriod.Start.Day())
		assert.EqualValues(t, 20, timePeriod.Start.Hour())
		assert.EqualValues(t, 30, timePeriod.Start.Minute())
		assert.EqualValues(t, 0, timePeriod.Start.Second())

		assert.EqualValues(t, now.Year(), timePeriod.End.Year())
		assert.EqualValues(t, now.Month(), timePeriod.End.Month())
		assert.EqualValues(t, now.Day(), timePeriod.End.Day())
		assert.EqualValues(t, now.Hour(), timePeriod.End.Hour())
		assert.EqualValues(t, now.Minute(), timePeriod.End.Minute())
		assert.EqualValues(t, now.Second(), timePeriod.End.Second())
	})

	t.Run("SimpleDate:SimpleDate", func(t *testing.T) {
		parseable := "09/01/1773:07/04/1776"
		timePeriod, err := Parse(parseable)
		assert.NoError(t, err)

		assert.EqualValues(t, 1773, timePeriod.Start.Year())
		assert.EqualValues(t, 9, timePeriod.Start.Month())
		assert.EqualValues(t, 1, timePeriod.Start.Day())

		assert.EqualValues(t, 1776, timePeriod.End.Year())
		assert.EqualValues(t, 7, timePeriod.End.Month())
		assert.EqualValues(t, 4, timePeriod.End.Day())
	})

	t.Run("SimpleDate:", func(t *testing.T) {
		parseable := "07/04/1776:"
		timePeriod, err := Parse(parseable)
		now := time.Now()
		assert.NoError(t, err)

		assert.EqualValues(t, 1776, timePeriod.Start.Year())
		assert.EqualValues(t, 7, timePeriod.Start.Month())
		assert.EqualValues(t, 4, timePeriod.Start.Day())

		assert.EqualValues(t, now.Year(), timePeriod.End.Year())
		assert.EqualValues(t, now.Month(), timePeriod.End.Month())
		assert.EqualValues(t, now.Day(), timePeriod.End.Day())
	})

	t.Run("returns error when format invalid", func(t *testing.T) {
		parseable := "djadk"
		_, err := Parse(parseable)

		assert.EqualError(t, err, "could not parse timerange string djadk")
	})

	t.Run("returns error when start time not parseable", func(t *testing.T) {
		parseable := "07/04/dksj:"

		_, err := Parse(parseable)

		assert.EqualError(t, err, "could not parse time string 07/04/dksj")
	})

	t.Run("returns error when end time not parseable", func(t *testing.T) {
		parseable := "07/04/1776:ab/04/1776"

		_, err := Parse(parseable)

		assert.EqualError(t, err, "could not parse time string ab/04/1776")
	})
}

func TestTimePeriod_Length(t *testing.T) {
	now := time.Now()
	tp := TimePeriod{
		Start: now.Add(-time.Minute * 10),
		End:   now,
	}

	assert.EqualValues(t, time.Minute*10, tp.Length())
}

func TestTimePeriod_Since(t *testing.T) {
	now := time.Now()

	t.Run("0", func(t *testing.T) {
		tp := TimePeriod{
			Start: now,
			End:   now.Add(time.Minute),
		}
		previousTimePeriod := TimePeriod{
			Start: now.Add(-time.Minute),
			End:   now,
		}

		since := tp.Since(previousTimePeriod)
		assert.EqualValues(t, 0, since)
	})

	t.Run("Positive", func(t *testing.T) {
		tp := TimePeriod{
			Start: now,
			End:   now.Add(time.Minute),
		}

		previousTimePeriod := TimePeriod{
			Start: now.Add(-time.Minute * 2),
			End:   now.Add(-time.Minute),
		}

		since := tp.Since(previousTimePeriod)

		assert.EqualValues(t, time.Minute, since)
	})
}

func TestTimePeriod_Advance(t *testing.T) {
	now := time.Now()

	t.Run("Advances by correct amount", func(t *testing.T) {
		tp := TimePeriod{
			Start: now,
			End:   now.Add(time.Minute),
		}

		tp = tp.Advance(1)

		assert.EqualValues(t, now.Add(time.Minute).UnixNano(), tp.Start.UnixNano())
		assert.EqualValues(t, now.Add(time.Minute*2).UnixNano(), tp.End.UnixNano())
	})
}
