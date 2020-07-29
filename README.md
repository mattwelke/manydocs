# manydocs

An experiment to create a database backed by Postgres or Google Cloud Bigtable (hereafter referred to as "Bigtable") that lets you store documents.

# Operations

## save doc

Save a document with the "save doc" operation, you can provide a set of "query prefixes". It stores the document so that it can be looked up by the document's automatically-generated ID, and it stores the document an extra time for each query prefix you provide. Then, you can query for the document along with any other document that query prefix matches. See the "query docs" operation below for more details on using query prefixes this way.

## get doc

Use the "get doc" operation to retrieve a single document you previously saved with the "save doc" operation. This operation uses the document's automatically-generated ID, returned with the output of the "save doc" operation, as input.

## query docs

Use the "query docs" operation to retrieve multiple documents at once that each match a particular "query prefix". See "save doc" operation above for more details on query prefixes.

## delete doc

Delete a document by its automatically-generated ID using the "delete doc" operation. The database deletes each copy of the document it previously saved. It tracks which saved copies must be deleted until you use the operation and then clears this state after deleting the document.

# Architecture

## Storage Engine

The database uses an underlying data store (hereafter referred to as "storage engine") to store each document saved. Right now, this is Postgres. In the future, this will also include Bigtable to enable high availability, high performance, and horizontal scalability.

The pattern of retrieving documents by ID equality or ID prefix allows each storage engine to efficiently perform "get doc" and "query docs" operations, provided there are enough Postgres or Bigtable nodes provisioned to handle the data stored and requests performed on the storage engine.

## Application Layer

To control the operations performed on the storage engine, an application accepts HTTP requests, parses operations, and performs as many SQL queries or Bigtable writes and reads as are necessary to complete the operation. In the future, this layer will also control IAM, setting `Cache-Control` header values (to enable global low latency read operations via a CDN such as Cloudflare).

# CI/CD

GitHub Actions runs upon commits to this repository's `master` branch to build and push the application layer to Google Cloud Run.

**WIP - Bigtable support is pending**
