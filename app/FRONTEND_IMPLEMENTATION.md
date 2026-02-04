# Expense Tracker App - API Integration Implementation

## Overview

This document outlines the complete implementation of the Expense Tracker Flutter app with full API integration to the StoreHouse backend. All endpoints from the Go backend have been mapped to the frontend with proper error handling, data parsing, and user interface components.

## Architecture

### Project Structure

```
app/lib/
├── main.dart                    # Main app entry and navigation
├── constants/
│   └── api_constants.dart       # API configuration and constants
├── services/
│   └── api_service.dart         # API client and data models
├── screens/
│   ├── dash_board_screen.dart   # Dashboard with summary statistics
│   ├── transactions_screen.dart # Transaction list with filtering
│   ├── add_transaction_screen.dart # Form to create transactions
│   ├── accounts_screen.dart     # Account management
│   ├── members_screen.dart      # Member management
│   └── groups_screen.dart       # Group management
└── theme/
    └── app_theme.dart           # Theme configuration (Teal/Material Design)
```

## Implemented Features

### 1. **Dashboard Screen**
- **Location:** `screens/dash_board_screen.dart`
- **API Endpoints Used:**
  - `GET /api/v1/transactions` - Fetch all transactions
  - `GET /api/v1/accounts` - Fetch all accounts
  - `GET /api/v1/members` - Fetch all members
- **Features:**
  - Real-time income/expense summary cards
  - Account, member, and transaction count statistics
  - Recent transactions list with type indicators
  - Pull-to-refresh functionality
  - Error handling with retry button

### 2. **Transactions Screen**
- **Location:** `screens/transactions_screen.dart`
- **API Endpoints Used:**
  - `GET /api/v1/transactions` - Get all transactions
  - `GET /api/v1/transactions/type/{type}` - Filter by type
  - `DELETE /api/v1/transactions/{id}` - Delete transaction
- **Features:**
  - Filter chips for transaction types (All, Receipts, Expenses, Withdrawals, Transfers)
  - Transaction cards with amount, date, and reference
  - Delete functionality with confirmation dialog
  - Pull-to-refresh support
  - Empty state handling

### 3. **Add Transaction Screen**
- **Location:** `screens/add_transaction_screen.dart`
- **API Endpoints Used:**
  - `POST /api/v1/transactions` - Create new transaction
  - `GET /api/v1/accounts` - Load available accounts
  - `GET /api/v1/members` - Load available members
- **Features:**
  - Transaction type dropdown (Receipts, Withdrawal, Expenses, Transfer)
  - Account selection (required)
  - Member selection (optional)
  - Amount input with currency formatting
  - Date picker with default to today
  - Reference number field
  - Notes field
  - Form validation
  - Success/error feedback

### 4. **Accounts Screen**
- **Location:** `screens/accounts_screen.dart`
- **API Endpoints Used:**
  - `GET /api/v1/accounts` - Get all accounts
  - `POST /api/v1/accounts` - Create new account
  - `DELETE /api/v1/accounts/{id}` - Deactivate account
- **Features:**
  - Account list with type color coding
  - Account type icons (Bank, Expense, Income, Asset, Liability)
  - Active/Inactive status badge
  - Add account dialog with form
  - Delete/Deactivate functionality
  - Account type selection
  - Local share percentage field

### 5. **Members Screen**
- **Location:** `screens/members_screen.dart`
- **API Endpoints Used:**
  - `GET /api/v1/members` - Get all members
  - `GET /api/v1/members/search?q={query}` - Search members
  - `POST /api/v1/members` - Create new member
  - `DELETE /api/v1/members/{id}` - Delete member
- **Features:**
  - Member list with name, phone, and email
  - Search functionality with query filtering
  - Add member dialog
  - Delete with confirmation
  - Phone and email validation
  - Full name and notes fields

### 6. **Groups Screen**
- **Location:** `screens/groups_screen.dart`
- **API Endpoints Used:**
  - `GET /api/v1/groups` - Get all groups
  - `POST /api/v1/groups` - Create new group
  - `DELETE /api/v1/groups/{id}` - Delete group
- **Features:**
  - Group list display
  - Add group dialog
  - Delete group with confirmation
  - Group name and description fields
  - Pull-to-refresh support

## API Service Implementation

### File: `services/api_service.dart`

#### Core Features:
- **Centralized API client** with error handling
- **Timeout management** (30 seconds default)
- **Request/Response handling** with JSON serialization
- **Exception handling** with custom `ApiException` class
- **Data model classes** for type safety:
  - `Transaction` - Transaction data model
  - `Account` - Account data model
  - `Member` - Member data model
  - `MembersGroup` - Group data model

#### Key Methods:

##### Transaction Endpoints
```dart
getAll Transactions()           // GET /transactions
getTransaction(id)              // GET /transactions/{id}
getTransactionsByAccount(id)    // GET /transactions/account/{id}
getTransactionsByMember(id)     // GET /transactions/member/{id}
getTransactionsByType(type)     // GET /transactions/type/{type}
getTransactionsByDateRange()    // GET /transactions/date-range
createTransaction(data)         // POST /transactions
updateTransaction(id, data)     // PUT /transactions/{id}
deleteTransaction(id)           // DELETE /transactions/{id}
```

##### Account Endpoints
```dart
getAllAccounts()                // GET /accounts
getAccount(id)                  // GET /accounts/{id}
createAccount(data)             // POST /accounts
updateAccount(id, data)         // PUT /accounts/{id}
deleteAccount(id)               // DELETE /accounts/{id}
```

