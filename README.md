# manydocs

An experiment to create a database backed by Postgres or Google Cloud Bigtable (Bigtable) that lets you store documents.

When a document is stored, optional "query keys" can be provided, which are each serialized into a prefix that the underlying data store can use to perform key prefix lookups. The data is stored multiple times with each key prefix. This technique allows a querying pattern as powerful as Bigtable, with the database scaling as high as it needs (when backed by Bigtable) to store the data coming in and serve the queries being performed, without the user having to manage GCP. When backed by Postgres, users can operate it less expensively than with Bigtable, with the caveat of it not being horizontally-scalable.

**WIP**
