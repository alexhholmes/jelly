# Jelly - Photo Sharing API

A Go-based photo sharing social media API built with minimal dependencies.

## Development TODO List

### Health Check Enhancements

#### Database Connectivity
- [ ] Add PostgreSQL connection health check
- [ ] Verify read/write operations work in health check
- [ ] Test connection pool health and metrics
- [ ] Add database migration status check
- [ ] Monitor database query performance metrics

#### External Service Dependencies
- [ ] Add S3 bucket connectivity check
- [ ] Verify S3 read/write permissions
- [ ] Test S3 authentication and credentials validity
- [ ] Add external API dependency checks (if any)
- [ ] Monitor network connectivity to dependent services

#### Application-Specific Checks
- [ ] Validate configuration is loaded correctly
- [ ] Add Redis/cache system connectivity (when implemented)
- [ ] Verify file system permissions for uploads/storage
- [ ] Test photo upload directory write permissions
- [ ] Add photo processing pipeline health checks

#### File Storage & Processing
- [ ] Implement S3 integration for photo storage
- [ ] Add file compression service integration
- [ ] Implement image format conversion service
- [ ] Add thumbnail generation service
- [ ] Implement file validation and virus scanning
- [ ] Add CDN integration for photo delivery
- [ ] Add check for to test photo upload and processing pipeline timings

#### Enhanced Health Check Response
- [ ] Replace simple `{"status": "ok"}` with detailed status
- [ ] Add timestamp and version information
- [ ] Include individual component status checks
- [ ] Add uptime and performance metrics
- [ ] Implement degraded service status handling

Example target response structure:
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.2.3",
  "checks": {
    "database": "ok",
    "s3_storage": "ok",
    "image_processing": "degraded",
    "disk_space": "ok",
    "memory": "ok"
  },
  "uptime": "72h30m15s",
  "metrics": {
    "photos_uploaded": 1234,
    "active_connections": 45,
    "avg_response_time": "120ms"
  }
}
```

#### Monitoring & Alerting
- [ ] Add structured logging for health check failures
- [ ] Implement health check metrics collection
- [ ] Add Prometheus metrics endpoint
- [ ] Create health check alerts for critical failures
- [ ] Add health check dashboard integration

#### Configuration Management
- [ ] Add environment-specific health check configurations
- [ ] Implement health check timeout settings
- [ ] Add health check retry logic
- [ ] Create health check endpoint versioning
- [ ] Add health check authentication (if needed)
