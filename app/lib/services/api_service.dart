import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:async';

class ApiService {
  static const String baseUrl = 'http://localhost:8080/api/v1';
  static const Duration timeout = Duration(seconds: 30);

  static Future<T> _makeRequest<T>(
    String endpoint, {
    String method = 'GET',
    dynamic body,
    required T Function(dynamic) parser,
  }) async {
    try {
      final url = Uri.parse('$baseUrl$endpoint');
      late http.Response response;

      final headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      };

      if (method == 'GET') {
        response = await http.get(url, headers: headers).timeout(timeout);
      } else if (method == 'POST') {
        response = await http.post(
          url,
          headers: headers,
          body: jsonEncode(body),
        ).timeout(timeout);
      } else if (method == 'PUT') {
        response = await http.put(
          url,
          headers: headers,
          body: jsonEncode(body),
        ).timeout(timeout);
      } else if (method == 'DELETE') {
        response = await http.delete(url, headers: headers).timeout(timeout);
      }

      if (response.statusCode >= 200 && response.statusCode < 300) {
        final decodedBody = jsonDecode(response.body);
        return parser(decodedBody);
      } else {
        throw ApiException(
          'Failed with status ${response.statusCode}',
          response.statusCode,
        );
      }
    } on TimeoutException {
      throw ApiException('Request timeout', 408);
    } catch (e) {
      throw ApiException('Error: $e', 500);
    }
  }

  // Transaction endpoints
  static Future<List<Transaction>> getAllTransactions() {
    return _makeRequest(
      '/transactions',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Transaction.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<Transaction> getTransaction(String id) {
    return _makeRequest(
      '/transactions/$id',
      parser: (data) => Transaction.fromJson(data),
    );
  }

  static Future<List<Transaction>> getTransactionsByAccount(String accountId) {
    return _makeRequest(
      '/transactions/account/$accountId',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Transaction.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<List<Transaction>> getTransactionsByMember(String memberId) {
    return _makeRequest(
      '/transactions/member/$memberId',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Transaction.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<List<Transaction>> getTransactionsByType(String type) {
    return _makeRequest(
      '/transactions/type/$type',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Transaction.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<List<Transaction>> getTransactionsByDateRange(
    DateTime startDate,
    DateTime endDate,
  ) {
    return _makeRequest(
      '/transactions/date-range?start=${startDate.toIso8601String()}&end=${endDate.toIso8601String()}',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Transaction.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<Transaction> createTransaction(Map<String, dynamic> data) {
    return _makeRequest(
      '/transactions',
      method: 'POST',
      body: data,
      parser: (data) => Transaction.fromJson(data),
    );
  }

  static Future<Transaction> updateTransaction(
    String id,
    Map<String, dynamic> data,
  ) {
    return _makeRequest(
      '/transactions/$id',
      method: 'PUT',
      body: data,
      parser: (data) => Transaction.fromJson(data),
    );
  }

  static Future<void> deleteTransaction(String id) {
    return _makeRequest(
      '/transactions/$id',
      method: 'DELETE',
      parser: (data) => null,
    );
  }

  // Account endpoints
  static Future<List<Account>> getAllAccounts() {
    return _makeRequest(
      '/accounts',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Account.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<Account> getAccount(String id) {
    return _makeRequest(
      '/accounts/$id',
      parser: (data) => Account.fromJson(data),
    );
  }

  static Future<Account> createAccount(Map<String, dynamic> data) {
    return _makeRequest(
      '/accounts',
      method: 'POST',
      body: data,
      parser: (data) => Account.fromJson(data),
    );
  }

  static Future<Account> updateAccount(String id, Map<String, dynamic> data) {
    return _makeRequest(
      '/accounts/$id',
      method: 'PUT',
      body: data,
      parser: (data) => Account.fromJson(data),
    );
  }

  static Future<void> deleteAccount(String id) {
    return _makeRequest(
      '/accounts/$id',
      method: 'DELETE',
      parser: (data) => null,
    );
  }

  // Member endpoints
  static Future<List<Member>> getAllMembers() {
    return _makeRequest(
      '/members',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Member.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<List<Member>> searchMembers(String query) {
    return _makeRequest(
      '/members/search?q=$query',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Member.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<Member> getMember(String id) {
    return _makeRequest(
      '/members/$id',
      parser: (data) => Member.fromJson(data),
    );
  }

  static Future<Member> getMemberByPhone(String phone) {
    return _makeRequest(
      '/members/phone/$phone',
      parser: (data) => Member.fromJson(data),
    );
  }

  static Future<Member> getMemberByEmail(String email) {
    return _makeRequest(
      '/members/email/$email',
      parser: (data) => Member.fromJson(data),
    );
  }

  static Future<List<Member>> getMembersByGroup(String groupId) {
    return _makeRequest(
      '/members/group/$groupId',
      parser: (data) {
        if (data is List) {
          return data.map((item) => Member.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<Member> createMember(Map<String, dynamic> data) {
    return _makeRequest(
      '/members',
      method: 'POST',
      body: data,
      parser: (data) => Member.fromJson(data),
    );
  }

  static Future<Member> updateMember(String id, Map<String, dynamic> data) {
    return _makeRequest(
      '/members/$id',
      method: 'PUT',
      body: data,
      parser: (data) => Member.fromJson(data),
    );
  }

  static Future<void> deleteMember(String id) {
    return _makeRequest(
      '/members/$id',
      method: 'DELETE',
      parser: (data) => null,
    );
  }

  // Group endpoints
  static Future<List<MembersGroup>> getAllGroups() {
    return _makeRequest(
      '/groups',
      parser: (data) {
        if (data is List) {
          return data.map((item) => MembersGroup.fromJson(item)).toList();
        }
        return [];
      },
    );
  }

  static Future<List<Map<String, dynamic>>> getGroupsWithMemberCount() {
    return _makeRequest(
      '/groups/with-count',
      parser: (data) {
        if (data is List) {
          return List<Map<String, dynamic>>.from(data);
        }
        return [];
      },
    );
  }

  static Future<MembersGroup> getGroup(String id) {
    return _makeRequest(
      '/groups/$id',
      parser: (data) => MembersGroup.fromJson(data),
    );
  }

  static Future<int> getGroupMemberCount(String id) {
    return _makeRequest(
      '/groups/$id/count',
      parser: (data) {
        if (data is Map && data.containsKey('count')) {
          return data['count'] as int;
        }
        return 0;
      },
    );
  }

  static Future<MembersGroup> getGroupByName(String name) {
    return _makeRequest(
      '/groups/name/$name',
      parser: (data) => MembersGroup.fromJson(data),
    );
  }

  static Future<MembersGroup> createGroup(Map<String, dynamic> data) {
    return _makeRequest(
      '/groups',
      method: 'POST',
      body: data,
      parser: (data) => MembersGroup.fromJson(data),
    );
  }

  static Future<MembersGroup> updateGroup(String id, Map<String, dynamic> data) {
    return _makeRequest(
      '/groups/$id',
      method: 'PUT',
      body: data,
      parser: (data) => MembersGroup.fromJson(data),
    );
  }

  static Future<void> deleteGroup(String id) {
    return _makeRequest(
      '/groups/$id',
      method: 'DELETE',
      parser: (data) => null,
    );
  }
}

class ApiException implements Exception {
  final String message;
  final int statusCode;

  ApiException(this.message, this.statusCode);

  @override
  String toString() => 'ApiException: $message (Status: $statusCode)';
}

// Models
class Transaction {
  final String id;
  final String? transactionRef;
  final DateTime transactionDate;
  final String transactionType;
  final double amount;
  final String? notes;
  final String debitAccountId;
  final String? memberId;
  final String createdBy;
  final DateTime createdAt;
  final DateTime updatedAt;

  Transaction({
    required this.id,
    this.transactionRef,
    required this.transactionDate,
    required this.transactionType,
    required this.amount,
    this.notes,
    required this.debitAccountId,
    this.memberId,
    required this.createdBy,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Transaction.fromJson(Map<String, dynamic> json) {
    return Transaction(
      id: json['id'] ?? '',
      transactionRef: json['transaction_ref'],
      transactionDate: DateTime.tryParse(json['transaction_date'] ?? '') ?? DateTime.now(),
      transactionType: json['transaction_type'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      notes: json['notes'],
      debitAccountId: json['debit_account_id'] ?? '',
      memberId: json['member_id'],
      createdBy: json['created_by'] ?? '',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'transaction_ref': transactionRef,
    'transaction_date': transactionDate.toIso8601String(),
    'transaction_type': transactionType,
    'amount': amount,
    'notes': notes,
    'debit_account_id': debitAccountId,
    'member_id': memberId,
  };
}

class Account {
  final String id;
  final String accountName;
  final String accountType;
  final double? localShare;
  final String? notes;
  final bool isActive;
  final String createdBy;
  final DateTime createdAt;
  final DateTime updatedAt;

  Account({
    required this.id,
    required this.accountName,
    required this.accountType,
    this.localShare,
    this.notes,
    required this.isActive,
    required this.createdBy,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Account.fromJson(Map<String, dynamic> json) {
    return Account(
      id: json['id'] ?? '',
      accountName: json['account_name'] ?? '',
      accountType: json['account_type'] ?? '',
      localShare: (json['local_share'])?.toDouble(),
      notes: json['notes'],
      isActive: json['is_active'] ?? true,
      createdBy: json['created_by'] ?? '',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'account_name': accountName,
    'account_type': accountType,
    'local_share': localShare,
    'notes': notes,
  };
}

class Member {
  final String id;
  final String fullName;
  final String phoneNumber;
  final String? email;
  final String? notes;
  final String? groupId;
  final DateTime createdAt;
  final DateTime updatedAt;

  Member({
    required this.id,
    required this.fullName,
    required this.phoneNumber,
    this.email,
    this.notes,
    this.groupId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Member.fromJson(Map<String, dynamic> json) {
    return Member(
      id: json['id'] ?? '',
      fullName: json['full_name'] ?? '',
      phoneNumber: json['phone_number'] ?? '',
      email: json['email'],
      notes: json['notes'],
      groupId: json['group_id'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'full_name': fullName,
    'phone_number': phoneNumber,
    'email': email,
    'notes': notes,
    'group_id': groupId,
  };
}

class MembersGroup {
  final String id;
  final String groupName;
  final String? description;
  final DateTime createdAt;
  final DateTime updatedAt;

  MembersGroup({
    required this.id,
    required this.groupName,
    this.description,
    required this.createdAt,
    required this.updatedAt,
  });

  factory MembersGroup.fromJson(Map<String, dynamic> json) {
    return MembersGroup(
      id: json['id'] ?? '',
      groupName: json['group_name'] ?? '',
      description: json['description'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() => {
    'group_name': groupName,
    'description': description,
  };
}
