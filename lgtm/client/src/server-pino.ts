import express from "express";
import "dotenv/config";

import { logHelper } from "./logger/helper";
import { requestLogger } from "./logger/middleware";

// Create Express app
const app = express();
const port = process.env.PORT || 3002;

// Middleware
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Request logging middleware using the same format as existing middleware
app.use(requestLogger);

// Routes

// Health check route
app.get("/health", (req, res) => {
  logHelper.info("Health check endpoint accessed", {
    endpoint: "/health",
    request_id: (req as any).request_id,
  });

  res.json({
    status: "OK",
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    service: "lgtm-client-pino",
  });
});

// Test different log levels
app.get("/test-logs", (req, res) => {
  const testId = Math.random().toString(36).substr(2, 9);

  logHelper.info("Testing different log levels", {
    testId,
    endpoint: "/test-logs",
    request_id: (req as any).request_id,
  });

  // Test various log levels
  logHelper.debug("Debug level log - detailed information", {
    testId,
    level: "debug",
    details: "This is a debug message with additional context",
  });

  logHelper.info("Info level log - general information", {
    testId,
    level: "info",
    message: "This is an informational message",
  });

  logHelper.warn("Warning level log - potential issue", {
    testId,
    level: "warn",
    warning: "This is a warning message",
    code: "WARN_001",
  });

  logHelper.error("Error level log - error occurred", {
    testId,
    level: "error",
    error: "This is an error message",
    code: "ERR_001",
    stack: new Error("Test error").stack,
  });

  res.json({
    success: true,
    message: "Test logs generated successfully",
    testId,
    timestamp: new Date().toISOString(),
    logLevels: ["debug", "info", "warn", "error"],
  });
});

// User management routes
app.get("/users", (req, res) => {
  const page = parseInt(req.query.page as string) || 1;
  const limit = parseInt(req.query.limit as string) || 10;

  logHelper.info("Users endpoint accessed", {
    endpoint: "/users",
    page,
    limit,
    request_id: (req as any).request_id,
  });

  // Simulate user data
  const users = Array.from({ length: limit }, (_, index) => ({
    id: (page - 1) * limit + index + 1,
    name: `User ${(page - 1) * limit + index + 1}`,
    email: `user${(page - 1) * limit + index + 1}@example.com`,
    status: Math.random() > 0.1 ? "active" : "inactive",
    createdAt: new Date().toISOString(),
  }));

  logHelper.debug("User data generated", {
    userCount: users.length,
    page,
    limit,
  });

  res.json({
    users,
    pagination: {
      page,
      limit,
      total: users.length,
    },
  });
});

app.get("/users/:id", (req, res) => {
  const userId = req.params.id;

  logHelper.info("Single user endpoint accessed", {
    endpoint: "/users/:id",
    userId,
    request_id: (req as any).request_id,
  });

  // Simulate user lookup
  if (userId === "999") {
    logHelper.warn("User not found", {
      userId,
      reason: "User ID 999 does not exist",
    });
    return res.status(404).json({
      error: "User not found",
      userId,
    });
  }

  const user = {
    id: userId,
    name: `User ${userId}`,
    email: `user${userId}@example.com`,
    status: "active",
    lastLogin: new Date().toISOString(),
    profile: {
      age: Math.floor(Math.random() * 50) + 18,
      city: ["New York", "London", "Tokyo", "Paris"][
        Math.floor(Math.random() * 4)
      ],
    },
  };

  logHelper.debug("User data retrieved", {
    userId,
    userStatus: user.status,
  });

  res.json(user);
});

