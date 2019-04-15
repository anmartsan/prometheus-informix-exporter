package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CustomMetrics struct {
	mutex   sync.Mutex
	metrics *prometheus.GaugeVec
}

func NewcustomMetrics() *CustomMetrics {

	return &CustomMetrics{
		metrics: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "informix",
			Name:      "custom_metrics",
			Help:      "Metricas por Chunks",
		}, []string{"informixserver", "metrica"}),
	}

}

func (c *CustomMetrics) Describe(ch chan<- *prometheus.Desc) {
	c.metrics.Describe(ch)
}

func (c *CustomMetrics) Collect(ch chan<- prometheus.Metric) {
	c.Scrape()
	c.metrics.Collect(ch)
}

func (c *CustomMetrics) Scrape() error {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	var err error

	for m, _ := range Instances.Servers {
		connect := "DSN=" + Instances.Servers[m].Informixserver
		for intentos := 0; intentos < 3; intentos++ {

			Instances.Servers[m].db, err = sql.Open("odbc", connect)
			err = Instances.Servers[m].db.Ping()
			if err != nil {
				time.Sleep(1 * time.Second)

			} else {
				break
			}
		}
		if err != nil {
			Instances.Servers = append(Instances.Servers[:m], Instances.Servers[m+1:]...)
			log.Println("Error en Open Database: ", err)
		}
	}
	defer func() {
		for m, _ := range Instances.Servers {
			Instances.Servers[m].db.Close()
		}
	}()

	for m, _ := range Instances.Servers {
		for k, _ := range Instances.Custom {
			result := getQuery(Instances.Servers[m], Instances.Custom[k].Query)
			c.metrics.WithLabelValues(Instances.Servers[m].Name, Instances.Custom[k].Response).Set(float64(result))
		}
	}
	return nil
}

func getQuery(Instancia Instance, Query string) float64 {

	var value float64
	rows, err := Instancia.db.Query(Query)
	if err != nil {
		log.Fatal("Error en Query: \n", err)
	}

	for rows.Next() {
		err := rows.Scan(&value)
		if err != nil {
			log.Fatal("Error en Scan", err)
		}
	}

	return float64(value)
}
