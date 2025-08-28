export const env = {
  rabbitmq: {
    url: process.env.RABBITMQ_URL!,
  },
} as const;
