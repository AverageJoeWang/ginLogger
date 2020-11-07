package ginLogger

// LogConfig
type LogConfig struct {
	Level        string `json:"level"`       //default debug
	Filename     string `json:"filename"`    //default "app.log"
	MaxSize      int    `json:"maxsize"`     //the maximum size in megabytes of the log file before it gets rotated. It defaults to 200 megabytes.
	MaxAge       int    `json:"max_age"`     //the maximum number of days to retain old log files,default 7 days
	MaxBackups   int    `json:"max_backups"` //is the maximum number of old log files to retain.
	TimeLocation string `json:"time_zone"`   //time zone, default "Asia/Chongqing
	TimeFormat   string `json:"time_format"` //default 2006-01-02 15:04:05.000000
	LogFormat 	string	`json:"log_format"`//there two mode such as json/console，default json
}

type commonParam struct {
	ResId     string `json:"res_id"`
	RequestId string `json:"request_id"`
}
