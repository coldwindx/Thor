package config

var Config = new(Configuration)

type Configuration struct {
    Application ApplicationConfiguration `mapstructure:"application" json:"app" yaml:"app"`
    Log         LogConfiguration         `mapstructure:"log" json:"log" yaml:"log"`
    Database    DatabaseConfiguration    `mapstructure:"database" json:"database" yaml:"database"`
    Redis       RedisConfiguration       `mapstructure:"redis" json:"redis" yaml:"redis"`
    Jwt         JwtConfiguration         `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
    RabbitMq    RabbitMqConfiguration    `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
}

type ApplicationConfiguration struct {
    Name string `mapstructure:"name" json:"name" yaml:"name"`
    Host string `mapstructure:"host" json:"host" yaml:"host"`
    Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type DatabaseConfiguration struct {
    Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
    Host                string `mapstructure:"host" json:"host" yaml:"host"`
    Port                int    `mapstructure:"port" json:"port" yaml:"port"`
    Database            string `mapstructure:"database" json:"database" yaml:"database"`
    UserName            string `mapstructure:"username" json:"username" yaml:"username"`
    Password            string `mapstructure:"password" json:"password" yaml:"password"`
    Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
    MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
    MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
    LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
    EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
    LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

type JwtConfiguration struct {
    Secret                  string `mapstructure:"secret" json:"secret" yaml:"secret"`
    JwtTtl                  int64  `mapstructure:"jwt_ttl" json:"jwt_ttl" yaml:"jwt_ttl"`
    JwtBlacklistGracePeriod int64  `mapstructure:"jwt_blacklist_grace_period" json:"jwt_blacklist_grace_period" yaml:"jwt_blacklist_grace_period"`
    RefreshGracePeriod      int64  `mapstructure:"refresh_grace_period" json:"refresh_grace_period" yaml:"refresh_grace_period"`
}

type LogConfiguration struct {
    Level      string `mapstructure:"level" json:"level" yaml:"level"`
    Dir        string `mapstructure:"dir" json:"dir" yaml:"dir"`
    Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`
    Format     string `mapstructure:"format" json:"format" yaml:"format"`
    ShowLine   bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
    MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
    MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
    MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
    Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

type RabbitMqConfiguration struct {
    Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
    Port     int    `mapstructure:"port" json:"port" yaml:"port"`
    User     string `mapstructure:"user" json:"user" yaml:"user"`
    Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type RedisConfiguration struct {
    Host     string `mapstructure:"host" json:"host" yaml:"host"`
    Port     int    `mapstructure:"port" json:"port" yaml:"port"`
    DB       int    `mapstructure:"db" json:"db" yaml:"db"`
    Password string `mapstructure:"password" json:"password" yaml:"password"`
}
