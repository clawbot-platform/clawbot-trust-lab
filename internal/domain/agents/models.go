package agents

type AgentRole string
type ActorType string

const (
	AgentRoleReviewer  AgentRole = "reviewer"
	AgentRoleOperator  AgentRole = "operator"
	AgentRoleAdversary AgentRole = "adversary"
)

const (
	ActorTypeHuman  ActorType = "human"
	ActorTypeAgent  ActorType = "agent"
	ActorTypeSystem ActorType = "system"
)

type RuntimeRef struct {
	Runtime string `json:"runtime"`
	Version string `json:"version"`
	Gateway string `json:"gateway"`
}

func KnownRoles() []AgentRole {
	return []AgentRole{
		AgentRoleReviewer,
		AgentRoleOperator,
		AgentRoleAdversary,
	}
}
