import { Channel } from "amqplib";
import { consumeJson } from "../util";

// Consumer Command: example.commands.q - create
export async function consumeCreateExampleCommand(channel: Channel) {
  await channel.assertQueue("example.commands.q", {
    durable: true,
    deadLetterExchange: "dlx",
  });
  channel.bindQueue("example.commands.q", "example.commands", "create.v1");

  consumeJson(channel, "example.commands.q", async (cmd) => {
    console.log(" [x] Received %s", cmd);
  });
}
