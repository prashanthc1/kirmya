import 'dart:async';
import 'dart:convert';
import 'dart:io';

class Conversation {
  final String id;
  final String type;
  final String title;
  final List<String> participants;
  final String updatedAt;
  final int unreadCount;
  final String lastMessagePreview;
  final bool isPinned;
  final bool isArchived;

  Conversation({
    required this.id,
    required this.type,
    required this.title,
    required this.participants,
    required this.updatedAt,
    required this.unreadCount,
    required this.lastMessagePreview,
    required this.isPinned,
    required this.isArchived,
  });

  factory Conversation.fromJson(Map<String, dynamic> json) {
    return Conversation(
      id: json['id'] ?? '',
      type: json['type'] ?? '',
      title: json['title'] ?? '',
      participants: List<String>.from(json['participants'] ?? []),
      updatedAt: json['updated_at'] ?? '',
      unreadCount: json['unread_count'] ?? 0,
      lastMessagePreview: json['last_message_preview'] ?? '',
      isPinned: json['is_pinned'] ?? false,
      isArchived: json['is_archived'] ?? false,
    );
  }
}

class Message {
  final String id;
  final String senderId;
  final String body;
  final String contentType;
  final String createdAt;
  final String? editedAt;
  final String? deletedAt;

  Message({
    required this.id,
    required this.senderId,
    required this.body,
    required this.contentType,
    required this.createdAt,
    this.editedAt,
    this.deletedAt,
  });

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: json['id'] ?? '',
      senderId: json['sender_id'] ?? '',
      body: json['body'] ?? '',
      contentType: json['content_type'] ?? 'text',
      createdAt: json['created_at'] ?? '',
      editedAt: json['edited_at'],
      deletedAt: json['deleted_at'],
    );
  }
}

class WSEvent {
  final String kind;
  final String conversationId;
  final String? senderId;
  final String? readerId;
  final String? id;
  final String? body;
  final String? contentType;
  final String? createdAt;
  final String? at;

  WSEvent({
    required this.kind,
    required this.conversationId,
    this.senderId,
    this.readerId,
    this.id,
    this.body,
    this.contentType,
    this.createdAt,
    this.at,
  });

  factory WSEvent.fromJson(Map<String, dynamic> json) {
    return WSEvent(
      kind: json['kind'] ?? '',
      conversationId: json['conversation_id'] ?? '',
      senderId: json['sender_id'],
      readerId: json['reader_id'],
      id: json['id'],
      body: json['body'],
      contentType: json['content_type'],
      createdAt: json['created_at'],
      at: json['at'],
    );
  }
}

typedef WSEventCallback = void Function(WSEvent event);

class KirmyaMessagingClient {
  final String apiBaseUrl; // e.g. "http://localhost:8080/api/v1"
  String _accessToken;
  final HttpClient _httpClient = HttpClient();

  // WebSocket State
  WebSocket? _webSocket;
  bool _shouldReconnect = false;
  int _reconnectDelayMs = 1000;
  Timer? _reconnectTimer;
  Timer? _pingTimer;

  // Stream controller for live events
  final StreamController<WSEvent> _eventStreamController = StreamController<WSEvent>.broadcast();

  KirmyaMessagingClient({
    required this.apiBaseUrl,
    required String accessToken,
  }) : _accessToken = accessToken;

  Stream<WSEvent> get events => _eventStreamController.stream;

  void updateToken(String token) {
    _accessToken = token;
    if (_webSocket != null) {
      // Reconnect with new token
      reconnect();
    }
  }

  // --- REST HTTP Client methods ---

