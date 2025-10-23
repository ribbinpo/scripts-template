import { LoggerProvider } from "@opentelemetry/sdk-logs";
import { OTLPLogExporter } from "@opentelemetry/exporter-logs-otlp-http";
import { BatchLogRecordProcessor } from "@opentelemetry/sdk-logs";
import { Resource } from "@opentelemetry/resources";
import { logs } from "@opentelemetry/api-logs";
import { trace } from "@opentelemetry/api";
import { hostname } from "os";

// Configure the resource with service information
const resource = new Resource({
  "service.name": process.env.SERVICE_NAME || "lgtm-client",
  "service.namespace": process.env.SERVICE_NAMESPACE || "lgtm-stack",
  "deployment.environment": process.env.NODE_ENV || "development",
  "host.name": hostname(),
});

// Configure the OTLP exporter to send logs to Alloy
const otlpExporter = new OTLPLogExporter({
  url: process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs",
  headers: {
    "Content-Type": "application/json",
  },
  timeoutMillis: 5000,
});

otlpExporter.export = ((originalExport) => {
  return (logs: any, resultCallback: any) => {
    console.log(
      `🔄 Exporting ${logs.length} logs to OTLP endpoint: ${
        process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs"
      }`
    );

    return originalExport(logs, (result: any) => {
      if (result.code === 0) {
        console.log("✅ OTLP logs exported successfully");
      } else {
        console.error("❌ OTLP export failed:", result);
      }
      resultCallback(result);
    });
  };
})(otlpExporter.export.bind(otlpExporter));

// Create and configure the logger provider
const loggerProvider = new LoggerProvider({
  resource: resource,
});

// Add the OTLP exporter to the logger provider
loggerProvider.addLogRecordProcessor(
  new BatchLogRecordProcessor(otlpExporter, {
    // Reduce batch timeout for faster testing
    scheduledDelayMillis: 1000,
    maxExportBatchSize: 10,
    maxQueueSize: 100,
  })
);

// Register the logger provider globally
logs.setGlobalLoggerProvider(loggerProvider);

// Helper function to get current trace ID
function getCurrentTraceId(): string | undefined {
  const span = trace.getActiveSpan();
  return span?.spanContext().traceId;
}

// Export the logger for use in the application
export const otelLogger = loggerProvider.getLogger("lgtm-client", "1.0.0");

// Export helper function for getting trace ID
export { getCurrentTraceId };

// Initialize telemetry
export function initTelemetry() {
  console.log("🔗 OpenTelemetry logging initialized");
  console.log(
    `📡 OTLP Endpoint: ${
      process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs"
    }`
  );
  console.log(`🏷️  Service: ${resource.attributes["service.name"]}`);
  console.log(`📦 Version: ${resource.attributes["service.version"]}`);
}

// Shutdown function for graceful cleanup
export async function shutdownTelemetry() {
  try {
    console.log("🔄 Flushing remaining logs...");
    await loggerProvider.forceFlush();
    await loggerProvider.shutdown();
    console.log("🔌 OpenTelemetry logging shutdown complete");
  } catch (error) {
    console.error("❌ Error shutting down OpenTelemetry:", error);
  }
}
