import { Channel, Options } from "amqplib";

export type Publish = (
  exchange: string,
  routingKey: string,
  msg: object,
  opts?: Options.Publish
) => Promise<void>;

export async function assertInfra(ch: Channel, commands: string[]) {
  // Exchanges
  await ch.assertExchange("domain.events", "topic", { durable: true });
  await ch.assertExchange("retry", "topic", { durable: true });
  await ch.assertExchange("dlx", "topic", { durable: true });
  for (const cmd of commands) {
    await ch.assertExchange(`${cmd}.commands`, "direct", { durable: true });
  }
}

export async function assertConsumerQueues(
  ch: Channel,
  name: string,
  bind: Array<{ rk: string }>,
  events = true
) {
  const base = events ? `${name}.events` : `${name}.commands`;
  const main = `${base}.q`;
  const retry = `${base}.retry.q`;
  const dlq = `${base}.dlq`;

  await ch.assertQueue(main, {
    durable: true,
    deadLetterExchange: "retry",
  });

  await ch.assertQueue(retry, {
    durable: true,
    messageTtl: 30_000, // backoff window (30s here)
    deadLetterExchange: "domain.events",
  });

  await ch.assertQueue(dlq, { durable: true });
  await ch.bindQueue(dlq, "dlx", "dead");

  // Bindings for events queues
  for (const { rk } of bind) {
    await ch.bindQueue(main, "domain.events", rk);
    await ch.bindQueue(retry, "retry", rk);
  }

  return { main, retry, dlq };
}

export function consumeJson<T>(
  ch: Channel,
  queue: string,
  handler: (msg: T, headers: any) => Promise<void>,
  opts = { prefetch: 10, maxRetries: 5 }
) {
  ch.prefetch(opts.prefetch);
  ch.consume(
    queue,
    async (m) => {
      if (!m) return;
      const headers = m.properties.headers || {};
      try {
        const parsed = JSON.parse(m.content.toString()) as T;
        await handler(parsed, headers);
        ch.ack(m);
      } catch (e) {
        const deaths = (headers["x-death"]?.[0]?.count ?? 0) as number;
        if (deaths >= (opts.maxRetries ?? 5)) {
          // publish to dlx with reason
          ch.publish("dlx", "dead", m.content, {
            contentType: m.properties.contentType,
            persistent: true,
            headers: { ...headers, error: (e as Error).message },
          });
          ch.ack(m); // prevent retry loop
        } else {
          ch.nack(m, false, false); // dead-letter to retry (then TTL sends back)
        }
      }
    },
    { noAck: false }
  );
}

export function publisher(ch: Channel): Publish {
  return async (exchange, routingKey, msg, opts) => {
    const payload = Buffer.from(JSON.stringify(msg));
    const ok = ch.publish(exchange, routingKey, payload, {
      contentType: "application/json",
      persistent: true,
      messageId: (msg as any).eventId || (msg as any).commandId,
      timestamp: Date.now(),
      ...opts,
    });
    if (!ok) await new Promise((r) => ch.once("drain", r));
  };
}
