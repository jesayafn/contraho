package main

const (
	projectIndexApi           = "/api/projects/search?qualifiers=TRK&ps=1&p=1"
	projectScrapeApi          = "/api/projects/search?qualifiers=TRK&ps=500&p="
	aplIndexApi               = "/api/projects/search?qualifiers=APP&ps=1&p=1"
	apltScrapeApi             = "/api/projects/search?qualifiers=APP&ps=500&p="
	projectScrapeAplApi       = "/api/measures/component_tree?ps=500&s=qualifier,name&metricKey=ncloc&strategy=children&p="
	projectIndexAplApi        = "/api/measures/component_tree?ps=1&s=qualifier,name&metricKey=ncloc&strategy=children"
	timeParseFormat           = "2006-01-02T15:04:05-0700"
	projectBranchesApi        = "/api/project_branches/list?project="
	ProjectBranchesLocApi     = "/api/measures/component?metricKeys=ncloc&component="
	ProjectUserPermissionsApi = "/api/permissions/users?projectKey="
	ProjectDateAnalysisApi    = "/api/project_analyses/search?ps=1&p=1&project="
)
