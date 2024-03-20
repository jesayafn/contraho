package main

type HttpRequestProperties struct {
	Path   string
	Params []string
}

type PagingProperties struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}

type ProjectSearchPage struct {
	Paging     PagingProperties `json:"paging"`
	Components []struct {
		Key              string `json:"key"`
		Name             string `json:"name"`
		Qualifier        string `json:"qualifier"`
		Visibility       string `json:"visibility"`
		LastAnalysisDate string `json:"lastAnalysisDate"`
		Revision         string `json:"revision"`
	} `json:"components"`
}

type ComponentSearchPage struct {
	Paging struct {
		PageSize int `json:"pageSize"`
		Total    int `json:"total"`
	} `json:"paging"`
}
type ComponentSearch struct {
	Components []struct {
		Key              string `json:"key"`
		Name             string `json:"name"`
		Qualifier        string `json:"qualifier"`
		Visibility       string `json:"visibility"`
		LastAnalysisDate string `json:"lastAnalysisDate"`
		Revision         string `json:"revision"`
	} `json:"components"`
}

type ProjectSearchList struct {
	Key                string
	Name               string
	Branch             string
	Loc                string
	Owner              string
	Email              string
	LastAnalysisBranch string
	LastAnalysisDate   string
	QualityGateId      string
	QualityGateName    string
	Language           string
	Bugs               string
	SecurityHotspots   string
	LineCoverage       string
	DuplicatedLines    string
	CodeSmells         string
	DebtInMinute       string
}

type AppList struct {
	Key              string
	Name             string
	Owner            string
	Email            string
	Language         string
	Loc              string
	Bugs             string
	SecurityHotspots string
	LineCoverage     string
	DuplicatedLines  string
	CodeSmells       string
	DebtInMinute     string
}

type ProjectBranchesList struct {
	Branches []struct {
		Name   string `json:"name"`
		IsMain bool   `json:"isMain"`
		Type   string `json:"type"`
		Status struct {
			QualityGateStatus string `json:"qualityGateStatus"`
		} `json:"status"`
		AnalysisDate      string `json:"analysisDate"`
		ExcludedFromPurge bool   `json:"excludedFromPurge"`
	} `json:"branches"`
}

type ProjectBranchesLoC struct {
	Branch string
	LoC    int
}

type ProjectMeasures struct {
	Component struct {
		Key       string `json:"key"`
		Name      string `json:"name"`
		Qualifier string `json:"qualifier"`
		Language  string `json:"language"`
		Path      string `json:"path"`
		Measures  []struct {
			Metric string `json:"metric"`
			Value  string `json:"value,omitempty"`
			Period struct {
				Value     string `json:"value"`
				BestValue bool   `json:"bestValue"`
			} `json:"period"`
		} `json:"measures"`
	} `json:"component"`
	Metrics []struct {
		Key                   string `json:"key"`
		Name                  string `json:"name"`
		Description           string `json:"description"`
		Domain                string `json:"domain"`
		Type                  string `json:"type"`
		HigherValuesAreBetter bool   `json:"higherValuesAreBetter"`
		Qualitative           bool   `json:"qualitative"`
		Hidden                bool   `json:"hidden"`
	} `json:"metrics"`
	Period struct {
		Mode      string `json:"mode"`
		Date      string `json:"date"`
		Parameter string `json:"parameter"`
	} `json:"period"`
}

const (
	ContentType = "application/json"
)

type ProjectPermissions struct {
	Paging PagingProperties `json:"paging"`
	Users  []struct {
		Login       string   `json:"login"`
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
		Avatar      string   `json:"avatar"`
	} `json:"users"`
}

type ProjectSearchOfApplicationPage struct {
	Paging   PagingProperties `json:"paging"`
	Projects []struct {
		Key        string `json:"key"`
		Name       string `json:"name"`
		Enabled    bool   `json:"enabled"`
		Selected   bool   `json:"selected"`
		Accessible bool   `json:"accessible"`
	} `json:"projects"`
}

