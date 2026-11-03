# Odigos Group By Trace Processor

The **odigosgroupbytrace** processor buffers spans until their trace is complete and then releases the whole trace as one batch, so processors after it (tail sampling, Insights) see full traces. Odigos installs it on the gateway when trace aggregation is needed, in place of the contrib `groupbytrace` processor.

Unlike contrib `groupbytrace`, it releases a trace as soon as it looks complete instead of always waiting the full window, holds traces with async (Kafka/queue) spans longer, moves spans instead of copying them, and rejects new data under memory pressure instead of growing until OOM.

## When a trace is released

| Trace | Released |
| ----- | -------- |
| Has a PRODUCER/CONSUMER span (async) | `async_wait` after its first span |
| Root span seen, no async spans | `root_settle` after the last span arrived, never later than `sync_wait` |
| No root span yet | `sync_wait` after its first span |
| Buffer full (`max_traces`) | Oldest trace released immediately |

## Configuration

```yaml
processors:
  odigosgroupbytrace:
    sync_wait: 30s               # longest an HTTP/RPC trace (no async spans) is held
    async_wait: 30s              # longest a trace with PRODUCER/CONSUMER spans is held
    root_settle: 7s              # release once the root arrived and no span came for this long; 0 disables
    max_traces: 200000           # buffer cap; past it the oldest trace is released early
    late_spans_max_memory_mib: 70 # memory for remembering released trace IDs (late_spans); 0 disables
```

All fields are optional; the values above are the defaults.

## Internal metrics

In Prometheus the names get the `otelcol_` prefix and `.` becomes `_`.

| Metric | Attributes | Meaning |
| ------ | ---------- | ------- |
| `odigos.groupbytrace.traces_in_memory` | | Traces currently buffered. |
| `odigos.groupbytrace.spans_in_memory` | | Spans currently buffered. |
| `odigos.groupbytrace.traces_released` | `class`, `reason` | Traces sent on, by why they were released. |
| `odigos.groupbytrace.spans_released` | `class`, `reason` | Spans sent on, by why their trace was released. |
| `odigos.groupbytrace.late_spans` | `class`, `reason` | Spans that arrived after their trace was already released, so the trace went out incomplete. |
| `odigos.memory_pressure.spans` | `component`, `action` | Spans rejected (retryable) because the Go memory limit was reached. |

- `class`: `sync` or `async`.
- `reason`: `settled`, `timeout`, `evicted`, `shutdown`; on `late_spans` also `long_running` for spans longer than the wait itself (watches, streams), which no wait increase would catch.
- Rising `late_spans` with `settled`/`timeout` means the waits are too short for the traffic.
