'use strict';

module.exports = {
    up: (queryInterface, Sequelize) => {
        return Promise.all([
            queryInterface.renameTable("stockadjustments", "stock_adjustments"),
            queryInterface.renameTable("stockmutations", "stock_mutations"),
            queryInterface.renameTable("ordertransactions", "order_transactions"),
            queryInterface.renameTable("orderinvoices", "order_invoices"),
            queryInterface.renameTable("returntransactions", "return_transactions"),

        ])
    },

    down: (queryInterface, Sequelize) => {
        /*
          Add reverting commands here.
          Return a promise to correctly handle asynchronicity.

          Example:
          return queryInterface.dropTable('users');
        */
        return Promise.all([])

    }
};
