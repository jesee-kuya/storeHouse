# StoreHouse Expense Tracker API Documentation

This document provides an overview of the REST API endpoints available in the StoreHouse Expense Tracker application.

## Base URL

All API endpoints are prefixed with `/api/v1`.

## Authentication

Currently, the API does not implement authentication. In a production environment, you would need to add JWT token-based authentication or session management.

## Endpoints

### Health Check

- **GET** `/health`
  - Returns the health status of the API
  - Response: `{"status": "ok"}`

### Accounts

- **GET** `/api/v1/accounts`
  - Get all accounts
  - Response: Array of Account objects

- **POST** `/api/v1/accounts`
  - Create a new account
  - Request Body:
    ```json
    {
      "account_name": "Church Building Fund",
      "account_type": "Asset",
      "local_share": 5000.00,
      "notes": "Fund for church building maintenance"
    }
    ```
  - Response: Created Account object

- **GET** `/api/v1/accounts/{id}`
  - Get account by ID
  - Response: Account object

- **PUT** `/api/v1/accounts/{id}`
  - Update account details
  - Request Body: Partial Account object
  - Response: Updated Account object

- **DELETE** `/api/v1/accounts/{id}`
  - Deactivate account (soft delete)
  - Response: Success message

### Transactions

- **GET** `/api/v1/transactions`
  - Get all transactions
  - Response: Array of Transaction objects

- **POST** `/api/v1/transactions`
  - Create a new transaction
  - Request Body:
    ```json
    {
      "transaction_type": "expenses",
      "amount": 150.00,
      "debit_account_id": "account-uuid",
      "notes": "Office supplies"
    }
    ```
  - Response: Created Transaction object

- **GET** `/api/v1/transactions/{id}`
  - Get transaction by ID
  - Response: Transaction object

- **GET** `/api/v1/transactions/ref/{ref}`
  - Get transaction by reference number
  - Response: Transaction object

- **GET** `/api/v1/transactions/account/{accountID}`
  - Get transactions for a specific account
  - Response: Array of Transaction objects

- **GET** `/api/v1/transactions/member/{memberID}`
  - Get transactions for a specific member
  - Response: Array of Transaction objects

- **GET** `/api/v1/transactions/type/{type}`
  - Get transactions by type (receipts, withdrawal, expenses, transfer)
  - Response: Array of Transaction objects

- **GET** `/api/v1/transactions/date-range?start_date={RFC3339}&end_date={RFC3339}`
  - Get transactions within a date range
  - Response: Array of Transaction objects

- **PUT** `/api/v1/transactions/{id}`
  - Update transaction details
  - Response: Updated Transaction object

- **DELETE** `/api/v1/transactions/{id}`
  - Delete a transaction
  - Response: Success message

### Members

- **GET** `/api/v1/members`
  - Get all members
  - Response: Array of Member objects

- **GET** `/api/v1/members/search?q={search_term}`
  - Search members by name, phone, or email
  - Response: Array of Member objects

- **POST** `/api/v1/members`
  - Create a new member
  - Request Body:
    ```json
    {
      "full_name": "John Doe",
      "phone_number": "+1234567890",
      "email": "john@example.com",
      "notes": "Regular member"
    }
    ```
  - Response: Created Member object

- **GET** `/api/v1/members/{id}`
  - Get member by ID
  - Response: Member object

- **GET** `/api/v1/members/phone/{phone}`
  - Get member by phone number
  - Response: Member object

- **GET** `/api/v1/members/email/{email}`
  - Get member by email
  - Response: Member object

- **GET** `/api/v1/members/group/{groupID}`
  - Get members in a specific group
  - Response: Array of Member objects

- **PUT** `/api/v1/members/{id}`
  - Update member details
  - Response: Updated Member object

- **DELETE** `/api/v1/members/{id}`
  - Delete a member
  - Response: Success message

### Users

- **GET** `/api/v1/users`
  - Get all users
  - Response: Array of User objects

- **GET** `/api/v1/users/active`
  - Get all active users
  - Response: Array of User objects

- **GET** `/api/v1/users/role/{role}`
  - Get users by role (Admin, Treasurer, Clerk)
  - Response: Array of User objects

