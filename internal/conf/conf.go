package conf

import (
	"log"
	"sync"

	"github.com/gurkankaymak/hocon"
)

var once sync.Once
var instance *conf

type conf struct {
	hocon         *hocon.Config
	Port          string
	KafkaHost     string
	KafkaPort     string
	VendorAPIHost string
	VendorAPIPort string
}

func GetConf() *conf {
	once.Do(func() {
		var c *conf = &conf{}
		c.setup()
		instance = c
	})
	return instance
}

func (c *conf) setup() {
	c.parseHOCONConfigFile()
	c.setPort()
	c.setKafkaHost()
	c.setKafkaPort()
	c.setVendorAPIHost()
	c.setVendorAPIPort()
}

func (c *conf) parseHOCONConfigFile() {
	hocon, err := hocon.ParseResource("application.conf")
	if err != nil {
		log.Panic("error while parsing configuration file: ", err)
	}

	log.Printf("configurations: %+v", *hocon)

	c.hocon = hocon
}

func (c *conf) setPort() {
	port := c.hocon.GetString("host.port")
	if len(port) == 0 {
		log.Panic("port environment variable not found")
	}

	c.Port = port
}

func (c *conf) setKafkaHost() {
	kafkaHost := c.hocon.GetString("kafka.host")
	if len(kafkaHost) == 0 {
		log.Panic("kafka host environment variable not found")
	}

	c.KafkaHost = kafkaHost
}

func (c *conf) setKafkaPort() {
	kafkaPort := c.hocon.GetString("kafka.port")
	if len(kafkaPort) == 0 {
		log.Panic("kafka port environment variable not found")
	}

	c.KafkaPort = kafkaPort
}

func (c *conf) setVendorAPIHost() {
	vendorAPIHost := c.hocon.GetString("vendorAPI.host")
	if len(vendorAPIHost) == 0 {
		log.Panic("vendor api host environment variable not found")
	}

	c.VendorAPIHost = vendorAPIHost
}

func (c *conf) setVendorAPIPort() {
	vendorAPIPort := c.hocon.GetString("vendorAPI.port")
	if len(vendorAPIPort) == 0 {
		log.Panic("vendor api port environment variable not found")
	}

	c.VendorAPIPort = vendorAPIPort
}
