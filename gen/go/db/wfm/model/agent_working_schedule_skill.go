//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type AgentWorkingScheduleSkill struct {
	ID                     int32 `sql:"primary_key"`
	DomainID               int64
	AgentWorkingScheduleID int64
	SkillID                int64
	Capacity               int16
}
