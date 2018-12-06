package setting

import (
	"bytes"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/chinx/wepub/crypts"
	_ "github.com/chinx/morph/bmp"
	_ "github.com/chinx/morph/gif"
	_ "github.com/chinx/morph/jpeg"
	_ "github.com/chinx/morph/png"
	_ "github.com/chinx/morph/tiff"
	_ "github.com/chinx/morph/webp"
)

type Options struct {
	HttpPort  int    `toml:"httpport"`
	HttpsPort int    `toml:"httpsport"`
	StaticDir string `toml:"staticdir"`
	Mysql     *Mysql `toml:"mysql"`
	Redis     *Redis `toml:"redis"`
}

type Mysql struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type Redis struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Session  int    `toml:"session"`
	Database int    `toml:"database"`
}

func LoadConfigFile(privateFile, configFile string) (*Options, error) {
	cryptData, err := ioutil.ReadFile(privateFile)
	if err != nil {
		return nil, err
	}
	err = crypts.DecryptPrivateKey(cryptData)
	if err != nil {
		return nil, err
	}

	confData, err := LoadDecryptFile(configFile)
	if err != nil {
		return nil, err
	}

	opt := &Options{}
	_, err = toml.DecodeReader(bytes.NewBuffer(confData), opt)
	if err != nil {
		return nil, err
	}
	return opt, nil
}

func LoadDecryptFile(filePath string) ([]byte, error) {
	cryptData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return crypts.AesDecrypt(cryptData)
}
