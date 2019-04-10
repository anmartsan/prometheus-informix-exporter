# Informix Exporter Prometheus

Prometheus exporter para varias metricas de Informix escrito en GO. 




### Pre-requisitos 📋

Es necesario Docker y docker-compose.



### Instalación 🔧

La instalacion se realizara mediante docker y una serie de ficheros de configuración.

En el fichero ./exporter/sqlhosts añadimos las instancias de Informix que queremos Monitorizar de la misma manera 
que si fuese el sqlhosts de Informix

./export/sqlhosts
```
#Server         Protocol         Host           Port

prueba		onsoctcp	192.168.1.50	1527
prueba2		onsoctcp	192.168.1.50	1530

```
En el fichero ./exporter/odbc.ini configuramos el ODBC


```
[ODBC]
UNICODE=UCS-2
[prueba]
Driver=/opt/IDS12/lib/cli/libifcli.so
Server=prueba
Database=sysmaster
TRANSLATIONDLL=/opt/IDS12/lib/esql/igo4a304.so
LogonID=informix
pwd=informix
[prueba2]
Driver=/opt/IDS12/lib/cli/libifcli.so
Server=prueba2
Database=sysmaster
TRANSLATIONDLL=/opt/IDS12/lib/esql/igo4a304.so
LogonID=informix
pwd=informix

```

El fichero ./exporter/config.yaml lo utilizara el exporter para leer los datos de configuracion
Ejemplo:

```
---
servers:
- name: pruebaids
  informixserver: prueba
  user: informix
  password: informix
- name: pruebaids2
  informixserver: prueba2
  user: informix
  password: informix
custom: 
- query: select tabid from systables where tabid=99 
  response: tabid


```

La configuracion de prometheus se encuentra en ./prometheus

Se puede cambiar el puerto donde se quiere que escuche el exporter.

```
- job_name: 'informix'

    # metrics_path defaults to '/metrics'
    # scheme defaults to 'http'.

    static_configs:
    - targets: ['ids_exporter:8080']
      #  - job_name: 'node'

```



## Arranque del sistema ⚙️

```
docker-compose up -d

```



## Autores ✒️



* **Antonio Martinez Sanchez-Escalonilla ** - [anmartsan](https://github.com/anmartsan)
    www.scmsi.es







---
⌨️ con ❤️ por [anmartsan](a.martinez@scmsi.es) 😊
