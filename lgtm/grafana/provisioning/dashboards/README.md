# Grafana Dashboards

This directory contains pre-configured Grafana dashboards that are automatically loaded when Grafana starts.

## Available Dashboards

### 1. System Overview
- **Purpose**: High-level system health and resource usage
- **Metrics**: CPU, Memory, Disk, Network, Docker container stats
- **Use Case**: System monitoring and capacity planning
- **Refresh Rate**: 10 seconds

### 2. Docker Services
- **Purpose**: Detailed Docker container monitoring
- **Metrics**: Container status, CPU/Memory usage, Network I/O, Disk I/O
- **Use Case**: Container performance analysis and troubleshooting
- **Refresh Rate**: 15 seconds

### 3. Application Metrics
- **Purpose**: Application-specific performance metrics
- **Metrics**: HTTP requests, response times, error rates, connections
- **Use Case**: Application performance monitoring and SLO tracking
- **Refresh Rate**: 15 seconds

### 4. Logs Overview
- **Purpose**: Log analysis and error tracking
- **Metrics**: Log volume, error rates, log level distribution
- **Use Case**: Log monitoring and error investigation
- **Refresh Rate**: 15 seconds

## Dashboard Features

### Service Separation
All dashboards use the service labels we configured in Alloy:
- `service`: General service identifier
- `container_name`: Docker container name
- `image`: Docker image name
- `service_name`: Extracted service name from image
- `environment`: Environment detection

### Metric Queries
The dashboards use PromQL queries that leverage our service labeling:
- Filter by service: `{service="client-app"}`
- Filter by container: `{container_name="grafana"}`
- Filter by image: `{image="grafana/grafana"}`

### Thresholds and Alerts
Key metrics have color-coded thresholds:
- **Green**: Normal operation
- **Yellow**: Warning level
- **Red**: Critical level

## Customization

### Adding New Panels
1. Edit the JSON dashboard files
2. Add new panels with appropriate PromQL queries
3. Use service labels for filtering: `{service="your-service-name"}`

### Creating Service-Specific Dashboards
1. Copy an existing dashboard JSON
2. Modify the title and tags
3. Update queries to filter for specific services
4. Add service-specific metrics

### Adding New Metrics
1. Ensure your application exposes metrics on `/metrics` endpoint
2. Add scraping configuration in `alloy/config.alloy`
3. Use appropriate service labels
4. Create panels in relevant dashboards

## Troubleshooting

### Dashboards Not Loading
- Check Grafana logs for provisioning errors
- Verify dashboard JSON syntax
- Ensure dashboard files are readable by Grafana

### Missing Metrics
- Verify Prometheus targets are up
- Check Alloy configuration for scraping rules
- Ensure service labels are properly applied

### Performance Issues
- Reduce dashboard refresh rates
- Limit time range queries
- Use appropriate aggregation functions

## Best Practices

1. **Service Naming**: Use consistent service naming conventions
2. **Label Strategy**: Apply meaningful labels for easy filtering
3. **Refresh Rates**: Set appropriate refresh rates based on metric volatility
4. **Time Ranges**: Use reasonable time ranges for queries
5. **Panel Layout**: Organize panels logically for easy navigation

## Next Steps

1. **Restart Grafana** to load the new dashboards
2. **Customize** dashboards for your specific needs
3. **Add Alerts** based on dashboard thresholds
4. **Create** service-specific dashboards
5. **Set up** automated dashboard provisioning