  Future<Map<String, String>> _headers() async {
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer $_accessToken',
    };
  }

  Future<dynamic> _get(String path) async {
    final uri = Uri.parse('$apiBaseUrl$path');
    final request = await _httpClient.getUrl(uri);
    (await _headers()).forEach((k, v) => request.headers.set(k, v));
    final response = await request.close();
    return _parseResponse(response);
  }

  Future<dynamic> _post(String path, Map<String, dynamic> body) async {
    final uri = Uri.parse('$apiBaseUrl$path');
    final request = await _httpClient.postUrl(uri);
    (await _headers()).forEach((k, v) => request.headers.set(k, v));
    request.write(jsonEncode(body));
    final response = await request.close();
    return _parseResponse(response);
  }

  Future<dynamic> _delete(String path) async {
    final uri = Uri.parse('$apiBaseUrl$path');
    final request = await _httpClient.deleteUrl(uri);
    (await _headers()).forEach((k, v) => request.headers.set(k, v));
    final response = await request.close();
    return _parseResponse(response);
  }

  Future<dynamic> _parseResponse(HttpClientResponse response) async {
    final body = await response.transform(utf8.decoder).join();
    if (response.statusCode >= 200 && response.statusCode < 300) {
      if (body.isEmpty) return null;
      final payload = jsonDecode(body);
      if (payload is Map && payload.containsKey('data')) {
        return payload['data'];
      }
      return payload;
    }
    throw HttpException('Request failed with status: ${response.statusCode}, body: $body');
  }

  // --- Messaging API ---

  Future<List<Conversation>> getConversations() async {
    final data = await _get('/conversations');
    final list = data['conversations'] as List;
    return list.map((item) => Conversation.fromJson(item)).toList();
  }

  Future<Conversation> startConversation(List<String> participantIds, {String title = ""}) async {
    final data = await _post('/conversations', {
      'participant_ids': participantIds,
      'title': title,
    });
    return Conversation.fromJson(data);
  }

  Future<List<Message>> getMessages(String conversationId, {String? searchQuery}) async {
    var path = '/conversations/$conversationId/messages';
    if (searchQuery != null && searchQuery.isNotEmpty) {
      path += '?q=${Uri.encodeComponent(searchQuery)}';
    }
    final data = await _get(path);
    final list = data['messages'] as List;
    return list.map((item) => Message.fromJson(item)).toList();
  }

  Future<Message> sendMessage(String conversationId, String body, {String contentType = 'text'}) async {
    final data = await _post('/conversations/$conversationId/messages', {
      'body': body,
      'content_type': contentType,
    });
    return Message.fromJson(data);
  }

  Future<bool> deleteMessage(String conversationId, String messageId) async {
    final data = await _delete('/conversations/$conversationId/messages/$messageId');
    return data != null && data['deleted'] == true;
  }

  Future<bool> markRead(String conversationId) async {
    final data = await _post('/conversations/$conversationId/read', {});
    return data != null && data['read'] == true;
  }

  Future<bool> sendTypingIndicator(String conversationId) async {
    final data = await _post('/conversations/$conversationId/typing', {});
    return data != null && data['ok'] == true;
  }

  Future<bool> archiveConversation(String conversationId, bool archive) async {
    final data = await _post('/conversations/$conversationId/archive', {'archive': archive});
    return data != null && data['archived'] == archive;
  }

  Future<bool> pinConversation(String conversationId, bool pin) async {
    final data = await _post('/conversations/$conversationId/pin', {'pin': pin});
    return data != null && data['pinned'] == pin;
  }

  // --- WebSocket Connection & Lifecycle ---

  void connectWebSocket() async {
    _shouldReconnect = true;
    final wsUrl = apiBaseUrl.replaceFirst('http', 'ws') + '/ws?token=${Uri.encodeComponent(_accessToken)}';

    try {
      print('[ws] Mobile client connecting to $wsUrl');
      _webSocket = await WebSocket.connect(wsUrl);

      print('[ws] Mobile client connected');
      _reconnectDelayMs = 1000; // reset backoff
      _startPing();

      _webSocket!.listen(
        (data) {
          _handleWSMessage(data);
        },
        onError: (err) {
          print('[ws] Mobile client error: $err');
          _reconnectLater();
        },
        onDone: () {
          print('[ws] Mobile client connection closed');
          _stopPing();
          if (_shouldReconnect) {
            _reconnectLater();
          }
        },
        cancelOnError: true,
      );
    } catch (e) {
      print('[ws] Mobile client connection failed: $e');
      _reconnectLater();
    }
  }

  void _handleWSMessage(dynamic data) {
    try {
      final payload = jsonDecode(data);
      if (payload['type'] == 'ping' || payload['type'] == 'pong') {
        return;
      }
      final ev = WSEvent.fromJson(payload);
      _eventStreamController.add(ev);
    } catch (e) {
      print('[ws] Error parsing incoming WS message: $e');
    }
  }

  void sendTypingWS(String conversationId) {
    if (_webSocket != null && _webSocket!.readyState == WebSocket.open) {
      _webSocket!.add(jsonEncode({
        'type': 'typing',
        'conversation_id': conversationId,
      }));
    }
  }

  void sendPingWS() {
    if (_webSocket != null && _webSocket!.readyState == WebSocket.open) {
      _webSocket!.add(jsonEncode({
        'type': 'ping',
      }));
    }
  }

  void reconnect() {
    close();
    connectWebSocket();
  }

  void close() {
    _shouldReconnect = false;
    _stopPing();
    if (_reconnectTimer != null) {
      _reconnectTimer!.cancel();
      _reconnectTimer = null;
    }
    if (_webSocket != null) {
      _webSocket!.close();
      _webSocket = null;
    }
  }

  void _reconnectLater() {
    if (_reconnectTimer != null) return;

    print('[ws] Mobile client reconnecting in $_reconnectDelayMs ms');
    _reconnectTimer = Timer(Duration(milliseconds: _reconnectDelayMs), () {
      _reconnectTimer = null;
      connectWebSocket();
    });

    _reconnectDelayMs = (_reconnectDelayMs * 2).clamp(1000, 30000);
  }

  void _startPing() {
    _stopPing();
    _pingTimer = Timer.periodic(Duration(seconds: 30), (timer) {
      sendPingWS();
    });
  }

  void _stopPing() {
    if (_pingTimer != null) {
      _pingTimer!.cancel();
      _pingTimer = null;
    }
  }
}
