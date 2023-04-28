package types

type ConfigFile struct {
	ServerBaseUrl string `json:"server_base_url"`
	ServerPort string `json:"server_port"`
	ServerAPIKey string `json:"server_api_key"`
	ServerCookieName string `json:"server_cookie_name"`
	ServerCookieSecret string `json:"server_cookie_secret"`
	ServerCookieSecretMessage string `json:"server_cookie_secret_message"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
	TimeZone string `json:"time_zone"`
	BoltDBPath string `json:"bolt_db_path"`
	BoltDBEncryptionKey string `json:"bolt_db_encryption_key"`
}

type AListResponse struct {
	UUIDS []string `json:"uuids"`
}

type RedisMultiCommand struct {
	Command string `json:"type"`
	Key string `json:"key"`
	Args string `json:"args"`
}

type SearchItem struct {
	UUID string
	Name string
}