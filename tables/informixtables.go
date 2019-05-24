package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/alexbrainman/odbc"
	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	configfile = flag.String("configfile", "config.yaml", "Configuration File")
	puerto     = flag.String("port", "8088", "Listen Port")
	Instances  *Configuration
)

type Instance struct {
	Name           string `yaml:"name"`
	Informixserver string `yaml:"informixserver"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	db             *sql.DB
}

type Configuration struct {
	Servers []Instance `yaml:"servers"`
	Custom  []struct {
		Query    string `yaml:"query"`
		Response string `yaml:"response"`
	} `yaml:"custom"`
}
type metric struct {
	Name string
	Help string
}

type Coleccion interface {
	Scrape() error
	Collect(chan<- prometheus.Metric)
	Describe(chan<- *prometheus.Desc)
}

type Exporter struct {
	m        sync.Mutex
	metricas []Coleccion
}

func NewExporter() *Exporter {

	e := &Exporter{

		metricas: []Coleccion{
			NewtablesMetrics(),
		},
	}

	return e

}

func (e *Exporter) scrape() {

	for _, m := range e.metricas {
		err := m.Scrape()
		if err != nil {
			log.Println("Error in scrape data")
		}

	}

}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.metricas {
		m.Describe(ch)
	}

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.m.Lock()
	defer e.m.Unlock()
	e.scrape()
	for _, m := range e.metricas {
		m.Collect(ch)
	}

}

func loadConfig(filename *string) (*Configuration, error) {

	bytes, err := ioutil.ReadFile(*filename)
	if err != nil {
		return &Configuration{}, err
	}

	var c Configuration
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {

		return &Configuration{}, err
	}

	return &c, nil
}

func main() {

	flag.Parse()
	var err error
	log.Println("Cargando fichero de configuracion:")
	Instances, err = loadConfig(configfile)
	if err != nil {
		log.Fatal("Error en  fichero Yaml:", err)

	}

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	log.Println("Arrancando servidor V0.2 en puerto:", *puerto)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+*puerto, nil))
	os.Exit(0)

}
