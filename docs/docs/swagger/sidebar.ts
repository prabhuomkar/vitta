import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebar: SidebarsConfig = {
  apisidebar: [
    {
      type: "doc",
      id: "swagger/vitta-api",
    },
    {
      type: "category",
      label: "Budgets",
      items: [
        {
          type: "doc",
          id: "swagger/create-group",
          label: "Create a new group",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "swagger/get-groups",
          label: "Get all existing groups",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "swagger/update-group",
          label: "Update an existing group",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "swagger/delete-group",
          label: "Delete an existing group",
          className: "api-method delete",
        },
        {
          type: "doc",
          id: "swagger/create-category",
          label: "Create a new category",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "swagger/get-categories",
          label: "Get all existing categories",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "swagger/update-category",
          label: "Update an existing category",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "swagger/delete-category",
          label: "Delete an existing category",
          className: "api-method delete",
        },
        {
          type: "doc",
          id: "swagger/get-budget",
          label: "Get Budget",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "swagger/set-budget",
          label: "Set Budget",
          className: "api-method put",
        },
      ],
    },
    {
      type: "category",
      label: "Payees",
      items: [
        {
          type: "doc",
          id: "swagger/create-payee",
          label: "Create a new payee",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "swagger/get-payees",
          label: "Get all existing payees",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "swagger/update-payee",
          label: "Update an existing payee",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "swagger/delete-payee",
          label: "Delete an existing payee",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Accounts",
      items: [
        {
          type: "doc",
          id: "swagger/create-account",
          label: "Create a new account",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "swagger/get-accounts",
          label: "Get all existing accounts",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "swagger/update-account",
          label: "Update an existing account",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "swagger/delete-account",
          label: "Delete an existing account",
          className: "api-method delete",
        },
      ],
    },
    {
      type: "category",
      label: "Transactions",
      items: [
        {
          type: "doc",
          id: "swagger/create-transaction",
          label: "Create a new transaction",
          className: "api-method post",
        },
        {
          type: "doc",
          id: "swagger/import-transactions",
          label: "Import transactions",
          className: "api-method put",
        },
        {
          type: "doc",
          id: "swagger/get-transactions",
          label: "Get all existing transactions",
          className: "api-method get",
        },
        {
          type: "doc",
          id: "swagger/update-transaction",
          label: "Update an existing transaction",
          className: "api-method patch",
        },
        {
          type: "doc",
          id: "swagger/delete-transaction",
          label: "Delete an existing transaction",
          className: "api-method delete",
        },
      ],
    },
  ],
};

export default sidebar.apisidebar;
