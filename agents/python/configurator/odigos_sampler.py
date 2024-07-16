from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.sampling import Sampler, Decision, ParentBased
from opentelemetry.sdk.trace.export import SimpleSpanProcessor, ConsoleSpanExporter
from threading import Lock

class OdigosSampler(Sampler):
    def __init__(self):
        self._condition = None
        self._lock = Lock()

    def should_sample(self, parent_context, trace_id, name, kind, attributes, links):
        with self._lock:
            if self._condition is None or self._condition(name, attributes):
                print('Recording and sampling')
                return Decision.RECORD_AND_SAMPLE
            else:
                print('droppimg')
                return Decision.DROP

    def get_description(self):
        return "OdigosSampler"

    def update_condition(self, new_condition):
        with self._lock:
            self._condition = new_condition