import Transport from 'winston-transport';
import { SeverityNumber } from '@opentelemetry/api-logs';
import { otelLogger } from './telemetry';

interface LogInfo {
  level: string;
  message: string;
  timestamp?: string;
  [key: string]: any;
}

interface OTLPTransportOptions {
  level?: string;
  silent?: boolean;
  name?: string;
  [key: string]: any;
}

class OTLPTransport extends Transport {
  constructor(opts: OTLPTransportOptions = {}) {
    super(opts);
  }

  log(info: LogInfo, callback: () => void) {
    setImmediate(() => {
      this.emit('logged', info);
    });

    // Convert Winston log level to OpenTelemetry severity number
    const severityNumber = this.mapLogLevelToSeverity(info.level);
    
    // Let OpenTelemetry handle timestamp automatically to avoid precision issues
    // The SDK will set the timestamp if not provided

    // Prepare log attributes
    const attributes = {
      'log.level': info.level,
      'service.name': 'lgtm-client',
      ...Object.keys(info)
        .filter(key => !['level', 'message', 'timestamp'].includes(key))
        .reduce((acc, key) => {
          acc[key] = info[key];
          return acc;
        }, {} as Record<string, any>)
    };

    // Emit log to OpenTelemetry
    try {
      otelLogger.emit({
        severityNumber,
        severityText: info.level.toUpperCase(),
        body: info.message,
        attributes,
      });
    } catch (error) {
      console.error('Error sending log to OTLP:', error);
    }

    callback();
  }

  private mapLogLevelToSeverity(level: string): SeverityNumber {
    switch (level.toLowerCase()) {
      case 'debug':
        return SeverityNumber.DEBUG;
      case 'info':
        return SeverityNumber.INFO;
      case 'warn':
        return SeverityNumber.WARN;
      case 'error':
        return SeverityNumber.ERROR;
      default:
        return SeverityNumber.INFO;
    }
  }
}

export default OTLPTransport;
