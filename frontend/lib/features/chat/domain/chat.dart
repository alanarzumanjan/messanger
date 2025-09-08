class Chat {
  final String id;
  final String title;
  String lastMessage;
  DateTime updatedAt;

  Chat({
    required this.id,
    required this.title,
    required this.lastMessage,
    required this.updatedAt,
  });
}
