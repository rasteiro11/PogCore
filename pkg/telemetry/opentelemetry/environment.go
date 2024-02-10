package opentelemetry

import (
	"context"

	"github.com/rasteiro11/PogCore/pkg/config"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envExporterEndpoint             = "OTEL_EXPORTER_OTLP_ENDPOINT"
	envExporterProtocol             = "OTEL_EXPORTER_OTLP_PROTOCOL"
	envExporterInsecure             = "OTEL_EXPORTER_OTLP_INSECURE"
	envExporterCertificateAuthority = "OTEL_EXPORTER_OTLP_CERTIFICATE"
	envExporterClientKey            = "OTEL_EXPORTER_OTLP_CLIENT_KEY"
	envExporterClientCertificate    = "OTEL_EXPORTER_OTLP_CLIENT_CERTIFICATE"
)

type envConfig struct {
	endpoint              string
	protocol              string
	insecure              bool
	certCA                string
	certClientKey         string
	certClientCertificate string
}

func newEnvConfig() envConfig {
	cfg := envConfig{
		endpoint: "localhost:4317",
		protocol: "grpc",
		insecure: true,
	}

	if val := config.Instance().String(envExporterEndpoint); val != "" {
		cfg.endpoint = val
	}

	if val := config.Instance().String(envExporterProtocol); val != "" {
		cfg.protocol = val
	}

	if ok := config.Instance().Bool(envExporterInsecure); ok {
		cfg.insecure = ok
	}

	cfg.certCA = config.Instance().String(envExporterCertificateAuthority)
	cfg.certClientKey = config.Instance().String(envExporterClientKey)
	cfg.certClientCertificate = config.Instance().String(envExporterClientCertificate)

	return cfg
}

func NewProcessorFromEnv() trace.SpanProcessor {
	cfg := newEnvConfig()

	if cfg.protocol == "grpc" {
		exporter := newExporterFromEnv(context.TODO(), cfg)
		return trace.NewBatchSpanProcessor(exporter)
	}

	return nil
}

func newExporterFromEnv(ctx context.Context, cfg envConfig) *otlptrace.Exporter {
	credentials := newExporterCredentialFromEnv(cfg)

	conn, err := grpc.DialContext(ctx, cfg.endpoint, grpc.WithTransportCredentials(credentials))
	if err != nil {
		logger.Global().Warnf("[main] gprc.DialContext() returned error=%+v\n", err)
		return nil
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		logger.Global().Warnf("[main] otlptracegrpc.New() returned error=%+v\n", err)
		return nil
	}

	return traceExporter
}

func newExporterCredentialFromEnv(cfg envConfig) credentials.TransportCredentials {
	insecure := insecure.NewCredentials()
	if cfg.insecure {
		return insecure
	}

	secure, err := credentials.NewClientTLSFromFile(cfg.certClientKey, cfg.certClientCertificate)
	if err != nil {
		return insecure
	}

	return secure
}

func NewResourceFromEnv() *resource.Resource {
	opts := []resource.Option{
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
		resource.WithOS(),
		resource.WithContainer(),
	}

	r, err := resource.New(context.TODO(), opts...)
	if err != nil {
		return nil
	}

	return r
}