- **POST** `/api/v1/users`
  - Create a new user
  - Request Body:
    ```json
    {
      "username": "admin",
      "email": "admin@church.com",
      "password": "SecurePass123",
      "full_name": "Church Administrator",
      "role": "Admin",
      "phone_number": "+1234567890"
    }
    ```
  - Response: Created User object

- **GET** `/api/v1/users/{id}`
  - Get user by ID
  - Response: User object

- **GET** `/api/v1/users/username/{username}`
  - Get user by username
  - Response: User object

- **GET** `/api/v1/users/email/{email}`
  - Get user by email
  - Response: User object

- **PUT** `/api/v1/users/{id}`
  - Update user details
  - Response: Updated User object

- **DELETE** `/api/v1/users/{id}`
  - Delete a user
  - Response: Success message

- **POST** `/api/v1/users/{id}/deactivate`
  - Deactivate a user
  - Response: Success message

- **POST** `/api/v1/users/authenticate`
  - Authenticate user
  - Request Body:
    ```json
    {
      "username": "admin",
      "password": "SecurePass123"
    }
    ```
  - Response: User object

- **POST** `/api/v1/users/{id}/change-password`
  - Change user password
  - Request Body:
    ```json
    {
      "old_password": "oldpass",
      "new_password": "newpass"
    }
    ```
  - Response: Success message

### Expenditures

- **GET** `/api/v1/expenditures`
  - Get all expenditures
  - Response: Array of Expenditure objects

- **POST** `/api/v1/expenditures`
  - Create a new expenditure
  - Request Body:
    ```json
    {
      "transaction_id": "transaction-uuid",
      "particulars": "Office supplies purchase",
      "bank_account_id": "account-uuid",
      "amount": 150.00
    }
    ```
  - Response: Created Expenditure object

- **GET** `/api/v1/expenditures/{id}`
  - Get expenditure by ID
  - Response: Expenditure object

- **GET** `/api/v1/expenditures/transaction/{transactionID}`
  - Get expenditures for a transaction
  - Response: Array of Expenditure objects

- **PUT** `/api/v1/expenditures/{id}`
  - Update expenditure details
  - Response: Updated Expenditure object

- **DELETE** `/api/v1/expenditures/{id}`
  - Delete an expenditure
  - Response: Success message

### Transfers

- **GET** `/api/v1/transfers`
  - Get all transfers
  - Response: Array of Transfer objects

- **POST** `/api/v1/transfers`
  - Create a new transfer
  - Request Body:
    ```json
    {
      "transaction_id": "transaction-uuid",
      "particulars": "Transfer to savings",
      "credit_account_id": "account-uuid",
      "amount": 1000.00
    }
    ```
  - Response: Created Transfer object

- **GET** `/api/v1/transfers/{id}`
  - Get transfer by ID
  - Response: Transfer object

- **GET** `/api/v1/transfers/transaction/{transactionID}`
  - Get transfers for a transaction
  - Response: Array of Transfer objects

- **GET** `/api/v1/transfers/credit-account/{accountID}`
  - Get transfers for a credit account
  - Response: Array of Transfer objects

- **GET** `/api/v1/transfers/credit-account/{accountID}/total`
  - Get total amount for a credit account
  - Response: `{"total": 5000.00}`

- **GET** `/api/v1/transfers/date-range?start_date={RFC3339}&end_date={RFC3339}`
  - Get transfers within a date range
  - Response: Array of Transfer objects

- **GET** `/api/v1/transfers/date-range/{accountID}?start_date={RFC3339}&end_date={RFC3339}`
  - Get total transfers for an account within a date range
  - Response: `{"total": 2000.00}`

- **PUT** `/api/v1/transfers/{id}`
  - Update transfer details
  - Response: Updated Transfer object

- **DELETE** `/api/v1/transfers/{id}`
  - Delete a transfer
  - Response: Success message

### Receipts

- **GET** `/api/v1/receipts`
  - Get all receipts
  - Response: Array of Receipt objects

- **POST** `/api/v1/receipts`
  - Create a new receipt
  - Request Body:
    ```json
    {
      "transaction_id": "transaction-uuid",
      "income_account_id": "account-uuid",
      "amount": 500.00
    }
    ```
  - Response: Created Receipt object