type ProjectAnalyses struct {
	Paging   PagingProperties `json:"paging"`
	Analyses []struct {
		Key                         string `json:"key"`
		Date                        string `json:"date"`
		ProjectVersion              string `json:"projectVersion"`
		BuildString                 string `json:"buildString"`
		Revision                    string `json:"revision,omitempty"`
		ManualNewCodePeriodBaseline bool   `json:"manualNewCodePeriodBaseline"`
		Events                      []struct {
			Key      string `json:"key"`
			Category string `json:"category"`
			Name     string `json:"name"`
		} `json:"events"`
		DetectedCI string `json:"detectedCI,omitempty"`
	} `json:"analyses"`
}

type NavigationGlobal struct {
	Version string `json:"version"`
	Edition string `json:"edition"`
	// GlobalPages []struct {
	// 	Key  string `json:"key"`
	// 	Name string `json:"name"`
	// } `json:"globalPages"`
	// Settings                map[string]string `json:"settings"`
	// LogoUrl                 string            `json:"logoUrl"`
	// LogoWidth               string            `json:"logoWidth"`
	// Qualifiers              []string          `json:"qualifiers"`
	// ProductionDatabase      bool              `json:"productionDatabase"`
	// BranchesEnabled         bool              `json:"branchesEnabled"`
	// ProjectImportFeature    bool              `json:"projectImportFeatureEnabled"`
	// RegulatoryReportFeature bool              `json:"regulatoryReportFeatureEnabled"`
	// CanAdmin                bool              `json:"canAdmin"`
	// Standalone              bool              `json:"standalone"`
}

type QualityGatesGetByProject struct {
	QualityGate struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Default bool   `json:"default"`
	} `json:"qualityGate"`
}

type NavigationComponent struct {
	Key  string `json:"key"`
	ID   string `json:"id"`
	Name string `json:"name"`
	// Description               string `json:"description"`
	// IsFavorite                bool   `json:"isFavorite"`
	// CanBrowseAllChildProjects bool   `json:"canBrowseAllChildProjects"`
	// Version                   string `json:"version"`
	// AnalysisDate              string `json:"analysisDate"`
	// Extensions                []struct {
	// 	Key  string `json:"key"`
	// 	Name string `json:"name"`
	// } `json:"extensions"`
	QualityProfiles []struct {
		Key      string `json:"key"`
		Name     string `json:"name"`
		Language string `json:"language"`
	} `json:"qualityProfiles"`
	// QualityGate struct {
	// 	Key       string `json:"key"`
	// 	Name      string `json:"name"`
	// 	IsDefault bool   `json:"isDefault"`
	// } `json:"qualityGate"`
	// Configuration struct {
	// 	ShowSettings        bool `json:"showSettings"`
	// 	ShowQualityProfiles bool `json:"showQualityProfiles"`
	// 	ShowQualityGates    bool `json:"showQualityGates"`
	// 	ShowLinks           bool `json:"showLinks"`
	// 	ShowPermissions     bool `json:"showPermissions"`
	// 	ShowHistory         bool `json:"showHistory"`
	// 	ShowUpdateKey       bool `json:"showUpdateKey"`
	// 	ShowBackgroundTasks bool `json:"showBackgroundTasks"`
	// 	CanBrowseProject    bool `json:"canBrowseProject"`
	// 	Extensions          []struct {
	// 		Key  string `json:"key"`
	// 		Name string `json:"name"`
	// 	} `json:"extensions"`
	// } `json:"configuration"`
	// Breadcrumbs []struct {
	// 	Key       string `json:"key"`
	// 	Name      string `json:"name"`
	// 	Qualifier string `json:"qualifier"`
	// } `json:"breadcrumbs"`
}

type QualityProfilesShow struct {
	Profile struct {
		Key          string `json:"key"`
		Name         string `json:"name"`
		Language     string `json:"language"`
		LanguageName string `json:"languageName"`
		// IsInherited               bool   `json:"isInherited"`
		// IsBuiltIn                 bool   `json:"isBuiltIn"`
		// IsDefault                 bool   `json:"isDefault"`
		// ActiveRuleCount           int    `json:"activeRuleCount"`
		// ActiveDeprecatedRuleCount int    `json:"activeDeprecatedRuleCount"`
		// ProjectCount              int    `json:"projectCount"`
		// RulesUpdatedAt            string `json:"rulesUpdatedAt"`
		// LastUsed                  string `json:"lastUsed"`
	} `json:"profile"`
}
