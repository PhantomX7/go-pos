'use strict';

module.exports = {
    up: (queryInterface, Sequelize) => {
        return Promise.all([
            queryInterface.addColumn("return_transactions", "invoice_id", {
                type: Sequelize.INTEGER,
                allowNull: false,
            }),
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
            queryInterface.removeColumn("return_transactions", "invoice_id")
        ])

    }
};
