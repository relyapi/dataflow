# DataFlow

**DataFlow** is a high-performance, flexible, and pluggable data ingestion tool that allows you to stream structured or
semi-structured data into various backends including:

- **Elasticsearch**
- **MongoDB**
- **MySQL**
- **Kafka**
- **ZincSearch**
- **Tencent COS**
- ...and more.

## 🚀 Features

- 🔌 **Pluggable Architecture**: Easily extend to support new data sources or sinks.
- ⚡ **High Throughput**: Optimized for efficient batch processing and concurrency.
- 🌊 **Streaming Support**: Real-time data delivery to multiple backends.
- 📦 **Flexible Input Format**: Supports JSON, NDJSON, CSV, and other formats.
- 🔒 **Reliable**: Fault-tolerant with retry mechanisms and error logging.
- 🧱 **Modular Design**: Choose only the components you need.

## 📦 Supported Backends

| Backend       | Supported | Notes                          |
|---------------|-----------|--------------------------------|
| Elasticsearch | ✅         | Bulk indexing supported        |
| MongoDB       | ✅         | Batched insert with upsert     |
| MySQL         | ✅         | Prepared statements & batching |
| Kafka         | ✅         | JSON or Avro output supported  |
| ZincSearch    | ✅         | Lightweight search backend     |
| Tencent COS   | ✅         | Streaming to object storage    |

