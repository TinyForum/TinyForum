package vo

type GetConfigYamlVO struct {
	Name    string `json:"name"`
	Format  string `json:"format"`
	Content string `json:"content"`
}
type GetConfigKvVO struct {
	Name   string                 `json:"name"`
	Format string                 `json:"format"`
	Kv     map[string]interface{} `json:"kv"`
}
