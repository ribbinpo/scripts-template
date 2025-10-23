import { randomUUID } from "crypto";
import express from "express";
import { logger } from "./logger";

export const requestLogger = (
  req: express.Request,
  res: express.Response,
  next: express.NextFunction
) => {
  const requestId = randomUUID();
  (req as any).request_id = requestId;

  const start = process.hrtime.bigint();

  res.on("finish", () => {
    const end = process.hrtime.bigint();
    const latencyMs = Number(end - start) / 1_000_000;

    const context = {
      method: req.method,
      route: req.originalUrl,
      status_code: res.statusCode,
      latency_ms: latencyMs.toFixed(2),
      request_id: requestId,
      ip: req.ip,
      user_agent: req.get("user-agent"),
    };

    if (res.statusCode >= 500) {
      logger.error("HTTP Request", context);
    } else if (res.statusCode >= 400) {
      logger.warn("HTTP Request", context);
    } else {
      logger.info("HTTP Request", context);
    }
  });

  next();
};
