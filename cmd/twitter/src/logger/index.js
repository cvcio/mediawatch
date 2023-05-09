const winston = require('winston');

const enumerateErrorFormat = winston.format((info) => {
  if (info instanceof Error) {
    Object.assign(info, { message: info.stack });
  }
  return info;
});

const logger = winston.createLogger({
  level: process.env.NODE_ENV === 'development' ? 'debug' : 'info',
  format: winston.format.combine(winston.format.timestamp({
    format: 'YYYY-MM-DD HH:mm:ss'
  }),
  enumerateErrorFormat(),
  process.env.NODE_ENV === 'development' ? winston.format.colorize() : winston.format.uncolorize(),
  winston.format.splat(),
  winston.format.printf(({ timestamp, level, message }) => `[${timestamp}] [${level}] ${message}`)),
  transports: [
    new winston.transports.Console({
      stderrLevels: ['error'],
      timestamp: true
    }),
  ],
});

module.exports = logger;
