import 'package:portal/src/stored/json.dart';

class Stored<T extends JSONSerializable> extends JSONSerializable {
  static const idJSONKey = "id";
  static const createdByJSONKey = "createdBy";
  static const createdAtJSONKey = "createdAt";
  static const modifiedByJSONKey = "modifiedBy";
  static const modifiedAtJSONKey = "modifiedAt";
  static const contentJSONKey = "content";

  final String id;
  final String createdBy;
  final DateTime createdAt;
  final String modifiedBy;
  final DateTime modifiedAt;
  final T content;

  Stored({
    required this.id,
    required this.createdBy,
    required this.createdAt,
    required this.modifiedBy,
    required this.modifiedAt,
    required this.content,
  });

  Stored.fromJSON(
    T Function(Map<String, dynamic>) contentParser,
    Map<String, dynamic> json,
  )   : id = json[idJSONKey],
        createdBy = json[createdByJSONKey],
        createdAt = DateTime.parse(json[createdAtJSONKey]),
        modifiedBy = json[modifiedByJSONKey],
        modifiedAt = DateTime.parse(json[modifiedAtJSONKey]),
        content = contentParser(json[contentJSONKey]);

  Stored.copy(Stored<T> original,
      {String? id,
      String? createdBy,
      DateTime? createdAt,
      String? modifiedBy,
      DateTime? modifiedAt,
      T? content})
      : id = id ?? original.id,
        createdBy = createdBy ?? original.createdBy,
        createdAt = createdAt ?? original.createdAt,
        modifiedBy = modifiedBy ?? original.modifiedBy,
        modifiedAt = modifiedAt ?? original.modifiedAt,
        content = content ?? original.content;

  @override
  Map<String, dynamic> toJSON() {
    return <String, dynamic>{
      Stored.idJSONKey: id,
      Stored.createdByJSONKey: createdBy,
      Stored.createdAtJSONKey: createdAt.toIso8601String(),
      Stored.modifiedByJSONKey: modifiedBy,
      Stored.modifiedAtJSONKey: modifiedAt.toIso8601String(),
      Stored.contentJSONKey: content.toJSON(),
    };
  }

  @override
  bool operator ==(Object other) =>
      other is Stored<T> &&
      other.runtimeType == runtimeType &&
      other.id == id &&
      other.createdBy == createdBy &&
      other.createdAt == createdAt &&
      other.modifiedBy == modifiedBy &&
      other.modifiedAt == modifiedAt &&
      other.content == content;

  @override
  int get hashCode =>
      Object.hash(id, createdBy, createdAt, modifiedBy, modifiedAt, content);
}