// Product routes
app.get("/products", (req, res) => {
  const category = req.query.category as string;
  const minPrice = parseFloat(req.query.minPrice as string) || 0;
  const maxPrice = parseFloat(req.query.maxPrice as string) || 1000;

  logHelper.info("Products endpoint accessed", {
    endpoint: "/products",
    filters: { category, minPrice, maxPrice },
    request_id: (req as any).request_id,
  });

  // Simulate product data
  const products = Array.from({ length: 20 }, (_, index) => ({
    id: index + 1,
    name: `Product ${index + 1}`,
    price: Math.floor(Math.random() * 900) + 10,
    category: ["Electronics", "Clothing", "Books", "Home", "Sports"][
      Math.floor(Math.random() * 5)
    ],
    inStock: Math.random() > 0.2,
    rating: (Math.random() * 5).toFixed(1),
  }));

  let filteredProducts = products;

  if (category) {
    filteredProducts = products.filter(
      (p) => p.category.toLowerCase() === category.toLowerCase()
    );
    logHelper.debug("Products filtered by category", {
      category,
      originalCount: products.length,
      filteredCount: filteredProducts.length,
    });
  }

  filteredProducts = filteredProducts.filter(
    (p) => p.price >= minPrice && p.price <= maxPrice
  );

  if (minPrice > 0 || maxPrice < 1000) {
    logHelper.debug("Products filtered by price range", {
      minPrice,
      maxPrice,
      filteredCount: filteredProducts.length,
    });
  }

  res.json({
    products: filteredProducts,
    filters: { category, minPrice, maxPrice },
    total: filteredProducts.length,
  });
});

// Order routes
app.post("/orders", (req, res) => {
  const { userId, products, totalAmount } = req.body;

  logHelper.info("Order creation attempt", {
    endpoint: "/orders",
    userId,
    productCount: products?.length || 0,
    totalAmount,
    request_id: (req as any).request_id,
  });

  // Validation
  if (!userId) {
    logHelper.warn("Order creation failed - missing userId", {
      error: "Missing required field: userId",
    });
    return res.status(400).json({
      error: "Missing required field: userId",
    });
  }

  if (!products || products.length === 0) {
    logHelper.warn("Order creation failed - no products", {
      error: "No products in order",
      userId,
    });
    return res.status(400).json({
      error: "Order must contain at least one product",
    });
  }

  // Simulate order creation
  const orderId = Math.floor(Math.random() * 10000) + 1;
  const order = {
    id: orderId,
    userId,
    products,
    totalAmount,
    status: "pending",
    createdAt: new Date().toISOString(),
    estimatedDelivery: new Date(
      Date.now() + 7 * 24 * 60 * 60 * 1000
    ).toISOString(),
  };

  logHelper.info("Order created successfully", {
    orderId,
    userId,
    totalAmount,
    productCount: products.length,
  });

  res.status(201).json(order);
});

