import express from "express";
import "dotenv/config";
import { initTelemetry, shutdownTelemetry } from "./winston/telemetry";
import { logger } from "./winston/logger";
import { requestLogger } from "./winston/middleware";

// Initialize OpenTelemetry logging
initTelemetry();

const app = express();
const port = process.env.PORT || 3001;

// Middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Routes

// Health check route
app.get("/health", (req, res) => {
  logger.info("Health check endpoint accessed");
  res.json({
    status: "OK",
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
  });
});

// OTLP log test route
app.get("/test-otlp", requestLogger, (req, res) => {
  logger.info("OTLP test endpoint accessed", {});

  logger.warn("OTLP warning test log", {
    warning: "This is a warning message",
    code: "WARN_001",
    stack: new Error("Test warning").stack,
  });

  logger.error("OTLP error test log", {
    error: "This is an error message",
    code: "ERR_001",
    stack: new Error("Test error").stack,
  });

  res.json({
    success: true,
    message: "OTLP test logs sent to Alloy",
    otlpEndpoint: process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs",
    timestamp: new Date().toISOString(),
  });
});

// Simple log generation route
app.get("/simple-generate-log", (req, res) => {
  const logLevel = (req.query.level as string) || "info";
  const message = (req.query.message as string) || "Simple log generated";
  const count = parseInt(req.query.count as string) || 1;

  logger.info("Log generation endpoint accessed", {
    logLevel,
    message,
    count,
    endpoint: "/simple-generate-log",
  });

  // Generate multiple logs based on count
  for (let i = 0; i < count; i++) {
    switch (logLevel.toLowerCase()) {
      case "error":
        logger.error(`${message} - Log ${i + 1}/${count}`, {
          logNumber: i + 1,
          totalLogs: count,
          generatedAt: new Date().toISOString(),
        });
        break;
      case "warn":
        logger.warn(`${message} - Log ${i + 1}/${count}`, {
          logNumber: i + 1,
          totalLogs: count,
          generatedAt: new Date().toISOString(),
        });
        break;
      case "debug":
        logger.debug(`${message} - Log ${i + 1}/${count}`, {
          logNumber: i + 1,
          totalLogs: count,
          generatedAt: new Date().toISOString(),
        });
        break;
      default:
        logger.info(`${message} - Log ${i + 1}/${count}`, {
          logNumber: i + 1,
          totalLogs: count,
          generatedAt: new Date().toISOString(),
        });
    }
  }

  res.json({
    success: true,
    message: `Generated ${count} log(s) with level: ${logLevel}`,
    timestamp: new Date().toISOString(),
    logsGenerated: count,
  });
});

// User simulation route
app.get("/user/:id", (req, res) => {
  const userId = req.params.id;

  logger.info("User endpoint accessed", {
    userId,
    endpoint: "/user/:id",
  });

  // Simulate some user data
  const userData = {
    id: userId,
    name: `User ${userId}`,
    email: `user${userId}@example.com`,
    lastLogin: new Date().toISOString(),
    status: "active",
  };

  res.json(userData);
});

// Products listing route
app.get("/products", (req, res) => {
  const page = parseInt(req.query.page as string) || 1;
  const limit = parseInt(req.query.limit as string) || 10;

  logger.info("Products endpoint accessed", {
    page,
    limit,
    endpoint: "/products",
  });

  // Simulate product data
  const products = Array.from({ length: limit }, (_, index) => ({
    id: (page - 1) * limit + index + 1,
    name: `Product ${(page - 1) * limit + index + 1}`,
    price: Math.floor(Math.random() * 1000) + 10,
    category: ["Electronics", "Clothing", "Books", "Home"][
      Math.floor(Math.random() * 4)
    ],
  }));

  res.json({
    products,
    pagination: {
      page,
      limit,
      total: products.length,
    },
  });
});

// Orders route
app.get("/orders", (req, res) => {
  const status = (req.query.status as string) || "all";

  logger.info("Orders endpoint accessed", {
    status,
    endpoint: "/orders",
  });

  // Simulate order data
  const orders = [
    {
      id: 1,
      status: "pending",
      amount: 99.99,
      createdAt: new Date().toISOString(),
    },
    {
      id: 2,
      status: "completed",
      amount: 149.5,
      createdAt: new Date().toISOString(),
    },
    {
      id: 3,
      status: "cancelled",
      amount: 75.25,
      createdAt: new Date().toISOString(),
    },
  ];

  const filteredOrders =
    status === "all"
      ? orders
      : orders.filter((order) => order.status === status);

  res.json({
    orders: filteredOrders,
    count: filteredOrders.length,
    filter: status,
  });
});