- **GET** `/api/v1/receipts/{id}`
  - Get receipt by ID
  - Response: Receipt object

- **GET** `/api/v1/receipts/transaction/{transactionID}`
  - Get receipts for a transaction
  - Response: Array of Receipt objects

- **GET** `/api/v1/receipts/account/{accountID}`
  - Get receipts for an account
  - Response: Array of Receipt objects

- **GET** `/api/v1/receipts/account/{accountID}/total`
  - Get total amount for an account
  - Response: `{"total": 3000.00}`

- **GET** `/api/v1/receipts/date-range?start_date={RFC3339}&end_date={RFC3339}`
  - Get receipts within a date range
  - Response: Array of Receipt objects

- **GET** `/api/v1/receipts/date-range/{accountID}?start_date={RFC3339}&end_date={RFC3339}`
  - Get total receipts for an account within a date range
  - Response: `{"total": 1500.00}`

- **PUT** `/api/v1/receipts/{id}`
  - Update receipt details
  - Response: Updated Receipt object

- **DELETE** `/api/v1/receipts/{id}`
  - Delete a receipt
  - Response: Success message

### Members Groups

- **GET** `/api/v1/groups`
  - Get all groups
  - Response: Array of Group objects

- **GET** `/api/v1/groups/with-count`
  - Get all groups with their member counts
  - Response: Array of Group objects with member counts

- **POST** `/api/v1/groups`
  - Create a new group
  - Request Body:
    ```json
    {
      "group_name": "Youth Ministry",
      "notes": "Young adults group"
    }
    ```
  - Response: Created Group object

- **GET** `/api/v1/groups/{id}`
  - Get group by ID
  - Response: Group object

- **GET** `/api/v1/groups/{id}/count`
  - Get member count for a group
  - Response: `{"member_count": 25}`

- **GET** `/api/v1/groups/name/{name}`
  - Get group by name
  - Response: Group object

- **PUT** `/api/v1/groups/{id}`
  - Update group details
  - Response: Updated Group object

- **DELETE** `/api/v1/groups/{id}`
  - Delete a group (only if no members exist)
  - Response: Success message

## Data Models

### Account
```json
{
  "id": "uuid",
  "account_name": "string",
  "account_type": "Bank|Expense|Income|Asset|liability",
  "local_share": "number",
  "notes": "string",
  "is_active": "boolean",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

### Transaction
```json
{
  "id": "uuid",
  "transaction_ref": "string",
  "transaction_date": "RFC3339 timestamp",
  "transaction_type": "receipts|withdrawal|expenses|transfer",
  "amount": "number",
  "notes": "string",
  "debit_account_id": "uuid",
  "member_id": "uuid",
  "created_by": "string",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

### Member
```json
{
  "id": "uuid",
  "full_name": "string",
  "phone_number": "string",
  "email": "string",
  "notes": "string",
  "group_id": "uuid",
  "created_by": "string",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

### User
```json
{
  "id": "uuid",
  "username": "string",
  "email": "string",
  "full_name": "string",
  "role": "Admin|Treasurer|Clerk",
  "phone_number": "string",
  "is_active": "boolean",
  "last_login": "RFC3339 timestamp",
  "created_at": "RFC3339 timestamp",
  "updated_at": "RFC3339 timestamp"
}
```

## Error Responses

All endpoints return appropriate HTTP status codes and error messages in case of failures:

```json
{
  "error": "Error description"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `500` - Internal Server Error

## Running the Application

1. Ensure you have PostgreSQL installed and running
2. Set the `DATABASE_URL` environment variable
3. Run the application:
   ```bash
   go run main.go
   ```
4. The API will be available at `http://localhost:8080`

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (required)
  Format: `postgres://username:password@localhost:5432/database_name`

## Development Notes

- The application uses Go Chi router for HTTP handling
- Database migrations are automatically applied on startup
- All timestamps should be in RFC3339 format
- Amounts are represented as float64 numbers
- Soft deletion is used for accounts and users (deactivation instead of deletion)
- Group deletion is prevented if members exist in the group