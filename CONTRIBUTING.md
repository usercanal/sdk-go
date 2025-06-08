# Contributing to UserCanal Go SDK

## Prerequisites
- Go 1.21+
- FlatBuffers compiler (`flatc`) - for schema changes only

## Code Contributions
- **Keep it simple** - avoid over-engineering
- **Follow Go conventions** - use `gofmt`, `golint`
- **Add tests** - especially for new public APIs
- **Update documentation** - README, examples, comments

## API Design
- Keep the public API simple and consistent
- Follow existing patterns:
  ```go
  client.Event(ctx, userID, eventName, properties)
  client.LogInfo(ctx, service, message, data)
  ```
- Maintain backward compatibility when possible

### Schema Changes
```bash
flatc --go -o internal/ schema/event.fbs
flatc --go -o internal/ schema/log.fbs
```

## Project Scope
- **Accept**: Bug fixes, performance improvements, analytics/logging features
- **Don't Accept**: Unrelated features, breaking changes without migration, heavy dependencies

## Getting Help
- Check [GitHub Issues](https://github.com/usercanal/sdk-go/issues)
- Read [ARCHITECTURE.md](ARCHITECTURE.md) for technical details
- Look at `/examples` directory