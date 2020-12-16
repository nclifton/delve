const fs = require("fs");
const path = require("path");
const Sequelize = require("sequelize");

let databases = [
  {
    name: "account",
    connectionString: process.env.ACCOUNT_POSTGRES_URL,
  },
  {
    name: "webhook",
    connectionString: process.env.WEBHOOK_POSTGRES_URL,
  },
  {
    name: "sms",
    connectionString: process.env.SMS_POSTGRES_URL,
  },
  {
    name: "mms",
    connectionString: process.env.MMS_POSTGRES_URL,
  },
];

const sequelize = {};
const db = {};
const models = {};

databases.forEach((databaseInfo) => {
  models[databaseInfo.name] = {};
  const isDevelopment =
    process.env.NODE_ENV === "development" || !process.env.NODE_ENV;
  const databaseOptions = {
    logging: isDevelopment ? console.log : false,
    pool: { maxConnections: 10, minConnections: 1 },
    dialectOptions: {},
  };
  if (
    process.env.DATABASE_SSL &&
    JSON.parse(process.env.DATABASE_SSL.toLowerCase())
  ) {
    databaseOptions.dialectOptions.ssl = true;
  }
  const connection = new Sequelize(
    databaseInfo.connectionString,
    databaseOptions
  );
  sequelize[databaseInfo.name] = connection;
  fs.readdirSync(path.join(__dirname, databaseInfo.name))
    .filter((file) => file.indexOf(".") !== 0 && file !== "index.js")
    .forEach((file) => {
      try {
        const model = connection.import(
          path.join(__dirname, databaseInfo.name, file)
        );
        models[databaseInfo.name][model.name] = model;
      } catch (error) {
        console.error("Model creation error: " + error);
      }
    });
  Object.keys(models[databaseInfo.name]).forEach((modelName) => {
    if ("associate" in models[databaseInfo.name][modelName]) {
      models[databaseInfo.name][modelName].associate(
        sequelize[databaseInfo.name].models
      );
    }
  });
});
db.sequelize = sequelize;
db.Sequelize = Sequelize;
module.exports = db;
