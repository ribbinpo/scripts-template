import pino from "pino";
import os from "os";

const logger = pino({
  level: process.env.LOG_LEVEL || "info",
  base: {
    service: process.env.SERVICE_NAME || "service-name",
    env: process.env.NODE_ENV || "prod",
    version: process.env.SERVICE_VERSION || "1.0.0",
    hostname: os.hostname(),
  },
  timestamp: pino.stdTimeFunctions.isoTime,
  formatters: {
    level(label) {
      return { level: label }; // keep level readable
    },
  },
  transport:
    process.env.NODE_ENV !== "production"
      ? {
          target: "pino-pretty",
          options: { colorize: true, translateTime: "SYS:standard" },
        }
      : undefined,
});

export default logger;
