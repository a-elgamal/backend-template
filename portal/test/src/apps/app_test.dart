import 'package:flutter_test/flutter_test.dart';
import 'package:portal/src/apps/app.dart';

void main() {
  group(
    "An App should",
    () {
      const expected = App(
        apiKey: "1",
        disabled: true,
      );

      test("parses JSON content successfully", () {
        final json = <String, dynamic>{
          App.apiKeyJSONKey: expected.apiKey,
          App.disabledJSONKey: expected.disabled,
        };
        final parsed = App.fromJSON(json);
        expect(parsed, equals(expected));
      });

      test("converts to JSON and reads it back", () {
        App parsed = App.fromJSON(expected.toJSON());
        expect(parsed, equals(expected));
      });

      test("copied objects are equal", () {
        final copy = App.copy(expected);
        expect(copy, equals(expected));
        expect(copy.hashCode, equals(expected.hashCode));
      });
    },
  );
}
