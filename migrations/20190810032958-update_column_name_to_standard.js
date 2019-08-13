'use strict';

module.exports = {
    up: (queryInterface, Sequelize) => {
        return Promise.all([
            queryInterface.renameColumn("users", "roleId", "role_id"),
            queryInterface.renameColumn("users", "updatedAt", "updated_at"),
            queryInterface.renameColumn("users", "createdAt", "created_at"),
            queryInterface.renameColumn("roles", "updatedAt", "updated_at"),
            queryInterface.renameColumn("roles", "createdAt", "created_at"),
            queryInterface.renameColumn("products", "updatedAt", "updated_at"),
            queryInterface.renameColumn("products", "createdAt", "created_at"),
            queryInterface.renameColumn("stockadjustments", "productId", "product_id"),
            queryInterface.renameColumn("stockadjustments", "updatedAt", "updated_at"),
            queryInterface.renameColumn("stockadjustments", "createdAt", "created_at"),
            queryInterface.renameColumn("invoices", "customerId", "customer_id"),
            queryInterface.renameColumn("invoices", "deletedAt", "deleted_at"),
            queryInterface.renameColumn("invoices", "updatedAt", "updated_at"),
            queryInterface.renameColumn("invoices", "createdAt", "created_at"),
            queryInterface.renameColumn("transactions", "invoiceId", "invoice_id"),
            queryInterface.renameColumn("transactions", "productId", "product_id"),
            queryInterface.renameColumn("transactions", "stockMutationId", "stock_mutation_id"),
            queryInterface.renameColumn("transactions", "updatedAt", "updated_at"),
            queryInterface.renameColumn("transactions", "createdAt", "created_at"),
            queryInterface.renameColumn("customers", "updatedAt", "updated_at"),
            queryInterface.renameColumn("customers", "createdAt", "created_at"),
            queryInterface.renameColumn("customers", "deletedAt", "deleted_at"),
            queryInterface.renameColumn("stockmutations", "productId", "product_id"),
            queryInterface.renameColumn("stockmutations", "updatedAt", "updated_at"),
            queryInterface.renameColumn("stockmutations", "createdAt", "created_at"),
            queryInterface.renameColumn("ordertransactions", "orderInvoiceId", "order_invoice_id"),
            queryInterface.renameColumn("ordertransactions", "productId", "product_id"),
            queryInterface.renameColumn("ordertransactions", "stockMutationId", "stock_mutation_id"),
            queryInterface.renameColumn("ordertransactions", "updatedAt", "updated_at"),
            queryInterface.renameColumn("ordertransactions", "createdAt", "created_at"),
            queryInterface.renameColumn("orderinvoices", "updatedAt", "updated_at"),
            queryInterface.renameColumn("orderinvoices", "createdAt", "created_at"),
            queryInterface.renameColumn("returntransactions", "transactionId", "transaction_id"),
            queryInterface.renameColumn("returntransactions", "stockMutationId", "stock_mutation_id"),
            queryInterface.renameColumn("returntransactions", "updatedAt", "updated_at"),
            queryInterface.renameColumn("returntransactions", "createdAt", "created_at"),
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
