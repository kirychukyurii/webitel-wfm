//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type WorkingConditionACL struct {
	ID      int32 `sql:"primary_key"`
	Dc      int64
	Grantor *int64
	Object  int32
	Subject int64
	Access  int16
}
