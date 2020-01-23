package gurl

import "fmt"

// Status returns the status of IOCat
func (io *IOCat) Status(id string) Status {
	switch v := io.Fail.(type) {
	case nil:
		return Status{
			ID:       id,
			Status:   "success",
			Duration: io.dur.Milliseconds(),
			Actual:   io.Body,
		}
	case *BadMatch:
		return Status{
			ID:       id,
			Status:   "failure",
			Duration: io.dur.Milliseconds(),
			Actual:   v.Actual,
			Expect:   v.Expect,
		}
	default:
		return Status{
			ID:       id,
			Status:   "failure",
			Duration: io.dur.Milliseconds(),
			Actual:   fmt.Sprint(io.Fail),
		}
	}
}
