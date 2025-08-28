import { Channel } from "amqplib";
import { assertConsumerQueues, consumeJson } from "../util";

// Consumer Event: example.subEvent.created.v1
export async function consumeExampleCreatedEvent(channel: Channel) {
  const { main: evQ } = await assertConsumerQueues(channel, "example", [
    { rk: "example.subEvent.created.v1" },
  ]);

  consumeJson(channel, evQ, async (msg, headers) => {
    console.log(" [x] Received %s", msg);
  });
}
