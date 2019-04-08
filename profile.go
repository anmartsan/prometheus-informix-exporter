package main

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	profileMetrics = map[string]metric{
		"pf_isamtot":          metric{Name: "pf_isamtot", Help: "Total ISAM"},
		"pf_isopens":          metric{Name: "pf_isopens", Help: "Total ISAM opens"},
		"pf_isreads":          metric{Name: "pf_isreads", Help: "Total ISAM reads"},
		"pf_iswrites":         metric{Name: "pf_iswrites", Help: "Total ISAM writes"},
		"pf_isrewrites":       metric{Name: "pf_isrewrites", Help: "Total ISAM updates"},
		"pf_isdeletes":        metric{Name: "pf_isdeletes", Help: "Total ISAM deletes"},
		"pf_iscommits":        metric{Name: "pf_iscommits", Help: "Total commits"},
		"pf_isrollbacks":      metric{Name: "pf_isrollbacks", Help: "Total rollbacks"},
		"pf_latchwts":         metric{Name: "pf_latchwts", Help: "Total latch waits"},
		"pf_buffwts":          metric{Name: "pf_buffwts", Help: "Total buffer waits"},
		"pf_lockreqs":         metric{Name: "pf_lockreqs", Help: "Total lock request"},
		"pf_lockwts":          metric{Name: "pf_lockwts", Help: "Total locks waits"},
		"pf_ckptwts":          metric{Name: "pf_ckptwts", Help: "Total checkpoint waits"},
		"pf_plgwrites":        metric{Name: "pf_plgwrites", Help: "Total physical log writes"},
		"pf_pagreads":         metric{Name: "pf_pagreads", Help: "Total page reads"},
		"pf_btradata":         metric{Name: "pf_btradata", Help: "Total pf_btradata"},
		"pf_rapgs_used":       metric{Name: "pf_rapgs_used", Help: "Total pf_rapgs_used"},
		"pf_seqscans":         metric{Name: "pf_seqscans", Help: "Total secuencial scans"},
		"pagreads_2K":         metric{Name: "pagreads_2K", Help: "Total paginas leidas 2k"},
		"bufreads_2K":         metric{Name: "bufreads_2K", Help: "Total buffer reads 2k"},
		"pagwrites_2K":        metric{Name: "pagwrites_2K", Help: "Total page writes 2k "},
		"bufwrites_2K":        metric{Name: "bufwrites_2K", Help: "Total buffer writes 2k"},
		"bufwaits_2K":         metric{Name: "bufwaits_2K", Help: "Total buffer waits 2k"},
		"pagreads_16K":        metric{Name: "pagreads_16K", Help: "Total page reads 16k"},
		"bufreads_16K":        metric{Name: "bufreads_16K", Help: "Total buffer reads 16k"},
		"pagwrites_16K":       metric{Name: "pagwrites_16K", Help: "Total page writes 16k"},
		"bufwrites_16K":       metric{Name: "bufwrites_16K", Help: "Total buffer writes 16k"},
		"bufwaits_16K":        metric{Name: "bufwaits_16K", Help: "Total buffer waits 16k"},
		"open_transactions":   metric{Name: "open_transactions", Help: "Open transactions"},
		"total_locks":         metric{Name: "total_locks", Help: "Total locks"},
		"locks_with_waiter":   metric{Name: "locks_with_waiter", Help: "Total locks with waiters"},
		"logs_without_backup": metric{Name: "logs_without_backup", Help: "Logs without backup"},
		"net_connects":        metric{Name: "net_connects", Help: "Number of connects"},
		"ckptotal":            metric{Name: "ckptotal", Help: "Total time chekcpoints"},
		"dskflush_per_sec":    metric{Name: "dskflush_per_sec", Help: "Total disk flush ckp"},
		"n_dirty_buffs":       metric{Name: "n_dirty_buffs", Help: "Total dirty buffers checkpoints"},
	}
)

type ProfileMetrics struct {
	mutex   sync.Mutex
	metrics map[string]*prometheus.GaugeVec
}

