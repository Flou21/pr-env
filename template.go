package main

type TemplateData struct {
	TemplateDeploymentName string
	TemplateNamespace      string
	TemplateAppLabel       string
	TemplateContainerName  string
	TemplateContainerImage string
	TemplateContainerTag   string
	TemplatePort           string
	TemplateMemory         string
	TemplateCpu            string
	TemplateDomainName     string
	PullRequestNumber      string
}
