import 'package:flutter/foundation.dart';
import '../../data/chat_repository.dart';
import '../../domain/chat.dart';
import '../../domain/message.dart';

class ChatProvider extends ChangeNotifier {
  final ChatRepository _repo;
  ChatProvider(this._repo);

  bool _loading = false;
  bool get loading => _loading;

  List<Chat> _chats = [];
  List<Chat> get chats => _chats;

  final Map<String, List<Message>> _messages = {};
  List<Message> messagesFor(String chatId) =>
      List.unmodifiable(_messages[chatId] ?? const []);

  Future<void> load() async {
    _loading = true;
    notifyListeners();
    _chats = await _repo.fetchChats();
    for (final c in _chats) {
      _messages[c.id] = await _repo.fetchMessages(c.id);
    }
    _loading = false;
    notifyListeners();
  }

  Future<void> send(String chatId, String text) async {
    if (text.trim().isEmpty) return;
    final msg = await _repo.sendMessage(chatId, text);
    _messages.putIfAbsent(chatId, () => []);
    _messages[chatId]!.add(msg);
    notifyListeners();

    // Через ~1с в репозитории появится "эхо" — подтянем заново чат
    await Future<void>.delayed(const Duration(milliseconds: 800));
    _messages[chatId] = await _repo.fetchMessages(chatId);
    _chats = await _repo.fetchChats();
    notifyListeners();
  }
}