func NewprofileMetrics() *ProfileMetrics {

	e := ProfileMetrics{metrics: map[string]*prometheus.GaugeVec{}}
	for key, _ := range profileMetrics {
		e.metrics[key] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "informix",
			Name:      key,
			Help:      key},
			[]string{"informixserver"})
	}
	return &e

}

func queryprofile(p *ProfileMetrics, Instancia Instance) {

	var (
		name             string
		value            float64
		open             float64
		locks            float64
		locksw           float64
		logs             float64
		timeckp          float64
		dskflush_per_sec float64
		n_dirty_buffs    float64
	)

	var err error
	rows, err := Instancia.db.Query("select name,value from sysshmhdr ")

	if err != nil {
		log.Fatal("Error in  Query: \n", err)
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&name, &value)

		if err != nil {
			log.Fatal("Error in Scan", err)
		}
		if _, ok := p.metrics[strings.TrimSpace(name)]; ok {
			p.metrics[strings.TrimSpace(name)].WithLabelValues(Instancia.Name).Set(value)

		}

	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows, err = Instancia.db.Query(`SELECT COUNT(*) as open_transactions, SUM(tx_nlocks) total_locks ,
	 (SELECT  COUNT(*) as locks_with_waiter FROM syslocks WHERE waiter IS NOT NULL AND
	 dbsname != 'sysmaster' AND tabname != 'sysdatabases') as locks_with_waiter  FROM systrans;
	`)

	if err != nil {
		log.Fatal("Error in Query: \n", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&open, &locks, &locksw)
		if err != nil {
			log.Fatal("Error in Scan", err)
		}
		p.metrics["open_transactions"].WithLabelValues(Instancia.Name).Set(open)
		p.metrics["total_locks"].WithLabelValues(Instancia.Name).Set(locks)
		p.metrics["locks_with_waiter"].WithLabelValues(Instancia.Name).Set(locksw)

	}

	rows, err = Instancia.db.Query(`select count(*) as logssinbackup from syslogs where is_backed_up=0 `)

	if err != nil {
		log.Fatal("Error in Query: \n", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&logs)
		if err != nil {
			log.Fatal("Error in Scan", err)
		}

		p.metrics["logs_without_backup"].WithLabelValues(Instancia.Name).Set(logs)

	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	rows, err = Instancia.db.Query(`select first 1  cp_time::decimal(4,2) as ckptotal,n_dirty_buffs,dskflush_per_sec  from syscheckpoint order by intvl desc `)

	if err != nil {
		log.Fatal("Error in Query: \n", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&timeckp, &n_dirty_buffs, &dskflush_per_sec)
		if err != nil {
			log.Fatal("Error in Scan", err)
		}

		p.metrics["ckptotal"].WithLabelValues(Instancia.Name).Set(timeckp)
		p.metrics["n_dirty_buffs"].WithLabelValues(Instancia.Name).Set(n_dirty_buffs)
		p.metrics["dskflush_per_sec"].WithLabelValues(Instancia.Name).Set(dskflush_per_sec)

	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}

func (p *ProfileMetrics) Scrape() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var err error

	for m, _ := range Instances.Servers {
		connect := "DSN=" + Instances.Servers[m].Informixserver
		Instances.Servers[m].db, err = sql.Open("odbc", connect)
		err = Instances.Servers[m].db.Ping()
		if err != nil {
			Instances.Servers = append(Instances.Servers[:m], Instances.Servers[m+1:]...)
			log.Println("Error in Open Database: ", err)
		}
	}

	defer func() {
		for m, _ := range Instances.Servers {
			Instances.Servers[m].db.Close()
		}
	}()
	for m, _ := range Instances.Servers {
		queryprofile(p, Instances.Servers[m])
	}
	return nil
}

func (p *ProfileMetrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range p.metrics {
		m.Describe(ch)
	}
}

func (p *ProfileMetrics) Collect(ch chan<- prometheus.Metric) {

	for _, m := range p.metrics {
		m.Collect(ch)
	}

}
