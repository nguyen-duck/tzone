# 📋 Logging & Error Handling Quick Reference

## 🎯 Quick Guide

This project now uses **comprehensive emoji-based logging** and **panic-free error handling** throughout.

---

## 📊 Logging Patterns

### Success Pattern
```go
log.Println("🔄 Starting operation...")
// ... do work ...
log.Println("✅ Operation completed successfully")
```

### Error Pattern
```go
log.Println("🔄 Starting operation...")
result, err := doSomething()
if err != nil {
    log.Printf("❌ Operation failed: %v", err)
    return fmt.Errorf("operation failed: %w", err)
}
log.Println("✅ Operation completed successfully")
```

### Warning Pattern
```go
if config.Optional == "" {
    log.Println("⚠️ Optional feature not configured, skipping...")
}
```

---

## 🎨 Emoji Reference

| Emoji | Use Case | Example |
|-------|----------|---------|
| 🚀 | App start | `log.Println("🚀 Starting application...")` |
| 🔄 | In progress | `log.Println("🔄 Connecting to database...")` |
| ✅ | Success | `log.Println("✅ Database connected")` |
| ❌ | Error | `log.Printf("❌ Failed to connect: %v", err)` |
| ⚠️ | Warning | `log.Println("⚠️ Config missing, using default")` |
| 🔧 | Setup | `log.Println("🔧 Initializing server...")` |
| 🔑 | Auth/Keys | `log.Println("🔑 Loading credentials...")` |
| 🌐 | Network | `log.Println("🌐 Server listening on :8080")` |
| 📝 | Info | `log.Println("📝 Using system environment")` |
| 🛑 | Fatal | `log.Println("🛑 Application terminated")` |
| 👋 | Shutdown | `log.Println("👋 Graceful shutdown")` |

---

## 🛡️ Error Handling Rules

### ❌ DON'T: Use panic
```go
// ❌ BAD
if err != nil {
    panic(err)
}
```

### ✅ DO: Return errors
```go
// ✅ GOOD
if err != nil {
    log.Printf("❌ Operation failed: %v", err)
    return fmt.Errorf("operation failed: %w", err)
}
```

### ✅ DO: Validate before operations
```go
// ✅ GOOD
if client == nil {
    log.Println("❌ Client is nil")
    return fmt.Errorf("client not initialized")
}

if input == "" {
    log.Println("❌ Input is empty")
    return fmt.Errorf("input cannot be empty")
}
```

### ✅ DO: Wrap errors with context
```go
// ✅ GOOD
result, err := database.Query(...)
if err != nil {
    return fmt.Errorf("failed to query database: %w", err)
}
```

---

## 📝 Common Patterns

### Database Connection
```go
log.Println("🔄 Connecting to database...")

if uri == "" {
    log.Println("❌ Database URI is empty")
    return nil, fmt.Errorf("database URI cannot be empty")
}

client, err := db.Connect(uri)
if err != nil {
    log.Printf("❌ Connection failed: %v", err)
    return nil, fmt.Errorf("failed to connect: %w", err)
}

log.Println("✅ Database connected successfully")
return client, nil
```

### Configuration Loading
```go
log.Println("🔄 Loading configuration...")

value := os.Getenv("CONFIG_VALUE")
if value == "" {
    value = "default"
    log.Printf("⚠️ CONFIG_VALUE not set, using default: %s", value)
} else {
    log.Printf("✅ CONFIG_VALUE: %s", value)
}

return value
```

### Graceful Degradation
```go
log.Println("🔄 Initializing optional feature...")

feature, err := InitOptionalFeature()
if err != nil {
    log.Printf("⚠️ Optional feature failed: %v", err)
    log.Println("⚠️ Continuing without optional feature...")
    feature = nil
} else {
    log.Println("✅ Optional feature initialized")
}

// Continue with or without feature
```

### Resource Cleanup
```go
func DoWork() error {
    log.Println("🔄 Starting work...")

    resource, err := acquireResource()
    if err != nil {
        log.Printf("❌ Failed to acquire resource: %v", err)
        return fmt.Errorf("resource acquisition failed: %w", err)
    }
    defer func() {
        if err := resource.Release(); err != nil {
            log.Printf("⚠️ Failed to release resource: %v", err)
        } else {
            log.Println("✅ Resource released")
        }
    }()

    // Do work with resource

    log.Println("✅ Work completed")
    return nil
}
```

