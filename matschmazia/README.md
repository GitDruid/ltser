# Matsch-Mazia

Matsch-Mazia specific tools.

**Ingestor** is a service that expose a REST API to receive JSONs with raw data from matsch-mazia sensor network. After some validation it will store the data in the database.

**Visualizer** is a tool to query the database. [*TEMPORARY*: it contains also some experimental logic to perform ADF test and Hampel filtering on series]