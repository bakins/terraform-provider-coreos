package google

import (
	"bytes"
	"fmt"

	"google.golang.org/api/compute/v1"

	"github.com/hashicorp/terraform/helper/resource"
)

// OperationWaitType is an enum specifying what type of operation
// we're waiting on.
type OperationWaitType byte

const (
	OperationWaitInvalid OperationWaitType = iota
	OperationWaitGlobal
	OperationWaitRegion
	OperationWaitZone
)

type OperationWaiter struct {
	Service *compute.Service
	Op      *compute.Operation
	Project string
	Region  string
	Zone    string
	Type    OperationWaitType
}

func (w *OperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var op *compute.Operation
		var err error

		switch w.Type {
		case OperationWaitGlobal:
			op, err = w.Service.GlobalOperations.Get(
				w.Project, w.Op.Name).Do()
		case OperationWaitRegion:
			op, err = w.Service.RegionOperations.Get(
				w.Project, w.Region, w.Op.Name).Do()
		case OperationWaitZone:
			op, err = w.Service.ZoneOperations.Get(
				w.Project, w.Zone, w.Op.Name).Do()
		default:
			return nil, "bad-type", fmt.Errorf(
				"Invalid wait type: %#v", w.Type)
		}

		if err != nil {
			return nil, "", err
		}

		return op, op.Status, nil
	}
}

func (w *OperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"PENDING", "RUNNING"},
		Target:  "DONE",
		Refresh: w.RefreshFunc(),
	}
}

// OperationError wraps compute.OperationError and implements the
// error interface so it can be returned.
type OperationError compute.OperationError

func (e OperationError) Error() string {
	var buf bytes.Buffer

	for _, err := range e.Errors {
		buf.WriteString(err.Message + "\n")
	}

	return buf.String()
}
