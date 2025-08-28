import { Publish } from "../util";

export async function emitExampleCreated(
  publish: Publish,
  example: { id: string; email: string }
) {
  await publish("domain.events", "example.subEvent.created.v1", {
    eventId: crypto.randomUUID(),
    occurredAt: new Date().toISOString(),
    example,
    source: "example",
  });
}
