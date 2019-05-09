# pgsql-backup-s3
[![](https://img.shields.io/docker/cloud/build/lnsp/pgsql-backup-s3.svg)](https://cloud.docker.com/repository/docker/lnsp/pgsql-backup-s3)

Small tool that performs a PostgreSQL backup and stores the result in a S3-compatible service.

## Environment configuration

| Variable       | Description                  | Default value        |
|----------------|------------------------------|----------------------|
| `HOST`         | Host parameter for `pg_dump` | `"localhost"`        |
| `PORT`         | Port parameter for `pg_dump` | `"5432"`             |
| `DATABASE`     | Database to dump             | none                 |
| `USER`         | Access role                  | `"root"`             |
| `PASSWORD`     | Access password              | `"root"`             |
| `PGDUMPBINARY` | Binary location of `pg_dump` | `"/usr/bin/pg_dump"` |
| `ACCESSKEY`    | Access key for S3 service    | none                 |
| `SECRETKEY`    | Secret key for S3 service    | none                 |
| `ENDPOINT`     | S3 endpoint name             | none                 |
| `BUCKET`       | S3 bucket identifier         | none                 |
| `PATH`         | Backup file path             | `""`                 |
| `PREFIX`       | Backup file prefix           | `""`                 |

## Run on Kubernetes

If you want to run the prebuilt container on Kubernetes, you are welcomed to do so. In the `/docs` folder you find a `CronJob` manifest example.

## Build from scratch

Since this project supports Go modules, building it from scratch is very straightforward.

```bash
git clone https://github.com/lnsp/pgsql-backup-s3.git
cd pgsql-backup-s3
go build
```
