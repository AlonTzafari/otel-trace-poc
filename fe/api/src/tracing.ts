import { NodeSDK } from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
// import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-proto";

export function startOtel() {
    const sdk = new NodeSDK({
        resource: resourceFromAttributes({
            [ATTR_SERVICE_NAME]: 'api'
        }),

        traceExporter: new OTLPTraceExporter({
            url: process.env.COLLECTOR || 'http://127.0.0.1:4317',
            // url: process.env.COLLECTOR || 'http://127.0.0.1:4318/v1/traces',
        }),

        instrumentations: [
            getNodeAutoInstrumentations()
        ],
    });

    sdk.start();
}
