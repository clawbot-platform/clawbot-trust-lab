package actors

import "clawbot-trust-lab/internal/domain/agents"

type ActorType string
type DelegationMode string

const (
	ActorTypeHuman ActorType = "human"
	ActorTypeAgent ActorType = "agent"
)

const (
	DelegationModeDirectHuman    DelegationMode = "direct_human"
	DelegationModeAgentAssisted  DelegationMode = "agent_assisted"
	DelegationModeFullyDelegated DelegationMode = "fully_delegated"
)

type PrincipalRef struct {
	PrincipalID   string `json:"principal_id"`
	PrincipalType string `json:"principal_type"`
}

type HumanActor struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      ActorType    `json:"type"`
	Principal PrincipalRef `json:"principal"`
	Tags      []string     `json:"tags"`
}

type AgentActor struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Type      ActorType         `json:"type"`
	Role      agents.AgentRole  `json:"role"`
	Runtime   agents.RuntimeRef `json:"runtime"`
	Principal PrincipalRef      `json:"principal"`
	Tags      []string          `json:"tags"`
}
