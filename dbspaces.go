package main

import (
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type chunkmetrics struct {
	path      string
	reads     int64
	writes    int64
	readtime  float64
	writetime float64
}

type dbsmetrics struct {
	name   string
	freemb float64
}

type DbspaceMetrics struct {
	mutex   sync.Mutex
	metrics *prometheus.GaugeVec
	space   *prometheus.GaugeVec
}

func NewdbspaceMetrics() *DbspaceMetrics {

	return &DbspaceMetrics{
		metrics: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "informix",
			Name:      "chunk_metrics",
			Help:      "Metricas por Chunks",
		}, []string{"informixserver", "chunk", "metrica"}),
		space: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "informix",
			Name:      "dbspace_metrics",
			Help:      "Metricas por Dbspace",
		}, []string{"informixserver", "dbspace", "metrica"}),
	}

}

func (d *DbspaceMetrics) Describe(ch chan<- *prometheus.Desc) {
	d.metrics.Describe(ch)
}

func (d *DbspaceMetrics) Collect(ch chan<- prometheus.Metric) {
	d.Scrape()
	d.metrics.Collect(ch)
}

func (d *DbspaceMetrics) Scrape() error {

	d.mutex.Lock()
	defer d.mutex.Unlock()
	var err error

	for m, _ := range Instances.Servers {
		connect := "DSN=" + Instances.Servers[m].Informixserver
		log.Println("Conectando a DSN", connect)
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
			log.Println("Cerrando DSN", m)
			Instances.Servers[m].db.Close()
		}
	}()

	c := []*chunkmetrics{}
	for m, _ := range Instances.Servers {
		log.Println("Ejecutando Querys:", m)
		c = getChunks(Instances.Servers[m])
		for i := range c {

			d.metrics.WithLabelValues(Instances.Servers[m].Name, c[i].path, "reads").Set(float64(c[i].reads))
			d.metrics.WithLabelValues(Instances.Servers[m].Name, c[i].path, "writes").Set(float64(c[i].writes))
			d.metrics.WithLabelValues(Instances.Servers[m].Name, c[i].path, "readtime").Set(c[i].readtime)
			d.metrics.WithLabelValues(Instances.Servers[m].Name, c[i].path, "writetime").Set(c[i].writetime)

		}
	}

	f := []*dbsmetrics{}
	for m, _ := range Instances.Servers {
		f = freeDbs(Instances.Servers[m])
		for i := range f {

			d.metrics.WithLabelValues(Instances.Servers[m].Name, f[i].name, "FreeMB").Set(float64(f[i].freemb))

		}
	}

	return nil
}

func getChunks(Instancia Instance) []*chunkmetrics {

	var (
		fname        string
		pagesread    int64
		pageswritten int64
		readtime     float64
		writetime    float64
	)
	var err error

	res := []*chunkmetrics{}
	c := new(chunkmetrics)

	rows, err := Instancia.db.Query("select fname,pagesread,pageswritten,readtime,writetime from syschktab ")

	if err != nil {
		log.Fatal("Error en Query: \n", err)
	}

	for rows.Next() {
		err := rows.Scan(&fname, &pagesread, &pageswritten, &readtime, &writetime)

		if err != nil {
			log.Fatal("Error en Scan", err)
		}
		c.path = strings.TrimSpace(fname)
		c.reads = pagesread
		c.writes = pageswritten
		c.readtime = readtime
		c.writetime = writetime
		res = append(res, c)
		c = new(chunkmetrics)

	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return res

}

func freeDbs(Instancia Instance) []*dbsmetrics {

	var (
		dbspace  string
		mblibres float64
	)
	var err error

	res := []*dbsmetrics{}
	c := new(dbsmetrics)

	rows, err := Instancia.db.Query(`select
	dbs.name[1,20] dbspace, 
	round(SUM(chk.nfree)*2/1024,2) MBlibres
	
	from sysdbspaces dbs, syschunks chk
	where dbs.dbsnum=chk.dbsnum and is_sbchunk <> 1
	group by dbs.name
	UNION
	select
	dbs.name[1,20] dbspace, 
	round(SUM(chk.udfree)*2/1024,2) MBlibres
	
	from sysdbspaces dbs, syschunks chk
	where dbs.dbsnum=chk.dbsnum and is_sbchunk=1
	group by dbs.name
	order by 1;
	`)

	if err != nil {
		log.Println("Error en Query: \n", err)
	}

	for rows.Next() {
		err := rows.Scan(&dbspace, &mblibres)

		if err != nil {
			log.Println("Error en Scan", err)
		}
		c.name = strings.TrimSpace(dbspace)
		c.freemb = mblibres

		res = append(res, c)
		c = new(dbsmetrics)

	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return res

}