---

## 🔍 Validation Pattern

```go
func ValidateInput(input string) error {
    log.Println("🔄 Validating input...")

    if input == "" {
        log.Println("❌ Validation failed: input is empty")
        return fmt.Errorf("input cannot be empty")
    }

    if len(input) > 100 {
        log.Println("❌ Validation failed: input too long")
        return fmt.Errorf("input must be 100 characters or less")
    }

    log.Println("✅ Input validated")
    return nil
}
```

---

## 🎯 Startup Logging Pattern

```go
func main() {
    log.Println("🚀 Starting Application...")
    log.Printf("📅 Date: %s", time.Now().Format("2006-01-02"))

    // Load config
    log.Println("🔄 Loading configuration...")
    cfg := loadConfig()
    log.Println("✅ Configuration loaded")

    // Validate config
    if err := cfg.Validate(); err != nil {
        log.Printf("❌ Configuration invalid: %v", err)
        log.Println("🛑 Application startup aborted")
        os.Exit(1)
    }
    log.Println("✅ Configuration validated")

    // Connect to services
    log.Println("🔄 Connecting to services...")
    // ... connections ...
    log.Println("✅ All services connected")

    // Start server
    log.Println("🌐 Starting HTTP server...")
    if err := startServer(); err != nil {
        log.Printf("❌ Server failed: %v", err)
        log.Println("🛑 Application terminated")
        os.Exit(1)
    }

    log.Println("👋 Application shutdown gracefully")
}
```

---

## 📊 Logging Levels

Use appropriate logging methods:

### Standard Log (Always visible)
```go
log.Println("🔄 Important operation")
log.Printf("✅ Result: %v", result)
```

### Structured Log (Production monitoring)
```go
slog.Info("Operation completed", "duration", elapsed, "result", result)
slog.Warn("Degraded mode", "reason", reason)
slog.Error("Operation failed", "error", err)
```

### Combined Pattern (Best practice)
```go
log.Printf("❌ Database query failed: %v", err)
slog.Error("Database query failed", "error", err, "query", queryType)
```

---

## 🚦 Status Indicator Pattern

```go
type Service struct {
    client *Client
}

func (s *Service) IsAvailable() bool {
    return s.client != nil
}

func (s *Service) DoWork() error {
    if !s.IsAvailable() {
        log.Println("⚠️ Service not available")
        return fmt.Errorf("service not initialized")
    }

    log.Println("🔄 Performing work...")
    // ... do work ...
    log.Println("✅ Work completed")
    return nil
}
```

---

## 💡 Tips

1. **Always log the START of an operation** with 🔄
2. **Always log the RESULT** with ✅ or ❌
3. **Use emoji consistently** for visual scanning
4. **Include context in error messages**
5. **Wrap errors with `%w`** to preserve error chains
6. **Validate inputs BEFORE operations**
7. **Clean up resources in defer blocks**
8. **Never panic in production code**
9. **Use warnings (⚠️) for non-fatal issues**
10. **Log configuration loading** for debugging

---

## 🎓 Error Message Best Practices

### ✅ GOOD Error Messages
```go
return fmt.Errorf("failed to connect to database: connection timeout after 30s")
return fmt.Errorf("invalid configuration: port must be between 1 and 65535, got %d", port)
return fmt.Errorf("failed to query table %s: %w", tableName, err)
```

### ❌ BAD Error Messages
```go
return errors.New("error")
return fmt.Errorf("failed")
return err  // without context
```

---

## 📈 Monitoring Integration

The logging format is designed to be:
- ✅ Human-readable (emoji visual indicators)
- ✅ Machine-parseable (structured format)
- ✅ Grep-friendly (consistent emoji prefixes)
- ✅ Cloud-ready (slog for structured logs)

### Example: Find all errors
```bash
grep "❌" app.log
```

### Example: Find all warnings
```bash
grep "⚠️" app.log
```

### Example: Monitor startup
```bash
tail -f app.log | grep -E "🚀|✅|❌"
```

---

## 🎉 Summary

- ✅ Use emojis for visual status
- ✅ Log operations start and end
- ✅ Never panic - always return errors
- ✅ Wrap errors with context
- ✅ Validate before operations
- ✅ Clean up resources properly
- ✅ Use warnings for non-fatal issues
- ✅ Make errors actionable

**Happy Logging!** 📊
