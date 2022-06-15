/**
 * @Author: Sun
 * @Description:
 * @File:  section
 * @Version: 1.0.0
 * @Date: 2022/5/29 09:58
 */

package pkg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

//Database:
//DBType: postgresql
//Username: dtcloud # 填写数据库账号
//Password: 123456  # 填写数据库密码
//Host: 122.51.164.176
//Port: 5432
//DBName: DTCloud #填写数据库名
type Config struct {
	Host               string `yaml:"Host"`
	Port               int    `yaml:"Port"`
	Prefix             string `yaml:"Prefix"`
	User               string `yaml:"User"`
	Pass               string `yaml:"Pass"`
	DataBase           string `yaml:"DataBase"`
	SetMaxIdleConns    int    `yaml:"SetMaxIdleConns"`
	SetMaxOpenConns    int    `yaml:"SetMaxOpenConns"`
	SetConnMaxLifetime int    `yaml:"SetConnMaxLifetime"`
	LogSpeed           int    `yaml:"Log_speed"`
}

//func NewGetConf() *DataBaseSettings {
//	var dbSettings *DataBaseSettings
//	//获取项目的执行路径
//	path, err := os.Getwd()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(path)
//	vip := viper.New()
//	vip.AddConfigPath(path + "/config") //设置读取的文件路径
//	vip.SetConfigName("config")         //设置读取的文件名
//	vip.SetConfigType("yaml")           //设置文件的类型
//	//尝试进行配置读取
//	if err := vip.ReadInConfig(); err != nil {
//		panic(err)
//	}
//
//	err = vip.Unmarshal(&dbSettings)
//	if err != nil {
//		panic(err)
//	}
//
//	return dbSettings
//}

// NewConf 获取配置信息
func NewConf() *Config {
	yamlFile, err := ioutil.ReadFile("/Users/Sun/DTCloud-code/DTCloudWebSocket/apps/config/config.yaml")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	c := new(Config)
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	fmt.Println(c)
	return c
}
