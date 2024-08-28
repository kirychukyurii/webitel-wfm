package service

import (
	"context"

	"github.com/webitel/webitel-wfm/infra/webitel/engine"
	"github.com/webitel/webitel-wfm/infra/webitel/logger"
	"github.com/webitel/webitel-wfm/internal/model"
	"github.com/webitel/webitel-wfm/pkg/compare"
	"github.com/webitel/webitel-wfm/pkg/werror"
)

type AgentAbsenceManager interface {
	CreateAgentAbsence(ctx context.Context, user *model.SignedInUser, in *model.AgentAbsence) (*model.AgentAbsence, error)
	ReadAgentAbsence(ctx context.Context, user *model.SignedInUser, search *model.SearchItem) (*model.AgentAbsence, error)
	UpdateAgentAbsence(ctx context.Context, user *model.SignedInUser, in *model.AgentAbsence) (*model.AgentAbsence, error)
	DeleteAgentAbsence(ctx context.Context, user *model.SignedInUser, agentId, id int64) error

	CreateAgentsAbsencesBulk(ctx context.Context, user *model.SignedInUser, agentIds []int64, in []*model.AgentAbsenceBulk) ([]*model.AgentAbsences, error)
	ReadAgentAbsences(ctx context.Context, user *model.SignedInUser, search *model.AgentAbsenceSearch) (*model.AgentAbsences, error)
	SearchAgentsAbsences(ctx context.Context, user *model.SignedInUser, search *model.AgentAbsenceSearch) ([]*model.AgentAbsences, error)
}

type AgentAbsence struct {
	store  AgentAbsenceManager
	audit  *logger.Audit
	engine *engine.Client
}

func NewAgentAbsence(store AgentAbsenceManager, audit *logger.Audit, engine *engine.Client) *AgentAbsence {
	return &AgentAbsence{
		store:  store,
		audit:  audit,
		engine: engine,
	}
}

func (a *AgentAbsence) CreateAgentAbsence(ctx context.Context, user *model.SignedInUser, in *model.AgentAbsence) (*model.AgentAbsence, error) {
	_, err := a.engine.Agent(ctx, in.Agent.Id)
	if err != nil {
		return nil, err
	}

	out, err := a.store.CreateAgentAbsence(ctx, user, in)
	if err != nil {
		return nil, err
	}

	if err := a.audit.Create(ctx, user, map[int64]any{out.Absence.Id: out}); err != nil {
		return nil, err
	}

	return out, nil
}

func (a *AgentAbsence) UpdateAgentAbsence(ctx context.Context, user *model.SignedInUser, in *model.AgentAbsence) (*model.AgentAbsence, error) {
	_, err := a.engine.Agent(ctx, in.Agent.Id)
	if err != nil {
		return nil, err
	}

	out, err := a.store.UpdateAgentAbsence(ctx, user, in)
	if err != nil {
		return nil, err
	}

	if err := a.audit.Update(ctx, user, map[int64]any{out.Absence.Id: out}); err != nil {
		return nil, err
	}

	return out, nil
}

func (a *AgentAbsence) DeleteAgentAbsence(ctx context.Context, user *model.SignedInUser, agentId, id int64) error {
	_, err := a.engine.Agent(ctx, agentId)
	if err != nil {
		return err
	}

	if err := a.store.DeleteAgentAbsence(ctx, user, agentId, id); err != nil {
		return err
	}

	if err := a.audit.Delete(ctx, user, map[int64]any{id: nil}); err != nil {
		return err
	}

	return nil
}

func (a *AgentAbsence) ReadAgentAbsences(ctx context.Context, user *model.SignedInUser, search *model.AgentAbsenceSearch) (*model.AgentAbsences, error) {
	_, err := a.engine.Agent(ctx, search.AgentIds[0])
	if err != nil {
		return nil, err
	}

	items, err := a.store.SearchAgentsAbsences(ctx, user, search)
	if err != nil {
		return nil, err
	}

	if len(items) > 1 {
		return nil, werror.NewDBEntityConflictError("service.agent_absence.read.conflict")
	}

	if len(items) == 0 {
		return nil, werror.NewDBNoRowsErr("service.agent_absence.read")
	}

	return items[0], nil
}

func (a *AgentAbsence) CreateAgentsAbsencesBulk(ctx context.Context, user *model.SignedInUser, agentIds []int64, in []*model.AgentAbsenceBulk) ([]*model.AgentAbsences, error) {
	agents, err := a.engine.Agents(ctx, &model.AgentSearch{Ids: agentIds})
	if err != nil {
		return nil, err
	}

	// Checks if signed user has read access to a desired set of agents.
	if ok := compare.ElementsMatch(agents, agentIds); !ok {
		return nil, werror.NewAuthForbiddenError("service.agent_absence.check_agents", "cc_agent", "read")
	}

	out, err := a.store.CreateAgentsAbsencesBulk(ctx, user, agents, in)
	if err != nil {
		return nil, err
	}

	records := make(map[int64]any)
	for _, item := range out {
		for _, absence := range item.Absence {
			records[absence.Id] = absence
		}
	}

	if err := a.audit.Create(ctx, user, records); err != nil {
		return nil, err
	}

	return out, nil
}

func (a *AgentAbsence) SearchAgentsAbsences(ctx context.Context, user *model.SignedInUser, search *model.AgentAbsenceSearch) ([]*model.AgentAbsences, bool, error) {
	s := &model.AgentSearch{
		SupervisorIds: search.SupervisorIds,
		TeamIds:       search.TeamIds,
		SkillIds:      search.SkillIds,
	}

	agents, err := a.engine.Agents(ctx, s)
	if err != nil {
		return nil, false, err
	}

	search.AgentIds = agents
	items, err := a.store.SearchAgentsAbsences(ctx, user, search)
	if err != nil {
		return nil, false, err
	}

	var next bool
	if len(items) == int(search.SearchItem.Limit()) {
		next = true
		items = items[:search.SearchItem.Limit()-1]
	}

	return items, next, nil
}
