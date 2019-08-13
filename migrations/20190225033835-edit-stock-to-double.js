"use strict";

module.exports = {
  up: (queryInterface, Sequelize) => {
    /*
      Add altering commands here.
      Return a promise to correctly handle asynchronicity.

      Example:
      return queryInterface.createTable('users', { id: Sequelize.INTEGER });
    */
    return Promise.all([
      queryInterface.changeColumn("Products", "stock", {
        type: Sequelize.DOUBLE,
        allowNull: false
      }),
      queryInterface.changeColumn("StockAdjustments", "amount", {
        type: Sequelize.DOUBLE,
        allowNull: false
      }),
      queryInterface.changeColumn("Transactions", "amount", {
        type: Sequelize.DOUBLE,
        allowNull: false
      }),
      queryInterface.changeColumn("StockMutations", "amount", {
        type: Sequelize.DOUBLE,
        allowNull: false
      }),
      queryInterface.changeColumn("OrderTransactions", "amount", {
        type: Sequelize.DOUBLE,
        allowNull: false
      }),
      queryInterface.changeColumn("ReturnTransactions", "amount", {
        type: Sequelize.DOUBLE,
        allowNull: false
      })
    ])
  },

  down: (queryInterface, Sequelize) => {
    /*
      Add reverting commands here.
      Return a promise to correctly handle asynchronicity.

      Example:
      return queryInterface.dropTable('users');
    */
   return Promise.all([
    queryInterface.changeColumn("Products", "stock", {
      type: Sequelize.INTEGER,
      allowNull: false
    }),
    queryInterface.changeColumn("StockAdjustments", "amount", {
      type: Sequelize.INTEGER,
      allowNull: false
    }),
    queryInterface.changeColumn("Transactions", "amount", {
      type: Sequelize.INTEGER,
      allowNull: false
    }),
    queryInterface.changeColumn("StockMutations", "amount", {
      type: Sequelize.INTEGER,
      allowNull: false
    }),
    queryInterface.changeColumn("OrderTransactions", "amount", {
      type: Sequelize.INTEGER,
      allowNull: false
    }),
    queryInterface.changeColumn("ReturnTransactions", "amount", {
      type: Sequelize.INTEGER,
      allowNull: false
    })
  ])
  }
};
