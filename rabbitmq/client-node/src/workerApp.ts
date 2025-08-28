import { consumeCreateExampleCommand } from "./handlers/example.command";
import { emitExampleCreated } from "./handlers/example.emit";
import { consumeExampleCreatedEvent } from "./handlers/example.event";
import rabbitmqConfig from "./rabbitmq.config";
import { assertInfra, publisher } from "./util";

export const runWorkerApp = async () => {
  const { channel } = await rabbitmqConfig.connect();
  await assertInfra(channel, []);
  const publish = publisher(channel);

  await consumeExampleCreatedEvent(channel);
  await consumeCreateExampleCommand(channel);
  await emitExampleCreated(publish, { id: "1", email: "test@test.com" });
};

runWorkerApp();
