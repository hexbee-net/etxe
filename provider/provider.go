package provider

type Provider interface {
	Schema() Schema
}

type Schema struct {
	ID             string
	Name           string
	Description    string
	ResourcesTypes map[string]SchemaResource
	DataSources    map[string]SchemaDataSource
}

type SchemaResource struct {
}

type SchemaDataSource struct {
}
