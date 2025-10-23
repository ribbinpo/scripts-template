import winston from "winston";
import "dotenv/config";
import OTLPTransport from "./winston-otlp-transport";

// Configure Winston logger with OTLP transport
export const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || "info",
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  defaultMeta: {
    service: "lgtm-client",
    environment: process.env.NODE_ENV || "development",
    version: "1.0.0",
  },
  transports: [
    // Console transport for development
    new winston.transports.Console({
      format: winston.format.combine(
        winston.format.colorize(),
        winston.format.simple()
      ),
    }),
    // File transport for local logging
    // new winston.transports.File({
    //   filename: "/var/log/app.log",
    //   format: winston.format.json(),
    // }),
    // OTLP transport to send logs to Alloy
    new OTLPTransport({
      level: process.env.OTLP_LOG_LEVEL || "info",
    }),
  ],
});
