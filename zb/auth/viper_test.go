package auth

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"testing"
)

func TestViper(t *testing.T) {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("$HOME")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			dir, _ := os.UserHomeDir()
			err := viper.WriteConfigAs(filepath.Join(dir, configName + "." + configType))
			if err != nil {
				t.Fatal(fmt.Errorf("Fatal error config file: %s \n", err))
			}
		} else {
			t.Fatal("Config file was found but another error was produced")
		}
	}

	m := viper.GetStringMapStringSlice("users")
	t.Logf("%+v", m)

	m["a"] = []string{"root", "password"}
	viper.Set("users", m)
	err = viper.WriteConfig()
	if err != nil {
		t.Fatal(err)
	}


	m = viper.GetStringMapStringSlice("users")
	t.Logf("%+v", m)
}
