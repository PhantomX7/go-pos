"use strict";
module.exports = {
  up: (queryInterface, Sequelize) => {
    return queryInterface.createTable("Transactions", {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: Sequelize.INTEGER
      },
      invoiceId: {
        type: Sequelize.INTEGER,
        allowNull: false
      },
      productId: {
        type: Sequelize.INTEGER,
        allowNull: false
      },
      stockMutationId: {
        type: Sequelize.INTEGER,
        allowNull: false
      },
      capital_price: {
        type: Sequelize.DOUBLE,
        allowNull: false
      },
      sell_price: {
        type: Sequelize.DOUBLE,
        allowNull: false
      },
      amount: {
        type: Sequelize.INTEGER,
        allowNull: false
      },
      total_sell_price: {
        type: Sequelize.DOUBLE,
        allowNull: false
      },
      profit: {
        type: Sequelize.DOUBLE,
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
      }
    });
  },
  down: (queryInterface, Sequelize) => {
    return queryInterface.dropTable("Transactions");
  }
};
