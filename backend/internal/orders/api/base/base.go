package base

// DBConfig holds the configuration details for connecting to a database and executing a SQL task.
type DBConfig struct {
	Hostname         string // The hostname of the database server.
	Port             uint16 // The port number on which the database server is listening.
	Charset          string // The character set to be used for the database connection.
	UserName         string // The username for authenticating with the database.
	Password         string // The password for authenticating with the database.
	Schema           string // The database schema to be used.
	DBType           string // The type of the database (e.g., MySQL, PostgreSQL).
	SQLType          string // The type of the SQL operation (e.g., SELECT, INSERT).
	SQL              string // The SQL query or command to be executed.
	OrderID          string // An identifier for the order related to the SQL task.
	TaskID           string // An identifier for the specific task to be executed.
	ExportFileFormat string // The format for exporting the data (e.g., CSV, JSON).
}

// ExportFile contains details about an exported file.
type ExportFile struct {
	FileName      string `json:"file_name"`      // The name of the exported file.
	FileSize      int64  `json:"file_size"`      // The size of the exported file in bytes.
	FilePath      string `json:"file_path"`      // The path where the exported file is stored.
	ContentType   string `json:"content_type"`   // The MIME type of the exported file.
	EncryptionKey string `json:"encryption_key"` // The key used to encrypt the file, if any.
	ExportRows    int64  `json:"export_rows"`    // The number of rows exported in the file.
	DownloadUrl   string `json:"download_url"`   // The URL to download the exported file.
}

// ReturnData contains the results and metadata of a SQL execution task.
type ReturnData struct {
	RollbackSQL     string `json:"rollback_sql"`      // The SQL command to rollback the operation.
	AffectedRows    int64  `json:"affected_rows"`     // The number of rows affected by the SQL execution.
	ExecuteCostTime string `json:"execute_cost_time"` // The time taken to execute the SQL command.
	BackupCostTime  string `json:"backup_cost_time"`  // The time taken to backup data, if applicable.
	ExecuteLog      string `json:"execute_log"`       // The log output of the SQL execution.
	ExportFile             // Embedded struct containing export file details.
	Error           string `json:"error"` // Error message, if any occurred during execution.
}
