package module

type MetaApiServerBean struct {
	Apiversion    string `yaml:"apiversion"`
	Name          string `yaml:"name"`
	Port          string `yaml:"port"`
	JwtKey        string `yaml:"jwtkey"`
	BboltDBPath   string `yaml:"bboltdbpath"`
	LogFPath      string `yaml:"logfpath"`
	HomeDir       string `yaml:"homedir"`
	LeaderIp      string `yaml:"leaderip"`
	LeaderPort    string `yaml:"leaderport"`
	LeaderMaxPort int64  `yaml:"leadermaxport"`
	LeaderMinPort int64  `yaml:"leaderminport"`
	ApiServerIp   string `yaml:"apiserverip"`
	ApiServerPort string `yaml:"apiserverport"`
	AccessToken   string `yaml:"accesstoken"`
}
