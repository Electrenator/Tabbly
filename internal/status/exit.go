package internal_status

const (
	OK = iota
	UNSPECIFIED_PRIMARY_FUNCTION_ERROR
	FILE_CREATION_ERROR
	FILE_OPEN_ERROR
	DB_CONNECT_ERROR
	DB_MIGRATION_ERROR
)