app.get("/orders/:id", (req, res) => {
  const orderId = req.params.id;

  logHelper.info("Order lookup attempt", {
    endpoint: "/orders/:id",
    orderId,
    request_id: (req as any).request_id,
  });

  // Simulate order lookup
  if (orderId === "999") {
    logHelper.warn("Order not found", {
      orderId,
      reason: "Order does not exist",
    });
    return res.status(404).json({
      error: "Order not found",
      orderId,
    });
  }

  const order = {
    id: orderId,
    userId: Math.floor(Math.random() * 100) + 1,
    products: [
      { id: 1, name: "Product 1", price: 29.99, quantity: 2 },
      { id: 2, name: "Product 2", price: 49.99, quantity: 1 },
    ],
    totalAmount: 109.97,
    status: ["pending", "processing", "shipped", "delivered"][
      Math.floor(Math.random() * 4)
    ],
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  logHelper.debug("Order retrieved successfully", {
    orderId,
    status: order.status,
    totalAmount: order.totalAmount,
  });

  res.json(order);
});

// Error simulation routes
app.get("/simulate-error", (req, res) => {
  const errorType = (req.query.type as string) || "generic";

  logHelper.info("Error simulation requested", {
    endpoint: "/simulate-error",
    errorType,
    request_id: (req as any).request_id,
  });

  switch (errorType) {
    case "auth":
      logHelper.error("Authentication error simulated", {
        errorType: "auth",
        errorCode: "AUTH_001",
        message: "Invalid authentication credentials",
      });
      return res.status(401).json({
        error: "Unauthorized",
        code: "AUTH_001",
        message: "Invalid authentication credentials",
      });

    case "validation":
      logHelper.warn("Validation error simulated", {
        errorType: "validation",
        errorCode: "VAL_001",
        message: "Invalid input data",
      });
      return res.status(400).json({
        error: "Bad Request",
        code: "VAL_001",
        message: "Invalid input data",
      });

    case "notfound":
      logHelper.warn("Not found error simulated", {
        errorType: "notfound",
        errorCode: "NOT_FOUND_001",
        message: "Resource not found",
      });
      return res.status(404).json({
        error: "Not Found",
        code: "NOT_FOUND_001",
        message: "Resource not found",
      });

    case "server":
      logHelper.error("Server error simulated", {
        errorType: "server",
        errorCode: "SERVER_001",
        message: "Internal server error",
        stack: new Error("Simulated server error").stack,
      });
      return res.status(500).json({
        error: "Internal Server Error",
        code: "SERVER_001",
        message: "Internal server error",
      });

    default:
      logHelper.error("Generic error simulated", {
        errorType: "generic",
        errorCode: "GEN_001",
        message: "Generic error occurred",
      });
      return res.status(400).json({
        error: "Bad Request",
        code: "GEN_001",
        message: "Generic error occurred",
      });
  }
});

// Performance test route
app.get("/performance-test", (req, res) => {
  const iterations = parseInt(req.query.iterations as string) || 100;

  logHelper.info("Performance test started", {
    endpoint: "/performance-test",
    iterations,
    request_id: (req as any).request_id,
  });

  const startTime = Date.now();

  // Simulate some work
  for (let i = 0; i < iterations; i++) {
    logHelper.debug("Performance test iteration", {
      iteration: i + 1,
      totalIterations: iterations,
    });
  }

  const endTime = Date.now();
  const duration = endTime - startTime;

  logHelper.info("Performance test completed", {
    iterations,
    duration: `${duration}ms`,
    averageTime: `${(duration / iterations).toFixed(2)}ms per iteration`,
  });

  res.json({
    success: true,
    iterations,
    duration: `${duration}ms`,
    averageTime: `${(duration / iterations).toFixed(2)}ms per iteration`,
    timestamp: new Date().toISOString(),
  });
});

// 404 handler
app.use("*", (req, res) => {
  logHelper.warn("Route not found", {
    method: req.method,
    path: req.path,
    ip: req.ip,
    userAgent: req.get("User-Agent"),
  });

  res.status(404).json({
    error: "Route not found",
    path: req.path,
    method: req.method,
    timestamp: new Date().toISOString(),
  });
});

// Global error handler
app.use(
  (
    err: Error,
    req: express.Request,
    res: express.Response,
    next: express.NextFunction
  ) => {
    logHelper.error("Unhandled application error", {
      error: err.message,
      stack: err.stack,
      method: req.method,
      path: req.path,
      ip: req.ip,
      userAgent: req.get("User-Agent"),
    });

    res.status(500).json({
      error: "Internal server error",
      message:
        process.env.NODE_ENV === "development"
          ? err.message
          : "Something went wrong",
      timestamp: new Date().toISOString(),
    });
  }
);

// Start server
const server = app.listen(port, () => {
  logHelper.info("Pino Express server started", {
    port,
    environment: process.env.NODE_ENV || "development",
    service: "lgtm-client-pino",
  });
  console.log(`🚀 Pino Express server running on http://localhost:${port}`);
  console.log(
    `📊 Logging with Pino - Level: ${process.env.LOG_LEVEL || "info"}`
  );
});

// Graceful shutdown
const gracefulShutdown = async (signal: string) => {
  logHelper.info("Received shutdown signal", {
    signal,
  });
  console.log(`\n🛑 Received ${signal}, shutting down gracefully...`);

  server.close(() => {
    logHelper.info("HTTP server closed");
    console.log("📴 HTTP server closed");
    process.exit(0);
  });

  // Force close after 10 seconds
  setTimeout(() => {
    logHelper.error("Force shutdown after timeout", {
      timeout: "10s",
    });
    console.error(
      "⏰ Could not close connections in time, forcefully shutting down"
    );
    process.exit(1);
  }, 10000);
};

process.on("SIGTERM", () => gracefulShutdown("SIGTERM"));
process.on("SIGINT", () => gracefulShutdown("SIGINT"));

export default app;
