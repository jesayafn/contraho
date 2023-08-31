package main

import "time"

type HttpRequestProperties struct {
	Path   string
	Params []string
}

type ProjectSearchPage struct {
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
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
	Key              string
	Name             string
	HighestBranch    string
	Loc              string
	Owner            string
	Email            string
	LastAnalysisDate time.Time
	// Qualifier        string
	// Visibility       string
	// LastAnalysisDate string
	// Revision         string
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
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	Users []struct {
		Login       string   `json:"login"`
		Name        string   `json:"name"`
		Email       string   `json:"email"`
		Permissions []string `json:"permissions"`
		Avatar      string   `json:"avatar"`
	} `json:"users"`
}
