# type-safe-mq

> 📦 Type-safe, language-agnostic message envelope format for distributed systems.

`type-safe-mq` is a lightweight protocol and library for wrapping and transmitting typed Protobuf messages across message queues (Redis Streams, Kafka, Pub/Sub, etc). It ensures message integrity by enforcing type-safe payloads and standardized envelope metadata like origin and timestamp.

---

## ✨ Features

- ✅ **Language support**: Python & Go
- 🔒 **Strong type guarantees** via Protobuf
- 🧩 **Structured envelope**: `payload` + `origin` + `timestamp`
- 📡 Suitable for **Redis Streams**, **gRPC**, **Pub/Sub**, and more
- 🚀 Lightweight and extensible
- 🌐 **Node.js** support is coming soon!

---

## 📦 Envelope Format

All messages follow this canonical envelope structure:

```json
{
  "payload": "<protobuf bytes>",
  "origin": "pod-xyz",
  "timestamp": 1716041999999
}