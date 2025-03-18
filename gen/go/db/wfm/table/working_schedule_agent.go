//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var WorkingScheduleAgent = newWorkingScheduleAgentTable("wfm", "working_schedule_agent", "")

type workingScheduleAgentTable struct {
	postgres.Table

	// Columns
	ID                postgres.ColumnInteger
	DomainID          postgres.ColumnInteger
	WorkingScheduleID postgres.ColumnInteger
	AgentID           postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
	DefaultColumns postgres.ColumnList
}

type WorkingScheduleAgentTable struct {
	workingScheduleAgentTable

	EXCLUDED workingScheduleAgentTable
}

// AS creates new WorkingScheduleAgentTable with assigned alias
func (a WorkingScheduleAgentTable) AS(alias string) *WorkingScheduleAgentTable {
	return newWorkingScheduleAgentTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new WorkingScheduleAgentTable with assigned schema name
func (a WorkingScheduleAgentTable) FromSchema(schemaName string) *WorkingScheduleAgentTable {
	return newWorkingScheduleAgentTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new WorkingScheduleAgentTable with assigned table prefix
func (a WorkingScheduleAgentTable) WithPrefix(prefix string) *WorkingScheduleAgentTable {
	return newWorkingScheduleAgentTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new WorkingScheduleAgentTable with assigned table suffix
func (a WorkingScheduleAgentTable) WithSuffix(suffix string) *WorkingScheduleAgentTable {
	return newWorkingScheduleAgentTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newWorkingScheduleAgentTable(schemaName, tableName, alias string) *WorkingScheduleAgentTable {
	return &WorkingScheduleAgentTable{
		workingScheduleAgentTable: newWorkingScheduleAgentTableImpl(schemaName, tableName, alias),
		EXCLUDED:                  newWorkingScheduleAgentTableImpl("", "excluded", ""),
	}
}

func newWorkingScheduleAgentTableImpl(schemaName, tableName, alias string) workingScheduleAgentTable {
	var (
		IDColumn                = postgres.IntegerColumn("id")
		DomainIDColumn          = postgres.IntegerColumn("domain_id")
		WorkingScheduleIDColumn = postgres.IntegerColumn("working_schedule_id")
		AgentIDColumn           = postgres.IntegerColumn("agent_id")
		allColumns              = postgres.ColumnList{IDColumn, DomainIDColumn, WorkingScheduleIDColumn, AgentIDColumn}
		mutableColumns          = postgres.ColumnList{DomainIDColumn, WorkingScheduleIDColumn, AgentIDColumn}
		defaultColumns          = postgres.ColumnList{IDColumn}
	)

	return workingScheduleAgentTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                IDColumn,
		DomainID:          DomainIDColumn,
		WorkingScheduleID: WorkingScheduleIDColumn,
		AgentID:           AgentIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
