package config

type App struct {
	Title         string `mapstructure:"title" json:"title" yaml:"title"`
	ListenAddress string `mapstructure:"listen_address" json:"listen_address" yaml:"listen_address"`
	Environment   string `mapstructure:"environment" json:"environment" yaml:"environment"`
	SECRET_KEY    string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
}

type Crontab struct {
	SyncDBMetas string `mapstructure:"sync_db_metas" json:"sync_db_metas" yaml:"sync_db_metas"`
}

type Log struct {
	Level   string `mapstructure:"level" json:"level" yaml:"level"`
	RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
}

type Database struct {
	Driver          string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	Database        string `mapstructure:"database" json:"database" yaml:"database"`
	UserName        string `mapstructure:"username" json:"username" yaml:"username"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	Charset         string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	ConnMaxLifetime int    `mapstructure:"conn_max_life_time" json:"conn_max_life_time" yaml:"conn_max_life_time"`
}

type Redis struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type RemoteDB struct {
	UserName string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type Das struct {
	MaxExecutionTime  uint64   `mapstructure:"max_execution_time" json:"max_execution_time" yaml:"max_execution_time"`
	DefaultReturnRows uint64   `mapstructure:"default_return_rows" json:"default_return_rows" yaml:"default_return_rows"`
	MaxReturnRows     uint64   `mapstructure:"max_return_rows" json:"max_return_rows" yaml:"max_return_rows"`
	AllowedUserAgents []string `mapstructure:"allowed_useragents" json:"allowed_useragents" yaml:"allowed_useragents"`
}

type Ghost struct {
	Path string   `mapstructure:"path" json:"path" yaml:"path"`
	Args []string `mapstructure:"args" json:"args" yaml:"args"`
}

type Notify struct {
	NoticeURL string `mapstructure:"notice_url" json:"notice_url" yaml:"notice_url"`
	Wechat    struct {
		Enable  bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
		Webhook string `mapstructure:"webhook" json:"webhook" yaml:"webhook"`
	}
	Mail struct {
		Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
		Username string `mapstructure:"username" json:"username" yaml:"username"`
		Password string `mapstructure:"password" json:"password" yaml:"password"`
		Host     string `mapstructure:"host" json:"host" yaml:"host"`
		Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	}
	DingTalk struct {
		Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
		Webhook  string `mapstructure:"webhook" json:"webhook" yaml:"webhook"`
		Keywords string `mapstructure:"keywords" json:"keywords" yaml:"keywords"`
	}
}

type Configuration struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Crontab  Crontab  `mapstructure:"crontab" json:"crontab" yaml:"crontab"`
	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	RemoteDB RemoteDB `mapstructure:"remotedb" json:"remotedb" yaml:"remotedb"`
	Das      Das      `mapstructure:"das" json:"das" yaml:"das"`
	Ghost    Ghost    `mapstructure:"ghost" json:"ghost" yaml:"ghost"`
	Notify   Notify   `mapstructure:"notify" json:"notify" yaml:"notify"`
}
