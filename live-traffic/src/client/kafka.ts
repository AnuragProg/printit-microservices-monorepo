import { CompressionCodecs, CompressionTypes, Kafka } from 'kafkajs';
import { SnappyCodec } from 'kafkajs-snappy-typescript';

CompressionCodecs[CompressionTypes.Snappy] = new SnappyCodec().codec;

const KAFKA_BROKER = process.env.KAFKA_BROKER || "localhost:9092"

const kafkaClient = new Kafka({
	brokers: [KAFKA_BROKER],
});



export default kafkaClient;
