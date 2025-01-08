package engine

import (
	"context"

	"github.com/webitel/webitel-go-kit/logging/wlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/webitel/webitel-wfm/infra/health"
	"github.com/webitel/webitel-wfm/infra/registry"
	"github.com/webitel/webitel-wfm/infra/shutdown"
	"github.com/webitel/webitel-wfm/infra/webitel"
	"github.com/webitel/webitel-wfm/pkg/werror"
)

var serviceName = "engine"

type Client struct {
	log  *wlog.Logger
	conn *grpc.ClientConn
}

func New(log *wlog.Logger, discovery registry.Discovery) (*Client, error) {
	conn, err := webitel.New(log, discovery, serviceName)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
	}, nil
}

func (c *Client) Shutdown(p *shutdown.Process) error {
	return c.conn.Close()
}

func (c *Client) HealthCheck(ctx context.Context) []health.CheckResult {
	state := c.conn.GetState()
	if state != connectivity.Idle && state != connectivity.Ready {
		return []health.CheckResult{{Name: serviceName, Err: werror.New("service is not ready", werror.WithValue("state", state.String()))}}
	}

	return []health.CheckResult{{Name: serviceName, Err: nil}}
}

func (c *Client) AgentService() *AgentService {
	return newAgentServiceClient(c)
}

func (c *Client) CalendarService() *CalendarService {
	return newCalendarServiceClient(c)
}

func (c *Client) TeamService() *TeamService {
	return newTeamServiceClient(c)
}
