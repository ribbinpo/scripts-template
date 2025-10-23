import Transport from "winston-transport";
import { SeverityNumber } from "@opentelemetry/api-logs";
import { otelLogger } from "./telemetry";
import { context } from "@opentelemetry/api";

interface LogInfo {
  message: string;
  timestamp: string;
  [key: string]: any;
}

interface LogHandler {
  level: string;
  message: string;
  timestamp: string;
  [key: string]: any;
}

interface OTLPTransportOptions {
  level?: string;
  silent?: boolean;
  name?: string;
  [key: string]: any;
}

class OTLPTransport extends Transport {
  constructor(opts: OTLPTransportOptions = {}) {
    super(opts);
  }

  private logHandler(log: LogHandler) {
    setImmediate(() => {
      this.emit("logged", log);
    });

    // Convert Winston log level to OpenTelemetry severity number
    const severityNumber = this.mapLogLevelToSeverity(log.level);

    const attributes: Record<string, any> = {};

    if (log.level === "error" && log.stack) {
      attributes["exception.stacktrace"] = log.stack;

      // 'name' (e.g., 'Error', 'TypeError') is also useful
      if (log.name) {
        attributes["exception.type"] = log.name;
      }
    }

    // Add additional attributes from the log info, excluding base format fields
    Object.keys(log)
      .filter(
        (key) =>
          ![
            "service",
            "version",
            "hostname",
            "level",
            "trace_id", // We get this from active context now
          ].includes(key)
      )
      .forEach((key) => {
        attributes[key] = log[key];
      });

    // Emit log to OpenTelemetry
    try {
      otelLogger.emit({
        timestamp: new Date(log.timestamp),
        severityNumber,
        severityText: log.level.toUpperCase(),
        body: log.message,
        attributes,
        context: context.active(),
      });
    } catch (error) {
      console.error("Error sending log to OTLP:", error);
    }
  }

  log(info: LogInfo, callback: () => void) {
    if (info.level === "debug") {
      this.debug(info);
    } else if (info.level === "info") {
      this.info(info);
    } else if (info.level === "warn") {
      this.warn(info);
    } else if (info.level === "error") {
      this.error(info);
    }
    callback();
  }

  private debug(info: LogInfo) {
    this.logHandler({ ...info, level: "debug" });
  }

  private info(info: LogInfo) {
    this.logHandler({ ...info, level: "info" });
  }

  private warn(info: LogInfo) {
    this.logHandler({ ...info, level: "warn" });
  }

  private error(info: LogInfo | Error) {
    if (info instanceof Error) {
      const { message, name, stack, ...rest } = info;
      this.logHandler({
        level: "error",
        message: message, // 'message' becomes the body
        stack: stack, // 'stack' will be handled
        name: name, // 'name' will be handled
        timestamp: new Date().toISOString(), // Use current time if it's just an error
        ...rest, // Spread any other custom properties on the error
      });
    } else {
      this.logHandler({ ...info, level: "error" });
    }
  }
  private mapLogLevelToSeverity(level: string): SeverityNumber {
    switch (level.toLowerCase()) {
      case "debug":
        return SeverityNumber.DEBUG;
      case "info":
        return SeverityNumber.INFO;
      case "warn":
        return SeverityNumber.WARN;
      case "error":
        return SeverityNumber.ERROR;
      default:
        return SeverityNumber.INFO;
    }
  }
}

export default OTLPTransport;
