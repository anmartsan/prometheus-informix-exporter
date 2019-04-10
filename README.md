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


_Finaliza con un ejemplo de cómo obtener datos del sistema o como usarlos para una pequeña demo_

## Ejecutando las pruebas ⚙️

_Explica como ejecutar las pruebas automatizadas para este sistema_

### Analice las pruebas end-to-end 🔩

_Explica que verifican estas pruebas y por qué_

```
Da un ejemplo
```

### Y las pruebas de estilo de codificación ⌨️

_Explica que verifican estas pruebas y por qué_

```
Da un ejemplo
```

## Deployment 📦

_Agrega notas adicionales sobre como hacer deploy_

## Construido con 🛠️

_Menciona las herramientas que utilizaste para crear tu proyecto_

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - El framework web usado
* [Maven](https://maven.apache.org/) - Manejador de dependencias
* [ROME](https://rometools.github.io/rome/) - Usado para generar RSS

## Contribuyendo 🖇️

Por favor lee el [CONTRIBUTING.md](https://gist.github.com/villanuevand/xxxxxx) para detalles de nuestro código de conducta, y el proceso para enviarnos pull requests.

## Wiki 📖

Puedes encontrar mucho más de cómo utilizar este proyecto en nuestra [Wiki](https://github.com/tu/proyecto/wiki)

## Versionado 📌

Usamos [SemVer](http://semver.org/) para el versionado. Para todas las versiones disponibles, mira los [tags en este repositorio](https://github.com/tu/proyecto/tags).

## Autores ✒️

_Menciona a todos aquellos que ayudaron a levantar el proyecto desde sus inicios_

* **Andrés Villanueva** - *Trabajo Inicial* - [villanuevand](https://github.com/villanuevand)
* **Fulanito Detal** - *Documentación* - [fulanitodetal](#fulanito-de-tal)

También puedes mirar la lista de todos los [contribuyentes](https://github.com/your/project/contributors) quíenes han participado en este proyecto. 

## Licencia 📄

Este proyecto está bajo la Licencia (Tu Licencia) - mira el archivo [LICENSE.md](LICENSE.md) para detalles

## Expresiones de Gratitud 🎁

* Comenta a otros sobre este proyecto 📢
* Invita una cerveza 🍺 a alguien del equipo. 
* Da las gracias públicamente 🤓.
* etc.



---
⌨️ con ❤️ por [Villanuevand](https://github.com/Villanuevand) 😊
