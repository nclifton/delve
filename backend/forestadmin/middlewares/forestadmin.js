const requireAll = require("require-all");
const chalk = require("chalk");
const path = require("path");
const Liana = require("forest-express-sequelize");
const models = require("../models");

module.exports = async function (app) {
  app.use(
    await Liana.init({
      modelsDir: path.join(__dirname, "../models"),
      configDir: path.join(__dirname, "../forest"),
      envSecret: process.env.FOREST_ENV_SECRET,
      authSecret: process.env.FOREST_AUTH_SECRET,
      sequelize: models.Sequelize,
      connections: [models.sequelize.account, models.sequelize.webhook],
    })
  );

  console.log(
    chalk.cyan(
      "Your admin panel is available here: https://app.forestadmin.com/projects"
    )
  );
};
