package main

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type TablesMetrics struct {
	mutex   sync.Mutex
	metrics map[string]*prometheus.GaugeVec
}

var (
	tableMetrics = map[string]metric{
		"extents":     metric{Name: "extents", Help: "Total extents"},
		"table_pages": metric{Name: "table_pages", Help: "Total pages"},
	}
)

func NewtablesMetrics() *TablesMetrics {

	e := TablesMetrics{metrics: map[string]*prometheus.GaugeVec{}}
	for key, _ := range tableMetrics {
		e.metrics[key] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "informix",
			Name:      key,
			Help:      key},
			[]string{"informixserver", "table"})
	}
	return &e
}

func querytables(p *TablesMetrics, Instancia Instance) error {

	var err error

	var (
		name  string
		value float64
		table string
		pages float64
	)

	rows, err := Instancia.db.Query(`select tabname[1,20], count(*)
	from sysmaster:sysextents e
	group by dbsname,tabname
	having count(*)>180
	order by 1,2
	`)

	if err != nil {
		log.Println("Error in  Query extents: \n", err)
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&name, &value)

		if err != nil {
			log.Println("Error in Scan", err)
			return err
		}
		//	if _, ok := p.metrics[strings.TrimSpace(name)]; ok {

		//	p.metrics[strings.TrimSpace(name)].WithLabelValues(Instancia.Name, name).Set(value)
		p.metrics["extents"].WithLabelValues(Instancia.Name, name).Set(value)

	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows, err = Instancia.db.Query(`select  tabname,size from sysextents
	group by 1,2
	having sum(size) > 14000000
	`)

	if err != nil {
		log.Println("Error in  Query pages: \n", err)
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&table, &pages)

		if err != nil {
			log.Println("Error in Scan", err)
			return err
		}
		//	if _, ok := p.metrics[strings.TrimSpace(name)]; ok {

		//	p.metrics[strings.TrimSpace(name)].WithLabelValues(Instancia.Name, name).Set(value)
		p.metrics["table_pages"].WithLabelValues(Instancia.Name, table).Set(pages)

	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func (p *TablesMetrics) Scrape() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var err error

	for m, _ := range Instances.Servers {
		connect := "DSN=" + Instances.Servers[m].Informixserver
		log.Println("Conectando a DSN", connect)
		for intentos := 0; intentos < 5; intentos++ {

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
			log.Println("Error in Open Database: ", err)
		}
	}

	defer func() {
		for m, _ := range Instances.Servers {
			log.Println("Cerrando DSN", m)
			Instances.Servers[m].db.Close()
		}
	}()
	for m, _ := range Instances.Servers {
		log.Println("Ejecutando Querys:", m)
		querytables(p, Instances.Servers[m])
	}
	return nil
}

func (p *TablesMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.metrics {
		m.Describe(ch)
	}
}

func (p *TablesMetrics) Collect(ch chan<- prometheus.Metric) {

	for _, m := range p.metrics {
		m.Collect(ch)
	}

}