// Statistics route
app.get("/stats", (req, res) => {
  logger.info("Statistics endpoint accessed", {
    endpoint: "/stats",
  });

  // Generate some random stats
  const stats = {
    totalUsers: Math.floor(Math.random() * 10000) + 1000,
    totalOrders: Math.floor(Math.random() * 5000) + 500,
    totalRevenue: (Math.random() * 100000 + 10000).toFixed(2),
    activeUsers: Math.floor(Math.random() * 1000) + 100,
    timestamp: new Date().toISOString(),
  };

  res.json(stats);
});

// Error simulation route
app.get("/simulate-error", (req, res) => {
  const errorType = (req.query.type as string) || "generic";

  logger.error("Error simulation endpoint accessed", {
    errorType,
    endpoint: "/simulate-error",
  });

  switch (errorType) {
    case "auth":
      logger.error("Authentication error simulated", { errorCode: 401 });
      return res.status(401).json({ error: "Unauthorized access" });
    case "notfound":
      logger.error("Not found error simulated", { errorCode: 404 });
      return res.status(404).json({ error: "Resource not found" });
    case "server":
      logger.error("Server error simulated", { errorCode: 500 });
      return res.status(500).json({ error: "Internal server error" });
    default:
      logger.error("Generic error simulated", { errorCode: 400 });
      return res.status(400).json({ error: "Bad request" });
  }
});

// Add this route after the existing routes
app.get("/test-direct-otlp", requestLogger, async (req, res) => {
  try {
    // Test direct OTLP HTTP call
    const testLog = {
      resourceLogs: [
        {
          resource: {
            attributes: [
              { key: "service.name", value: { stringValue: "lgtm-client" } },
              { key: "service.version", value: { stringValue: "1.0.0" } },
            ],
          },
          scopeLogs: [
            {
              scope: { name: "test-scope" },
              logRecords: [
                {
                  body: { stringValue: "Direct OTLP test log" },
                  severityText: "INFO",
                  severityNumber: 9,
                  timeUnixNano: (Date.now() * 1000000).toString(),
                  attributes: [
                    { key: "test.direct", value: { stringValue: "true" } },
                    {
                      key: "endpoint",
                      value: { stringValue: "/test-direct-otlp" },
                    },
                  ],
                },
              ],
            },
          ],
        },
      ],
    };

    const response = await fetch(
      process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(testLog),
      }
    );

    if (response.ok) {
      res.json({
        success: true,
        message: "Direct OTLP test successful",
        status: response.status,
        endpoint: process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs",
      });
    } else {
      const errorText = await response.text();
      res.status(500).json({
        success: false,
        message: "Direct OTLP test failed",
        status: response.status,
        error: errorText,
      });
    }
  } catch (error) {
    console.error("Direct OTLP test error:", error);
    res.status(500).json({
      success: false,
      message: "Direct OTLP test error",
      error: (error as any).message,
    });
  }
});

// 404 handler
app.use("*", (req, res) => {
  logger.warn("Route not found", {
    method: req.method,
    path: req.path,
    ip: req.ip,
  });

  res.status(404).json({
    error: "Route not found",
    path: req.path,
    method: req.method,
  });
});

// Error handler
app.use(
  (
    err: Error,
    req: express.Request,
    res: express.Response,
    next: express.NextFunction
  ) => {
    logger.error("Unhandled error", {
      error: err.message,
      stack: err.stack,
      method: req.method,
      path: req.path,
    });

    res.status(500).json({
      error: "Internal server error",
      message:
        process.env.NODE_ENV === "development"
          ? err.message
          : "Something went wrong",
    });
  }
);

// Start server
const server = app.listen(port, () => {
  logger.info("Server started", {
    port,
    environment: process.env.NODE_ENV || "development",
    otlpEndpoint: process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs",
    timestamp: new Date().toISOString(),
  });
  console.log(`🚀 Server running on http://localhost:${port}`);
  console.log(
    `📡 OTLP logs sending to: ${
      process.env.OTLP_ENDPOINT || "http://localhost:4318/v1/logs"
    }`
  );
});

// Graceful shutdown
const gracefulShutdown = async (signal: string) => {
  logger.info("Received shutdown signal", { signal });
  console.log(`\n🛑 Received ${signal}, shutting down gracefully...`);

  server.close(async () => {
    logger.info("HTTP server closed");
    console.log("📴 HTTP server closed");

    // Shutdown telemetry
    await shutdownTelemetry();

    process.exit(0);
  });

  // Force close after 10 seconds
  setTimeout(() => {
    console.error(
      "⏰ Could not close connections in time, forcefully shutting down"
    );
    process.exit(1);
  }, 10000);
};

process.on("SIGTERM", () => gracefulShutdown("SIGTERM"));
process.on("SIGINT", () => gracefulShutdown("SIGINT"));

export default app;
