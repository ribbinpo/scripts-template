import logger from "./index";

/**
 * Core log builder for consistent fields
 * @param {object} fields
 * @returns {object} merged fields
 */
function base(fields = {}) {
  return {
    ts: new Date().toISOString(),
    ...fields,
  };
}

export const logHelper = {
  trace(msg: string, context = {}) {
    logger.trace(base({ msg, ...context }));
  },

  debug(msg: string, context = {}) {
    logger.debug(base({ msg, ...context }));
  },

  info(msg: string, context = {}) {
    logger.info(base({ msg, ...context }));
  },

  warn(msg: string, context = {}) {
    logger.warn(base({ msg, ...context }));
  },

  error(errOrMsg: Error | string, context = {}) {
    if (errOrMsg instanceof Error) {
      logger.error(
        base({
          msg: errOrMsg.message,
          error: errOrMsg.name,
          stack: errOrMsg.stack,
          ...context,
        })
      );
    } else {
      logger.error(base({ msg: errOrMsg, ...context }));
    }
  },

  fatal(errOrMsg: Error | string, context = {}) {
    if (errOrMsg instanceof Error) {
      logger.fatal(
        base({
          msg: errOrMsg.message,
          error: errOrMsg.name,
          stack: errOrMsg.stack,
          ...context,
        })
      );
    } else {
      logger.fatal(base({ msg: errOrMsg, ...context }));
    }
  },
};
