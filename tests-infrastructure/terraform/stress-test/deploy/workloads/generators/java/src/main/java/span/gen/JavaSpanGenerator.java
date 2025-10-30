package span.gen;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.api.common.AttributeKey;
import io.opentelemetry.api.common.Attributes;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;
import java.util.logging.Level;

public class JavaSpanGenerator {
    private static final Logger logger = Logger.getLogger(JavaSpanGenerator.class.getName());

    public static void main(String[] args) {
        // Get configuration from environment variables
        int spansPerSec = getEnvInt("SPANS_PER_SEC", 1000);
        int spanBytes = getEnvInt("SPAN_BYTES", 1000);

        // Log startup information
        logger.info("Starting Java span generator with " + spansPerSec + " spans/sec, " + spanBytes + " bytes per span");
        logger.info("OTEL_SERVICE_NAME: " + System.getenv("OTEL_SERVICE_NAME"));
        logger.info("OTEL_RESOURCE_ATTRIBUTES: " + System.getenv("OTEL_RESOURCE_ATTRIBUTES"));

        // Create payload
        String payload = "x".repeat(spanBytes);

        // Get tracer from GlobalOpenTelemetry (auto-instrumentation will configure this)
        Tracer tracer = GlobalOpenTelemetry.getTracer("instrumentation-scope-name", "2.0.0");

        // Create scheduled executor for span generation
        ScheduledExecutorService executor = Executors.newScheduledThreadPool(1);

        final int[] totalSpans = {0};

        // Add shutdown hook for graceful termination
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            logger.info("Shutdown signal received, stopping span generator...");
            executor.shutdown();
            try {
                if (!executor.awaitTermination(5, TimeUnit.SECONDS)) {
                    executor.shutdownNow();
                }
            } catch (InterruptedException e) {
                executor.shutdownNow();
            }
            logger.info("Java span generator stopped. Total spans generated: " + totalSpans[0]);
        }));

        // Schedule span generation every second
        executor.scheduleAtFixedRate(() -> {
            try {
                for (int i = 0; i < spansPerSec; i++) {
                    Span span = tracer.spanBuilder("load-span")
                        .setAttribute(AttributeKey.stringKey("payload"), payload)
                        .setAttribute(AttributeKey.stringKey("lang"), "java")
                        .setAttribute(AttributeKey.stringKey("gen"), "java-span-gen")
                        .setAttribute(AttributeKey.longKey("payload_size"), (long) spanBytes)
                        .startSpan();

                    // Simulate some work
                    try {
                        // Add small delay to reduce CPU usage
                        if (i % 100 == 0) {
                            Thread.sleep(1);
                        }
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        return; // Exit the task if interrupted
                    } finally {
                        span.end();
                    }
                }

                totalSpans[0] += spansPerSec;
                logger.info("Generated " + spansPerSec + " spans in this second (total: " + totalSpans[0] + ")");

            } catch (Exception e) {
                logger.log(Level.SEVERE, "Error generating spans", e);
            }
        }, 0, 1, TimeUnit.SECONDS);

        logger.info("Span generator running... Press Ctrl+C to stop");

        // Keep the application running indefinitely
        while (!executor.isShutdown()) {
            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                break;
            }
        }
    }

    private static int getEnvInt(String name, int defaultValue) {
        try {
            String value = System.getenv(name);
            if (value != null && !value.trim().isEmpty()) {
                return Integer.parseInt(value.trim());
            }
        } catch (NumberFormatException e) {
            logger.warning("Invalid value for " + name + ", using default: " + defaultValue);
        }
        return defaultValue;
    }
}
