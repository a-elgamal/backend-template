import 'package:flutter_test/flutter_test.dart';
import 'package:portal/src/stored/json.dart';
import 'package:portal/src/stored/stored.dart';

class Simple extends JSONSerializable {
  final String body;

  const Simple(this.body);

  Simple.fromJSON(Map<String, dynamic> json) : body = json["body"];

  @override
  Map<String, dynamic> toJSON() => {"body": body};

  @override
  bool operator ==(Object other) =>
      other is Simple && other.runtimeType == runtimeType && other.body == body;

  @override
  int get hashCode => body.hashCode;
}

void main() {
  group(
    "A Stored should",
    () {
      final expected = Stored(
        id: "1",
        createdBy: "admin@kyny.com",
        createdAt: DateTime.now().add(const Duration(days: -1)),
        modifiedBy: "admin2@kynzy.com",
        modifiedAt: DateTime.now(),
        content: const Simple("body"),
      );

      test("parses JSON content successfully", () {
        final json = <String, dynamic>{
          Stored.idJSONKey: expected.id,
          Stored.createdByJSONKey: expected.createdBy,
          Stored.createdAtJSONKey: expected.createdAt.toIso8601String(),
          Stored.modifiedByJSONKey: expected.modifiedBy,
          Stored.modifiedAtJSONKey: expected.modifiedAt.toIso8601String(),
          Stored.contentJSONKey: <String, dynamic>{
            "body": expected.content.body,
          }
        };
        Stored<Simple> parsed = Stored.fromJSON(Simple.fromJSON, json);
        expect(parsed, equals(expected));
      });

      test("converts to JSON and reads it back", () {
        Stored<Simple> parsed =
            Stored.fromJSON(Simple.fromJSON, expected.toJSON());
        expect(parsed, equals(expected));
      });

      test("copied objects are equal", () {
        final copy = Stored<Simple>.copy(expected);
        expect(copy, equals(expected));
        expect(copy.hashCode, equals(expected.hashCode));
      });
    },
  );
}
