/// API Configuration Constants
class ApiConstants {
  // Base Configuration
  static const String apiBaseUrl = 'http://localhost:8080/api/v1';
  static const String hostName = 'localhost';
  static const int port = 8080;
  
  // Connection Settings
  static const Duration connectionTimeout = Duration(seconds: 30);
  static const Duration responseTimeout = Duration(seconds: 30);
  
  // Endpoints
  static const String transactionsEndpoint = '/transactions';
  static const String accountsEndpoint = '/accounts';
  static const String membersEndpoint = '/members';
  static const String groupsEndpoint = '/groups';
  static const String healthCheckEndpoint = '/health';
  
  // Transaction Types
  static const String transactionTypeReceipts = 'receipts';
  static const String transactionTypeWithdrawal = 'withdrawal';
  static const String transactionTypeExpenses = 'expenses';
  static const String transactionTypeTransfer = 'transfer';
  
  // Account Types
  static const String accountTypeBank = 'Bank';
  static const String accountTypeExpense = 'Expense';
  static const String accountTypeIncome = 'Income';
  static const String accountTypeAsset = 'Asset';
  static const String accountTypeLiability = 'liability';
}

/// HTTP Status Codes
class HttpStatus {
  static const int success = 200;
  static const int created = 201;
  static const int accepted = 202;
  static const int badRequest = 400;
  static const int unauthorized = 401;
  static const int forbidden = 403;
  static const int notFound = 404;
  static const int conflict = 409;
  static const int serverError = 500;
  static const int serviceUnavailable = 503;
}

/// Error Messages
class ErrorMessages {
  static const String networkError = 'Network error. Please check your connection.';
  static const String serverError = 'Server error. Please try again later.';
  static const String timeoutError = 'Request timeout. Please try again.';
  static const String invalidData = 'Invalid data provided.';
  static const String unauthorized = 'Unauthorized access. Please login.';
  static const String notFound = 'Resource not found.';
  static const String defaultError = 'An unexpected error occurred.';
}
