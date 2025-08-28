interface IMessage<T> {
  eventId: string;
  occurredAt: Date;
  source: string;
  data: T;
  reasonInit: string;
}
