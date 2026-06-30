package config

const (
	SchemaVersion = 1

	ModeObserve = "observe"
	ModeGuard   = "guard"
	ModeStrict  = "strict"
)

type Project struct {
	SchemaVersion int      `json:"schemaVersion"`
	ProjectName   string   `json:"projectName"`
	ProjectRoot   string   `json:"projectRoot"`
	Mode          string   `json:"mode"`
	PackageRunner string   `json:"packageRunner"`
	Framework     string   `json:"framework"`
	TruthFiles    []string `json:"truthFiles"`
	Checks        []Check  `json:"checks"`
	Agents        []Agent  `json:"agents"`
	Tools         Tools    `json:"tools"`
	Policies      Policies `json:"policies"`
	Skills        []string `json:"skills"`
	Notes         []string `json:"notes,omitempty"`
}

type Check struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

type Agent struct {
	Name    string `json:"name"`
	Adapter string `json:"adapter"`
	Enabled bool   `json:"enabled"`
}

type Tools struct {
	RTKFirst           bool `json:"rtkFirst"`
	SQZFallback        bool `json:"sqzFallback"`
	CodeGraphFirst     bool `json:"codeGraphFirst"`
	PackageRunnerGuard bool `json:"packageRunnerGuard"`
}

type Policies struct {
	BlockSecretWrites        bool `json:"blockSecretWrites"`
	BlockDestructiveGit      bool `json:"blockDestructiveGit"`
	BlockUnsafeDeployActions bool `json:"blockUnsafeDeployActions"`
	RequireHandoffOnStop     bool `json:"requireHandoffOnStop"`
	RequireChecksOnStop      bool `json:"requireChecksOnStop"`
	DesignDetectors          bool `json:"designDetectors"`
}

type RuntimeState struct {
	LastStage                 string   `json:"lastStage,omitempty"`
	TouchedFiles              []string `json:"touchedFiles,omitempty"`
	PendingChecks             []string `json:"pendingChecks,omitempty"`
	CodeGraphSeen             bool     `json:"codeGraphSeen,omitempty"`
	HandoffSaved              bool     `json:"handoffSaved,omitempty"`
	SecurityReviewRequired    bool     `json:"securityReviewRequired,omitempty"`
	SecurityReviewComplete    bool     `json:"securityReviewComplete,omitempty"`
	ProductionApproval        bool     `json:"productionApproval,omitempty"`
	LastCompactionSummaryPath string   `json:"lastCompactionSummaryPath,omitempty"`
}
