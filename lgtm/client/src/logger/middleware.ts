import { randomUUID } from "crypto";
import { logHelper } from "./helper";
import express from "express";

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

    const level =
      res.statusCode >= 500 ? "error" : res.statusCode >= 400 ? "warn" : "info";

    logHelper[level]("HTTP Request", {
      method: req.method,
      route: req.originalUrl,
      status_code: res.statusCode,
      latency_ms: latencyMs.toFixed(2),
      request_id: requestId,
      ip: req.ip,
      user_agent: req.get("user-agent"),
      trace_id: req.headers["x-trace-id"] || null,
      user_id: (req as any).user?.id || null,
    });
  });

  next();
};
