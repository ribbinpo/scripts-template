import { env } from "@/config/env";
import amqp from "amqplib";

// RabbitMQ connection
let amqpClient: amqp.ChannelModel | null = null;
let channel: amqp.Channel | null = null;

export const connect = async () => {
  amqpClient = await amqp.connect(env.rabbitmq.url);
  channel = await amqpClient.createChannel();
  return { amqpClient, channel };
};

export const disconnect = async () => {
  await channel?.close();
  await amqpClient?.close();
};

const rabbitmqConfig = { connect, disconnect };

export default rabbitmqConfig;
