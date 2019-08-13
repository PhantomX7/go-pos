"use strict";
module.exports = {
  up: (queryInterface, Sequelize) => {
    return queryInterface.createTable("Invoices", {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: Sequelize.INTEGER
      },
      image_url: {
        type: Sequelize.STRING
      },
      customerId: {
        type: Sequelize.INTEGER,
        allowNull: false
      },
      total_capital: {
        type: Sequelize.DOUBLE,
        allowNull: false
      },
      total_sell_price: {
        type: Sequelize.DOUBLE,
        allowNull: false
      },
      total_profit: {
        type: Sequelize.DOUBLE,
        allowNull: false
      },
      description: { type: Sequelize.STRING },
      payment_type: { type: Sequelize.STRING },
      payment_status: {
        type: Sequelize.BOOLEAN,
        allowNull: false
      },
      date: {
        type: Sequelize.DATE,
        allowNull: false
      },
      createdAt: {
        allowNull: false,
        type: Sequelize.DATE
      },
      updatedAt: {
        allowNull: false,
        type: Sequelize.DATE
      },
      deletedAt: {
        type: Sequelize.DATE
      }
    });
  },
  down: (queryInterface, Sequelize) => {
    return queryInterface.dropTable("Invoices");
  }
};
