//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type ForecastCalculation struct {
	ID          int32 `sql:"primary_key"`
	DomainID    int64
	CreatedAt   time.Time
	CreatedBy   int64
	UpdatedAt   time.Time
	UpdatedBy   int64
	Name        string
	Description *string
	Procedure   string
	Args        *string
}