##### Member Endpoints
```dart
getAllMembers()                 // GET /members
searchMembers(query)            // GET /members/search
getMember(id)                   // GET /members/{id}
getMemberByPhone(phone)         // GET /members/phone/{phone}
getMemberByEmail(email)         // GET /members/email/{email}
getMembersByGroup(groupId)      // GET /members/group/{groupId}
createMember(data)              // POST /members
updateMember(id, data)          // PUT /members/{id}
deleteMember(id)                // DELETE /members/{id}
```

##### Group Endpoints
```dart
getAllGroups()                  // GET /groups
getGroupsWithMemberCount()      // GET /groups/with-count
getGroup(id)                    // GET /groups/{id}
getGroupMemberCount(id)         // GET /groups/{id}/count
getGroupByName(name)            // GET /groups/name/{name}
createGroup(data)               // POST /groups
updateGroup(id, data)           // PUT /groups/{id}
deleteGroup(id)                 // DELETE /groups/{id}
```

## Design Language & UI/UX

### Color Scheme
- **Primary Color:** Teal (Colors.teal)
- **Secondary Colors:**
  - Green for income/positive amounts
  - Red for expenses/negative amounts
  - Purple for groups
  - Blue for bank accounts
  - Orange for liabilities

### Components
- **Cards:** White background with subtle gray borders (opacity 0.2)
- **Buttons:** Teal background with white text
- **Icons:** Material Design icons with color coding
- **Forms:** OutlineInputBorder with BorderRadius.circular(8)
- **Status Badges:** Color-coded badges for active/inactive states

### Navigation
- Bottom navigation bar with 6 tabs
- Dynamic app bar title based on current screen
- Material Design 3 patterns throughout

## Dependencies

The following packages were added to `pubspec.yaml`:

```yaml
dependencies:
  flutter:
    sdk: flutter
  cupertino_icons: ^1.0.8
  http: ^1.1.0              # HTTP client for API calls
  intl: ^0.19.0             # Internationalization and formatting
```

## Error Handling

### API Exception Handling
```dart
class ApiException implements Exception {
  final String message;
  final int statusCode;
  
  ApiException(this.message, this.statusCode);
}
```

### Error Scenarios Handled
1. **Network Timeouts** - 30-second timeout with user feedback
2. **Server Errors** - HTTP status code validation (200-299 success range)
3. **Connection Failures** - Exception wrapping with descriptive messages
4. **Empty Responses** - Default empty lists returned
5. **Parsing Errors** - Graceful fallback to default values

### UI Error States
- Error dialogs with retry buttons
- SnackBar notifications for operations
- Loading indicators during async operations
- Empty state messages

## Data Flow

### Typical Request Flow:
1. User triggers action (tap, form submission)
2. Screen calls ApiService method
3. ApiService constructs HTTP request with headers and body
4. Request sent with timeout protection
5. Response parsed and models instantiated
6. Data returned to screen via Future
7. UI updates via setState or FutureBuilder
8. Errors shown in SnackBars or dialogs

### State Management
- Uses Flutter's built-in state management (StatefulWidget)
- FutureBuilder for async data loading
- setState for local state updates
- Future-based data loading pattern

## Testing the Implementation

### Prerequisites
1. Backend server running on `http://localhost:8080`
2. Database properly initialized with migrations
3. Flutter SDK and dependencies installed

### Testing Steps

1. **Dashboard Screen**
   - Verify summary cards display correct totals
   - Confirm recent transactions appear
   - Test pull-to-refresh

2. **Transactions Screen**
   - Test filter chips (all, receipts, expenses, etc.)
   - Verify transaction list updates with filters
   - Test delete transaction with confirmation
   - Verify pagination if applicable

3. **Add Transaction Screen**
   - Test form validation
   - Verify dropdown options load from API
   - Test date picker
   - Confirm successful creation with feedback

4. **Accounts Screen**
   - Verify account list loads
   - Test add account dialog
   - Test account type selection
   - Verify delete/deactivate

5. **Members Screen**
   - Test member list
   - Verify search functionality
   - Test add member
   - Verify delete with confirmation

6. **Groups Screen**
   - Verify group list loads
   - Test add group
   - Test delete group

## Configuration

To change the backend URL, update `ApiConstants` in `lib/constants/api_constants.dart`:

```dart
static const String apiBaseUrl = 'http://localhost:8080/api/v1';
```

For production, update to your production server URL.

## Future Enhancements

1. **Authentication** - Add user login/authentication endpoints
2. **Pagination** - Implement pagination for large lists
3. **Filtering** - Add advanced filtering options
4. **Reports** - Add data export and reporting features
5. **Offline Support** - Local caching with SQLite
6. **Real-time Updates** - WebSocket support for live data
7. **Image Uploads** - Support for receipt images
8. **Advanced Search** - Full-text search capabilities

## Troubleshooting

### Common Issues

1. **"Connection refused" error**
   - Ensure backend is running on port 8080
   - Check if firewall is blocking connection

2. **"Empty list returned"**
   - Verify database has data
   - Check API endpoint parameters
   - Review network tab in DevTools

3. **"Form validation fails"**
   - Ensure all required fields are filled
   - Check field value formats match backend requirements

4. **"Timeout errors"**
   - Increase timeout duration in ApiConstants if needed
   - Check network performance
   - Review backend response times

## References

- [Flutter Documentation](https://flutter.dev/docs)
- [HTTP Package](https://pub.dev/packages/http)
- [Intl Package](https://pub.dev/packages/intl)
- Backend API Documentation: See `API_DOCUMENTATION.md`

