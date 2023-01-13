package kafka_service

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/nei7/ntube/internal"
	"github.com/nei7/ntube/internal/datastruct"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

const otelName = "github.com/nei7/ntube/internal/kafka"

const (
	VideoEventCreated = "videos.event.created"
	VideoEventDeleted = "video.event.deleted"
)

type Video struct {
	producer  *kafka.Producer
	topicName string
}

type event struct {
	Type  string
	Value datastruct.Video
}

func NewVideo(producer *kafka.Producer, topicName string) *Video {
	return &Video{
		producer,
		topicName,
	}
}

func (v *Video) Created(ctx context.Context, video datastruct.Video) error {
	return v.publish(ctx, "Video.Created", VideoEventCreated, video)
}

func (v *Video) Delete(ctx context.Context, id string) error {
	return v.publish(ctx, "Video.Deleted", VideoEventDeleted, datastruct.Video{ID: id})
}

func (v *Video) publish(ctx context.Context, spanName, msgType string, video datastruct.Video) error {
	_, span := otel.Tracer(otelName).Start(ctx, spanName)
	defer span.End()

	span.SetAttributes(attribute.KeyValue{
		Key:   semconv.MessagingSystemKey,
		Value: attribute.StringValue("kafka"),
	}, attribute.KeyValue{
		Key:   semconv.MessagingDestinationKey,
		Value: attribute.StringValue(v.topicName),
	})

	var b bytes.Buffer

	evt := event{
		Type:  msgType,
		Value: video,
	}

	if err := json.NewEncoder(&b).Encode(&evt); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Encoder")
	}

	if err := v.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &v.topicName,
			Partition: kafka.PartitionAny,
		},
		Value: b.Bytes(),
	}, nil); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "product.Producer")
	}

	return nil
}
