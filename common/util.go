package common

import (
	"strconv"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

const (
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	Stamp       = "Jan _2 15:04:05"
	StampMilli  = "Jan _2 15:04:05.000"
	StampMicro  = "Jan _2 15:04:05.000000"
	StampNano   = "Jan _2 15:04:05.000000000"
)

func StringfyNullTimeToRFC3339(OldDate pq.NullTime) string {
	if OldDate.Valid {
		return OldDate.Time.Format(RFC3339)
	}
	return ""
}

func StringfyDateToRFC3339(OldDate time.Time) string {
	if OldDate.IsZero() {
		return ""
	}
	return OldDate.Format(RFC3339)

}

func GetMapInt64Value(m map[string]interface{}, key string, value *int64) {
	if v, ok := m[key].(int64); ok {
		*value = v
	}
	if v, ok := m[key].(float64); ok {
		*value = int64(v)
	}

	if v, ok := m[key].(string); ok {
		s, _ := strconv.ParseInt(v, 10, 64)
		*value = s
	}

}

func GetMapFloat64Value(m map[string]interface{}, key string, value *float64) {
	if v, ok := m[key].(float64); ok {
		*value = v
	}

	if v, ok := m[key].(decimal.Decimal); ok {
		f, _ := v.Float64()
		*value = f
	}
	if v, ok := m[key].(decimal.NullDecimal); ok {
		f, _ := v.Decimal.Float64()
		*value = f
	}
}

func GetMapStringValue(m map[string]interface{}, key string, value *string) {
	if v, ok := m[key].(string); ok {
		*value = v
	}
	if v, ok := m[key].(float64); ok {
		*value = strconv.FormatFloat(v, 'f', -1, 64)
	}

	if v, ok := m[key].(int64); ok {
		*value = strconv.FormatInt(v, 64)
	}
	if v, ok := m[key].(time.Time); ok {
		*value = StringfyDateToRFC3339(v)
	}

}

func GetMapBoolValue(m map[string]interface{}, key string, value *bool) {
	if v, ok := m[key].(bool); ok {
		*value = v
	}

}
